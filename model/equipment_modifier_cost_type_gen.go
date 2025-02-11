// Code generated from "enum.go.tmpl" - DO NOT EDIT.

/*
 * Copyright ©1998-2022 by Richard A. Wilkes. All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, version 2.0. If a copy of the MPL was not distributed with
 * this file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * This Source Code Form is "Incompatible With Secondary Licenses", as
 * defined by the Mozilla Public License, version 2.0.
 */

package model

import (
	"strings"

	"github.com/richardwilkes/toolbox/i18n"
)

// Possible values.
const (
	OriginalEquipmentModifierCostType EquipmentModifierCostType = iota
	BaseEquipmentModifierCostType
	FinalBaseEquipmentModifierCostType
	FinalEquipmentModifierCostType
	LastEquipmentModifierCostType = FinalEquipmentModifierCostType
)

// AllEquipmentModifierCostType holds all possible values.
var AllEquipmentModifierCostType = []EquipmentModifierCostType{
	OriginalEquipmentModifierCostType,
	BaseEquipmentModifierCostType,
	FinalBaseEquipmentModifierCostType,
	FinalEquipmentModifierCostType,
}

// EquipmentModifierCostType describes how an Equipment Modifier's cost is applied.
type EquipmentModifierCostType byte

// EnsureValid ensures this is of a known value.
func (enum EquipmentModifierCostType) EnsureValid() EquipmentModifierCostType {
	if enum <= LastEquipmentModifierCostType {
		return enum
	}
	return 0
}

// Key returns the key used in serialization.
func (enum EquipmentModifierCostType) Key() string {
	switch enum {
	case OriginalEquipmentModifierCostType:
		return "to_original_cost"
	case BaseEquipmentModifierCostType:
		return "to_base_cost"
	case FinalBaseEquipmentModifierCostType:
		return "to_final_base_cost"
	case FinalEquipmentModifierCostType:
		return "to_final_cost"
	default:
		return EquipmentModifierCostType(0).Key()
	}
}

// String implements fmt.Stringer.
func (enum EquipmentModifierCostType) String() string {
	switch enum {
	case OriginalEquipmentModifierCostType:
		return i18n.Text("to original cost")
	case BaseEquipmentModifierCostType:
		return i18n.Text("to base cost")
	case FinalBaseEquipmentModifierCostType:
		return i18n.Text("to final base cost")
	case FinalEquipmentModifierCostType:
		return i18n.Text("to final cost")
	default:
		return EquipmentModifierCostType(0).String()
	}
}

// AltString returns the alternate string.
func (enum EquipmentModifierCostType) AltString() string {
	switch enum {
	case OriginalEquipmentModifierCostType:
		return i18n.Text("\"+5\", \"-5\", \"+10%\", \"-10%\", \"x3.2\"")
	case BaseEquipmentModifierCostType:
		return i18n.Text("\"x2\", \"+2 CF\", \"-0.2 CF\"")
	case FinalBaseEquipmentModifierCostType:
		return i18n.Text("\"+5\", \"-5\", \"+10%\", \"-10%\", \"x3.2\"")
	case FinalEquipmentModifierCostType:
		return i18n.Text("\"+5\", \"-5\", \"+10%\", \"-10%\", \"x3.2\"")
	default:
		return EquipmentModifierCostType(0).AltString()
	}
}

// MarshalText implements the encoding.TextMarshaler interface.
func (enum EquipmentModifierCostType) MarshalText() (text []byte, err error) {
	return []byte(enum.Key()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (enum *EquipmentModifierCostType) UnmarshalText(text []byte) error {
	*enum = ExtractEquipmentModifierCostType(string(text))
	return nil
}

// ExtractEquipmentModifierCostType extracts the value from a string.
func ExtractEquipmentModifierCostType(str string) EquipmentModifierCostType {
	for _, enum := range AllEquipmentModifierCostType {
		if strings.EqualFold(enum.Key(), str) {
			return enum
		}
	}
	return 0
}
