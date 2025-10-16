package handler

import (
	"avenger/internal/domain"
	"avenger/internal/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type InventoryHandler struct {
	service service.InventoryService
}

func NewInventoryHandler(s service.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: s}
}

func (h *InventoryHandler) GetAll(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	data, err := h.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *InventoryHandler) GetByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, _ := strconv.Atoi(p.ByName("id"))
	data, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if data == nil {
		http.Error(w, "not found", 404)
		return
	}

	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *InventoryHandler) Create(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var inv domain.Inventory
	if err := json.NewDecoder(r.Body).Decode(&inv); err != nil {
		http.Error(w, "invalid body", 400)
		return
	}

	if inv.Status != "active" && inv.Status != "broken" {
		http.Error(w, "invalid status", 400)
		return
	}

	if err := h.service.Create(inv); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "created"}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *InventoryHandler) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, _ := strconv.Atoi(p.ByName("id"))
	var inv domain.Inventory
	if err := json.NewDecoder(r.Body).Decode(&inv); err != nil {
		http.Error(w, "invalid body", 400)
		return
	}
	if err := h.service.Update(id, inv); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "updated"}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *InventoryHandler) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, _ := strconv.Atoi(p.ByName("id"))
	if err := h.service.Delete(id); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "deleted"}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
