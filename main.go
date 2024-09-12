package main

import (
	"git_control/config"
	"git_control/handlers"
	"git_control/logInterface"
	"git_control/router"
	"log/slog"
	"net/http"
)

func main() {
	conf := config.MustLoad()

	logger := logInterface.SetLogger(conf.Env)
	logger.Info("Север запущен с конфигом: ", slog.String("env", conf.Env))

	r := router.GetRouter()
	r.Get("/branches", handlers.GetList(logger))
	r.Get("/set/{branch}", handlers.SetBranch(logger))
	r.Get("/update", handlers.UpdateBranch(logger))

	server := &http.Server{
		Addr:         conf.HttpServer.Address,
		Handler:      r,
		ReadTimeout:  conf.HttpServer.Timeout,
		WriteTimeout: conf.HttpServer.Timeout,
		IdleTimeout:  conf.HttpServer.IdleTimeout,
	}

	err := server.ListenAndServe()
	if err != nil {
		logger.Error("server stopped with error: ", err)
		return
	}

}
