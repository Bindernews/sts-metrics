package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	stms "github.com/bindernews/sts-msr"
	"github.com/bindernews/sts-msr/tonoauth"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %+v", err)
	}
	r := gin.Default()
	if err := setupControllers(r); err != nil {
		log.Fatalln(err)
	}
	r.Run("localhost:8080")
}

func setupControllers(r *gin.Engine) (err error) {
	srv := &stms.Services{}
	if err = srv.LoadDefaults(); err != nil {
		return
	}
	ctrlMain := stms.MainController{Srv: srv}
	if err = ctrlMain.Init(r); err != nil {
		return
	}
	// Setup the oauth controller, we can support multiple oauth
	ctrlOauth := tonoauth.NewOauthController(os.Getenv("BASE_URL"))
	ctrlOauth.AddProviders(
		tonoauth.NewGithubProvider(tonoauth.NewProviderOptsFromEnv("GH_")))
	if err = ctrlOauth.Init(r); err != nil {
		return
	}
	// Setup charts view
	ctrlCharts := stms.StatsView{Srv: srv}
	ctrlCharts.DefaultCharts()
	if err = ctrlCharts.Init(r.Group("/stats")); err != nil {
		return
	}
	return
}
