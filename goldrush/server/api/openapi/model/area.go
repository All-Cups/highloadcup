// Code generated by go-swagger; DO NOT EDIT.

package model

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"bytes"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Area area
//
// swagger:model area
type Area struct {

	// pos x
	// Required: true
	// Minimum: 0
	PosX *int64 `json:"posX"`

	// pos y
	// Required: true
	// Minimum: 0
	PosY *int64 `json:"posY"`

	// size x
	// Minimum: 1
	SizeX int64 `json:"sizeX,omitempty"`

	// size y
	// Minimum: 1
	SizeY int64 `json:"sizeY,omitempty"`
}

// UnmarshalJSON unmarshals this object while disallowing additional properties from JSON
func (m *Area) UnmarshalJSON(data []byte) error {
	var props struct {

		// pos x
		// Required: true
		// Minimum: 0
		PosX *int64 `json:"posX"`

		// pos y
		// Required: true
		// Minimum: 0
		PosY *int64 `json:"posY"`

		// size x
		// Minimum: 1
		SizeX int64 `json:"sizeX,omitempty"`

		// size y
		// Minimum: 1
		SizeY int64 `json:"sizeY,omitempty"`
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&props); err != nil {
		return err
	}

	m.PosX = props.PosX
	m.PosY = props.PosY
	m.SizeX = props.SizeX
	m.SizeY = props.SizeY
	return nil
}

// Validate validates this area
func (m *Area) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validatePosX(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePosY(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSizeX(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSizeY(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Area) validatePosX(formats strfmt.Registry) error {

	if err := validate.Required("posX", "body", m.PosX); err != nil {
		return err
	}

	if err := validate.MinimumInt("posX", "body", int64(*m.PosX), 0, false); err != nil {
		return err
	}

	return nil
}

func (m *Area) validatePosY(formats strfmt.Registry) error {

	if err := validate.Required("posY", "body", m.PosY); err != nil {
		return err
	}

	if err := validate.MinimumInt("posY", "body", int64(*m.PosY), 0, false); err != nil {
		return err
	}

	return nil
}

func (m *Area) validateSizeX(formats strfmt.Registry) error {

	if swag.IsZero(m.SizeX) { // not required
		return nil
	}

	if err := validate.MinimumInt("sizeX", "body", int64(m.SizeX), 1, false); err != nil {
		return err
	}

	return nil
}

func (m *Area) validateSizeY(formats strfmt.Registry) error {

	if swag.IsZero(m.SizeY) { // not required
		return nil
	}

	if err := validate.MinimumInt("sizeY", "body", int64(m.SizeY), 1, false); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Area) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Area) UnmarshalBinary(b []byte) error {
	var res Area
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}