package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/danigunawan/rpg-api/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
