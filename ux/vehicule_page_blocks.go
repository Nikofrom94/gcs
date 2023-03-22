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
	"github.com/richardwilkes/unison"
)

func createVehiculePageTopBlock(vehicule *model.Vehicule, targetMgr *TargetMgr) (page *VehiculePage, modifiedFunc func()) {
	page = NewVehiculePage(vehicule)
	var top *unison.Panel
	top, modifiedFunc = createVehiculePageFirstRow(vehicule, targetMgr)
	page.AddChild(top)
	page.AddChild(createVehiculePageSecondRow(vehicule, targetMgr))
	return page, modifiedFunc
}

func createVehiculePageFirstRow(vehicule *model.Vehicule, targetMgr *TargetMgr) (top *unison.Panel, modifiedFunc func()) {
	right := unison.NewPanel()
	right.SetLayout(&unison.FlexLayout{
		Columns:  3,
		HSpacing: 1,
		VSpacing: 1,
		HAlign:   unison.FillAlignment,
		VAlign:   unison.FillAlignment,
	})
	right.AddChild(NewVehiculeIdentityPanel(vehicule, targetMgr))
	miscPanel := NewVehiculeMiscPanel(vehicule, targetMgr)
	right.AddChild(miscPanel)
	//right.AddChild(NewPointsPanel(entity, targetMgr))
	right.AddChild(NewVehiculeDescriptionPanel(vehicule, targetMgr))

	top = unison.NewPanel()
	portraitPanel := NewVehiculePortraitPanel(vehicule)
	top.SetLayout(&VehiculePortraitLayout{
		portrait: portraitPanel,
		rest:     right,
	})
	top.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
	})
	top.AddChild(portraitPanel)
	top.AddChild(right)

	return top, miscPanel.UpdateModified
}

func createVehiculePageSecondRow(vehicule *model.Vehicule, targetMgr *TargetMgr) *unison.Panel {
	p := unison.NewPanel()
	p.SetLayout(&unison.FlexLayout{
		Columns:  1,
		HSpacing: 1,
		VSpacing: 1,
		HAlign:   unison.FillAlignment,
		VAlign:   unison.FillAlignment,
	})
	p.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
	})

	p.AddChild(NewVehiculePrimaryAttrPanel(vehicule, targetMgr))

	return p
}
