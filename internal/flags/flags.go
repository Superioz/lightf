// simple package to hold all constant values
// for flags.
// Normally we'd only use the flags in one package
// but as we have client AND server in one project,
// we have to supply them centralized somehow.
package flags

const (
	Verbose = "verbose"
	Short   = "short"
	Length  = "length"
	Charset = "charset"
	Server  = "server"
	Context = "context"
	Token   = "token"
	Timeout = "timeout"

	V = "v"
	S = "s"
	L = "l"
	C = "c"
	T = "t"
)
