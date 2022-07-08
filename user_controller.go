package main

import (
	"goweb/framework"
	"time"
)

func UserLoginController(c *framework.Context) error {
	time.Sleep(5 * time.Second)
	c.JsonAndCode(200, "ok, UserLoginController")
	return nil
}
