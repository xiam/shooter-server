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

package entity

import (
	"github.com/xiam/g"
	"math"
	"github.com/xiam/shooter-server/diff"
)

var compressTypes = true

type Entity struct {
	Width     float64
	Height    float64
	Position  [2]float64
	Direction [2]float64
	Speed     float64
	Life      int
	Id        string
	Kind      int
	Model     string
	dataMap   map[string]interface{}
	poly      *g.Poly
	*diff.Diff
}

func NewEntity() *Entity {
	self := &Entity{}
	self.Diff = diff.NewDiff()
	self.dataMap = map[string]interface{}{}
	self.poly = g.NewPoly(
		g.NewPoint(0, 0),
		g.NewPoint(0, 0),
		g.NewPoint(0, 0),
		g.NewPoint(0, 0),
	)
	return self
}

func truncate(i float64) float64 {
	return float64(int(i*1000)) / 1000.0
}

func (self *Entity) UpdateDataMap() {
	self.dataMap["w"] = uint64(self.Width)
	self.dataMap["h"] = uint64(self.Height)
	if compressTypes == true {
		self.dataMap["p"] = []float32{float32(truncate(self.Position[0])), float32(truncate(self.Position[1]))}
		self.dataMap["d"] = []float32{float32(truncate(self.Direction[0])), float32(truncate(self.Direction[1]))}
		self.dataMap["s"] = float32(self.Speed)
	} else {
		self.dataMap["p"] = self.Position
		self.dataMap["d"] = self.Direction
		self.dataMap["s"] = self.Speed
	}
	self.dataMap["k"] = self.Kind
	self.dataMap["m"] = self.Model
}

func (self *Entity) DataMap() *map[string]interface{} {
	self.UpdateDataMap()
	return &self.dataMap
}

func (self *Entity) Serialize() (buf []byte) {
	self.UpdateDataMap()
	self.Diff.SetData(&self.dataMap)
	return self.Diff.Serialize()
}

func (self *Entity) SetId(id string) {
	self.Id = id
}

func (self *Entity) SetSpeed(s float64) {
	self.Speed = s
}

func (self *Entity) Poly() *g.Poly {
	mw := self.Width / 2.0
	mh := self.Height / 2.0

	n := [2]float64{self.Direction[1], -self.Direction[0]}

	ox, oy := self.Position[0], self.Position[1]

	ax, ay := ox+self.Direction[0]*mw+n[0]*mh, oy+self.Direction[1]*mw+n[1]*mh
	bx, by := ax-n[0]*self.Height, ay-n[1]*self.Height
	cx, cy := bx-self.Direction[0]*self.Width, by-self.Direction[1]*self.Width
	dx, dy := cx+n[0]*self.Height, cy+n[1]*self.Height

	self.poly.Points[0].Set(ax, ay)
	self.poly.Points[1].Set(bx, by)
	self.poly.Points[2].Set(cx, cy)
	self.poly.Points[3].Set(dx, dy)

	return self.poly
}

func (self *Entity) SetPosition(x float64, y float64) {
	self.Position[0], self.Position[1] = x, y
}

func (self *Entity) SetDirection(x float64, y float64) {
	var d float64
	d = math.Sqrt(x*x + y*y)
	if d > 0 {
		//self.Direction = [2]float64{x / d, y / d}
		self.Direction[0], self.Direction[1] = x/d, y/d
	} else {
		//self.Direction = [2]float64{0, 0}
		self.Direction[0], self.Direction[1] = 0.0, 0.0
	}
}

func (self *Entity) Tick() {
	self.Position[0], self.Position[1] = self.Position[0]+self.Direction[0]*self.Speed, self.Position[1]+self.Direction[1]*self.Speed
}
