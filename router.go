package main

import (
	"goweb/framework"
	"goweb/framework/middleware"
	"time"
)

func registerRouter(core *framework.Core) {
	core.Use(middleware.Recovery(), middleware.Cost())
	core.Get("/", middleware.Test1(), middleware.Test2(), func(c *framework.Context) error {
		return c.JsonAndCode(200, "index page")
	})
	
	core.Get("/shut", func(c *framework.Context) error {
		foo, _ := c.QueryString("foo", "def")
		// 等待10s才结束执行
		time.Sleep(10 * time.Second)
		// 输出结果
		c.SetOkStatus().Json("ok, shut controller: " + foo)
		return nil
	})

	// 需求1+2:HTTP方法+静态路由匹配
	core.Get("/user/login", middleware.Test3(), func(c *framework.Context) error {
		time.Sleep(1 * time.Second)
		return c.JsonAndCode(200, "login")
	})

	// 需求3:批量通用前缀
	subjectApi := core.Group("/subject")
	{
		// 需求4:动态路由
		subjectApi.Delete("/:id", SubjectDelController)
		subjectApi.Put("/:id", SubjectUpdateController)
		subjectApi.Get("/:id", SubjectGetController)
		subjectApi.Get("/list/all", SubjectListController)

		subjectInnerApi := subjectApi.Group("/info")
		{
			subjectInnerApi.Get("/name", SubjectNameController)
		}
	}
}
