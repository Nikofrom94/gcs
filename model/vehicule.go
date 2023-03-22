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
	"bytes"
	"context"
	"io/fs"
	"strconv"

	"github.com/google/uuid"
	"github.com/richardwilkes/gcs/v5/model/jio"
	"github.com/richardwilkes/json"
	"github.com/richardwilkes/toolbox/errs"
)

/*
var (
	_ eval.VariableResolver = &Vehicule{}
	_ ListProvider          = &Vehicule{}
)
*/
// VehiculeProvider provides a way to retrieve a (possibly nil) Vehicule.
type VehiculerProvider interface {
	Vehicule() *Vehicule
}

// VehiculeHeader holds the Vehicule data that is written to disk.
type VehiculeHeader struct {
	Type          string         `json:"type"`
	Version       int            `json:"version"`
	ID            uuid.UUID      `json:"id"`
	Profile       *Profile       `json:"profile,omitempty"`
	Notes         []*Note        `json:"notes,omitempty"`
	CreatedOn     jio.Time       `json:"created_date"`
	ModifiedOn    jio.Time       `json:"modified_date"`
	SheetSettings *SheetSettings `json:"settings,omitempty"`
}

func NewVehiculeHeader() *VehiculeHeader {
	v := &VehiculeHeader{
		Type:      "",
		ID:        NewUUID(),
		CreatedOn: jio.Now(),
		Profile:   &Profile{},
	}
	v.SheetSettings = GlobalSettings().SheetSettings().Clone(nil)
	v.ModifiedOn = v.CreatedOn
	return v
}

// VehiculeAttributes holds the Vehicule main attributes
type VehiculeAttributes struct {
	TL    int                     `json:"tl"`
	ST    VehiculeAttributeInt    `json:"st"`
	HP    VehiculeAttributeInt    `json:"hp"`
	HT    VehiculeAttributeInt    `json:"ht"`
	Hnd   VehiculeAttributeInt    `json:"hnd"`
	SR    VehiculeAttributeInt    `json:"sr"`
	Move  VehiculeAttributeInt    `json:"move"`
	LWt   VehiculeAttributeInt    `json:"lwt"`
	Load  VehiculeAttributeInt    `json:"load"`
	SM    VehiculeAttributeInt    `json:"sm"`
	Occ   VehiculeAttributeString `json:"occ"`
	DR    VehiculeAttributeInt    `json:"dr"`
	Range VehiculeAttributeInt    `json:"range"`
	Cost  VehiculeAttributeInt    `json:"cost"`
}

type VehiculeAttributeCore struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayname"`
	Code        string `json:"code"`
}

type VehiculeAttributeString struct {
	VehiculeAttributeCore
	Value string `json:"value"`
}

type VehiculeAttributeInt struct {
	VehiculeAttributeCore
	Value int `json:"value"`
}

func (v *VehiculeAttributeInt) ToString() string {
	return strconv.Itoa(v.Value) + v.Code
}

func (v *VehiculeAttributeInt) SetValue(strVal string) {
	if iVal, err := strconv.Atoi(strVal); err == nil {
		v.Value = iVal
	}
}

func NewVehiculeAttributeInt(name string, displayname string, defaultvalue int, code string) *VehiculeAttributeInt {
	return &VehiculeAttributeInt{
		VehiculeAttributeCore: VehiculeAttributeCore{
			Name:        name,
			DisplayName: displayname,
			Code:        code,
		},
		Value: defaultvalue,
	}
}
func NewVehiculeAttributeString(name string, displayname string, defaultvalue string, code string) *VehiculeAttributeString {
	return &VehiculeAttributeString{
		VehiculeAttributeCore: VehiculeAttributeCore{
			Name:        name,
			DisplayName: displayname,
			Code:        code,
		},
		Value: defaultvalue,
	}
}

/*
// MarshalJSON implements json.Marshaler.

	func (a *VehiculeAttributeInt) MarshalJSON() ([]byte, error) {
		data := struct {
			VehiculeAttributeInt
		}{
			*a,
		}

		return json.Marshal(&data)
	}

// UnmarshalJSON implements json.Unmarshaler.

	func (a *VehiculeAttributeInt) UnmarshalJSON(data []byte) error {
		a = &VehiculeAttributeInt{}
		if err := json.Unmarshal(data, &a); err != nil {
			return err
		}
		return nil
	}
*/
func NewVehiculeAttributes() *VehiculeAttributes {
	attr := &VehiculeAttributes{
		TL:    0,
		ST:    *NewVehiculeAttributeInt("ST", "Strength", 0, ""),
		HP:    *NewVehiculeAttributeInt("HP", "Hit Points", 0, ""),
		HT:    *NewVehiculeAttributeInt("HT", "Health", 0, ""),
		Hnd:   *NewVehiculeAttributeInt("Hnd", "Handling", 0, ""),
		SR:    *NewVehiculeAttributeInt("Hnd", "Stability Rating", 0, ""),
		Move:  *NewVehiculeAttributeInt("Move", "Move", 0, ""),
		LWt:   *NewVehiculeAttributeInt("Lwt", "Loading Weight", 0, ""),
		Load:  *NewVehiculeAttributeInt("Load", "Load", 0, ""),
		SM:    *NewVehiculeAttributeInt("SM", "Size Modifier", 1, ""),
		Occ:   *NewVehiculeAttributeString("Occ", "Occupancy", "1", ""),
		DR:    *NewVehiculeAttributeInt("Hnd", "Handling", 0, ""),
		Range: *NewVehiculeAttributeInt("Hnd", "Handling", 0, ""),
		Cost:  *NewVehiculeAttributeInt("Hnd", "Handling", 0, ""),
	}
	return attr
}

