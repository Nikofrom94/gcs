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
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

// VehiculeDescriptionPanel holds the contents of the description block on the vehicule sheet.
type VehiculeDescriptionPanel struct {
	unison.Panel
	vehicule  *model.Vehicule
	targetMgr *TargetMgr
	prefix    string
}

// NewVehiculeDescriptionPanel creates a new description panel.
func NewVehiculeDescriptionPanel(vehicule *model.Vehicule, targetMgr *TargetMgr) *VehiculeDescriptionPanel {
	d := &VehiculeDescriptionPanel{
		vehicule:  vehicule,
		targetMgr: targetMgr,
		prefix:    targetMgr.NextPrefix(),
	}
	d.Self = d
	d.SetLayout(&unison.FlexLayout{
		Columns:  3,
		HSpacing: 4,
	})
	d.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  2,
		HAlign: unison.FillAlignment,
		HGrab:  true,
	})
	d.SetBorder(unison.NewCompoundBorder(&TitledBorder{Title: i18n.Text("Description")}, unison.NewEmptyBorder(unison.Insets{
		Top:    1,
		Left:   2,
		Bottom: 1,
		Right:  2,
	})))
	d.DrawCallback = func(gc *unison.Canvas, rect unison.Rect) {
		gc.DrawRect(rect, unison.ContentColor.Paint(gc, rect, unison.Fill))
	}
	d.AddChild(d.createColumn1())
	d.AddChild(d.createColumn2())
	d.AddChild(d.createColumn3())
	return d
}

func (d *VehiculeDescriptionPanel) createColumn1() *unison.Panel {
	column := createColumn()
	/*
		title := i18n.Text("Gender")
		genderField := NewStringPageField(d.targetMgr, d.prefix+"gender", title,
			func() string { return d.vehicule.Profile.Gender },
			func(s string) { d.vehicule.Profile.Gender = s })
		column.AddChild(NewPageLabelWithRandomizer(title,
			i18n.Text("Randomize the gender using the current ancestry"), func() {
				d.vehicule.Profile.Gender = d.vehicule.Ancestry().RandomGender(d.vehicule.Profile.Gender)
				SetTextAndMarkModified(genderField.Field, d.vehicule.Profile.Gender)
			}))
		genderField.ClientData()[SkipDeepSync] = true
		column.AddChild(genderField)
	*/
	title := i18n.Text("Age")
	ageField := NewStringPageField(d.targetMgr, d.prefix+"age", title,
		func() string { return d.vehicule.Profile.Age },
		func(s string) { d.vehicule.Profile.Age = s })
	column.AddChild(NewPageLabelWithRandomizer(title,
		i18n.Text("Randomize the age using the current ancestry"), func() {

		}))
	ageField.ClientData()[SkipDeepSync] = true
	column.AddChild(ageField)

	title = i18n.Text("Launchday")
	birthdayField := NewStringPageField(d.targetMgr, d.prefix+"birthday", title,
		func() string { return d.vehicule.Profile.Birthday },
		func(s string) { d.vehicule.Profile.Birthday = s })
	column.AddChild(NewPageLabelWithRandomizer(title,
		i18n.Text("Randomize the birthday using the current calendar"), func() {

		}))
	birthdayField.ClientData()[SkipDeepSync] = true
	column.AddChild(birthdayField)
	/*
		title = i18n.Text("Religion")
		column.AddChild(NewPageLabelEnd(title))
		religionField := NewStringPageField(d.targetMgr, d.prefix+"religion", title,
			func() string { return d.vehicule.Profile.Religion },
			func(s string) { d.vehicule.Profile.Religion = s })
		religionField.ClientData()[SkipDeepSync] = true
		column.AddChild(religionField)
	*/
	return column
}

