package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bindernews/sts-msr/web"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const usage = `Usage of %s:
  -c, --config string
      Config file (default "%s")
`

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %+v", err)
	}

	cmd := MainCmd{}
	if err := cmd.setup(); err != nil {
		log.Fatalln(err)
	}

	server := &http.Server{Addr: cmd.srv.Config.Listen, Handler: cmd.r}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, os.Interrupt, syscall.SIGTERM)
	<-sigC
	server.Shutdown(context.Background())
}

type MainCmd struct {
	Config string
	r      *gin.Engine
	srv    *web.Services
}

func (m *MainCmd) setup() (err error) {
	flag.StringVar(&m.Config, "config", "config.toml", "")
	flag.StringVar(&m.Config, "c", "config.toml", "")
	flag.Usage = func() {
		fmt.Printf(usage,
			os.Args[0],
			flag.Lookup("config").DefValue,
		)
	}
	flag.Parse()

	m.srv = new(web.Services)
	m.r = gin.Default()
	if err = m.srv.LoadDefaults(); err != nil {
		return
	}
	if err = m.srv.Config.LoadFile(m.Config, false); err != nil {
		return
	}
	ctrlMain := web.MainController{Srv: m.srv}
	if err = ctrlMain.Init(m.r); err != nil {
		return
	}
	return
}
