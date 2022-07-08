package main

import (
	"context"
	"fmt"
	"goweb/framework"
	"time"
)

func FooControllerHandler(ctx *framework.Context) error {
	durationCtx, cancel := context.WithTimeout(ctx.BaseContext(), 3*time.Second)
	defer cancel()
	finish := make(chan struct{}, 1)
	panicChan := make(chan interface{}, 1)
	go func() {
		if p := recover(); p != nil {
			panicChan <- p
		}
		time.Sleep(2 * time.Second)
		ctx.JsonAndCode(200, map[string]interface{}{"data": "test"})
		finish <- struct{}{}
	}()
	select {
	case <-durationCtx.Done():
		ctx.WriterMux().Lock()
		defer ctx.WriterMux().Unlock()
		ctx.JsonAndCode(500, "time out")
		ctx.SetHasTimeout()
	case <-finish:
		fmt.Println("finish")
	case p := <-panicChan:
		ctx.WriterMux().Lock()
		defer ctx.WriterMux().Unlock()
		fmt.Println(p)
		ctx.JsonAndCode(500, "panic")
	}
	return nil
}
