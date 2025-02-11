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
	"fmt"

	"github.com/richardwilkes/gcs/v5/model"
	"github.com/richardwilkes/gcs/v5/model/fxp"
	"github.com/richardwilkes/rpgtools/dice"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

// Rebuildable defines the methods a rebuildable panel should provide.
type Rebuildable interface {
	unison.Paneler
	fmt.Stringer
	Rebuild(full bool)
}

// Syncer should be called to sync an object's UI state to its model.
type Syncer interface {
	Sync()
}

// DeepSync does a depth-first traversal of the panel and all of its descendents and calls Sync() on any Syncer objects
// it finds.
func DeepSync(panel unison.Paneler) {
	p := panel.AsPanel()
	for _, child := range p.Children() {
		DeepSync(child)
	}
	if syncer, ok := p.Self.(Syncer); ok {
		syncer.Sync()
	}
}

// ModifiableRoot marks the root of a modifable tree of components, typically a Dockable.
type ModifiableRoot interface {
	MarkModified(src unison.Paneler)
}

// MarkModified looks for a ModifiableRoot, starting at the panel. If found, it then called MarkModified() on it.
func MarkModified(panel unison.Paneler) {
	p := panel.AsPanel()
	for p != nil {
		if modifiable, ok := p.Self.(ModifiableRoot); ok {
			modifiable.MarkModified(panel)
			break
		}
		p = p.Parent()
	}
}

func addNameLabelAndField(parent *unison.Panel, fieldData *string) {
	addLabelAndStringField(parent, i18n.Text("Name"), "", fieldData)
}

func addSpecializationLabelAndField(parent *unison.Panel, fieldData *string) {
	addLabelAndStringField(parent, i18n.Text("Specialization"), "", fieldData)
}

func addPageRefLabelAndField(parent *unison.Panel, fieldData *string) {
	addLabelAndStringField(parent, i18n.Text("Page Reference"), pageRefTooltipText(), fieldData)
}

func addNotesLabelAndField(parent *unison.Panel, fieldData *string) {
	addLabelAndMultiLineStringField(parent, i18n.Text("Notes"), "", fieldData)
}

func addVTTNotesLabelAndField(parent *unison.Panel, fieldData *string) {
	addLabelAndMultiLineStringField(parent, i18n.Text("VTT Notes"),
		i18n.Text("Any notes for VTT use; see the instructions for your VVT to determine if/how these can be used"),
		fieldData)
}

func addUserDescLabelAndField(parent *unison.Panel, fieldData *string) {
	addLabelAndMultiLineStringField(parent, i18n.Text("User Description"),
		i18n.Text("Additional notes for your own reference. These only exist in character sheets and will be removed if transferred to a data list or template"),
		fieldData)
}

func addTechLevelRequired(parent *unison.Panel, fieldData **string, includeField bool) {
	tl := i18n.Text("Tech Level")
	var field *StringField
	if includeField {
		wrapper := addFlowWrapper(parent, tl, 2)
		field = NewStringField(nil, "", tl, func() string {
			if *fieldData == nil {
				return ""
			}
			return **fieldData
		}, func(value string) {
			if *fieldData == nil {
				return
			}
			**fieldData = value
			MarkModified(parent)
		})
		field.Tooltip = unison.NewTooltipWithText(techLevelInfo())
		if *fieldData == nil {
			field.SetEnabled(false)
		}
		field.SetMinimumTextWidthUsing("12^")
		wrapper.AddChild(field)
		parent = wrapper
	} else {
		parent.AddChild(NewFieldLeadingLabel(tl))
	}
	last := *fieldData
	required := last != nil
	parent.AddChild(NewCheckBox(nil, "", i18n.Text("Required"),
		func() unison.CheckState { return unison.CheckStateFromBool(required) },
		func(state unison.CheckState) {
			if required = state == unison.OnCheckState; required {
				if last == nil {
					var data string
					last = &data
				}
				*fieldData = last
				if field != nil {
					field.SetEnabled(true)
				}
			} else {
				last = *fieldData
				*fieldData = nil
				if field != nil {
					field.SetEnabled(false)
				}
			}
		}))
}

