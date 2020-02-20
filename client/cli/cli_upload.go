package cli

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/spf13/cobra"
	"github.com/testy/lightf/client/context"
	"github.com/testy/lightf/internal/flags"
	"github.com/testy/lightf/pkg/slog"
)

func runUploadCmd(cmd *cobra.Command, args []string) {
	cfg := contextConfig // local copy

	if len(args) == 0 {
		slog.Fatalf("You have to specify a file to upload.")
	}
	path := args[0]

	server, err := cmd.Flags().GetString(flags.Context)
	if err != nil || server == "" {
		server = cfg.CurrentServer
	}
	ctx, err := cfg.GetContext(server)

	if ctx == nil {
		// we have to create a custom one with flags
		server, err := cmd.Flags().GetString(flags.Server)
		if err != nil || server == "" {
			slog.Fatal(fmt.Errorf("Without context you have to use --server <address>"))
		}
		token, err := cmd.Flags().GetString(flags.Token)
		if err != nil || token == "" {
			slog.Fatal(fmt.Errorf("Without context you have to use --token <token>"))
		}

		ctx = &context.Context{
			Address: server,
			Token:   token,
		}
	}

	if timeout, err := cmd.Flags().GetInt64(flags.Timeout); err != nil {
		http.DefaultClient.Timeout = time.Duration(timeout) * time.Second
	}
	res, err := postFile(http.DefaultClient, path, ctx)

	if os.IsNotExist(err) {
		slog.Fatalf("This file does not exist! path=%q", path)
	}
	if err != nil {
		slog.Fatal(err)
	}
	b, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	slog.Infof("Status: %d", res.StatusCode)
	slog.Infoln(string(b))
}

// post file to http server
func postFile(cl *http.Client, file string, ctx *context.Context) (*http.Response, error) {
	buf := &bytes.Buffer{}
	write := multipart.NewWriter(buf)

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	mime, err := mimetype.DetectReader(f)
	if err != nil {
		return nil, err
	}
	f.Seek(0, 0)

	form, err := createFormFile(write, "fileId", f.Name(), mime.String())
	if err != nil {
		return nil, err
	}
	cType := write.FormDataContentType()

	_, err = io.Copy(form, f)
	if err != nil {
		return nil, err
	}

	write.Close()
	req, err := http.NewRequest("POST", ctx.Address, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", cType)
	req.Header.Set("Token", ctx.Token)

	res, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// wrapper method for ".CreatePart".
// normally we could use CreateFormFile, but somehow the function
// always sets the type to "application/octet-stream"
func createFormFile(w *multipart.Writer, fieldname string, filename string, contenttype string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name=%q;  filename=%q`, fieldname, filename))
	h.Set("Content-Type", contenttype)
	return w.CreatePart(h)
}
