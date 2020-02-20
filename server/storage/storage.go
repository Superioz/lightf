// package to manage the files storage
package storage

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/testy/lightf/pkg/slog"
)

// A storager is just a routine
// moving files from one place to another.
// In our case its purpose is to store the uploaded
// files and to archive old files eventually.
type Storager struct {
	FilesPath     string
	ArchivePath   string
	AddressFormat string

	Settings *Settings

	expireMap map[string]int64
	lock      *sync.Mutex
}

type Settings struct {
	ArchiveOnStartup bool
	ArchiveAuto      bool
	ArchiveAutoDelay time.Duration
}

// creates a new storager instance with given parameters.
func S(files string, archive string, addr string, settings *Settings) *Storager {
	// create paths and panic if something went wrong
	// I know duplicated code, but file handling in Go sucks
	if err := os.MkdirAll(files, os.ModePerm); err != nil {
		slog.Fatalf("could not create dir: %v", files, err)
	}
	if err := os.MkdirAll(archive, os.ModePerm); err != nil {
		slog.Fatalf("could not create dir: %v", archive, err)
	}

	return &Storager{
		FilesPath:     files,
		ArchivePath:   archive,
		AddressFormat: addr,
		Settings:      settings,
		lock:          &sync.Mutex{},
		expireMap:     make(map[string]int64),
	}
}

// takes the name and its content ("r") and creates
// a file in the /storage folder with its data.
// Uses a lock to not clash with the archiver.
func (s *Storager) Store(name string, expire int64, r io.Reader) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	// read data from reader
	dat, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("%s/%s", s.FilesPath, name)
	err = ioutil.WriteFile(path, dat, 0644)
	if err != nil {
		return "", err
	}
	s.expireMap[name] = time.Now().Unix() + expire
	return fmt.Sprintf(s.AddressFormat, name), nil
}

// Starts the archive routine.
func (s *Storager) StartupArchiver() {
	if s.Settings.ArchiveOnStartup {
		s.Archive()
	}

	if !s.Settings.ArchiveAuto {
		// archiver is not enabled, only on startup (eventually)
		return
	}
	go func(s *Storager) {
		for {
			time.Sleep(s.Settings.ArchiveAutoDelay * time.Second)
			s.Archive()
		}
	}(s)
}

// Takes old/expired files from the /files dir
// and moves them to the /archive directory.
func (s *Storager) Archive() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	expired := make([]string, 0)
	now := time.Now().Unix()

	filepath.Walk(s.FilesPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		name := info.Name()

		// check if either name is NOT present in map
		// or if he should EXPIRE
		expire, ok := s.expireMap[name]
		if expire == 0 || now >= expire {
			if ok {
				// delete from map if exists. it did expire, so ...
				delete(s.expireMap, name)
			}
			expired = append(expired, name)
		}
		return nil
	})

	for _, f := range expired {
		path := fmt.Sprintf("%s/%s", s.FilesPath, f)
		if _, err := os.Stat(path); err != nil {
			// file does not exist, just continue
			continue
		}

		err := os.Rename(path, fmt.Sprintf("%s/%s", s.ArchivePath, f))
		if err != nil {
			return err
		}
	}
	return nil
}
