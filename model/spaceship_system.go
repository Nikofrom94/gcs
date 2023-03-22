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
	"github.com/richardwilkes/json"
	"github.com/richardwilkes/toolbox/i18n"
)

// SpaceshipSystemData holds the Hitlocation data that gets written to disk.
type SpaceshipSystemData struct {
	SysID       string `json:"id"`
	TL          int    `json:"tl,omitempty"`
	ChoiceName  string `json:"choice_name"`
	TableName   string `json:"table_name"`
	Slots       int    `json:"slots,omitempty"`
	HitPenalty  int    `json:"hit_penalty,omitempty"`
	DRBonus     int    `json:"dr_bonus,omitempty"`
	Description string `json:"description,omitempty"`
	IsCore      bool   `json:"is_core"`
	IsHE        bool   `json:"is_he"` // is high energy system
	SM          int    `json:"sm,omitempty"`
	LoadedMass  int    `json:"loaded_mass,omitempty"`
	Length      int    `json:"length,omitempty"`
	dST         int    `json:"dst,omitempty"`
	HP          int    `json:"hp,omitempty"`
	Hnd         int    `json:"hnd,omitempty"`
	SR          int    `json:"sr,omitempty"`
	Workspaces  int    `json:"workspaces,omitempty"`
}

func NewSpaceshipSystemData() *SpaceshipSystemData {
	s := SpaceshipSystemData{}
	s.IsCore = false
	s.IsHE = false
	s.Workspaces = 0
	return &s
}

// SpaceshipSystem holds a single hit location.
type SpaceshipSystem struct {
	SpaceshipSystemData
	RollRange   string
	KeyPrefix   string
	owningTable *Body
}

func NewSpaceshipSystem() *SpaceshipSystem {
	s := SpaceshipSystem{
		SpaceshipSystemData: *NewSpaceshipSystemData(),
	}
	return &s
}

// NewSpaceshipLocation creates a new hit location.
func NewSpaceshipLocation(entity *Entity, keyPrefix string) *SpaceshipSystem {
	return &SpaceshipSystem{
		SpaceshipSystemData: SpaceshipSystemData{
			SysID:      "id",
			ChoiceName: i18n.Text("untitled choice"),
			TableName:  i18n.Text("untitled location"),
		},
		KeyPrefix: keyPrefix,
	}
}

// Clone a copy of this.
func (h *SpaceshipSystem) Clone() *SpaceshipSystem {
	clone := *h
	return &clone
}

// MarshalJSON implements json.Marshaler.
func (h *SpaceshipSystem) MarshalJSON() ([]byte, error) {
	type calc struct {
		RollRange string         `json:"roll_range"`
		DR        map[string]int `json:"dr,omitempty"`
	}
	data := struct {
		SpaceshipSystemData
		Calc calc `json:"calc"`
	}{
		SpaceshipSystemData: h.SpaceshipSystemData,
		Calc: calc{
			RollRange: h.RollRange,
		},
	}

	return json.Marshal(&data)
}

// UnmarshalJSON implements json.Unmarshaler.
func (h *SpaceshipSystem) UnmarshalJSON(data []byte) error {
	h.SpaceshipSystemData = SpaceshipSystemData{}
	if err := json.Unmarshal(data, &h.SpaceshipSystemData); err != nil {
		return err
	}

	return nil
}

// ID returns the ID.
func (h *SpaceshipSystem) ID() string {
	return h.SysID
}

// SetID sets the ID, sanitizing it in the process (i.e. it may be changed from what you set -- read it back if you want
// to be sure of what it gets set to.
func (h *SpaceshipSystem) SetID(value string) {
	h.SysID = SanitizeID(value, false, ReservedIDs...)
}

// OwningTable returns the owning table.
func (h *SpaceshipSystem) OwningTable() *Body {
	return h.owningTable
}
