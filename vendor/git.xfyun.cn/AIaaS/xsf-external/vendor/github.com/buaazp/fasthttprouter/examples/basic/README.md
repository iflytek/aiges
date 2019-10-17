# Example of FastHttpRouter

These examples show you the usage of `fasthttprouter`. You can easily build a web application with it. Or you can make your own midwares such as custom logger, metrics, or any one you want.

### Basic example

This is just a quick introduction, view the [GoDoc](http://godoc.org/github.com/buaazp/fasthttprouter) for details.

Let's start with a trivial example:

```go
package main

import (
	"fmt"
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

// Index is the index handler
func Index(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "Welcome!\n")
}

// Hello is the Hello handler
func Hello(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "hello, %s!\n", ctx.UserValue("name"))
}

// MultiParams is the multi params handler
func MultiParams(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "hi, %s, %s!\n", ctx.UserValue("name"), ctx.UserValue("word"))
}

// QueryArgs is used for uri query args test #11:
// if the req uri is /ping?name=foo, output: Pong! foo
// if the req uri is /piNg?name=foo, redirect to /ping, output: Pong!
func QueryArgs(ctx *fasthttp.RequestCtx) {
	name := ctx.QueryArgs().Peek("name")
	fmt.Fprintf(ctx, "Pong! %s\n", string(name))
}

func main() {
	router := fasthttprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)
	router.GET("/multi/:name/:word", MultiParams)
	router.GET("/ping", QueryArgs)

	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}
```
