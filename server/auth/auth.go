// package to hold information and methods for
// authentication
package auth

// represents the auth.yml config structure
// which only consists of one array of authenticated
// users.
type Config struct {
	Users []*User `mapstructure:"users" json:"users"`
}

// user represents a client connecting with a token
// to this server.
type User struct {
	Name    string   `mapstructure:"name" json:"name"`
	Token   string   `mapstructure:"token" json:"token"`
	Types   []string `mapstructure:"types" json:"types"`
	Size    int64    `mapstructure:"size" json:"size"`
	Expire  int64    `mapstructure:"expire" json:"expire"`
	FileKey string   `mapstructure:"fileKey" json:"fileKey"`
}

var (
	tokenMap map[string]*User
)

// setups the tokenMap by taking all
// auth users defined in the config and puts
// them into it.
// This makes requests faster.
func (c *Config) CollectUsers() {
	tokenMap = make(map[string]*User)
	for _, u := range c.Users {
		tokenMap[u.Token] = u
	}
}

// returns an user with given token or nil
// if no such user exists.
func (c *Config) GetUser(token string) *User {
	return tokenMap[token]
}
