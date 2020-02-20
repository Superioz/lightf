package token

import (
	"fmt"
	"math/rand"
	"time"
)

// our generator object.
// we could simply use global functions, but we
// want to open up the possibility to pass around
// a object pointer.
//
// also: the fields are exported, as they are not really
// necessary to be protected.
type Generator struct {
	Charset string `json:"charset"`
	Length  int    `json:"length"`

	// the security settings for this generator
	settings *Settings

	// only for generation. can be edited via "Rand()"
	random *rand.Rand
}

// settings for our generator, which includes
// minimum values.
type Settings struct {
	// for security, we don't want less than x characters ...
	MinCharsetLength int

	// for security as well
	MinTokenLength int

	// minimum expiration length in seconds
	// <=0 excluded, meaning that
	// it only gets active, IF an expiration greater
	// than zero is set.
	MinExpireLength int64
}

// our New() method.
func G(charset string, length int, settings *Settings) *Generator {
	return &Generator{
		Charset:  charset,
		Length:   length,
		settings: settings,
	}
}

// sets the random instance for our generator
func (g *Generator) Rand(r *rand.Rand) {
	g.random = r
}

// executes the generation and creates a token.
// this algorithm is pretty simple and shouldn't be used
// for REALLY HIGH level of security.
// But for most purposes (random names, tokens) it is more than
// enough. (still depends on the length and the charset though)
func (g *Generator) DoExpire(expire int64) (*token, error) {
	if g.random == nil {
		// create own random instance
		g.random = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	l := len(g.Charset)
	if l < g.settings.MinCharsetLength {
		return nil, fmt.Errorf("Can't generate token with given charset, expected len >= %d, got %d",
			g.settings.MinCharsetLength, l)
	}
	if g.Length < g.settings.MinTokenLength {
		return nil, fmt.Errorf("Can't generate token with given length, expected len >= %d, got %d",
			g.settings.MinTokenLength, g.Length)
	}
	if expire > 0 && expire < g.settings.MinExpireLength {
		return nil, fmt.Errorf("Can't generate token with given expire, expected time >= %d, got %d",
			g.settings.MinExpireLength, expire)
	}

	// append random char from charset to string
	b := make([]byte, 0)
	for i := 0; i < g.Length; i++ {
		b = append(b, g.Charset[g.random.Intn(l)])
	}
	return &token{
		Str:        string(b),
		Creation:   time.Now().UnixNano(),
		Expiration: expire,
	}, nil
}

func (g *Generator) Do() (*token, error) {
	return g.DoExpire(0)
}
