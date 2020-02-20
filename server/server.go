package server

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/testy/lightf/internal/env"
	"github.com/testy/lightf/pkg/token"
	"github.com/testy/lightf/server/auth"
	"github.com/testy/lightf/server/ratelimit"
	"github.com/testy/lightf/server/route"
	"github.com/testy/lightf/server/storage"
)

var (
	// our router we use for mux
	// and creating our webserver
	router *gin.Engine

	// the config for authorized users
	authConfig *auth.Config

	// our config provider - which is in this case viper
	configProvider *viper.Viper

	// objects for upload
	gen   *token.Generator
	store *storage.Storager
)

func init() {
	// ================================
	// get viper to set default envs
	// ================================
	viper.AutomaticEnv()

	viper.SetDefault(env.VerboseLogging, false)

	viper.SetDefault(env.GinMode, "release")
	viper.SetDefault(env.GinHost, "0.0.0.0")
	viper.SetDefault(env.GinPort, "8081")
	viper.SetDefault(env.GinMaxEventsPerSecond, 1000)
	viper.SetDefault(env.GinMaxBurstSize, 20)

	viper.SetDefault(env.AuthConfigPath, "/etc/lightf")
	viper.SetDefault(env.StoragePath, "/var/lib/lightf/storage")
	viper.SetDefault(env.StorageArchiveOnStartup, false)
	viper.SetDefault(env.StorageArchiveAuto, true)
	viper.SetDefault(env.StorageArchiveAutoDelay, 60)
	viper.SetDefault(env.StorageAddress, fmt.Sprintf("%s:%d/f/%v", viper.GetString(env.GinHost), viper.GetInt(env.GinPort), "%s"))

	viper.SetDefault(env.NameGenMinLength, 3)
	viper.SetDefault(env.NameGenMinCharsetLength, 3)
	viper.SetDefault(env.NameGenLength, 7)
	viper.SetDefault(env.NameGenCharset, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	// ================================
	// get the auth.yml config
	// ================================
	configProvider = viper.New()
	configProvider.SetConfigName("auth")
	configProvider.SetConfigType("yml")

	configProvider.AddConfigPath(viper.GetString(env.AuthConfigPath))
	configProvider.AddConfigPath(".")

	if err := configProvider.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
		}
		// actually we can ignore both cases.
		// if no config is set, the user will run into a message
		// telling him, that no context could be found, and that he should
		// use the context flags. No error.
	}
	configProvider.Unmarshal(&authConfig)
	authConfig.CollectUsers()

	// ================================
	// startup webserver and add routes
	// ================================
	gin.SetMode(viper.GetString(env.GinMode))
	router = gin.New()
	router.Use(gin.Logger(), gin.Recovery(), ratelimit.Throttle(
		viper.GetInt(env.GinMaxEventsPerSecond),
		viper.GetInt(env.GinMaxBurstSize),
	))

	gen = token.G(viper.GetString(env.NameGenCharset), viper.GetInt(env.NameGenLength), &token.Settings{
		MinCharsetLength: viper.GetInt(env.NameGenMinCharsetLength),
		MinTokenLength:   viper.GetInt(env.NameGenMinLength),
	})

	spath := viper.GetString(env.StoragePath)
	store = storage.S(fmt.Sprintf("%s/files", spath), fmt.Sprintf("%s/archive", spath),
		viper.GetString(env.StorageAddress), &storage.Settings{
			ArchiveOnStartup: viper.GetBool(env.StorageArchiveOnStartup),
			ArchiveAuto:      viper.GetBool(env.StorageArchiveAuto),
			ArchiveAutoDelay: time.Duration(viper.GetInt(env.StorageArchiveAutoDelay)),
		})

	router.POST("/upload", route.Upload(authConfig, gen, store))
}

// starts the webserver, nothing else.
func Run() error {
	// TODO cleanup storage to archive (only if configured so)
	store.StartupArchiver()

	addr := fmt.Sprintf("%v:%d", viper.GetString(env.GinHost), viper.GetInt(env.GinPort))
	return router.Run(addr)
}
