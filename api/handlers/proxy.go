package handlers

import (
	"encoding/json"

	middleware "github.com/AldoRuizP/proxy-app/api/middlewares"
	"github.com/kataras/iris"
)

// HandleRedirection should redirect traffic
func HandleRedirection(app *iris.Application) {
	app.Get("/ping", middleware.ProxyMiddleware, proxyHandler)
}

func proxyHandler(c iris.Context) {
	res, err := json.Marshal(middleware.Que)
	if err != nil {
		c.JSON(iris.Map{"status": 400, "result": "parse error"})
		return
	}
	parsed := string(res)
	c.JSON(iris.Map{"result": parsed})
}
