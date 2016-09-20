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
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	//"github.com/xiam/backstream"
	"fmt"
	"github.com/xiam/g"
	"io"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"github.com/xiam/shooter-server/ship"
	"time"
)

const (
	playerMaxSpeed  = 5.0
	playerTurnSpeed = math.Pi / 24
	playerMaxLife   = 32
	playerNearValue = 500
)

var playerShowLog = false

var (
	ErrNoWebsocket = errors.New(`Don't have any websocket connection.`)
)

func init() {
	if os.Getenv("DEBUG") != "" {
		playerShowLog = true
	}
	rand.Seed(time.Now().Unix())
}

type player struct {
	ws     *websocket.Conn
	output chan []byte
	addr   net.Addr

	points     uint64
	hitValue   int
	killValue  int
	bulletType int
	sector     *sector
	compressor io.Writer
	*ship.Ship
	*control
	shootTicks uint

	speedFactor float64
	destroying  bool
}

func newPlayer(ws *websocket.Conn) *player {
	self := &player{}
	self.ws = ws

	// Bot players do not use a websocket.
	if self.ws != nil {
		self.output = make(chan []byte, 256)
		//self.compressor = backstream.NewWriter(self, 0)
		self.addr = ws.RemoteAddr()
	} else {
		self.output = nil
	}

	// Creating controller.
	self.control = NewControl()

	self.speedFactor = 1.0

	// Creating a ship.
	self.Ship = ship.NewShip()

	self.Ship.Entity.Kind = (1 + rand.Intn(8))

	switch self.Ship.Entity.Kind {
	case 1, 2, 3, 4:
		self.Ship.Entity.Width = 80
		self.Ship.Entity.Height = 120
	case 5, 6, 7, 8:
		self.Ship.Entity.Width = 80
		self.Ship.Entity.Height = 110
	}

	self.Ship.Entity.Model = fmt.Sprintf("ship-%d", self.Ship.Entity.Kind)

	self.hitValue = 2

	self.Life = playerMaxLife

	self.bulletType = BULLET_1X
	self.shootTicks = 0

	// Default values.
	self.SetSpeed(0)
	self.SetDirection(0, -1)
	self.SetPosition(0, 0)

	return self
}

func (self *player) log(s string) {
	if self.addr != nil {
		log.Printf("%s: %s\n", self.addr, s)
	} else {
		log.Printf("bot: %s\n", s)
	}
}

// Determines is a player is near other player.
func (self *player) isNear(other *player) bool {
	xdiff := math.Abs(self.Position[0] - other.Position[0])
	ydiff := math.Abs(self.Position[1] - other.Position[1])
	mdiff := math.Max(xdiff, ydiff)
	if mdiff > playerNearValue {
		return false
	}
	return true
}

func (self *player) Write(p []byte) (n int, err error) {

	if self.ws == nil {
		return 0, ErrNoWebsocket
	}

	err = self.ws.WriteMessage(websocket.TextMessage, p)

	if playerShowLog == true {
		log.Printf("%s <- %s: %v\n", self.ws.RemoteAddr(), p, err)
	}

	return len(p), err
}

// Serializes player.
func (self *player) Serialize() (buf []byte) {
	data := self.Ship.DataMap()
	(*data)["L"] = self.Life
	(*data)["P"] = self.points
	(*data)["N"] = self.control.Name
	self.Ship.Diff.SetData(data)
	return self.Ship.Diff.Serialize()
}

func (self *player) addLife(delta int) {
	self.Life = int(math.Min(float64(self.Life+delta), float64(playerMaxLife)))
}

// Adds points to a player
func (self *player) addPoints(delta int) {
	self.points = self.points + uint64(delta)
}

func (self *player) newBullet() *fire {
	f := newFire(self.bulletType)
	f.player = self
	return f
}

func (self *player) shoot() {
	if self.sector == nil {
		return
	}

	self.shootTicks++

	if self.shootTicks%3 != 1 {
		return
	}

	b := self.newBullet()
	b.SetPosition(
		self.Position[0]+self.Direction[0]*(self.Width/2.0)*1.1,
		self.Position[1]+self.Direction[1]*(self.Width/2.0)*1.1,
	)
	b.SetDirection(self.Direction[0], self.Direction[1])
	b.SetSpeed(float64(playerMaxSpeed) * 3.0 * self.speedFactor)
	self.sector.addFire(b)
}

func (self *player) hit(other *player, val int) {
	if self.Life > 1 {
		self.Life = self.Life - val
	} else {
		if other != nil {
			other.addPoints(self.killValue)
		}
		self.destroy()
	}
}

func (self *player) ident() {
	if self.Id != "" {
		self.write(identFn(self.Id))
	}
}

func (self *player) sameAs(other *player) bool {
	if self.Id != "" {
		return self.Id == other.Id
	}
	return false
}

func (self *player) write(data []byte) {
	if self.ws == nil {
		return
	}
	if self.output == nil {
		return
	}
	self.output <- data
}

func (self *player) notice(other *player) {
	data := other.DataMap()
	buf, err := json.Marshal(data)
	if err == nil {
		if other.ws == nil {
			self.write(createFn("ship-ai", other.Id, buf))
		} else {
			self.write(createFn("ship", other.Id, buf))
		}
	}
}

func (self *player) collidesWithPlayer(other *player) ([]*g.Point, error) {
	a := self.Poly()
	b := other.Poly()
	points, err := g.PolyIntersectsPoly(a, b)
	if err == nil {
		return points, nil
	}
	return points, err
}

