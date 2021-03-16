package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/new_web/api"
	"github.com/new_web/config"
	"github.com/new_web/pkg/middleware"
	"github.com/new_web/pkg/util"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const Version = "0.0.1"

func Run(c *cli.Context) error {
	configFile := c.String("config")
	conf, err := config.NewConfig(configFile)
	if err != nil {
		logrus.Fatal(err)
	}
	_, cancelFunc := context.WithCancel(context.Background())
	wgAll := sync.WaitGroup{}
	InitGlobal(conf)
	err = util.InitLoggger(conf.LogPath, conf.LogReserveDay)
	if err != nil {
		panic(err)
	}
	client, err := api.NewApiClient(conf)
	if err != nil {
		panic(err)
	}
	if conf.LaunchModule.LaunchApi {
		logrus.Info("launch Api services")
		server := gin.New()
		// 加载日志中间件
		server.Use(middleware.Logger())
		client.LoadRouter(server)
		srv := &http.Server{Addr: fmt.Sprintf("%s:%d", conf.ListenServer, conf.ListenPort), Handler: server}
		defer func() {
			if err := srv.Close(); err != nil {
				logrus.Fatal("server shutdown:", err)
			}
		}()
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logrus.Fatalf("listen: %s \n", err)
			}
		}()
	} else {
		logrus.Warn("Api Lauch disabled")
	}
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	logrus.Info("Wait Stop Signal...")
	<-quit
	cancelFunc()
	wgAll.Wait()
	logrus.Info("Shutdown Server ...")
	logrus.Info("Server exited")
	return nil
}

func main() {
	app := cli.App{
		Name:    "vue web service: sky",
		Version: Version,
		Action:  Run,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "config, c", Value: "../config/sky.conf", Usage: "Custom configuration file path"},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

//初始化全局配置
func InitGlobal(conf *config.Config) {
	config.ApiConf = conf
}
