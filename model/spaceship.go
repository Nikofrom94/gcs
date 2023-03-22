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

type Spaceship struct {
	Vehicule
	Hulls []*SpaceshipHull `json:"hulls,omitempty"`
}

func NewSpaceship() *Spaceship {
	s := Spaceship{}
	s.Vehicule = *NewVehicule()
	s.CreateSpaceshipHulls()
	return &s
}

func (s *Spaceship) CreateSpaceshipHulls() {
	s.Hulls = make([]*SpaceshipHull, 3)
	s.Hulls[0] = NewSpaceshipHull("Front Hull")
	s.Hulls[1] = NewSpaceshipHull("Central Hull")
	s.Hulls[2] = NewSpaceshipHull("Rear Hull")
}
