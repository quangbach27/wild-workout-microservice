package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/quangbach27/wild-workout-microservice/internal/common/logs"
	"github.com/quangbach27/wild-workout-microservice/internal/common/server"
	"github.com/quangbach27/wild-workout-microservice/internal/trainer/ports"
	"github.com/quangbach27/wild-workout-microservice/internal/trainer/service"
)

func main() {
	logs.Init()

	ctx := context.Background()

	application := service.NewApplication(ctx)

	serverType := strings.ToLower(os.Getenv("SERVER_TO_RUN"))
	switch serverType {
	case "http":
		server.RunHTTPServer(func(router chi.Router) http.Handler {
			return ports.HandlerFromMux(
				ports.NewHttpServer(application),
				router,
			)
		})
	default:
		panic(fmt.Sprintf("server type '%s' is not supported", serverType))
	}
}
