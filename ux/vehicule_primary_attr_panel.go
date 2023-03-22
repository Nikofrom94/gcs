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
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

// PrimaryAttrPanel holds the contents of the primary attributes block on the sheet.
type VehiculePrimaryAttrPanel struct {
	unison.Panel
	targetMgr *TargetMgr
	prefix    string
	crc       uint64
}

// NewPrimaryAttrPanel creates a new primary attributes panel.
func NewVehiculePrimaryAttrPanel(vehicule *model.Vehicule, targetMgr *TargetMgr) *VehiculePrimaryAttrPanel {
	p := &VehiculePrimaryAttrPanel{
		targetMgr: targetMgr,
		prefix:    targetMgr.NextPrefix(),
	}
	p.Self = p
	p.SetLayout(&unison.FlexLayout{
		Columns:  3,
		HSpacing: 1,
		VSpacing: 1,
	})
	p.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
	})
	p.SetBorder(unison.NewCompoundBorder(&TitledBorder{Title: i18n.Text("Vehicule Statistics")}, unison.NewEmptyBorder(unison.Insets{
		Top:    1,
		Left:   2,
		Bottom: 1,
		Right:  2,
	})))
	/*
		p.DrawCallback = func(gc *unison.Canvas, rect unison.Rect) {
			gc.DrawRect(rect, unison.ContentColor.Paint(gc, rect, unison.Fill))
		}
	*/
	p.rebuild(vehicule)
	return p
}

func (p *VehiculePrimaryAttrPanel) rebuild(vehicule *model.Vehicule) {
	p.RemoveAllChildren()
	p.AddChild(p.createFirstColumn(vehicule))
	p.AddChild(p.createSecondColumn(vehicule))
	p.AddChild(p.createThirdColumn(vehicule))
}

// first column for vehicule statistics/attributes : ST/HP
func (p *VehiculePrimaryAttrPanel) createFirstColumn(vehicule *model.Vehicule) *unison.Panel {
	column := createColumn()

	title := i18n.Text(vehicule.ST.DisplayName)
	column.AddChild(NewPageLabelEnd(title))
	stField := NewStringPageField(p.targetMgr, p.prefix+"st", title,
		func() string { return vehicule.ST.ToString() },
		func(s string) { vehicule.ST.SetValue(s) })
	stField.ClientData()[SkipDeepSync] = true
	column.AddChild(stField)

	title = i18n.Text(vehicule.HP.DisplayName)
	column.AddChild(NewPageLabelEnd(title))
	hpField := NewStringPageField(p.targetMgr, p.prefix+"hp", title,
		func() string { return vehicule.HP.ToString() },
		func(s string) { vehicule.HP.SetValue(s) })
	hpField.ClientData()[SkipDeepSync] = true
	column.AddChild(hpField)

	return column
}

// second column for vehicule statistics/attributes : Hnd/SR
func (p *VehiculePrimaryAttrPanel) createSecondColumn(vehicule *model.Vehicule) *unison.Panel {
	column := createColumn()

	title := i18n.Text(vehicule.Hnd.DisplayName)
	column.AddChild(NewPageLabelEnd(title))
	hndField := NewStringPageField(p.targetMgr, p.prefix+"hnd", title,
		func() string { return vehicule.Hnd.ToString() },
		func(s string) { vehicule.Hnd.SetValue(s) })
	hndField.ClientData()[SkipDeepSync] = true
	column.AddChild(hndField)

	title = i18n.Text(vehicule.SR.DisplayName)
	column.AddChild(NewPageLabelEnd(title))
	srField := NewStringPageField(p.targetMgr, p.prefix+"sr", title,
		func() string { return vehicule.SR.ToString() },
		func(s string) { vehicule.SR.SetValue(s) })
	srField.ClientData()[SkipDeepSync] = true
	column.AddChild(srField)

	return column
}

// third column for vehicule statistics/attributes : Move/Lwt
func (p *VehiculePrimaryAttrPanel) createThirdColumn(vehicule *model.Vehicule) *unison.Panel {
	column := createColumn()

	title := i18n.Text(vehicule.Move.DisplayName)
	column.AddChild(NewPageLabelEnd(title))
	moveField := NewStringPageField(p.targetMgr, p.prefix+"move", title,
		func() string { return vehicule.Move.ToString() },
		func(s string) { vehicule.Move.SetValue(s) })
	moveField.ClientData()[SkipDeepSync] = true
	column.AddChild(moveField)

	title = i18n.Text(vehicule.LWt.DisplayName)
	column.AddChild(NewPageLabelEnd(title))
	lwtField := NewStringPageField(p.targetMgr, p.prefix+"lwt", title,
		func() string { return vehicule.LWt.ToString() },
		func(s string) { vehicule.LWt.SetValue(s) })
	lwtField.ClientData()[SkipDeepSync] = true
	column.AddChild(lwtField)

	return column
}
