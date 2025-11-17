package models

import (
	"time"

	"github.com/google/uuid"
)

type AssignedAttribute struct {
	AttributeID   uuid.UUID
	ObjectId      uuid.UUID
	VersionId     uuid.UUID
	DataType      string
	TextValue     *string
	RichTextValue *string
	BooleanValue  *bool
	IntegerValue  *int
	FloatValue    *float64
	DateValue     *time.Time
	AttributeType string
	AttributeName string
	IsMandatory   bool
}

type Attribute struct {
	AttributeId       uuid.UUID  `json:"attributeId" db:"AttributeId"`
	AttributeName     string     `json:"attributeName" db:"AttributeName"`
	AttributeType     string     `json:"attributeType" db:"AttributeType"`
	IsMandatory       bool       `json:"isMandatory" db:"IsMandatory"`
	IsSynchronised    bool       `json:"isSynchronised" db:"IsSynchronised"`
	VisioSyncName     string     `json:"visioSyncName" db:"VisioSyncName"`
	Description       string     `json:"description" db:"Description"`
	TooltipText       string     `json:"tooltipText" db:"TooltipText"`
	TextDefaultValue  *string    `json:"textDefaultValue,omitempty" db:"TextDefaultValue"`
	TextRowCount      *int       `json:"textRowCount,omitempty" db:"TextRowCount"`
	IntDefaultValue   *int64     `json:"intDefaultValue,omitempty" db:"IntDefaultValue"`
	IntLowerLimit     *int64     `json:"intLowerLimit,omitempty" db:"IntLowerLimit"`
	IntUpperLimit     *int64     `json:"intUpperLimit,omitempty" db:"IntUpperLimit"`
	FloatDefaultValue *float64   `json:"floatDefaultValue,omitempty" db:"FloatDefaultValue"`
	FloatLowerLimit   *float64   `json:"floatLowerLimit,omitempty" db:"FloatLowerLimit"`
	FloatUpperLimit   *float64   `json:"floatUpperLimit,omitempty" db:"FloatUpperLimit"`
	DateDefaultValue  *time.Time `json:"dateDefaultValue,omitempty" db:"DateDefaultValue"`
	BoolDefaultValue  *bool      `json:"boolDefaultValue,omitempty" db:"BoolDefaultValue"`
	AutoIdPrefix      *string    `json:"autoIdPrefix,omitempty" db:"AutoIdPrefix"`
	AutoIdSuffix      *string    `json:"autoIdSuffix,omitempty" db:"AutoIdSuffix"`
	AutoIdPadding     *int       `json:"autoIdPadding,omitempty" db:"AutoIdPadding"`
	AutoIdStartValue  *int       `json:"autoIdStartValue,omitempty" db:"AutoIdStartValue"`
	AutoIdNextValue   *int       `json:"autoIdNextValue,omitempty" db:"AutoIdNextValue"`
	ListDefaultValue  *int       `json:"listDefaultValue,omitempty" db:"ListDefaultValue"`
	ListType          *uint8     `json:"listType,omitempty" db:"ListType"`
	ListValues        *string    `json:"listValues,omitempty" db:"ListValues"`
	IsCalculated      bool       `json:"isCalculated" db:"IsCalculated"`
}

type AssignAttributeToObjectTypeRequest struct {
	ObjectTypeId       int       `json:"objectTypeId"`
	RelationTypeId     uuid.UUID `json:"relationTypeId"`
	AttributeGroupId   uuid.UUID `json:"attributeGroupId"`
	AttributeGroupName string    `json:"attributeGroupName"`
	AttributeId        uuid.UUID `json:"attributeId"`
}
