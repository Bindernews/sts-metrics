package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	stms "github.com/bindernews/sts-msr"
	"github.com/bindernews/sts-msr/tonoauth"
)

const usage = `Usage of %s:
  -c, --config string
      Config file (default "%s")
`

var (
	optConfig = flag.String("config", "config.toml", "")
)

func init() {
	flag.StringVar(optConfig, "c", "config.toml", "")

	flag.Usage = func() {
		fmt.Printf(usage,
			os.Args[0],
			flag.Lookup("config").DefValue,
		)
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %+v", err)
	}
	flag.Parse()

	srv := new(stms.Services)
	r := gin.Default()
	if err := setup(srv, r); err != nil {
		log.Fatalln(err)
	}
	if err := r.Run(srv.Config.Listen...); err != nil {
		log.Fatalln(err)
	}
}

func setup(srv *stms.Services, r *gin.Engine) (err error) {
	if err = srv.LoadDefaults(); err != nil {
		return
	}
	if err = srv.Config.LoadFile(*optConfig, false); err != nil {
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
