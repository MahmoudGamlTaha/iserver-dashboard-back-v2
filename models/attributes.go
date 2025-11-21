package models

import (
	"time"

	"github.com/google/uuid"
)

type General struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
type AssignedAttribute struct {
	AttributeID   uuid.UUID  `json:"attributeId" db:"AttributeId"`
	ObjectId      uuid.UUID  `json:"objectId" db:"ObjectId"`
	VersionId     uuid.UUID  `json:"versionId" db:"VersionId"`
	DataType      string     `json:"dataType" db:"DataType"`
	TextValue     *string    `json:"textValue,omitempty" db:"TextValue"`
	RichTextValue *string    `json:"richTextValue,omitempty" db:"RichTextValue"`
	BooleanValue  *bool      `json:"booleanValue,omitempty" db:"BooleanValue"`
	IntegerValue  *int       `json:"integerValue,omitempty" db:"IntegerValue"`
	FloatValue    *float64   `json:"floatValue,omitempty" db:"FloatValue"`
	DateValue     *time.Time `json:"dateValue,omitempty" db:"DateValue"`
	AttributeType string     `json:"attributeType" db:"AttributeType"`
	AttributeName string     `json:"attributeName" db:"AttributeName"`
	IsMandatory   bool       `json:"isMandatory" db:"IsMandatory"`
	ObjectTypeId  int        `json:"objectTypeId" db:"objectTypeId"`
}
type ObjectTypeAssignedAttribute struct {
	AttributeId         uuid.UUID `json:"attributeId" db:"AttributeId"`
	AttributeName       string    `json:"attributeName" db:"AttributeName"`
	AttributeType       string    `json:"attributeType" db:"AttributeType"`
	Description         string    `json:"description" db:"Description"`
	TooltipText         string    `json:"tooltipText" db:"TooltipText"`
	TextDefaultValue    *string   `json:"textDefaultValue,omitempty" db:"TextDefaultValue"`
	IntDefaultValue     *int64    `json:"intDefaultValue,omitempty" db:"IntDefaultValue"`
	BoolDefaultValue    *bool     `json:"boolDefaultValue,omitempty" db:"BoolDefaultValue"`
	ListType            *uint8    `json:"listType,omitempty" db:"ListType"`
	ListDefaultValue    *int      `json:"listDefaultValue,omitempty" db:"ListDefaultValue"`
	ListValues          *string   `json:"listValues,omitempty" db:"ListValues"`
	AttributeGroupId    uuid.UUID `json:"attributeGroupId" db:"AttributeGroupId"`
	ObjectTypeId        int       `json:"objectTypeId" db:"ObjectTypeId"`
	RelationTypeId      uuid.UUID `json:"relationTypeId" db:"RelationTypeId"`
	GeneralType         *string   `json:"generalType" db:"GeneralType"`
	SequenceWithinGroup int       `json:"sequenceWithinGroup" db:"SequenceWithinGroup"`
	AttributeGroupName  string    `json:"attributeGroupName" db:"AttributeGroupName"`
	ObjectTypeName      *string   `json:"objectTypeName" db:"ObjectTypeName"`
}

type ObjectInstanceAttribute struct {
	AssignedAttribute        []ObjectTypeAssignedAttribute
	AssignedAttributesValues []AssignedAttribute
	General
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
type AttributeAssignment struct {
	AttributeId         uuid.UUID `json:"attributeId" db:"AttributeId"`
	AttributeName       string    `json:"attributeName" db:"AttributeName"`
	AttributeType       string    `json:"attributeType" db:"AttributeType"`
	Description         string    `json:"description" db:"Description"`
	TooltipText         string    `json:"tooltipText" db:"tooltipText"`
	AttributeGroupId    uuid.UUID `json:"attributeGroupId" db:"AttributeGroupId"`
	ObjectTypeId        int       `json:"objectTypeId" db:"ObjectTypeId"`
	RelationTypeId      uuid.UUID `json:"relationTypeId" db:"RelationTypeId"`
	GeneralType         *string   `json:"generalType" db:"GeneralType"`
	SequenceWithinGroup int       `json:"sequenceWithinGroup" db:"SequenceWithinGroup"`
	AttributeGroupName  string    `json:"attributeGroupName" db:"AttributeGroupName"`
	ObjectTypeName      *string   `json:"objectTypeName" db:"ObjectTypeName"`
}

type UnassignAttributeFromObjectTypeRequest struct {
	ObjectTypeId     int        `json:"objectTypeId"`
	RelationTypeId   *uuid.UUID `json:"relationTypeId"`
	AttributeGroupId uuid.UUID  `json:"attributeGroupId"`
	AttributeId      uuid.UUID  `json:"attributeId"`
}
