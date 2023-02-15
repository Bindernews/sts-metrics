package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	stms "github.com/bindernews/sts-msr"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %+v", err)
	}
	r := gin.Default()
	controller, err := stms.NewStsMs()
	if err != nil {
		log.Fatal(err)
	}
	controller.Register(&r.RouterGroup)
	r.Run()
}
