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

// ObjectHandler handles HTTP requests for objects
type ObjectHandler struct {
	service *services.ObjectService
}

// NewObjectHandler creates a new ObjectHandler
func NewObjectHandler(service *services.ObjectService) *ObjectHandler {
	return &ObjectHandler{service: service}
}

// CreateObject handles POST /api/objects
func (h *ObjectHandler) CreateObject(w http.ResponseWriter, r *http.Request) {
	var req models.CreateObjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	object, err := h.service.CreateObject(req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create object", err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, object)
}

// GetObjectByID handles GET /api/objects/{id}
func (h *ObjectHandler) GetObjectByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid object ID", err.Error())
		return
	}

	object, err := h.service.GetObjectByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Object not found", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, object)
}

// GetAllObjects handles GET /api/objects
func (h *ObjectHandler) GetAllObjects(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	response, err := h.service.GetAllObjects(page, pageSize)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve objects", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

// UpdateObject handles PUT /api/objects/{id}
func (h *ObjectHandler) UpdateObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid object ID", err.Error())
		return
	}

	var req models.UpdateObjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	object, err := h.service.UpdateObject(id, req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update object", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, object)
}

// DeleteObject handles DELETE /api/objects/{id}
func (h *ObjectHandler) DeleteObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid object ID", err.Error())
		return
	}

	if err := h.service.DeleteObject(id); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete object", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, models.SuccessResponse{
		Message: "Object deleted successfully",
	})
}

// GetLibraries handles GET /api/objects/libraries
func (h *ObjectHandler) GetLibraries(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	response, err := h.service.GetLibraries(page, pageSize)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve libraries", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

// GetObjectsByTypeID handles GET /api/objects/type/{typeId}
func (h *ObjectHandler) GetObjectsByTypeID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	typeID, err := strconv.Atoi(vars["typeId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid object type ID", err.Error())
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	response, err := h.service.GetObjectsByTypeID(typeID, page, pageSize)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve objects by type", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}
func (h *ObjectHandler) GetHierarchyFolder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	objectID, err := uuid.Parse(vars["objectID"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid object ID", err.Error())
		return
	}
	profileID, _ := strconv.Atoi(r.URL.Query().Get("profileID"))
	isFolder, _ := strconv.Atoi(r.URL.Query().Get("isFolder"))
	response, err := h.service.GetHierarchyFolder(objectID, profileID, isFolder == 1)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve hierarchy folder", err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, response)
}
