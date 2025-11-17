package models

// EATag represents the EA_Tags table in the database
type EATag struct {
	ID     int    `json:"id" db:"id"`
	NameAr string `json:"name_ar" db:"name_ar"`
	NameEn string `json:"name_en" db:"name_en"`
}

// CreateEATagRequest represents the request body for creating a new EA tag
type CreateEATagRequest struct {
	NameAr string `json:"name_ar" validate:"required"`
	NameEn string `json:"name_en" validate:"required"`
}

// UpdateEATagRequest represents the request body for updating an EA tag
type UpdateEATagRequest struct {
	NameAr *string `json:"name_ar,omitempty"`
	NameEn *string `json:"name_en,omitempty"`
}

// EATagDimention represents the EA_Tags_Dimentions table in the database
type EATagDimention struct {
	ID           int `json:"id" db:"id"`
	EATagID      int `json:"ea_tag_id" db:"ea_tag_id"`
	ObjectTypeID int `json:"object_type_id" db:"object_type_id"`
}

// AssignObjectTypeToDimentionRequest represents the request for assigning object type to dimension
type AssignObjectTypeToDimentionRequest struct {
	ObjectTypeID int `json:"object_type_id" validate:"required"`
	EAID         int `json:"ea_tag_id" validate:"required"`
}

type AssignObjectTypeToDimentionResponse struct {
	ObjectTypeID int `json:"object_type_id" validate:"required"`
	EAID         int `json:"ea_tag_id" validate:"required"`
}