func (d *VehiculeDescriptionPanel) createColumn2() *unison.Panel {
	column := createColumn()

	title := i18n.Text("Leight")
	heightField := NewVehiculeLeightPageField(d.targetMgr, d.prefix+"leight", title, d.vehicule.SheetSettings.DefaultLengthUnits,
		func() model.Length { return d.vehicule.Profile.Height },
		func(v model.Length) { d.vehicule.Profile.Height = v }, 0, model.Length(fxp.Max), true)
	column.AddChild(NewPageLabelWithRandomizer(title,
		i18n.Text("Randomize the height using the current ancestry"), func() {

		}))
	heightField.ClientData()[SkipDeepSync] = true
	column.AddChild(heightField)
	/*
		title = i18n.Text("Weight")
		weightField := NewWeightPageField(d.targetMgr, d.prefix+"weight", title, d.vehicule,
			func() model.Weight { return d.vehicule.Profile.Weight },
			func(v model.Weight) { d.vehicule.Profile.Weight = v }, 0, model.Weight(fxp.Max), true)
		column.AddChild(NewPageLabelWithRandomizer(title,
			i18n.Text("Randomize the weight using the current ancestry"), func() {
			}))
		weightField.ClientData()[SkipDeepSync] = true
		column.AddChild(weightField)
	*/
	title = i18n.Text("Size")
	column.AddChild(NewPageLabelEnd(title))
	field := NewIntegerPageField(d.targetMgr, d.prefix+"size", title,
		func() int { return d.vehicule.Profile.AdjustedSizeModifier() },
		func(v int) { d.vehicule.Profile.SetAdjustedSizeModifier(v) }, -99, 99, true, false)
	field.HAlign = unison.StartAlignment
	column.AddChild(field)

	title = i18n.Text("TL")
	column.AddChild(NewPageLabelEnd(title))
	tlField := NewStringPageField(d.targetMgr, d.prefix+"tl", title,
		func() string { return d.vehicule.Profile.TechLevel },
		func(s string) { d.vehicule.Profile.TechLevel = s })
	tlField.Tooltip = unison.NewTooltipWithText(techLevelInfo())
	column.AddChild(tlField)

	return column
}

func (d *VehiculeDescriptionPanel) createColumn3() *unison.Panel {
	column := createColumn()

	// Loaded Weight field
	title := i18n.Text(d.vehicule.LWt.DisplayName)
	lwtField := NewStringPageField(d.targetMgr, d.prefix+"Lwt", title,
		func() string { return d.vehicule.LWt.ToString() },
		func(s string) { d.vehicule.LWt.SetValue(s) })
	column.AddChild(NewPageLabelWithRandomizer(title,
		i18n.Text("Randomize the lwt using the current ancestry"), func() {

		}))
	lwtField.ClientData()[SkipDeepSync] = true
	column.AddChild(lwtField)

	// Load field
	title = i18n.Text(d.vehicule.Load.DisplayName)
	loadField := NewStringPageField(d.targetMgr, d.prefix+"load", title,
		func() string { return d.vehicule.Load.ToString() },
		func(s string) { d.vehicule.Load.SetValue(s) })
	column.AddChild(NewPageLabelWithRandomizer(title,
		i18n.Text("Randomize the load using the current ancestry"), func() {

		}))
	loadField.ClientData()[SkipDeepSync] = true
	column.AddChild(loadField)

	// Occupancy field
	title = i18n.Text(d.vehicule.Occ.DisplayName)
	occField := NewStringPageField(d.targetMgr, d.prefix+"occ", title,
		func() string { return d.vehicule.Occ.Value },
		func(s string) { d.vehicule.Occ.Value = s })
	column.AddChild(NewPageLabelWithRandomizer(title,
		i18n.Text("Randomize the occ using the current ancestry"), func() {

		}))
	occField.ClientData()[SkipDeepSync] = true
	column.AddChild(occField)

	// Range field
	title = i18n.Text(d.vehicule.Range.DisplayName)
	rangeField := NewStringPageField(d.targetMgr, d.prefix+"range", title,
		func() string { return d.vehicule.Range.ToString() },
		func(s string) { d.vehicule.Range.SetValue(s) })
	column.AddChild(NewPageLabelWithRandomizer(title,
		i18n.Text("Randomize the range using the current ancestry"), func() {

		}))
	rangeField.ClientData()[SkipDeepSync] = true
	column.AddChild(rangeField)

	return column
}
