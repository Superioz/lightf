package context

import "testing"

func TestHasAnyContext(t *testing.T) {
	cfg := &Config{}
	if cfg.HasAnyContext() {
		t.Errorf("Config has context while initializing, whut.")
	}
	cfg.Servers = append(cfg.Servers, &Context{})
	if !cfg.HasAnyContext() {
		t.Errorf("Config does not have context after appending.")
	}
}

func TestHasContext(t *testing.T) {
	cfg := &Config{}
	if cfg.HasAnyContext() {
		t.Errorf("Config has context while initializing, whut.")
	}
	cfg.Servers = append(cfg.Servers, &Context{Name: "foo"})
	if !cfg.HasContext("foo") {
		t.Errorf("Config does not have context after appending.")
	}
}

func TestCurrentContext(t *testing.T) {
	cfg := &Config{
		CurrentServer: "foo",
	}
	_, err := cfg.CurrentContext()
	if err == nil {
		t.Errorf("Found context even though we did not add one.")
	}

	cfg.Servers = append(cfg.Servers, &Context{
		Name: "foo",
	})
	_, err = cfg.CurrentContext()
	if err != nil {
		t.Errorf("Did not found context even though we should find one.")
	}
}

func TestSetCurrentServer(t *testing.T) {
	cfg := &Config{
		CurrentServer: "foo",
	}
	res := cfg.SetCurrentServer("bar")
	if !res {
		t.Errorf("Could not detect a current-server change, expected one.")
	}

	if cfg.CurrentServer != "bar" {
		t.Errorf("Value of current-server did not change properly.")
	}
}
