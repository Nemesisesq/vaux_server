package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/nemesisesq/vaux_server/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
