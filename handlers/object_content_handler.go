package handlers

import (
	"encoding/json"
	"enterprise-architect-api/models"
	"enterprise-architect-api/services"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// ObjectContentHandler handles HTTP requests for object contents
type ObjectContentHandler struct {
	service *services.ObjectContentService
}

// NewObjectContentHandler creates a new ObjectContentHandler
func NewObjectContentHandler(service *services.ObjectContentService) *ObjectContentHandler {
	return &ObjectContentHandler{service: service}
}

// CreateObjectContent handles POST /api/object-contents
func (h *ObjectContentHandler) CreateObjectContent(w http.ResponseWriter, r *http.Request) {
	var req models.CreateObjectContentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	objectContent, err := h.service.CreateObjectContent(req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create object content", err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, objectContent)
}

// GetObjectContentByID handles GET /api/object-contents/{id}
func (h *ObjectContentHandler) GetObjectContentByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid object content ID", err.Error())
		return
	}

	objectContent, err := h.service.GetObjectContentByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Object content not found", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, objectContent)
}

// GetAllObjectContents handles GET /api/object-contents
func (h *ObjectContentHandler) GetAllObjectContents(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	response, err := h.service.GetAllObjectContents(page, pageSize)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve object contents", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

// UpdateObjectContent handles PUT /api/object-contents/{id}
func (h *ObjectContentHandler) UpdateObjectContent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid object content ID", err.Error())
		return
	}

	var req models.UpdateObjectContentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	objectContent, err := h.service.UpdateObjectContent(id, req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update object content", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, objectContent)
}

// DeleteObjectContent handles DELETE /api/object-contents/{id}
func (h *ObjectContentHandler) DeleteObjectContent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid object content ID", err.Error())
		return
	}

	if err := h.service.DeleteObjectContent(id); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete object content", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, models.SuccessResponse{
		Message: "Object content deleted successfully",
	})
}

func (h *ObjectContentHandler) GetDashboardStatistics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	libraryId, err := uuid.Parse(vars["libraryId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid library ID", err.Error())
		return
	}

	objectContents, err := h.service.DashboardCount(libraryId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve object contents", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, objectContents)
}
