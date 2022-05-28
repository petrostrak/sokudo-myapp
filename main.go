// sudo lsof -i -P -n | grep LISTEN
// to see running listened connections
package main

import (
	"myapp/data"
	"myapp/handlers"
	"myapp/middleware"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/petrostrak/sokudo"
)

type application struct {
	App        *sokudo.Sokudo
	Handlers   *handlers.Handlers
	Models     data.Models
	Middleware *middleware.Middleware
	wg         sync.WaitGroup
}

func main() {
	s := initApplication()
	go s.listenForShutdown()
	err := s.App.ListenAndServe()
	s.App.ErrorLog.Println(err)
}

func (a *application) shutdown() {
	// put any clean up task here

	// block until the WaitGroup is empty
	a.wg.Wait()
}

func (a *application) listenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	s := <-quit
	a.App.InfoLog.Println("Received signal", s.String())

	a.shutdown()

	os.Exit(0)
}