func addHitLocationChoicePopup(parent *unison.Panel, entity *model.Entity, prefix string, fieldData *string) *unison.PopupMenu[*model.HitLocationChoice] {
	choices, current := model.HitLocationChoices(entity, prefix, *fieldData)
	popup := addPopup(parent, choices, &current)
	popup.SelectionChangedCallback = func(p *unison.PopupMenu[*model.HitLocationChoice]) {
		if choice, ok := p.Selected(); ok {
			*fieldData = choice.Key
			MarkModified(parent)
		}
	}
	return popup
}

func addAttributeChoicePopup(parent *unison.Panel, entity *model.Entity, prefix string, fieldData *string, flags model.AttributeFlags) *unison.PopupMenu[*model.AttributeChoice] {
	choices, current := model.AttributeChoices(entity, prefix, flags, *fieldData)
	popup := addPopup(parent, choices, &current)
	popup.SelectionChangedCallback = func(p *unison.PopupMenu[*model.AttributeChoice]) {
		if choice, ok := p.Selected(); ok {
			*fieldData = choice.Key
			MarkModified(parent)
		}
	}
	return popup
}

func addDifficultyLabelAndFields(parent *unison.Panel, entity *model.Entity, difficulty *model.AttributeDifficulty) {
	wrapper := addFlowWrapper(parent, i18n.Text("Difficulty"), 3)
	addAttributeChoicePopup(wrapper, entity, "", &difficulty.Attribute, model.TenFlag)
	wrapper.AddChild(NewFieldTrailingLabel("/"))
	addPopup(wrapper, model.AllDifficulty, &difficulty.Difficulty)
}

func addTagsLabelAndField(parent *unison.Panel, fieldData *[]string) {
	addLabelAndListField(parent, i18n.Text("Tags"), i18n.Text("tags"), fieldData)
}

