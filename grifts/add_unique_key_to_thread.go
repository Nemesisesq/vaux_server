package grifts

import (
	. "github.com/markbates/grift/grift"
)

var _ = Desc("add_unique_key_to_thread", "Task Description")
var _ = Add("add_unique_key_to_thread", func(c *Context) error {
	return nil
})
