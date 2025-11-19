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

type AttributeHandler struct {
	service *services.AttributeService
}

func NewAttributeHandler(service *services.AttributeService) *AttributeHandler {
	return &AttributeHandler{service: service}
}

func (ah *AttributeHandler) GetAttributeForObject(w http.ResponseWriter, r *http.Request) {
	objectID := mux.Vars(r)["objectID"]
	if objectID == "" {
		respondWithError(w, http.StatusBadRequest, "objectID is required", "objectID is required")
		return
	}

	attributes, err := ah.service.GetAttributeForObject(uuid.Must(uuid.Parse(objectID)))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error retrieving attributes", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, attributes)
}

// CreateAttribute handles POST /api/attributes
func (ah *AttributeHandler) CreateAttribute(w http.ResponseWriter, r *http.Request) {
	var attribute models.Attribute
	if err := json.NewDecoder(r.Body).Decode(&attribute); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	if err := ah.service.CreateAttribute(&attribute); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create attribute", err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, attribute)
}

// GetAttributeByID handles GET /api/attributes/{id}
func (ah *AttributeHandler) GetAttributeByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid attribute ID", "ID is required")
		return
	}

	attribute, err := ah.service.GetAttributeByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Attribute not found", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, attribute)
}

// GetAllAttributes handles GET /api/attributes
func (ah *AttributeHandler) GetAllAttributes(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	response, err := ah.service.GetAllAttributes(page, pageSize)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve attributes", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

// UpdateAttribute handles PUT /api/attributes/{id}
func (ah *AttributeHandler) UpdateAttribute(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid attribute ID", "ID is required")
		return
	}

	var attribute models.Attribute
	if err := json.NewDecoder(r.Body).Decode(&attribute); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	updatedAttribute, err := ah.service.UpdateAttribute(id, &attribute)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update attribute", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, updatedAttribute)
}

// DeleteAttribute handles DELETE /api/attributes/{id}
func (ah *AttributeHandler) DeleteAttribute(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid attribute ID", "ID is required")
		return
	}

	if err := ah.service.DeleteAttribute(id); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete attribute", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, models.SuccessResponse{
		Message: "Attribute deleted successfully",
	})
}

// AssignAttributeToObjectType handles POST /api/attributes/assign-to-object-type
func (ah *AttributeHandler) AssignAttributeToObjectType(w http.ResponseWriter, r *http.Request) {
	var req models.AssignAttributeToObjectTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	if err := ah.service.AssignAttributeToObjectType(&req); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to assign attribute to object type", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, models.SuccessResponse{
		Message: "Attribute assigned to object type successfully",
	})
}

func (ah *AttributeHandler) GetAttributeAssignments(w http.ResponseWriter, r *http.Request) {
	objectTypeId, err := strconv.Atoi(r.URL.Query().Get("objectTypeId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid objectTypeId", "objectTypeId must be a valid integer")
		return
	}

	relationTypeId := uuid.Nil
	if relationTypeIdStr := r.URL.Query().Get("relationTypeId"); relationTypeIdStr != "" {
		var err error
		relationTypeId, err = uuid.Parse(relationTypeIdStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid relationTypeId", err.Error())
			return
		}
	}
	assignments, err := ah.service.GetAttributeAssignments(objectTypeId, relationTypeId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get attribute assignments", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, assignments)
}
