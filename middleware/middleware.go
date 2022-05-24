package middleware

import (
	"github.com/petrostrak/sokudo"
	"myapp/data"
)

type Middleware struct {
	App    *sokudo.Sokudo
	Models data.Models
}
