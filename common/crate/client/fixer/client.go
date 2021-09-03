package fixer

import (
	"xrate/common/crate"
	"xrate/config"
)

func newClient(cfg config.Client) *client {
	return &client{cfg: cfg}
}

func init() {
	crate.Register("fixer", newClient(config.GetClient()))
}