type VehiculeInterface interface {
	NoteList() []*Note
	SetNoteList(list []*Note)
	ReadVehiculeFromFile(fileSystem fs.FS, filePath string) (*VehiculeInterface, error)
	Save(filePath string) error
	UnmarshalJSON(data []byte) error
	MarshalJSON() ([]byte, error)
}

// Vehicule holds the base information for vehicule
type Vehicule struct {
	VehiculeHeader
	VehiculeAttributes
}

// NewVehiculeFromFile loads an Vehicule from a file.
func NewVehiculeFromFile(fileSystem fs.FS, filePath string) (*Vehicule, error) {
	var vehicule Vehicule
	if err := jio.LoadFromFS(context.Background(), fileSystem, filePath, &vehicule); err != nil {
		return nil, errs.NewWithCause(invalidFileDataMsg(), err)
	}
	if err := CheckVersion(vehicule.Version); err != nil {
		return nil, err
	}
	return &vehicule, nil
}

// NewVehicule creates a new Vehicule.
func NewVehicule() *Vehicule {
	vehicule := &Vehicule{
		VehiculeHeader:     *NewVehiculeHeader(),
		VehiculeAttributes: *NewVehiculeAttributes(),
	}
	return vehicule
}

// Vehicule implements VehiculeProvider.
func (v *Vehicule) Vehicule() *Vehicule {
	return v
}

// Save the Vehicule to a file as JSON.
func (v *Vehicule) Save(filePath string) error {
	return jio.SaveToFile(context.Background(), filePath, v)
}

// MarshalJSON implements json.Marshaler.
func (v *Vehicule) MarshalJSON() ([]byte, error) {
	v.VehiculeHeader.Version = CurrentDataVersion
	data := struct {
		VehiculeHeader
		VehiculeAttributes
	}{
		v.VehiculeHeader,
		v.VehiculeAttributes,
	}

	return json.Marshal(&data)
}

// UnmarshalJSON implements json.Unmarshaler.
func (v *Vehicule) UnmarshalJSON(data []byte) error {
	v = NewVehicule()
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if err := json.Unmarshal(data, &v.VehiculeAttributes); err != nil {
		return err
	}

	return nil
}

// DiscardCaches discards the internal caches.
func (v *Vehicule) DiscardCaches() {

}

// Recalculate the statistics.
func (v *Vehicule) Recalculate() {
	v.DiscardCaches()
}

// Ancestry returns the current Ancestry.
func (v *Vehicule) Ancestry() *Ancestry {
	var anc *Ancestry
	/*
		Traverse(func(t *Trait) bool {
			if t.Container() && t.ContainerType == RaceContainerType {
				if anc = LookupAncestry(t.Ancestry, GlobalSettings().Libraries()); anc != nil {
					return true
				}
			}
			return false
		}, true, false,...)
		if anc == nil {
			if anc = LookupAncestry(DefaultAncestry, GlobalSettings().Libraries()); anc == nil {
				jot.Fatal(1, "unable to load default ancestry (Human)")
			}
		}*/
	return anc
}

// NoteList implements ListProvider
func (v *Vehicule) NoteList() []*Note {
	return v.Notes
}

// SetNoteList implements ListProvider
func (v *Vehicule) SetNoteList(list []*Note) {
	for _, one := range list {
		one.SetOwningEntity(nil)
	}
	v.Notes = list
}

// CRC64 computes a CRC-64 value for the canonical disk format of the data. The ModifiedOn field is ignored for this
// calculation.
func (v *Vehicule) CRC64() uint64 {
	var buffer bytes.Buffer
	saved := v.ModifiedOn
	v.ModifiedOn = jio.Time{}
	defer func() { v.ModifiedOn = saved }()
	if err := jio.Save(context.Background(), &buffer, v); err != nil {
		return 0
	}
	return CRCBytes(0, buffer.Bytes())
}

// TraitList implements ListProvider : N/A for Vehicule
func (v *Vehicule) TraitList() []*Trait {
	return nil
}

// SetTraitList implements ListProvider
func (v *Vehicule) SetTraitList(list []*Trait) {

}

// CarriedEquipmentList implements ListProvider: N/A for Vehicule
func (v *Vehicule) CarriedEquipmentList() []*Equipment {
	return nil
}

// SetCarriedEquipmentList implements ListProvider: N/A for Vehicule
func (v *Vehicule) SetCarriedEquipmentList(list []*Equipment) {

}

// OtherEquipmentList implements ListProvider: N/A for Vehicule
func (v *Vehicule) OtherEquipmentList() []*Equipment {
	return nil
}

// SetOtherEquipmentList implements ListProvider
func (v *Vehicule) SetOtherEquipmentList(list []*Equipment) {

}

// SkillList implements ListProvider: N/A for Vehicule
func (v *Vehicule) SkillList() []*Skill {
	return nil
}

// SetSkillList implements ListProvider
func (v *Vehicule) SetSkillList(list []*Skill) {

}

// SpellList implements ListProvider
func (v *Vehicule) SpellList() []*Spell {
	return nil
}

// SetSpellList implements ListProvider
func (v *Vehicule) SetSpellList(list []*Spell) {

}

func (v *Vehicule) Entity() *Entity {
	return nil
}
