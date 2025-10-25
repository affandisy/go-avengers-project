package main

import (
	"avenger/internal/handler"
	"avenger/internal/middleware"
	"avenger/internal/repository"
	"avenger/internal/service"
	"avenger/pkg/db"
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found")
	}

	connInv := db.InitPostgres()
	defer func() {
		if err := connInv.Close(); err != nil {
			slog.Error("Failed to close inventory database connection", slog.Any("error", err))
		}
	}()

	connUserRecipe := db.InitPostgresGORM()
	sqlDB, _ := connUserRecipe.DB()
	defer func() {
		if err := sqlDB.Close(); err != nil {
			slog.Error("Failed to close user/recipe database connection", slog.Any("error", err))
		}
	}()

	// Initialize repositories
	repoInv := repository.NewInventoryRepository(connInv)
	userRepo := repository.NewUserRepository(connUserRecipe)
	recipeRepo := repository.NewRecipeRepository(connUserRecipe)

	// Initialize services
	svcInv := service.NewInventoryService(repoInv)
	userSvc := service.NewUserService(userRepo)
	recipeSvc := service.NewRecipeService(recipeRepo)

	// Initialize handlers
	inventoryHandler := handler.NewInventoryHandler(svcInv)
	authHandler := handler.NewAuthHandler(userSvc)
	recipeHandler := handler.NewRecipeHandler(recipeSvc)

	router := httprouter.New()

	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, err any) {
		slog.Error("PANIC occurred", slog.Any("error", err), slog.String("method", r.Method), slog.String("path", r.URL.Path), slog.String("remote_addr", r.RemoteAddr))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "Internal server error", "errors": null}`))
	}

	// Custom 404 handler
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Warn("Endpoint not found",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
		)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Endpoint not found", "errors": null}`))
	})

	// Custom method not allowed handler
	router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Warn("Method not allowed",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
		)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"message": "Method not allowed", "errors": null}`))
	})

	// ========== INVENTORY ROUTES (Public) ==========
	router.GET("/inventories", inventoryHandler.GetAll)
	router.GET("/inventories/:id", inventoryHandler.GetByID)
	router.POST("/inventories", inventoryHandler.Create)
	router.PUT("/inventories/:id", inventoryHandler.Update)
	router.DELETE("/inventories/:id", inventoryHandler.Delete)

	// ========== AUTH ROUTES (Public) ==========
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)

	// ========== RECIPE ROUTES ==========
	// Public: Anyone can view recipes
	router.Handler("GET", "/recipes", wrapHandler(recipeHandler.GetAll))

	// Protected: Only superadmin can create recipes
	router.Handler("POST", "/recipes", wrapHandler(
		middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			params := httprouter.ParamsFromContext(r.Context())
			recipeHandler.Create(w, r, params)
		}, "superadmin"),
	))

	// Protected: Only superadmin can delete recipes
	router.Handler("DELETE", "/recipes/:id", wrapHandler(
		middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			params := httprouter.ParamsFromContext(r.Context())
			recipeHandler.Delete(w, r, params)
		}, "superadmin"),
	))

	// Create HTTP server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		slog.Info("🚀 Server starting on http://localhost:8080")
		log.Println("Server running on localhost:8080")
		log.Println("=====================================")
		log.Println("📚 Available Endpoints:")
		log.Println("  POST   /register          - Register new user")
		log.Println("  POST   /login             - Login and get token")
		log.Println("  GET    /inventories       - Get all inventories")
		log.Println("  GET    /inventories/:id   - Get inventory by ID")
		log.Println("  POST   /inventories       - Create inventory")
		log.Println("  PUT    /inventories/:id   - Update inventory")
		log.Println("  DELETE /inventories/:id   - Delete inventory")
		log.Println("  GET    /recipes           - Get all recipes (public)")
		log.Println("  POST   /recipes           - Create recipe (superadmin)")
		log.Println("  DELETE /recipes/:id       - Delete recipe (superadmin)")
		log.Println("=====================================")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start:", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	slog.Info("Server exited properly")
}

func wrapHandler(f interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch h := f.(type) {
		case func(http.ResponseWriter, *http.Request, httprouter.Params):
			params := httprouter.ParamsFromContext(r.Context())
			h(w, r, params)
		case http.HandlerFunc:
			h(w, r)
		default:
			slog.Error("Invalid handler type", slog.Any("type", h))
			http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
		}
	})
}
