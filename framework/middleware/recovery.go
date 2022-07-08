package middleware

import (
	"fmt"
	"goweb/framework"
)

func Recovery() framework.ControllerHandler {
	return func(c *framework.Context) error {
		fmt.Println("middleware pre recovery")

		defer func() {
			if p := recover(); p != nil {
				fmt.Println("recovery panic")
				c.JsonAndCode(500, p)
			}
		}()

		c.Next()
		fmt.Println("middleware post recovery")
		return nil
	}
}
