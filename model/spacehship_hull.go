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
	"context"
	"io/fs"
	"sort"
	"strings"

	"github.com/richardwilkes/gcs/v5/model/jio"
	"github.com/richardwilkes/rpgtools/dice"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/txt"
)

const (
	hullTypeListTypeKey = "hull_type"
)

/*
go:embed embedded_data
var embeddedFS embed.FS
*/
// SpaceshipHull holds a set of hit locations.
type SpaceshipHull struct {
	Name         string             `json:"name,omitempty"`
	Roll         *dice.Dice         `json:"roll"`
	Systems      []*SpaceshipSystem `json:"systems,omitempty"`
	KeyPrefix    string             `json:"-"`
	owningSystem *SpaceshipSystem
	systemLookup map[string]*SpaceshipSystem
}

type SpaceshipHullData struct {
	Type    string `json:"type"`
	Version int    `json:"version"`
	*SpaceshipHull
}

func NewSpaceshipHull(name string) *SpaceshipHull {
	sh := SpaceshipHull{
		Name:         name,
		Roll:         &dice.Dice{},
		Systems:      []*SpaceshipSystem{},
		KeyPrefix:    "",
		owningSystem: &SpaceshipSystem{},
		systemLookup: map[string]*SpaceshipSystem{},
	}
	return &sh
}

// NewSpaceshipHullFromFile loads a SpaceshipHull from a file.
func NewSpaceshipHullFromFile(fileSystem fs.FS, filePath string) (*SpaceshipHull, error) {
	var data struct {
		SpaceshipHullData
		OldSpaceshipSystems *SpaceshipHull `json:"hit_locations"`
	}
	if err := jio.LoadFromFS(context.Background(), fileSystem, filePath, &data); err != nil {
		return nil, errs.NewWithCause(invalidFileDataMsg(), err)
	}

	if err := CheckVersion(data.Version); err != nil {
		return nil, err
	}

	return data.SpaceshipHull, nil
}

// Clone a copy of this.
func (b *SpaceshipHull) Clone() *SpaceshipHull {
	clone := &SpaceshipHull{
		Name:    b.Name,
		Roll:    dice.New(b.Roll.String()),
		Systems: make([]*SpaceshipSystem, len(b.Systems)),
	}
	for i, one := range b.Systems {
		clone.Systems[i] = one.Clone()
	}
	return clone
}

// Save writes the SpaceshipHull to the file as JSON.
func (b *SpaceshipHull) Save(filePath string) error {
	return jio.SaveToFile(context.Background(), filePath, &SpaceshipHullData{
		Type:          bodyTypeListTypeKey,
		Version:       CurrentDataVersion,
		SpaceshipHull: b,
	})
}

// OwningLocation returns the owning hit location, or nil if this is the top-level body.
func (h *SpaceshipHull) OwningLocation() *SpaceshipSystem {
	return h.owningSystem
}

// SetOwningLocation sets the owning SpaceshipSystem.
func (h *SpaceshipHull) SetOwningLocation(sys *SpaceshipSystem) {
	h.owningSystem = sys
	if sys != nil {
		h.Name = ""
	}
}

// AddLocation adds a SpaceshipSystem to the end of list.
func (h *SpaceshipHull) AddLocation(sys *SpaceshipSystem) {
	h.Systems = append(h.Systems, sys)
}

// RemoveLocation removes a SpaceshipSystem.
func (h *SpaceshipHull) RemoveLocation(sys *SpaceshipSystem) {
	for i, one := range h.Systems {
		if one == sys {
			copy(h.Systems[i:], h.Systems[i+1:])
			h.Systems[len(h.Systems)-1] = nil
			h.Systems = h.Systems[:len(h.Systems)-1]
		}
	}
}

// UniqueSpaceshipSystems returns the list of unique hit locations.
func (h *SpaceshipHull) UniqueSpaceshipSystems(entity *Entity) []*SpaceshipSystem {
	locations := make([]*SpaceshipSystem, 0, len(h.systemLookup))
	for _, v := range h.systemLookup {
		locations = append(locations, v)
	}
	sort.Slice(locations, func(i, j int) bool {
		if txt.NaturalLess(locations[i].ChoiceName, locations[j].ChoiceName, false) {
			return true
		}
		if strings.EqualFold(locations[i].ChoiceName, locations[j].ChoiceName) {
			return txt.NaturalLess(locations[i].ID(), locations[j].ID(), false)
		}
		return false
	})
	return locations
}

// LookupLocationByID returns the SpaceshipSystem that matches the given ID.
func (h *SpaceshipHull) LookupLocationByID(entity *Entity, idStr string) *SpaceshipSystem {

	return h.systemLookup[idStr]
}
