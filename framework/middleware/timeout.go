package middleware

import (
	"context"
	"fmt"
	"goweb/framework"
	"time"
)

func TimeoutHandler(d time.Duration) framework.ControllerHandler {
	return func(ctx *framework.Context) error {
		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)
		durationCtx, cancel := context.WithTimeout(ctx.BaseContext(), d)
		defer cancel()
		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()

			ctx.Next()
			finish <- struct{}{}
		}()
		select {
		case <-durationCtx.Done():
			ctx.JsonAndCode(500, "time out")
			ctx.SetHasTimeout()
		case <-finish:
			fmt.Println("finish")
		case p := <-panicChan:
			fmt.Println(p)
			ctx.JsonAndCode(500, "panic")
		}

		return nil
	}

}
