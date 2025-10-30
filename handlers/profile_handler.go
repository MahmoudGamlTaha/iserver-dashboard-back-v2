package handlers

import (
	"encoding/json"
	"enterprise-architect-api/models"
	"enterprise-architect-api/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// ProfileHandler handles HTTP requests for profiles
type ProfileHandler struct {
	service *services.ProfileService
}

// NewProfileHandler creates a new ProfileHandler
func NewProfileHandler(service *services.ProfileService) *ProfileHandler {
	return &ProfileHandler{service: service}
}

// CreateProfile handles POST /api/profiles
func (h *ProfileHandler) CreateProfile(w http.ResponseWriter, r *http.Request) {
	var req models.CreateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	profile, err := h.service.CreateProfile(req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create profile", err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, profile)
}

// GetProfileByID handles GET /api/profiles/{id}
func (h *ProfileHandler) GetProfileByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid profile ID", err.Error())
		return
	}

	profile, err := h.service.GetProfileByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Profile not found", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, profile)
}

// GetAllProfiles handles GET /api/profiles
func (h *ProfileHandler) GetAllProfiles(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	response, err := h.service.GetAllProfiles(page, pageSize)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve profiles", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

// UpdateProfile handles PUT /api/profiles/{id}
func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid profile ID", err.Error())
		return
	}

	var req models.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	profile, err := h.service.UpdateProfile(id, req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update profile", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, profile)
}

// DeleteProfile handles DELETE /api/profiles/{id}
func (h *ProfileHandler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid profile ID", err.Error())
		return
	}

	if err := h.service.DeleteProfile(id); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete profile", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, models.SuccessResponse{
		Message: "Profile deleted successfully",
	})
}

