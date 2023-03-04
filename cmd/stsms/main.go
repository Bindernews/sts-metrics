package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bindernews/sts-msr/web"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	srv := new(web.Services)
	r := gin.Default()
	if err := setup(srv, r); err != nil {
		log.Fatalln(err)
	}

	server := &http.Server{
		Addr:    srv.Config.Listen,
		Handler: r,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()

	adminRouter := gin.Default()
	adminRouter.POST("/stop", func(c *gin.Context) {
		server.Shutdown(context.Background())
	})
	if err := adminRouter.Run(srv.Config.AdminListen); err != nil {
		log.Fatalln(err)
	}
}

func setup(srv *web.Services, r *gin.Engine) (err error) {
	if err = srv.LoadDefaults(); err != nil {
		return
	}
	if err = srv.Config.LoadFile(*optConfig, false); err != nil {
		return
	}
	if srv.Config.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	ctrlMain := web.MainController{Srv: srv}
	if err = ctrlMain.Init(r); err != nil {
		return
	}
	// Setup the oauth controller, we can support multiple oauth
	ctrlOauth := web.NewOauthController(os.Getenv("BASE_URL"))
	ctrlOauth.AddProviders(
		web.NewGithubProvider(web.NewProviderOptsFromEnv("GH_")))
	if err = ctrlOauth.Init(r); err != nil {
		return
	}
	// Setup charts view
	ctrlCharts := web.StatsView{Srv: srv}
	ctrlCharts.DefaultCharts()
	if err = ctrlCharts.Init(r.Group("/stats2")); err != nil {
		return
	}
	return
}
