package handlers

import (
	"encoding/json"
	"enterprise-architect-api/models"
	"enterprise-architect-api/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// EATagHandler handles HTTP requests for EA tags
type EATagHandler struct {
	service *services.EATagService
}

// NewEATagHandler creates a new EATagHandler
func NewEATagHandler(service *services.EATagService) *EATagHandler {
	return &EATagHandler{service: service}
}

// CreateEATag handles POST /api/ea-tags
func (h *EATagHandler) CreateEATag(w http.ResponseWriter, r *http.Request) {
	var req models.CreateEATagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	tag, err := h.service.CreateEATag(req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create EA tag", err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, tag)
}

// GetEATagByID handles GET /api/ea-tags/{id}
func (h *EATagHandler) GetEATagByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid EA tag ID", err.Error())
		return
	}

	tag, err := h.service.GetEATagByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "EA tag not found", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, tag)
}

// GetAllEATags handles GET /api/ea-tags
func (h *EATagHandler) GetAllEATags(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	response, err := h.service.GetAllEATags(page, pageSize)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve EA tags", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

// UpdateEATag handles PUT /api/ea-tags/{id}
func (h *EATagHandler) UpdateEATag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid EA tag ID", err.Error())
		return
	}

	var req models.UpdateEATagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	tag, err := h.service.UpdateEATag(id, req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update EA tag", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, tag)
}

// DeleteEATag handles DELETE /api/ea-tags/{id}
func (h *EATagHandler) DeleteEATag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid EA tag ID", err.Error())
		return
	}

	if err := h.service.DeleteEATag(id); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete EA tag", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, models.SuccessResponse{
		Message: "EA tag deleted successfully",
	})
}

// AssignObjectTypeToDimention handles POST /api/ea-tags/assign-dimension
func (h *EATagHandler) AssignObjectTypeToDimention(w http.ResponseWriter, r *http.Request) {
	var req models.AssignObjectTypeToDimentionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	dimention, err := h.service.AssignObjectTypeToDimention(req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to assign object type to dimension", err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, dimention)
}
func (h *EATagHandler) GetEAObjectTypesAssignedToDimension(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	objectTypeID, err := strconv.Atoi(vars["objectTypeID"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid object type ID", err.Error())
		return
	}

	objectTypes, err := h.service.GetEAObjectTypesAssignedToDimension(int64(objectTypeID))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve EA object types assigned to dimension", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, objectTypes)
}
