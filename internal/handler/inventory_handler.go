package handler

import (
	"avenger/internal/domain"
	"avenger/internal/service"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator"
	"github.com/julienschmidt/httprouter"
)

type InventoryHandler struct {
	service  service.InventoryService
	validate *validator.Validate
}

func NewInventoryHandler(s service.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: s, validate: validator.New()}
}

func (h *InventoryHandler) GetAll(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	data, err := h.service.GetAll()
	if err != nil {
		slog.Error("GetAll inventory error", slog.Any("error", err))
		writeError(w, http.StatusInternalServerError, "Failed to retrieve inventories", nil)
		return
	}

	writeJSON(w, http.StatusOK, Response{
		Message: "success",
		Data:    data,
	})
}

func (h *InventoryHandler) GetByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idStr := p.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "Invalid ID parameter", map[string]string{
			"id": "ID must be a postivie integer",
		})
		return
	}

	data, err := h.service.GetByID(id)
	if err != nil {
		slog.Error("GetByID inventory error", slog.Int("id", id), slog.Any("error", err))
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "Inventory not found", nil)
			return
		}
	}

	writeJSON(w, http.StatusOK, Response{
		Message: "success",
		Data:    data,
	})
}

func (h *InventoryHandler) Create(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var inv domain.Inventory
	if err := json.NewDecoder(r.Body).Decode(&inv); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body", map[string]string{
			"body": "Request body must be valid JSON",
		})
		return
	}

	// Validate
	if err := h.validate.Struct(inv); err != nil {
		slog.Warn("Create inventory validation failed", slog.Any("error", err))
		writeError(w, http.StatusBadRequest, "Validation failed", formatValidationErrors(err))
		return
	}

	id, err := h.service.Create(inv)
	if err != nil {
		slog.Error("Create inventory error", slog.Any("error", err))

		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "already exists") {
			writeError(w, http.StatusConflict, "inventory already exists", nil)
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to create inventory", nil)
		return
	}

	writeJSON(w, http.StatusCreated, Response{
		Message: "inventory created successfully",
		Data:    map[string]any{"id": id},
	})

}

func (h *InventoryHandler) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idStr := p.ByName("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil || idInt <= 0 {
		writeError(w, http.StatusBadRequest, "Invalid ID parameter", map[string]string{
			"id": "ID must be a positive integer",
		})
		return
	}

	var inv domain.Inventory
	if err := json.NewDecoder(r.Body).Decode(&inv); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body", map[string]string{
			"body": "Request body must be valid JSON",
		})
		return
	}

	if err := h.validate.Struct(inv); err != nil {
		slog.Warn("Update inventory validation failed", slog.Int("id", idInt), slog.Any("error", err))
		writeError(w, http.StatusBadRequest, "Validation failed", formatValidationErrors(err))
		return
	}

	if err := h.service.Update(idInt, inv); err != nil {
		slog.Error("Update inventory error", slog.Int("id", idInt), slog.Any("error", err))

		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "Inventory not found", nil)
			return
		}

		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "already exists") {
			writeError(w, http.StatusConflict, "Inventory code already exists", nil)
			return
		}

		writeError(w, http.StatusInternalServerError, "Failed to update inventory", nil)
		return
	}

	writeJSON(w, http.StatusOK, Response{
		Message: "Inventory updated successfully",
		Data:    map[string]any{"id": idInt},
	})
}

func (h *InventoryHandler) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idStr := p.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "Invalid ID parameter", map[string]string{
			"id": "ID must be a positive integer",
		})
		return
	}

	if err := h.service.Delete(id); err != nil {
		slog.Error("Delete inventory error", slog.Int("id", id), slog.Any("error", err))

		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "Inventory not found", nil)
			return
		}

		writeError(w, http.StatusInternalServerError, "Failed to delete inventory", nil)
		return
	}

	writeJSON(w, http.StatusOK, Response{
		Message: "Inventory deleted successfully",
		Data:    map[string]any{"id": id},
	})

}
