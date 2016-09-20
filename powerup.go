// Copyright 2014 José Carlos Nieto
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/xiam/g"
	"math/rand"
	"github.com/xiam/shooter-server/item"
)

const BONUS_POINTS = 50
const LIFE_POINTS = 20

const (
	POWERUP_BEAM_1X = iota
	POWERUP_BEAM_2X
	POWERUP_BEAM_3X

	POWERUP_BONUS_1X
	POWERUP_BONUS_2X
	POWERUP_BONUS_3X

	POWERUP_RECOVER_1X
	POWERUP_RECOVER_2X
	POWERUP_RECOVER_3X
	POWERUP_RECOVER_FULL

	POWERUP_SPEED_1X
	POWERUP_SPEED_2X
	POWERUP_SPEED_3X

	POWERUP_LIMIT
)

type powerup struct {
	*item.Item
	sector *sector
}

func newPowerup() (self *powerup) {
	self = &powerup{}
	self.Item = item.NewItem()

	self.Item.Kind = rand.Intn(POWERUP_LIMIT)

	switch self.Item.Kind {
	case POWERUP_BEAM_1X:
		self.Item.Model = "item-beam-1x"
		self.Item.Width = 34
		self.Item.Height = 33
	case POWERUP_BEAM_2X:
		self.Item.Model = "item-beam-2x"
		self.Item.Width = 34
		self.Item.Height = 33
	case POWERUP_BEAM_3X:
		self.Item.Model = "item-beam-3x"
		self.Item.Width = 34
		self.Item.Height = 33
	case POWERUP_BONUS_1X:
		self.Item.Model = "item-bonus-1x"
		self.Item.Width = 31
		self.Item.Height = 30
	case POWERUP_BONUS_2X:
		self.Item.Model = "item-bonus-2x"
		self.Item.Width = 31
		self.Item.Height = 30
	case POWERUP_BONUS_3X:
		self.Item.Model = "item-bonus-3x"
		self.Item.Width = 31
		self.Item.Height = 30
	case POWERUP_RECOVER_1X:
		self.Item.Model = "item-recover-1x"
		self.Item.Width = 22
		self.Item.Height = 21
	case POWERUP_RECOVER_2X:
		self.Item.Model = "item-recover-2x"
		self.Item.Width = 22
		self.Item.Height = 21
	case POWERUP_RECOVER_3X:
		self.Item.Model = "item-recover-3x"
		self.Item.Width = 22
		self.Item.Height = 21
	case POWERUP_RECOVER_FULL:
		self.Item.Model = "item-recover-full"
		self.Item.Width = 22
		self.Item.Height = 21
	case POWERUP_SPEED_1X:
		self.Item.Model = "item-speed-1x"
		self.Item.Width = 19
		self.Item.Height = 30
	case POWERUP_SPEED_2X:
		self.Item.Model = "item-speed-2x"
		self.Item.Width = 19
		self.Item.Height = 30
	case POWERUP_SPEED_3X:
		self.Item.Model = "item-speed-3x"
		self.Item.Width = 19
		self.Item.Height = 30
	}
	return self
}

func (self *powerup) destroy() {
	self.sector.broadcast(destroyFn(self.Id))
	self.sector.gbgPowerup <- self
}

func (self *powerup) Tick() {
	if self.Life <= 0 {
		self.destroy()
		return
	}
	if self.sector != nil {
		poly := self.Poly()
		for player := range self.sector.players {
			if player != nil {
				_, err := g.PolyIntersectsPoly(poly, player.Poly())
				if err == nil {
					// Intercepts player
					switch self.Item.Kind {

					case POWERUP_BEAM_1X:
						player.bulletType = BULLET_1X
					case POWERUP_BEAM_2X:
						player.bulletType = BULLET_2X
					case POWERUP_BEAM_3X:
						player.bulletType = BULLET_3X

					case POWERUP_BONUS_1X:
						player.addPoints(BONUS_POINTS * 1)
					case POWERUP_BONUS_2X:
						player.addPoints(BONUS_POINTS * 2)
					case POWERUP_BONUS_3X:
						player.addPoints(BONUS_POINTS * 3)

					case POWERUP_RECOVER_1X:
						player.addLife(LIFE_POINTS * 1)
					case POWERUP_RECOVER_2X:
						player.addLife(LIFE_POINTS * 2)
					case POWERUP_RECOVER_3X:
						player.addLife(LIFE_POINTS * 3)
					case POWERUP_RECOVER_FULL:
						player.addLife(LIFE_POINTS * 10)

					case POWERUP_SPEED_1X:
						player.speedFactor = 1.0
					case POWERUP_SPEED_2X:
						player.speedFactor = 2.0
					case POWERUP_SPEED_3X:
						player.speedFactor = 3.0

					}
					self.Life = self.Life - 1
					return
				}
			}
		}
	}
}
