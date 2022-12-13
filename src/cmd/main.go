package main

import (
	"context"
	"ehdw/smartiko-test/src/config"
	"ehdw/smartiko-test/src/db"
	"ehdw/smartiko-test/src/logic/queue"
	"ehdw/smartiko-test/src/logic/repo"
	"ehdw/smartiko-test/src/logic/server"
	"ehdw/smartiko-test/src/logic/service"
	"ehdw/smartiko-test/src/util"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	lg "github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	InitLogger()
	config.Init()

	//TODO use the pool
	q := queue.NewQueue()
	defer q.Stop()
	pdb := db.ConnectPGDatabase(ctx, true)
	svc := service.NewService(repo.NewRepository(pdb), q)
	r := mux.NewRouter()
	server.MakeServer(r, svc)
	util.CreateHealthCheck(r)
	go func() {
		srv := http.Server{
			Handler:      r,
			ReadTimeout:  time.Minute,
			WriteTimeout: time.Minute,
			Addr:         fmt.Sprintf("0.0.0.0:%d", config.Config().ServerPort),
		}
		lg.WithField("port", config.Config().ServerPort).Info("Listening on port")
		lg.Error(srv.ListenAndServe().Error())
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		lg.Info("Catched shutdown signal")
		cancel()
	}()

	lg.Info("Initialised smartiko-test server")
	<-ctx.Done()
}

func InitLogger() {
	lg.SetLevel(lg.DebugLevel)
}
