package wserver

import (
	"context"
	"net/http"
)

//needed a constructor for WebServer add a struct to struct
//and initialize it with NewWebServer (possibly including all handlers to it to store enpoints, verbs, handlers)
//Then .Start starts the server, stop stops it actually

type WebServer struct {
	httpServer *http.Server
}

//remake it so that it accepts all/most of the parameters
func (ws *WebServer) Start(sPort string, vHandler map[string]http.Handler) error {
	ws.httpServer = &http.Server{
		Addr:    ":" + sPort,
		Handler: nil,
	}
	for key, value := range vHandler {
		http.Handle(key, value)
	}
	return ws.httpServer.ListenAndServe()
}

func (ws *WebServer) Stop() error {
	return ws.httpServer.Shutdown(context.Background())
}
