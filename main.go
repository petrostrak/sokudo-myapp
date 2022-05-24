// sudo lsof -i -P -n | grep LISTEN
// to see running listened connections
package main

import (
	"github.com/petrostrak/sokudo"
	"myapp/data"
	"myapp/handlers"
	"myapp/middleware"
)

type application struct {
	App        *sokudo.Sokudo
	Handlers   *handlers.Handlers
	Models     data.Models
	Middleware *middleware.Middleware
}

func main() {
	s := initApplication()
	s.App.ListenAndServe()
}
