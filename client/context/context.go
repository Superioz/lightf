// package to manage and configure
// contexts for our client.
// A context is a specific server with its
// authentication in a settings file.
package context

import "fmt"

// the representation of the context config
// as struct.
// holds multiple servers to be able to connect to them
type Config struct {
	Servers       []*Context `mapstructure:"servers" json:"servers"`
	CurrentServer string     `mapstructure:"current-server" json:"current-server"`
}

// one context entry for the Servers slice inside
// the struct.
type Context struct {
	Name    string `mapstructure:"name" json:"name"`
	Address string `mapstructure:"address" json:"address"`
	Token   string `mapstructure:"token" json:"token"`
}

// simply checks if any context is present.
func (c *Config) HasAnyContext() bool {
	return len(c.Servers) != 0
}

// checks if a given context is present.
func (c *Config) HasContext(ctx string) bool {
	for _, s := range c.Servers {
		if s.Name == ctx {
			return true
		}
	}
	return false
}

// returns the context by following the name
// given by "serv".
// If no current-server found, return error.
func (c *Config) GetContext(serv string) (*Context, error) {
	var ctx *Context
	for _, c := range c.Servers {
		if c.Name == serv {
			ctx = c
		}
	}
	if ctx == nil {
		return nil, fmt.Errorf("No context found for name %v", serv)
	}
	return ctx, nil
}

// returns the context by following the name
// given by "current-server" field.
// If no current-server found, return error.
func (c *Config) CurrentContext() (*Context, error) {
	name := c.CurrentServer

	return c.GetContext(name)
}

// sets the current server to given server
// and returns if the server has changed.
// That way, we can write to the config file only
// if neccessary
func (c *Config) SetCurrentServer(server string) bool {
	if c.CurrentServer != server {
		c.CurrentServer = server
		return true
	}
	return false
}
