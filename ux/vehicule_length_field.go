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

package ux

import (
	"github.com/richardwilkes/gcs/v5/model"
	"github.com/richardwilkes/gcs/v5/model/fxp"
	"github.com/richardwilkes/unison"
)

// VehiculeLengthField is field that holds a length value.
type VehiculeLengthField = NumericField[model.Length]

// NewVehiculeLengthField creates a new field that holds a fixed-point number.
func NewVehiculeLengthField(targetMgr *TargetMgr, targetKey, undoTitle string, lengthUnit model.LengthUnits, get func() model.Length, set func(model.Length), min, max model.Length, noMinWidth bool) *LengthField {
	var getPrototypes func(min, max model.Length) []model.Length
	if !noMinWidth {
		getPrototypes = func(min, max model.Length) []model.Length {
			if min == model.Length(fxp.Min) {
				min = model.Length(-fxp.One)
			}
			min = model.Length(fxp.Int(min).Trunc() + fxp.One - 1)
			if max == model.Length(fxp.Max) {
				max = model.Length(fxp.One)
			}
			max = model.Length(fxp.Int(max).Trunc() + fxp.One - 1)
			return []model.Length{min, model.Length(fxp.Two - 1), max}
		}
	}
	format := func(value model.Length) string {
		return lengthUnit.Format(value)
	}
	extract := func(s string) (model.Length, error) {
		return model.LengthFromString(s, lengthUnit)
	}
	f := NewNumericField[model.Length](targetMgr, targetKey, undoTitle, getPrototypes, get, set, format, extract, min, max)
	f.RuneTypedCallback = f.DefaultRuneTyped
	return f
}

func NewVehiculeLeightPageField(targetMgr *TargetMgr, targetKey, undoTitle string, lengthUnit model.LengthUnits, get func() model.Length, set func(model.Length), min, max model.Length, noMinWidth bool) *VehiculeLengthField {
	field := NewVehiculeLengthField(targetMgr, targetKey, undoTitle, lengthUnit, get, set, min, max, noMinWidth)
	field.Font = model.PageFieldPrimaryFont
	field.FocusedBorder = unison.NewLineBorder(unison.AccentColor, 0, unison.Insets{Bottom: 1}, false)
	field.UnfocusedBorder = unison.NewLineBorder(unison.ControlEdgeColor, 0, unison.Insets{Bottom: 1}, false)
	field.SetBorder(field.UnfocusedBorder)
	field.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
	})
	return field
}
