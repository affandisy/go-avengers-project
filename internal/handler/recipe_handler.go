package handler

import (
	"avenger/internal/domain"
	"avenger/internal/service"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

type RecipeHandler struct {
	service  service.RecipeService
	validate *validator.Validate
}

func NewRecipeHandler(s service.RecipeService) *RecipeHandler {
	return &RecipeHandler{service: s, validate: validator.New()}
}

func (h *RecipeHandler) GetAll(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	data, err := h.service.GetAll()
	if err != nil {
		slog.Error("GetAll recipes error", slog.Any("error", err))
		writeError(w, http.StatusInternalServerError, "Failed to retrieve recipes", nil)
		return
	}

	writeJSON(w, http.StatusOK, Response{
		Message: "success",
		Data:    data,
	})
}

func (h *RecipeHandler) Create(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var rec domain.Recipe
	if err := json.NewDecoder(r.Body).Decode(&rec); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body", map[string]string{
			"body": "Request body must be valid JSON",
		})
		return
	}

	if err := h.service.Create(&rec); err != nil {
		slog.Error("Create recipe error", slog.Any("error", err))
		writeError(w, http.StatusInternalServerError, "Failed to create recipe", nil)
		return
	}

	writeJSON(w, http.StatusCreated, Response{
		Message: "Recipe created successfully",
		Data:    rec,
	})
}

func (h *RecipeHandler) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idStr := p.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "Invalid ID parameter", map[string]string{
			"id": "ID must be a positive integer",
		})
		return
	}

	if err := h.service.Delete(id); err != nil {
		slog.Error("Delete recipe error", slog.Int("id", id), slog.Any("error", err))
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "Recipe not found", nil)
			return
		}

		writeError(w, http.StatusInternalServerError, "Failed to delete recipe", nil)
		return
	}

	writeJSON(w, http.StatusOK, Response{
		Message: "Recipe deleted successfully",
		Data: map[string]any{
			"id": id,
		},
	})
}
