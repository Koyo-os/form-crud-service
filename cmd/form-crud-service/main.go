package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/Koyo-os/form-crud-service/internal/app"
	"github.com/Koyo-os/form-crud-service/internal/repo"
)

func getRepo() *repo.Repository {
	
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()



	
}