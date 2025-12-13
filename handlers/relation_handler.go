package handlers

import (
	"encoding/json"
	"enterprise-architect-api/models"
	"enterprise-architect-api/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type RelationHandler struct {
	service *services.RelationService
}

func NewRelationHandler(service *services.RelationService) *RelationHandler {
	return &RelationHandler{service: service}
}

// CreateRelation handles POST /api/relations
func (h *RelationHandler) CreateRelation(w http.ResponseWriter, r *http.Request) {
	var req models.CreateRelationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	// Assuming 62 is the default user if not provided, similar to other handlers
	// But the model has CreatedBy as int, and the request has it.
	// If the user didn't send it, it will be 0.
	// We might want to enforce a default or validation here if needed,
	// but for now relying on validation or the request data.

	relation, err := h.service.CreateRelation(req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create relationship", err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, relation)
}

// GetRelationsByObjectID handles GET /api/objects/{id}/relations
func (h *RelationHandler) GetRelationsByObjectID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	objectID := vars["id"]

	relations, err := h.service.GetRelationsByObjectID(objectID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve relationships", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, relations)
}

// GetAllRelationTypes handles GET /api/relation-types
func (h *RelationHandler) GetAllRelationTypes(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	response, err := h.service.GetAllRelationTypes(page, pageSize)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve relation types", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}
