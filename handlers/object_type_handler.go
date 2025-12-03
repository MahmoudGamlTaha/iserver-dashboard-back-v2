package handlers

import (
	"encoding/json"
	"enterprise-architect-api/models"
	"enterprise-architect-api/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// ObjectTypeHandler handles HTTP requests for object types
type ObjectTypeHandler struct {
	service *services.ObjectTypeService
}

// NewObjectTypeHandler creates a new ObjectTypeHandler
func NewObjectTypeHandler(service *services.ObjectTypeService) *ObjectTypeHandler {
	return &ObjectTypeHandler{service: service}
}

// CreateObjectType handles POST /api/object-types
func (h *ObjectTypeHandler) CreateObjectType(w http.ResponseWriter, r *http.Request) {
	var req models.CreateObjectTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	objectType, err := h.service.CreateObjectType(req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create object type", err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, objectType)
}

// GetObjectTypeByID handles GET /api/object-types/{id}
func (h *ObjectTypeHandler) GetObjectTypeByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid object type ID", err.Error())
		return
	}

	objectType, err := h.service.GetObjectTypeByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Object type not found", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, objectType)
}

// GetAllObjectTypes handles GET /api/object-types
func (h *ObjectTypeHandler) GetAllObjectTypes(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	response, err := h.service.GetAllObjectTypes(page, pageSize)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve object types", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

// SearchObjectTypes handles GET /api/object-types/search
func (h *ObjectTypeHandler) SearchObjectTypes(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	response, err := h.service.SearchObjectTypesByName(name, page, pageSize)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to search object types", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

// UpdateObjectType handles PUT /api/object-types/{id}
func (h *ObjectTypeHandler) UpdateObjectType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid object type ID", err.Error())
		return
	}

	var req models.UpdateObjectTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	objectType, err := h.service.UpdateObjectType(id, req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update object type", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, objectType)
}

// DeleteObjectType handles DELETE /api/object-types/{id}
func (h *ObjectTypeHandler) DeleteObjectType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid object type ID", err.Error())
		return
	}

	if err := h.service.DeleteObjectType(id); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete object type", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, models.SuccessResponse{
		Message: "Object type deleted successfully",
	})
}

// GetFolderRepositoryTree handles GET /api/object-types/folder-tree
func (h *ObjectTypeHandler) GetFolderRepositoryTree(w http.ResponseWriter, r *http.Request) {
	response, err := h.service.GetFolderRepositoryTree()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve folder repository tree", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

// AddFolderToTree handles POST /api/object-types/folder-tree
func (h *ObjectTypeHandler) AddFolderToTree(w http.ResponseWriter, r *http.Request) {
	var req models.AddFolderToTreeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	folderTypeHierarchyId, err := h.service.AddFolderToTree(req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to add folder to tree", err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message":               "Folder added to tree successfully",
		"folderTypeHierarchyId": folderTypeHierarchyId,
	})
}

// AssignObjectTypeToFolder handles POST /api/object-types/folder-assignments
func (h *ObjectTypeHandler) AssignObjectTypeToFolder(w http.ResponseWriter, r *http.Request) {
	var req models.FolderObjectTypes
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	if err := h.service.AssignObjectTypeToFolder(req); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to assign object type to folder", err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, models.SuccessResponse{
		Message: "Object type assigned to folder successfully",
	})
}

// GetAvailableTypesForFolder handles GET /api/object-types/folder-assignments/{folderObjectTypeId}
func (h *ObjectTypeHandler) GetAvailableTypesForFolder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	folderObjectTypeId, err := strconv.Atoi(vars["folderObjectTypeId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid folder object type ID", err.Error())
		return
	}

	folderObjectTypes, err := h.service.GetAvailableTypesForFolder(folderObjectTypeId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve available types for folder", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, folderObjectTypes)
}

// DeleteObjectTypeFromFolder handles DELETE /api/object-types/folder-assignments/{folderObjectTypeId}/{objectTypeId}
func (h *ObjectTypeHandler) DeleteObjectTypeFromFolder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	folderObjectTypeId, err := strconv.Atoi(vars["folderObjectTypeId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid folder object type ID", err.Error())
		return
	}

	objectTypeId, err := strconv.Atoi(vars["objectTypeId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid object type ID", err.Error())
		return
	}

	if err := h.service.DeleteObjectTypeFromFolder(folderObjectTypeId, objectTypeId); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete object type from folder", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, models.SuccessResponse{
		Message: "Object type removed from folder successfully",
	})
}
