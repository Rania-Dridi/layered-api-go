package main

import (
	"log"
	"net/http"

	"github.com/Rania-Dridi/layered-api-go/internal/handler"
	"github.com/Rania-Dridi/layered-api-go/internal/model"
	"github.com/Rania-Dridi/layered-api-go/internal/repository"
	"github.com/Rania-Dridi/layered-api-go/internal/service"
)

func main() {
	repo := repository.NewMemory([]model.Course{
		{ID: 1, Title: "Introduction to SQL", Capacity: 2},
		{ID: 2, Title: "REST API design", Capacity: 30},
	})
	svc := service.New(repo)
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler.Router(svc)))
}
