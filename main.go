package main

import (
	"enterprise-architect-api/config"
	"enterprise-architect-api/handlers"
	"enterprise-architect-api/middleware"
	"enterprise-architect-api/repositories"
	"enterprise-architect-api/services"
	"enterprise-architect-api/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := utils.ConnectDB(utils.DBConfig{
		Server:   cfg.Database.Server,
		Port:     cfg.Database.Port,
		Database: cfg.Database.Database,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Successfully connected to database")

	// Initialize repositories
	objectRepo := repositories.NewObjectRepository(db)
	objectTypeRepo := repositories.NewObjectTypeRepository(db)
	profileRepo := repositories.NewProfileRepository(db)
	objectContentRepo := repositories.NewObjectContentRepository(db)
	folderRepo := repositories.NewFolderRepository(db)
	attributeRepo := repositories.NewAttributeRepository(db)
	reportConfigRepo := repositories.NewReportConfigRepository(db)

	// Initialize services
	objectService := services.NewObjectService(objectRepo)
	objectTypeService := services.NewObjectTypeService(objectTypeRepo)
	profileService := services.NewProfileService(profileRepo)
	objectContentService := services.NewObjectContentService(objectContentRepo)
	folderService := services.NewFolderService(folderRepo)
	attributeService := services.NewAttributeService(attributeRepo)
	fileObjectsService := services.NewFileObjectsService()
	eaTagService := services.NewEATagService(reportConfigRepo)

	// Initialize handlers
	objectHandler := handlers.NewObjectHandler(objectService)
	objectTypeHandler := handlers.NewObjectTypeHandler(objectTypeService)
	profileHandler := handlers.NewProfileHandler(profileService)
	objectContentHandler := handlers.NewObjectContentHandler(objectContentService)
	folderHandler := handlers.NewFolderHandler(folderService)
	attributeHandler := handlers.NewAttributeHandler(attributeService)
	fileObjectsHandler := handlers.NewFileObjectsHandler(fileObjectsService)
	eaTagHandler := handlers.NewEATagHandler(eaTagService)

	// Setup router
	router := mux.NewRouter()

	// API routes
	api := router.PathPrefix("/api").Subrouter()
	api.Use(loggingMiddleware)
	// Object routes
	api.HandleFunc("/objects", objectHandler.GetAllObjects).Methods("GET")
	api.HandleFunc("/objects", objectHandler.CreateObject).Methods("POST")
	api.HandleFunc("/objects/libraries", objectHandler.GetLibraries).Methods("GET")
	api.HandleFunc("/objects/hierarchy/{objectID}", objectHandler.GetHierarchyFolder).Methods("GET")
	api.HandleFunc("/objects/type/{typeId}", objectHandler.GetObjectsByTypeID).Methods("GET")
	api.HandleFunc("/objects/{id}", objectHandler.GetObjectByID).Methods("GET")
	api.HandleFunc("/objects/{id}", objectHandler.UpdateObject).Methods("PUT")
	api.HandleFunc("/objects/{id}", objectHandler.DeleteObject).Methods("DELETE")
	api.HandleFunc("/objects/{objectTypeID}/{libraryID}", objectHandler.GetObjectsByObjectTypeIDAndLibraryID).Methods("GET")

	// ObjectType routes
	api.HandleFunc("/object-types", objectTypeHandler.GetAllObjectTypes).Methods("GET")
	api.HandleFunc("/object-types", objectTypeHandler.CreateObjectType).Methods("POST")
	api.HandleFunc("/object-types/{id}", objectTypeHandler.GetObjectTypeByID).Methods("GET")
	api.HandleFunc("/object-types/{id}", objectTypeHandler.UpdateObjectType).Methods("PUT")
	api.HandleFunc("/object-types/{id}", objectTypeHandler.DeleteObjectType).Methods("DELETE")

	// Profile routes
	api.HandleFunc("/profiles", profileHandler.GetAllProfiles).Methods("GET")
	api.HandleFunc("/profiles", profileHandler.CreateProfile).Methods("POST")
	api.HandleFunc("/profiles/{id}", profileHandler.GetProfileByID).Methods("GET")
	api.HandleFunc("/profiles/{id}", profileHandler.UpdateProfile).Methods("PUT")
	api.HandleFunc("/profiles/{id}", profileHandler.DeleteProfile).Methods("DELETE")

	// ObjectContent routes
	api.HandleFunc("/object-contents", objectContentHandler.GetAllObjectContents).Methods("GET")
	api.HandleFunc("/object-contents", objectContentHandler.CreateObjectContent).Methods("POST")
	api.HandleFunc("/object-contents/{id}", objectContentHandler.GetObjectContentByID).Methods("GET")
	api.HandleFunc("/object-contents/{id}", objectContentHandler.UpdateObjectContent).Methods("PUT")
	api.HandleFunc("/object-contents/{id}", objectContentHandler.DeleteObjectContent).Methods("DELETE")

	// Folder routes
	api.HandleFunc("/folders/object-type/{libraryId}", folderHandler.GetObjectTypeFolders).Methods("GET")
	api.HandleFunc("/folders/{folderId}/contents", folderHandler.GetFoldersByLibrary).Methods("GET")

	// Dashboard routes
	api.HandleFunc("/dashboard/object-counts/{libraryId}", objectContentHandler.GetDashboardStatistics).Methods("GET")
	api.HandleFunc("/dashboard/object-counts-grouped/{libraryId}", objectContentHandler.GetDashboardStatisticsGrouped).Methods("GET")

	// Attribute routes
	api.HandleFunc("/attributes", attributeHandler.GetAllAttributes).Methods("GET")
	api.HandleFunc("/attributes", attributeHandler.CreateAttribute).Methods("POST")
	api.HandleFunc("/attributes/assignments", attributeHandler.GetAttributeAssignments).Methods("GET")
	api.HandleFunc("/attributes/assign-to-object-type", attributeHandler.AssignAttributeToObjectType).Methods("POST")
	api.HandleFunc("/attributes/unassign-from-object-type", attributeHandler.UnassignAttributeFromObjectType).Methods("DELETE")
	api.HandleFunc("/attributes/object/{objectID}", attributeHandler.GetAttributeForObject).Methods("GET")
	api.HandleFunc("/attributes/{id}", attributeHandler.GetAttributeByID).Methods("GET")
	api.HandleFunc("/attributes/{id}", attributeHandler.UpdateAttribute).Methods("PUT")
	api.HandleFunc("/attributes/{id}", attributeHandler.DeleteAttribute).Methods("DELETE")

	// File conversion routes
	api.HandleFunc("/convert-visio", fileObjectsHandler.ConvertVisioToSVGHandler).Methods("POST")

	// EA Tags routes
	api.HandleFunc("/ea-tags", eaTagHandler.GetAllEATags).Methods("GET")
	api.HandleFunc("/ea-tags", eaTagHandler.CreateEATag).Methods("POST")
	api.HandleFunc("/ea-tags/{id}", eaTagHandler.GetEATagByID).Methods("GET")
	api.HandleFunc("/ea-tags/{id}", eaTagHandler.UpdateEATag).Methods("PUT")
	api.HandleFunc("/ea-tags/{id}", eaTagHandler.DeleteEATag).Methods("DELETE")
	api.HandleFunc("/ea-tags/assign-dimension", eaTagHandler.AssignObjectTypeToDimention).Methods("POST")
	api.HandleFunc("/ea-tags/assigned-dimension/{objectTypeID}", eaTagHandler.GetEAObjectTypesAssignedToDimension).Methods("GET")

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Start server
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting server on %s", serverAddr)

	// Wrap router with CORS handler
	handler := middleware.CorsMiddleware()(router)

	if err := http.ListenAndServe(serverAddr, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("API Request: %s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
