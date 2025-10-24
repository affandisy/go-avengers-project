package handler

import (
	"avenger/internal/domain"
	"avenger/internal/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type RecipeHandler struct {
	service service.RecipeService
}

func NewRecipeHandler(s service.RecipeService) *RecipeHandler {
	return &RecipeHandler{service: s}
}

func (h *RecipeHandler) GetAll(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	data, _ := h.service.GetAll()
	json.NewEncoder(w).Encode(data)
}

func (h *RecipeHandler) Create(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var rec domain.Recipe
	if err := json.NewDecoder(r.Body).Decode(&rec); err != nil {
		http.Error(w, "invalid body", 400)
		return
	}
	h.service.Create(&rec)
	json.NewEncoder(w).Encode(map[string]string{"message": "recipe created"})
}

func (h *RecipeHandler) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, _ := strconv.Atoi(p.ByName("id"))
	h.service.Delete(id)
	json.NewEncoder(w).Encode(map[string]string{"message": "recipe deleted"})
}
