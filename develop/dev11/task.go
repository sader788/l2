package main

import (
	"L2/develop/dev11/internal/event"
	"L2/develop/dev11/internal/validator"
	"L2/develop/dev11/pkg/handlers"
	"L2/develop/dev11/pkg/logger"
	"github.com/julienschmidt/httprouter"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	l := logger.NewLogger()

	validator, err := validator.NewValidator(&event.Event{})
	if err != nil {
		return //TODO
	}

	router := httprouter.New()
	handler := handlers.NewHandler(validator, l)
	handler.Register(router)

	err = startHttpServer(router)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func startHttpServer(r *httprouter.Router) error {
	//servAddr := cfg.HttpServerHost + ":" + strconv.Itoa(cfg.HttpServerPort)

	listener, err := net.Listen("tcp", "0.0.0.0:80")

	if err != nil {
		return err
	}

	server := http.Server{
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	//l.Info("httpserver: http server started " + servAddr)
	log.Fatalln(server.Serve(listener))
	return nil
}