func (self *player) update() {
	b := self.Serialize()
	if b != nil {
		chunk := updateFn(self.Id, b)
		if self.sector != nil {
			for p, _ := range self.sector.players {
				if p.ws != nil && p.output != nil {
					if self.isNear(p) == true {
						p.write(chunk)
					}
				}
			}
		}
	}
}

func (self *player) destroy() {
	if self.destroying == false {

		self.destroying = true
		//self.log("Attempt to destroy.")

		// Announcing this ship is destroyed.
		if self.sector != nil {
			// Adding to top scores.
			highScores.Add(self.control.Name, self.points)

			// Last update
			self.update()
			//self.log("Updated")

			//self.log("Removing from sector.")
			self.sector.broadcast(destroyFn(self.Id))
			self.sector.gbgPlayer <- self
			self.Life = 0

			//self.log("Removed")
		}

		if self.ws != nil {
			topScores := highScores.GetTop()
			jsonScores, _ := json.Marshal(topScores)
			// Writing directly on the compressor
			//self.log("Writing scores.")
			self.write(scoresFn(jsonScores))
			//self.log("Wrote.")
		}

		//self.log("Destroyed")
	}
}

func (self *player) reader() {

	for {
		_, message, err := self.ws.ReadMessage()
		if err != nil {
			break
		}
		json.Unmarshal(message, self.control)
		if playerShowLog == true {
			log.Printf("%s -> %s\n", self.ws.RemoteAddr(), message)
		}
	}
	//self.log("Exiting reader.")
	//self.close()

	self.destroy()
}

func (self *player) writer() {
	var start, diff, sleep int64
	var buf []byte

	writing := true

	for writing {
		buf = make([]byte, 0, 1024*10)

		start = time.Now().UnixNano()

		loop := true
		for loop {
			select {
			case message := <-self.output:
				buf = append(buf, message...)
				buf = append(buf, '\n')
			default:
				loop = false
			}
		}

		if len(buf) > 0 {
			_, err := self.Write(buf)
			if err != nil {
				writing = false
			}
		}

		diff = time.Now().UnixNano() - start
		sleep = fpsn - diff

		//fmt.Printf("sleep: %d, diff: %d\n", sleep, diff)

		if sleep < 0 {
			continue
		}

		time.Sleep(time.Duration(sleep) * time.Nanosecond)
	}

	/*
		for message := range self.output {
			//_, err := self.compressor.Write(message)
			_, err := self.Write(message)
			if err != nil {
				//self.log("Got write error.")
				break
			}
		}
	*/

	//self.log("Exiting writer.")

	self.close()
}

func (self *player) close() {
	if self.ws != nil {
		//self.log("Closing websocket.")
		self.ws.Close()
		self.ws = nil
		//self.log("Websocket closed.")
	}

	if self.output != nil {
		//self.log("Closing channel.")
		close(self.output)
		self.output = nil
		//self.log("Channel closed.")
	}
}

func (self *player) isFree() bool {
	var poly *g.Poly

	if self.sector != nil {
		for other := range self.sector.players {
			if self.sameAs(other) == false {
				poly = self.Poly()
				_, err := g.PolyIntersectsPoly(poly, other.Poly())
				if err == nil {
					return false
				}
			}
		}
	}
	return true
}

func (self *player) correct() {
	for self.isFree() == false {
		area := int64(math.Max(float64(self.sector.bounds[0])/5, float64(self.sector.bounds[1]/5)))
		self.Position[0], self.Position[1] = float64(rand.Int63n(area)-area/2), float64(rand.Int63n(area)-area/2)
	}
}

func (self *player) Tick() {
	var t float64

	t = playerTurnSpeed

	if self.control.Y > 0 {
		self.Speed = -playerMaxSpeed
	} else if self.control.Y < 0 {
		self.Speed = playerMaxSpeed
	} else {
		self.Speed = 0.0
	}

	self.Speed = self.Speed * self.speedFactor

	if self.control.X < 0 {
		t = -t
	} else if self.control.X == 0 {
		t = 0.0
	}

	x := self.Direction[0]
	y := self.Direction[1]

	self.Direction[0] = x*math.Cos(t) - y*math.Sin(t)
	self.Direction[1] = x*math.Sin(t) + y*math.Cos(t)

	// Attempt to move.
	self.Position[0] = self.Position[0] + self.Direction[0]*self.Speed
	self.Position[1] = self.Position[1] + self.Direction[1]*self.Speed

	// Boundary checking
	if int64(self.Position[0]) > self.sector.bounds[0] {
		self.Position[0] = float64(self.sector.bounds[0])
	}

	if int64(self.Position[0]) < -self.sector.bounds[0] {
		self.Position[0] = float64(-self.sector.bounds[0])
	}

	if int64(self.Position[1]) > self.sector.bounds[1] {
		self.Position[1] = float64(self.sector.bounds[1])
	}

	if int64(self.Position[1]) < -self.sector.bounds[1] {
		self.Position[1] = float64(-self.sector.bounds[1])
	}

	// Collision check
	poly := self.Poly()
	for other, _ := range self.sector.players {
		if self.sameAs(other) == false {
			_, err := g.PolyIntersectsPoly(poly, other.Poly())
			if err == nil {
				self.hit(nil, 1)
			}
		}
	}

	if self.control.S > 0 {
		self.shoot()
	}

}
