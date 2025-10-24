package main

import (
	"avenger/internal/handler"
	"avenger/internal/middleware"
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

	connInv := db.InitPostgres()

	repoInv := repository.NewInventoryRepository(connInv)
	svcInv := service.NewInventoryService(repoInv)
	handlerInv := handler.NewInventoryHandler(svcInv)

	router := httprouter.New()
	router.GET("/inventories", handlerInv.GetAll)
	router.GET("/inventories/:id", handlerInv.GetByID)
	router.POST("/inventories", handlerInv.Create)
	router.PUT("/inventories/:id", handlerInv.Update)
	router.DELETE("/inventories/:id", handlerInv.Delete)

	connUserRecipe := db.InitPostgresGORM()

	userRepo := repository.NewUserRepository(connUserRecipe)
	userSvc := service.NewUserService(userRepo)
	authHandler := handler.NewAuthHandler(userSvc)

	recipeRepo := repository.NewRecipeRepository(connUserRecipe)
	recipeSvc := service.NewRecipeService(recipeRepo)
	recipeHandler := handler.NewRecipeHandler(recipeSvc)

	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)

	router.Handler("GET", "/recipes", wrap(recipeHandler.GetAll))
	router.Handler("POST", "/recipes", wrap(middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())
		recipeHandler.Create(w, r, params)
	}, "superadmin")))
	router.Handler("DELETE", "/recipes/:id", wrap(middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())
		recipeHandler.Delete(w, r, params)
	}, "superadmin")))

	fmt.Println("Server running on localhost:8080")

	listandserv := http.ListenAndServe(":8080", router)
	log.Fatal(listandserv)
}

func wrap(f interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch h := f.(type) {
		case func(http.ResponseWriter, *http.Request, httprouter.Params):
			params := httprouter.ParamsFromContext(r.Context())
			h(w, r, params)
		case http.HandlerFunc:
			h(w, r)
		default:
			http.Error(w, "invalid handler", 500)
		}
	})
}
