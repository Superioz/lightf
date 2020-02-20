package env

const (
	// default value for verbose logging, can be overriden
	// in client by command options.
	//
	// true = verbose logging enabled
	// false = ... disabled
	VerboseLogging = "LIGHTF_VERBOSE_LOGGING"

	// ======================================================================
	// client
	// ======================================================================

	// default value for shortened version in "version"
	// command.
	//
	// true = only short version "1.0.0"
	// false = long version "lightf v1.0.0"
	VersionShortened = "LIGHTF_VERSION_SHORTENED"

	// default settings for token generation
	// can be partially overriden with command options
	TokenGenMinLength        = "LIGHTF_TOKEN_MIN_LENGTH"
	TokenGenMinCharsetLength = "LIGHTF_TOKEN_MIN_CHARSET_LENGTH"
	TokenGenMinExpireLength  = "LIGHTF_TOKEN_MIN_EXPIRE_LENGTH"
	TokenGenDefaultLength    = "LIGHTF_TOKEN_DEFAULT_LENGTH"
	TokenGenDefaultCharset   = "LIGHTF_TOKEN_DEFAULT_CHARSET"

	// the path where to find the context.yml config
	// defaults to ".", current working directory
	ContextConfigPath = "LIGHTF_CONTEXT_CONF"

	UploadTimeoutInSeconds = "LIGHTF_UPLOAD_DEFAULT_TIMEOUT_IN_SECONDS"

	// ======================================================================
	// server
	// ======================================================================

	// the mode for the gin engine
	// release = in production
	// debug = during development
	GinMode = "LIGHTF_GIN_MODE"
	GinHost = "LIGHTF_GIN_HOST"
	GinPort = "LIGHTF_GIN_PORT"

	GinMaxEventsPerSecond = "LIGHTF_RATE_LIMIT"
	GinMaxBurstSize       = "LIGHTF_RATE_MAX_BURST_SIZE"

	// the path to the auth.yml config
	AuthConfigPath = "LIGHTF_AUTH_CONF"

	// the path to where the files will be stored
	// two folders will be created inside:
	// /files - for all files that get uploaded
	// /archive - for the archived files
	StoragePath             = "LIGHTF_STORAGE"
	StorageArchiveOnStartup = "LIGHTF_STORAGE_ARCHIVE_ON_STARTUP"
	StorageArchiveAuto      = "LIGHTF_STORAGE_ARCHIVE_AUTO"
	StorageArchiveAutoDelay = "LIGHTF_STORAGE_ARCHIVE_AUTO_DELAY"
	StorageAddress          = "LIGHTF_STORAGE_ADDRESS"

	// configuration for our file name generator
	NameGenMinLength        = "LIGHTF_NAME_MIN_LENGTH"
	NameGenMinCharsetLength = "LIGHTF_NAME_MIN_CHARSET_LENGTH"
	NameGenLength           = "LIGHTF_NAME_DEFAULT_LENGTH"
	NameGenCharset          = "LIGHTF_NAME_DEFAULT_CHARSET"
)
