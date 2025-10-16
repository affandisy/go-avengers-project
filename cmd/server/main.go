package main

import (
	"avenger/internal/handler"
	"avenger/internal/repository"
	"avenger/internal/service"
	"avenger/pkg/db"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("file .env not found")
	}

	conn := db.InitPostgres()

	repo := repository.NewInventoryRepository(conn)
	svc := service.NewInventoryService(repo)
	h := handler.NewInventoryHandler(svc)

	router := httprouter.New()
	router.GET("/inventories", h.GetAll)
	router.GET("/inventories/:id", h.GetByID)
	router.POST("/inventories", h.Create)
	router.PUT("/inventories/:id", h.Update)
	router.DELETE("/inventories/:id", h.Delete)

	fmt.Println("Server running on localhost:8080")

	listandserv := http.ListenAndServe(":8080", router)
	log.Fatal(listandserv)
}
