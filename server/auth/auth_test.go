package auth

import "testing"

func TestCollectUsers(t *testing.T) {
	cfg := &Config{
		Users: []*User{
			&User{
				Token: "bar",
			},
		},
	}

	if cfg.GetUser("bar") != nil {
		t.Errorf("found user, expecting nil")
	}

	cfg.CollectUsers()
	if cfg.GetUser("bar") == nil {
		t.Errorf("got nil, expected user from token")
	}
}
