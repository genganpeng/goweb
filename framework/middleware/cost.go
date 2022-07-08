package middleware

import (
	"fmt"
	"goweb/framework"
	"time"
)

func Cost() framework.ControllerHandler {
	return func(c *framework.Context) error {
		startTime := time.Now()
		c.Next()
		fmt.Println(fmt.Printf("request uri: %v, cost time: %v", c.GetRequest().RequestURI, time.Now().Sub(startTime)))
		return nil
	}

}
