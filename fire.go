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
	"fmt"
	"github.com/xiam/g"
	"math"
	"github.com/xiam/shooter-server/bullet"
)

const (
	bulletLifeSpan = 30
)

const (
	BULLET_1X = iota
	BULLET_2X
	BULLET_3X
)

type fire struct {
	*bullet.Bullet
	sector   *sector
	player   *player
	hitValue int
}

func newFire(kind int) *fire {
	self := &fire{}
	self.Bullet = bullet.NewBullet()
	self.Bullet.Kind = kind
	self.Bullet.Model = fmt.Sprintf("beam-%dx", self.Bullet.Kind+1)
	self.Bullet.Life = bulletLifeSpan
	self.Bullet.Diff.Ignore["p"] = true
	self.hitValue = 1
	switch kind {
	case BULLET_2X:
		self.hitValue = 2
	}
	return self
}

func (self *fire) isNear(other *player) bool {
	xdiff := math.Abs(self.Position[0] - other.Position[0])
	ydiff := math.Abs(self.Position[1] - other.Position[1])
	mdiff := math.Max(xdiff, ydiff)
	if mdiff > playerNearValue {
		return false
	}
	return true
}

func (self *fire) update() {
	b := self.Serialize()
	if b != nil {
		chunk := updateFn(self.Id, b)
		//self.sector.broadcast(chunk)
		for p, _ := range self.sector.players {
			if p.ws != nil {
				if self.isNear(p) == true {
					p.write(chunk)
				}
			}
		}
	}
}

func (self *fire) destroy() {
	self.sector.broadcast(destroyFn(self.Id))
	self.sector.gbgFire <- self
}

func (self *fire) Tick() {
	self.Life = self.Life - 1
	if self.Life > 0 {
		self.Position[0] = self.Position[0] + self.Direction[0]*self.Speed
		self.Position[1] = self.Position[1] + self.Direction[1]*self.Speed
	} else {
		self.destroy()
		return
	}
	if self.sector != nil {
		poly := self.Poly()
		for player := range self.sector.players {
			if player != nil {
				_, err := g.PolyIntersectsPoly(poly, player.Poly())
				if err == nil {
					player.hit(self.player, self.hitValue)
					if self.player != nil {
						self.player.addPoints(player.hitValue)
					}
					self.Life = 0
				}
			}
		}
	}
}
