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

// Report report
//
// swagger:model report
type Report struct {

	// area
	// Required: true
	Area *Area `json:"area"`

	// amount
	// Required: true
	Amount Amount `json:"amount"`

	// Histogram, key is depth ("1", "2", …).
	AmountPerDepth map[string]Amount `json:"amountPerDepth,omitempty"`
}

// UnmarshalJSON unmarshals this object while disallowing additional properties from JSON
func (m *Report) UnmarshalJSON(data []byte) error {
	var props struct {

		// area
		// Required: true
		Area *Area `json:"area"`

		// amount
		// Required: true
		Amount Amount `json:"amount"`

		// Histogram, key is depth ("1", "2", …).
		AmountPerDepth map[string]Amount `json:"amountPerDepth,omitempty"`
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&props); err != nil {
		return err
	}

	m.Area = props.Area
	m.Amount = props.Amount
	m.AmountPerDepth = props.AmountPerDepth
	return nil
}

// Validate validates this report
func (m *Report) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateArea(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateAmount(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateAmountPerDepth(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Report) validateArea(formats strfmt.Registry) error {

	if err := validate.Required("area", "body", m.Area); err != nil {
		return err
	}

	if m.Area != nil {
		if err := m.Area.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("area")
			}
			return err
		}
	}

	return nil
}

func (m *Report) validateAmount(formats strfmt.Registry) error {

	if err := m.Amount.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("amount")
		}
		return err
	}

	return nil
}

func (m *Report) validateAmountPerDepth(formats strfmt.Registry) error {

	if swag.IsZero(m.AmountPerDepth) { // not required
		return nil
	}

	for k := range m.AmountPerDepth {

		if val, ok := m.AmountPerDepth[k]; ok {
			if err := val.Validate(formats); err != nil {
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *Report) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Report) UnmarshalBinary(b []byte) error {
	var res Report
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
