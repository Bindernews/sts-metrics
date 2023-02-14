package main

import (
	"github.com/gin-gonic/gin"

	stms "github.com/bindernews/sts-metrics-server"
)

func main() {
	r := gin.Default()
	controller := stms.NewStsMs()
	controller.Register(&r.RouterGroup)
	r.Run()
}
