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
	"github.com/richardwilkes/gcs/v5/svg"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
)

// EditTraitModifier displays the editor for a trait modifier.
func EditTraitModifier(owner Rebuildable, modifier *model.TraitModifier) {
	displayEditor[*model.TraitModifier, *model.TraitModifierEditData](owner, modifier, svg.GCSTraitModifiers,
		"md:Help/Interface/Trait Modifiers", nil, initTraitModifierEditor)
}

func initTraitModifierEditor(e *editor[*model.TraitModifier, *model.TraitModifierEditData], content *unison.Panel) func() {
	addNameLabelAndField(content, &e.editorData.Name)
	addNotesLabelAndField(content, &e.editorData.LocalNotes)
	addVTTNotesLabelAndField(content, &e.editorData.VTTNotes)
	if !e.target.Container() {
		content.AddChild(unison.NewPanel())
		addInvertedCheckBox(content, i18n.Text("Enabled"), &e.editorData.Disabled)
		costLabel := i18n.Text("Cost")
		wrapper := addFlowWrapper(content, costLabel, 3)
		addDecimalField(wrapper, nil, "", costLabel, "", &e.editorData.Cost, -fxp.MaxBasePoints, fxp.MaxBasePoints)
		costTypePopup := addCostTypePopup(wrapper, e)
		affectsPopup := addPopup(wrapper, model.AllAffects, &e.editorData.Affects)
		levels := addLabelAndDecimalField(content, nil, "", i18n.Text("Level"), "", &e.editorData.Levels, 0, fxp.Thousand)
		adjustFieldBlank(levels, !e.target.HasLevels())
		total := NewNonEditableField(func(field *NonEditableField) {
			enabled := true
			switch costTypePopup.SelectedIndex() - 1 {
			case -1:
				field.Text = e.editorData.Cost.Mul(e.editorData.Levels).StringWithSign() + model.PercentageTraitModifierCostType.String()
			case int(model.PercentageTraitModifierCostType):
				field.Text = e.editorData.Cost.StringWithSign() + model.PercentageTraitModifierCostType.String()
			case int(model.PointsTraitModifierCostType):
				field.Text = e.editorData.Cost.StringWithSign()
			case int(model.MultiplierTraitModifierCostType):
				field.Text = model.MultiplierTraitModifierCostType.String() + e.editorData.Cost.String()
				affectsPopup.Select(model.TotalAffects)
				enabled = false
			default:
				jot.Errorf("unhandled cost type popup index: %d", costTypePopup.SelectedIndex())
				field.Text = e.editorData.Cost.StringWithSign() + model.PercentageTraitModifierCostType.String()
			}
			affectsPopup.SetEnabled(enabled)
			field.MarkForLayoutAndRedraw()
		})
		insets := total.Border().Insets()
		total.SetLayoutData(&unison.FlexLayoutData{
			MinSize: unison.NewSize(total.Font.SimpleWidth((-fxp.MaxBasePoints*2).String())+insets.Left+insets.Right, 0),
		})
		content.AddChild(NewFieldLeadingLabel(i18n.Text("Total")))
		content.AddChild(total)
		costTypePopup.SelectionChangedCallback = func(popup *unison.PopupMenu[string]) {
			index := popup.SelectedIndex()
			if index == 0 {
				e.editorData.CostType = model.PercentageTraitModifierCostType
				if e.editorData.Levels < fxp.One {
					levels.SetText("1")
				}
			} else {
				e.editorData.CostType = model.AllTraitModifierCostType[index-1]
				e.editorData.Levels = 0
			}
			adjustFieldBlank(levels, index != 0)
			MarkModified(wrapper)
		}
	}
	addTagsLabelAndField(content, &e.editorData.Tags)
	addPageRefLabelAndField(content, &e.editorData.PageRef)
	if !e.target.Container() {
		content.AddChild(newFeaturesPanel(e.target.Entity, e.target, &e.editorData.Features))
	}
	return nil
}

func addCostTypePopup(parent *unison.Panel, e *editor[*model.TraitModifier, *model.TraitModifierEditData]) *unison.PopupMenu[string] {
	popup := unison.NewPopupMenu[string]()
	popup.AddItem(i18n.Text("% per level"))
	for _, one := range model.AllTraitModifierCostType {
		popup.AddItem(one.String())
	}
	if e.target.HasLevels() {
		popup.SelectIndex(0)
	} else {
		popup.Select(e.editorData.CostType.String())
	}
	parent.AddChild(popup)
	return popup
}
