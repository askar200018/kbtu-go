package main

import (
	"context"
	"hw-6/internal/http"
	"hw-6/internal/store/inmemory"
	"log"
)

func main() {
	store := inmemory.NewDB()

	srv := http.NewServer(context.Background(), ":8080", store)
	if err := srv.Run(); err != nil {
		log.Println(err)
	}

	srv.WaitForGracefulTermination()
}
