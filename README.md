### WEB控制台

web-console库实现一个自包含前端页面的DB组件控制台，引入者导入库后，结合具体的DB类型进行实例操作后，结合具体使用的web框架，将组件暴露的handler注入到路由框架中

### Supported DB
---
* MySQL
* Redis
* MongoDB (TODO)

### Getting started
---
### MySQL控制台接入

##### 1.download module

```
go get git.qpaas.com/go-components/webconsole@v1.0.0
```
##### 2.import and use
   * Gin框架集成
   ---
```
package main

import (
	"fmt"

	"git.qpaas.com/go-components/webconsole/pkg/common"
	"git.qpaas.com/go-components/webconsole/pkg/console"
	"github.com/gin-gonic/gin"
)

func main() {
	// initial mysql console
	mysqlConsole := console.NewMySQLConsole()

	// initial router engine and registry console route
	r := gin.Default()

	// user set console route path. eg: "/console/mysql"
	consoleRoute := r.Group("/console") 
	{
		// *any: 采用Gin的模糊匹配
		consoleRoute.Any("/mysql/*any", func(ctx *gin.Context) {
			// user custom business logic
			// ...

			// 基础用法
			// construct handler options object
			opt := &common.HandlerOptions{
				// set connect configure (required)
				Conn: common.ConnConfig{
					IP:       "127.0.0.01",
					Port:     3306,
					UserName: "root",
					Password: "root",
				},
			}

			// 进阶用法
			// construct handler options object
			opt := &common.HandlerOptions{
				// set connect configure (required)
				Conn: common.ConnConfig{
					IP:       "127.0.0.01",
					Port:     3306,
					UserName: "root",
					Password: "root",
				},

				// set query option (optional)
				QueryOpt: common.QueryOptions{
					Timeout: 15, // sql execute timeout (unit: s), default 15s
				},

				// IsIgnoreSystemIntercept
				// 是否启用系统内置的拦截器
				// 默认是开启，开启后系统会根据内置的命令白名单(只读命令)进行拦截
				// 当系统拦截器不满足业务需求时，用户可以关闭，通过QueryBeforeHook钩子实现自定义的
				// 拦截逻辑
				IsIgnoreSystemIntercept: false

				// set SQL Type allow execute (optional)
				AllowSQLType: nil, // default readOnly SQL, eg: select|use|show|explain

				// callback func execute before query (optional)
				// if before callback return error
				// query will aborted
				QueryBeforeHook: func(pha *common.PrevHookArgs) error {
					fmt.Printf("schema: %s; sql:%s\n", pha.Schema, pha.SQL)
					return nil
				},

				// callback func execute after query (optional)
				QueryAfterHook: func(pha *common.PostHookArgs) {
					fmt.Printf("schema:%s; sql:%s; queryTime: %d; err: %s\n",
						pha.Schema, pha.SQL, pha.QueryDuration, pha.Err)
				},
			}

			// registry console handler
			console.Handler(ctx.Writer, ctx.Request, "/console/mysql", mysqlConsole, opt)
		})
	}

	fmt.Println("server listen on port 9099")
	r.Run(":9099")
}
```
   * Iris框架集成
   因Iris框架对于模糊匹配的支持没有Gin友好,相比于集成到Gin,集成到Iris的步骤略为繁琐
   ---
```
package main

import (
	"fmt"

	"github.com/kataras/iris/v12/context"

	"git.qpaas.com/go-components/webconsole/pkg/common"
	"git.qpaas.com/go-components/webconsole/pkg/console"
	"github.com/kataras/iris/v12"
)

func main() {
	// initial mysql console
	mysqlConsole := console.NewMySQLConsole()

	// initial router engine and registry console route
	irisR := iris.Default()

	// user set console route path. eg: "/console/mysql/"
	// the end "/" is not ignore
	consoleRoute := irisR.Party("/console")

	// 注册处理console首页的handler
	consoleRoute.Get("/mysql/", func(ctx *context.Context) {
		console.HandlerStaticFile(ctx.ResponseWriter(), ctx.Request(), "/console/mysql", mysqlConsole)
	})

	// 注册处理console其他静态资源的handler eg: css, js, img and so on
	consoleRoute.Get("/mysql/{any:path}", func(ctx *context.Context) {
		console.HandlerStaticFile(ctx.ResponseWriter(), ctx.Request(), "/console/mysql", mysqlConsole)
	})

	// 注册处理console控制台API请求的handler
	consoleRoute.Post("/mysql/", func(ctx *context.Context) {
		// user custom business logic
		// ...

		// 基础用法
		// construct handler options object
		opt := &common.HandlerOptions{
			// set connect configure (required)
			Conn: common.ConnConfig{
				IP:       "127.0.0.1",
				Port:     3306,
				UserName: "root",
				Password: "root",
			},
		}

		// 进阶用法
		// construct handler options object
		opt := &common.HandlerOptions{
			// set connect configure (required)
			Conn: common.ConnConfig{
				IP:       "127.0.0.1",
				Port:     3306,
				UserName: "root",
				Password: "root",
			},

			// set query option (optional)
			QueryOpt: common.QueryOptions{
				Timeout: 15, // sql execute timeout (unit: s), default 15s
			},

			// set SQL Type allow execute (optional)
			AllowSQLType: nil, // default readOnly SQL, eg: select|use|show|explain

			// IsIgnoreSystemIntercept
			// 是否启用系统内置的拦截器
			// 默认是开启，开启后系统会根据内置的命令白名单(只读命令)进行拦截
			// 当系统拦截器不满足业务需求时，用户可以关闭，通过QueryBeforeHook钩子实现自定义的
			// 拦截逻辑
			IsIgnoreSystemIntercept: false

			// callback func execute before query (optional)
			// if before callback return error
			// query will aborted
			QueryBeforeHook: func(pha *common.PrevHookArgs) error {
				fmt.Printf("schema: %s; sql:%s\n", pha.Schema, pha.SQL)
				return nil
			},

			// callback func execute after query (optional)
			QueryAfterHook: func(pha *common.PostHookArgs) {
				fmt.Printf("schema:%s; sql:%s; queryTime: %d; err: %s\n",
					pha.Schema, pha.SQL, pha.QueryDuration, pha.Err)
			},
		}

		console.HandlerAPI(ctx.ResponseWriter(), ctx.Request(), mysqlConsole, opt)
	})

	fmt.Println("iris engine listen on 9099")
	irisR.Run(iris.Addr(":9099"), iris.WithoutPathCorrectionRedirection)
}
```
##### 3.access web cosole
```
浏览器键入地址：http://hostIP:9099/console/mysql/
```
![web console page preview](/web-console-demo.png)

### Lib some default action explian
---
* MySQL
  * if SQL is empty， `desc table` as default SQL.
  * if Select SQL is not set limit, lib will append limit 100 to sql to avoid query set too big.
  * if AllowSQLType not set, lib will use default white list(Select、Show、Explain、Desc) for sql valid.
  * if Is IsIgnoreSystemIntercept set true, lib will not check sql and use can set custom check login in beforeQuery hook func.
  * if sql execute timeout is not set, 15 seconds will be set.
  * user can use before query hook to complete some business logic, if before hook func return error, query will be aborted not continue. 