func addLabelAndListField(parent *unison.Panel, labelText, pluralForTooltip string, fieldData *[]string) {
	tooltip := fmt.Sprintf(i18n.Text("Separate multiple %s with commas"), pluralForTooltip)
	label := NewFieldLeadingLabel(labelText)
	if tooltip != "" {
		label.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	parent.AddChild(label)
	field := NewMultiLineStringField(nil, "", labelText,
		func() string { return model.CombineTags(*fieldData) },
		func(value string) {
			*fieldData = model.ExtractTags(value)
			parent.MarkForLayoutAndRedraw()
			MarkModified(parent)
		})
	if tooltip != "" {
		field.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	field.AutoScroll = false
	parent.AddChild(field)
}

func addLabelAndStringField(parent *unison.Panel, labelText, tooltip string, fieldData *string) *StringField {
	label := NewFieldLeadingLabel(labelText)
	if tooltip != "" {
		label.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	parent.AddChild(label)
	return addStringField(parent, labelText, tooltip, fieldData)
}

func addStringField(parent *unison.Panel, labelText, tooltip string, fieldData *string) *StringField {
	field := NewStringField(nil, "", labelText,
		func() string { return *fieldData },
		func(value string) {
			*fieldData = value
			MarkModified(parent)
		})
	if tooltip != "" {
		field.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	parent.AddChild(field)
	return field
}

func addLabelAndMultiLineStringField(parent *unison.Panel, labelText, tooltip string, fieldData *string) {
	label := NewFieldLeadingLabel(labelText)
	if tooltip != "" {
		label.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	parent.AddChild(label)
	field := NewMultiLineStringField(nil, "", labelText,
		func() string { return *fieldData },
		func(value string) {
			*fieldData = value
			parent.MarkForLayoutAndRedraw()
			MarkModified(parent)
		})
	if tooltip != "" {
		field.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	field.AutoScroll = false
	parent.AddChild(field)
}

func addIntegerField(parent *unison.Panel, targetMgr *TargetMgr, targetKey, labelText, tooltip string, fieldData *int, min, max int) *IntegerField {
	field := NewIntegerField(targetMgr, targetKey, labelText,
		func() int { return *fieldData },
		func(value int) {
			*fieldData = value
			MarkModified(parent)
		}, min, max, false, false)
	if tooltip != "" {
		field.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	parent.AddChild(field)
	return field
}

func addLabelAndDecimalField(parent *unison.Panel, targetMgr *TargetMgr, targetKey, labelText, tooltip string, fieldData *fxp.Int, min, max fxp.Int) *DecimalField {
	label := NewFieldLeadingLabel(labelText)
	if tooltip != "" {
		label.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	parent.AddChild(label)
	return addDecimalField(parent, targetMgr, targetKey, labelText, tooltip, fieldData, min, max)
}

func addDecimalField(parent *unison.Panel, targetMgr *TargetMgr, targetKey, labelText, tooltip string, fieldData *fxp.Int, min, max fxp.Int) *DecimalField {
	field := NewDecimalField(targetMgr, targetKey, labelText,
		func() fxp.Int { return *fieldData },
		func(value fxp.Int) {
			*fieldData = value
			MarkModified(parent)
		}, min, max, false, false)
	if tooltip != "" {
		field.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	parent.AddChild(field)
	return field
}

func addWeightField(parent *unison.Panel, targetMgr *TargetMgr, targetKey, labelText, tooltip string, entity *model.Entity, fieldData *model.Weight, noMinWidth bool) *WeightField {
	field := NewWeightField(targetMgr, targetKey, labelText, entity,
		func() model.Weight { return *fieldData },
		func(value model.Weight) {
			*fieldData = value
			MarkModified(parent)
		}, 0, model.Weight(fxp.Max), noMinWidth)
	if tooltip != "" {
		field.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	parent.AddChild(field)
	return field
}

func addCheckBox(parent *unison.Panel, labelText string, fieldData *bool) *CheckBox {
	checkBox := NewCheckBox(nil, "", labelText,
		func() unison.CheckState { return unison.CheckStateFromBool(*fieldData) },
		func(state unison.CheckState) { *fieldData = state == unison.OnCheckState })
	parent.AddChild(checkBox)
	return checkBox
}

func addInvertedCheckBox(parent *unison.Panel, labelText string, fieldData *bool) {
	parent.AddChild(NewCheckBox(nil, "", labelText,
		func() unison.CheckState { return unison.CheckStateFromBool(!*fieldData) },
		func(state unison.CheckState) { *fieldData = state == unison.OffCheckState }))
}

func addFlowWrapper(parent *unison.Panel, labelText string, count int) *unison.Panel {
	parent.AddChild(NewFieldLeadingLabel(labelText))
	wrapper := unison.NewPanel()
	wrapper.SetLayout(&unison.FlexLayout{
		Columns:  count,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
		VAlign:   unison.MiddleAlignment,
	})
	parent.AddChild(wrapper)
	return wrapper
}

func addLabelAndNullableDice(parent *unison.Panel, labelText, tooltip string, fieldData **dice.Dice) *StringField {
	var data string
	if *fieldData != nil {
		data = (*fieldData).String()
	}
	label := NewFieldLeadingLabel(labelText)
	parent.AddChild(label)
	field := NewStringField(nil, "", labelText,
		func() string { return data },
		func(value string) {
			data = value
			if value == "" {
				*fieldData = nil
			} else {
				*fieldData = dice.New(data)
			}
			MarkModified(parent)
		})
	if tooltip != "" {
		label.Tooltip = unison.NewTooltipWithText(tooltip)
		field.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	parent.AddChild(field)
	return field
}

func addLabelAndPopup[T comparable](parent *unison.Panel, labelText, tooltip string, choices []T, fieldData *T) *unison.PopupMenu[T] {
	label := NewFieldLeadingLabel(labelText)
	if tooltip != "" {
		label.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	parent.AddChild(label)
	return addPopup[T](parent, choices, fieldData)
}

func addPopup[T comparable](parent *unison.Panel, choices []T, fieldData *T) *unison.PopupMenu[T] {
	popup := unison.NewPopupMenu[T]()
	for _, one := range choices {
		popup.AddItem(one)
	}
	popup.Select(*fieldData)
	popup.SelectionChangedCallback = func(p *unison.PopupMenu[T]) {
		if item, ok := p.Selected(); ok {
			*fieldData = item
			MarkModified(parent)
		}
	}
	parent.AddChild(popup)
	return popup
}

func addBoolPopup(parent *unison.Panel, trueChoice, falseChoice string, fieldData *bool) *unison.PopupMenu[string] {
	popup := unison.NewPopupMenu[string]()
	popup.AddItem(trueChoice)
	popup.AddItem(falseChoice)
	if *fieldData {
		popup.SelectIndex(0)
	} else {
		popup.SelectIndex(1)
	}
	popup.SelectionChangedCallback = func(p *unison.PopupMenu[string]) {
		*fieldData = p.SelectedIndex() == 0
		MarkModified(parent)
	}
	parent.AddChild(popup)
	return popup
}

func addHasPopup(parent *unison.Panel, has *bool) {
	addBoolPopup(parent, i18n.Text("has"), i18n.Text("doesn't have"), has)
}

func adjustFieldBlank(field unison.Paneler, blank bool) {
	panel := field.AsPanel()
	panel.SetEnabled(!blank)
	if blank {
		panel.DrawOverCallback = func(gc *unison.Canvas, _ unison.Rect) {
			var ink unison.Ink
			if f, ok := panel.Self.(*unison.Field); ok {
				ink = f.BackgroundInk
			} else {
				ink = unison.DefaultFieldTheme.BackgroundInk
			}
			r := panel.ContentRect(false)
			gc.DrawRect(r, ink.Paint(gc, r, unison.Fill))
		}
	} else {
		panel.DrawOverCallback = nil
	}
}

func adjustPopupBlank[T comparable](popup *unison.PopupMenu[T], blank bool) {
	popup.SetEnabled(!blank)
	if blank {
		popup.DrawOverCallback = func(gc *unison.Canvas, _ unison.Rect) {
			unison.DrawRoundedRectBase(gc, popup.ContentRect(false), popup.CornerRadius, 1, popup.BackgroundInk, popup.EdgeInk)
		}
	} else {
		popup.DrawOverCallback = nil
	}
}

func addNameCriteriaPanel(parent *unison.Panel, strCriteria *model.StringCriteria, hSpan int, includeEmptyFiller bool) {
	prefix := i18n.Text("whose name")
	addStringCriteriaPanel(parent, prefix, prefix, i18n.Text("Name Qualifier"), strCriteria, hSpan, includeEmptyFiller)
}

func addSpecializationCriteriaPanel(parent *unison.Panel, strCriteria *model.StringCriteria, hSpan int, includeEmptyFiller bool) {
	prefix := i18n.Text("and whose specialization")
	addStringCriteriaPanel(parent, prefix, prefix, i18n.Text("Specialization Qualifier"), strCriteria, hSpan, includeEmptyFiller)
}

func addTagCriteriaPanel(parent *unison.Panel, strCriteria *model.StringCriteria, hSpan int, includeEmptyFiller bool) {
	addStringCriteriaPanel(parent, i18n.Text("and at least one tag"), i18n.Text("and all tags"), i18n.Text("Tag Qualifier"), strCriteria, hSpan, includeEmptyFiller)
}

func addNotesCriteriaPanel(parent *unison.Panel, strCriteria *model.StringCriteria, hSpan int, includeEmptyFiller bool) {
	prefix := i18n.Text("and whose notes")
	addStringCriteriaPanel(parent, prefix, prefix, i18n.Text("Notes Qualifier"), strCriteria, hSpan, includeEmptyFiller)
}

func addStringCriteriaPanel(parent *unison.Panel, prefix, notPrefix, undoTitle string, strCriteria *model.StringCriteria, hSpan int, includeEmptyFiller bool) (*unison.PopupMenu[string], *StringField) {
	if includeEmptyFiller {
		parent.AddChild(unison.NewPanel())
	}
	panel := unison.NewPanel()
	panel.SetLayout(&unison.FlexLayout{
		Columns:  2,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
		VAlign:   unison.MiddleAlignment,
	})
	panel.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  hSpan,
		HAlign: unison.FillAlignment,
		HGrab:  true,
	})
	var criteriaField *StringField
	popup := unison.NewPopupMenu[string]()
	for _, one := range model.PrefixedStringCompareTypeChoices(prefix, notPrefix) {
		popup.AddItem(one)
	}
	popup.SelectIndex(model.ExtractStringCompareTypeIndex(string(strCriteria.Compare)))
	popup.SelectionChangedCallback = func(p *unison.PopupMenu[string]) {
		strCriteria.Compare = model.AllStringCompareTypes[p.SelectedIndex()]
		adjustFieldBlank(criteriaField, strCriteria.Compare == model.AnyString)
		MarkModified(panel)
	}
	panel.AddChild(popup)
	criteriaField = addStringField(panel, undoTitle, "", &strCriteria.Qualifier)
	adjustFieldBlank(criteriaField, strCriteria.Compare == model.AnyString)
	parent.AddChild(panel)
	return popup, criteriaField
}

func addLevelCriteriaPanel(parent *unison.Panel, targetMgr *TargetMgr, targetKey string, numCriteria *model.NumericCriteria, hSpan int, includeEmptyFiller bool) {
	addNumericCriteriaPanel(parent, targetMgr, targetKey, i18n.Text("and whose level"), i18n.Text("Level Qualifier"),
		numCriteria, 0, fxp.Thousand, hSpan, false, includeEmptyFiller)
}

func addNumericCriteriaPanel(parent *unison.Panel, targetMgr *TargetMgr, targetKey, prefix, undoTitle string, numCriteria *model.NumericCriteria, min, max fxp.Int, hSpan int, integerOnly, includeEmptyFiller bool) (popup *unison.PopupMenu[string], field unison.Paneler) {
	if includeEmptyFiller {
		parent.AddChild(unison.NewPanel())
	}
	panel := unison.NewPanel()
	panel.SetLayout(&unison.FlexLayout{
		Columns:  2,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
		VAlign:   unison.MiddleAlignment,
	})
	panel.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  hSpan,
		HAlign: unison.FillAlignment,
		HGrab:  true,
	})
	popup = unison.NewPopupMenu[string]()
	for _, one := range model.PrefixedNumericCompareTypeChoices(prefix) {
		popup.AddItem(one)
	}
	popup.SelectIndex(model.ExtractNumericCompareTypeIndex(string(numCriteria.Compare)))
	popup.SelectionChangedCallback = func(p *unison.PopupMenu[string]) {
		numCriteria.Compare = model.AllNumericCompareTypes[p.SelectedIndex()]
		adjustFieldBlank(field, numCriteria.Compare == model.AnyNumber)
		MarkModified(panel)
	}
	panel.AddChild(popup)
	if integerOnly {
		field = NewIntegerField(targetMgr, targetKey, undoTitle,
			func() int { return fxp.As[int](numCriteria.Qualifier) },
			func(value int) {
				numCriteria.Qualifier = fxp.From(value)
				MarkModified(panel)
			}, fxp.As[int](min), fxp.As[int](max), false, false)
		panel.AddChild(field)
	} else {
		field = addDecimalField(panel, targetMgr, targetKey, undoTitle, "", &numCriteria.Qualifier, min, max)
	}
	adjustFieldBlank(field, numCriteria.Compare == model.AnyNumber)
	parent.AddChild(panel)
	return popup, field
}

func addWeightCriteriaPanel(parent *unison.Panel, targetMgr *TargetMgr, targetKey string, entity *model.Entity, weightCriteria *model.WeightCriteria) {
	popup := unison.NewPopupMenu[string]()
	for _, one := range model.PrefixedNumericCompareTypeChoices(i18n.Text("which")) {
		popup.AddItem(one)
	}
	popup.SelectIndex(model.ExtractNumericCompareTypeIndex(string(weightCriteria.Compare)))
	parent.AddChild(popup)
	field := addWeightField(parent, targetMgr, targetKey, i18n.Text("Weight Qualifier"), "", entity,
		&weightCriteria.Qualifier, false)
	popup.SelectionChangedCallback = func(p *unison.PopupMenu[string]) {
		weightCriteria.Compare = model.AllNumericCompareTypes[p.SelectedIndex()]
		adjustFieldBlank(field, weightCriteria.Compare == model.AnyNumber)
		MarkModified(parent)
	}
	adjustFieldBlank(field, weightCriteria.Compare == model.AnyNumber)
	parent.SetLayout(&unison.FlexLayout{
		Columns:  len(parent.Children()),
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
}

func addQuantityCriteriaPanel(parent *unison.Panel, targetMgr *TargetMgr, targetKey string, numCriteria *model.NumericCriteria) {
	choices := []string{
		i18n.Text("exactly"),
		i18n.Text("at least"),
		i18n.Text("at most"),
	}
	var numType string
	switch numCriteria.Compare {
	case model.AtLeastNumber:
		numType = choices[1]
	case model.AtMostNumber:
		numType = choices[2]
	default:
		numType = choices[0]
	}
	popup := unison.NewPopupMenu[string]()
	for _, one := range choices {
		popup.AddItem(one)
	}
	popup.Select(numType)
	popup.SelectionChangedCallback = func(p *unison.PopupMenu[string]) {
		switch p.SelectedIndex() {
		case 0:
			numCriteria.Compare = model.EqualsNumber
		case 1:
			numCriteria.Compare = model.AtLeastNumber
		case 2:
			numCriteria.Compare = model.AtMostNumber
		}
		MarkModified(parent)
	}
	parent.AddChild(popup)
	parent.AddChild(NewIntegerField(targetMgr, targetKey, i18n.Text("Quantity Criteria"),
		func() int { return fxp.As[int](numCriteria.Qualifier) },
		func(value int) {
			numCriteria.Qualifier = fxp.From(value)
			MarkModified(parent)
		}, 0, 9999, false, false))
}

func addLeveledAmountPanel(parent *unison.Panel, targetMgr *TargetMgr, targetKey, title string, amount *model.LeveledAmount) {
	parent.AddChild(NewDecimalField(targetMgr, targetKey, i18n.Text("Amount"),
		func() fxp.Int { return amount.Amount },
		func(value fxp.Int) {
			amount.Amount = value
			MarkModified(parent)
		}, fxp.Min, fxp.Max, true, false))
	addCheckBox(parent, title, &amount.PerLevel)
}

func addTemplateChoices(parent *unison.Panel, targetmgr *TargetMgr, targetKey string, picker **model.TemplatePicker) {
	if *picker == nil {
		*picker = &model.TemplatePicker{}
	}
	last := (*picker).Type
	wrapper := addFlowWrapper(parent, i18n.Text("Template Choices"), 3)
	templatePickerTypePopup := addPopup(wrapper, model.AllTemplatePickerType, &(*picker).Type)
	text := i18n.Text("Template Choice Quantifier")
	popup, field := addNumericCriteriaPanel(wrapper, targetmgr, targetKey, "", text, &(*picker).Qualifier, fxp.Min,
		fxp.Max, 1, false, false)
	templatePickerTypePopup.SelectionChangedCallback = func(p *unison.PopupMenu[model.TemplatePickerType]) {
		if item, ok := p.Selected(); ok {
			(*picker).Type = item
			if last == model.NotApplicableTemplatePickerType && item != model.NotApplicableTemplatePickerType {
				(*picker).Qualifier.Qualifier = fxp.One
				field.(Syncer).Sync()
			}
			last = item
			adjustFieldBlank(field, item == model.NotApplicableTemplatePickerType || (*picker).Qualifier.Compare == model.AnyNumber)
			adjustPopupBlank(popup, item == model.NotApplicableTemplatePickerType)
			MarkModified(parent)
		}
	}
	adjustFieldBlank(field, (*picker).Type == model.NotApplicableTemplatePickerType)
}

// WrapWithSpan wraps a number of children with a single panel that request to fill in span number of columns.
func WrapWithSpan(span int, children ...unison.Paneler) *unison.Panel {
	wrapper := unison.NewPanel()
	wrapper.SetLayout(&unison.FlexLayout{
		Columns:  len(children),
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	wrapper.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  span,
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
		HGrab:  true,
	})
	for _, child := range children {
		wrapper.AddChild(child)
	}
	return wrapper
}

func pageRefTooltipText() string {
	return i18n.Text(`A reference to the book and page the item appears on e.g. B22 would refer to "Basic Set", page 22`)
}

func techLevelInfo() string {
	return i18n.Text(`TL0: Stone Age (Prehistory)
TL1: Bronze Age (3500 B.C.+)
TL2: Iron Age (1200 B.C.+)
TL3: Medieval (600 A.D.+)
TL4: Age of Sail (1450+)
TL5: Industrial Revolution (1730+)
TL6: Mechanized Age (1880+)
TL7: Nuclear Age (1940+)
TL8: Digital Age (1980+)
TL9: Microtech Age (2025+?)
TL10: Robotic Age (2070+?)
TL11: Age of Exotic Matter
TL12: Anything Goes`)
}
