package main

import (
	"github.com/ewgRa/gocsfixer/fixers"
	"fmt"
)

func NewCsFixerConfig(recommend, lint, fix bool, csFixer fixers.CsFixer) *CsFixerConfig {
	return &CsFixerConfig{recommend: recommend, lint: lint, fix: fix, csFixer: csFixer}
}

type CsFixerConfig struct {
	recommend bool
	lint bool
	fix bool
	csFixer fixers.CsFixer
}

func (c *CsFixerConfig) Recommend() bool {
	return c.recommend
}

func (c *CsFixerConfig) Lint() bool {
	return c.lint
}

func (c *CsFixerConfig) Fix() bool {
		return c.fix
}

func (c *CsFixerConfig) String() string {
	return fmt.Sprintf("Config, fix: %v, csFixer: %v", c.fix, c.csFixer)
}
