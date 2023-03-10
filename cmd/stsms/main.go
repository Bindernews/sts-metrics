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

	server := &http.Server{Addr: srv.Config.Listen, Handler: r}
	go safeServe(server)

	adminRouter := gin.Default()
	adminServer := &http.Server{Addr: srv.Config.AdminListen, Handler: adminRouter}
	adminRouter.POST("/stop", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
		ctx := context.Background()
		server.Shutdown(ctx)
		adminServer.Shutdown(ctx)
	})
	safeServe(adminServer)
}

func safeServe(server *http.Server) {
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
	ctrlMain := web.MainController{Srv: srv}
	if err = ctrlMain.Init(r); err != nil {
		return
	}
	return
}
