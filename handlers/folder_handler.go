package handlers

import (
	"enterprise-architect-api/services"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// FolderHandler handles HTTP requests for folders
type FolderHandler struct {
	service *services.FolderService
}

// NewFolderHandler creates a new FolderHandler
func NewFolderHandler(service *services.FolderService) *FolderHandler {
	return &FolderHandler{service: service}
}

// GetObjectTypeFolders handles GET /api/folders/object-type/{libraryId}
func (h *FolderHandler) GetObjectTypeFolders(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	libraryID, err := uuid.Parse(vars["libraryId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid library ID", err.Error())
		return
	}

	folders, err := h.service.GetObjectTypeFolders(libraryID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve object type folders", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, folders)
}

// GetFoldersByLibrary handles GET /api/folders/{folderId}/contents
func (h *FolderHandler) GetFoldersByLibrary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	folderID, err := uuid.Parse(vars["folderId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid folder ID", err.Error())
		return
	}

	// Get profileId from query parameter
	profileIDStr := r.URL.Query().Get("profileId")
	if profileIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "Profile ID is required", "profileId query parameter is missing")
		return
	}

	profileID, err := strconv.Atoi(profileIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid profile ID", err.Error())
		return
	}

	contents, err := h.service.GetFoldersByLibrary(folderID, profileID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve folder contents", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, contents)
}

