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
	"log"
	"math"
	"math/rand"
	"time"
)

const (
	agentSleep    = 1000 / 48 // Milliseconds.
	agentMaxTicks = 1000      // Ticks before self-destruction.
	largeDistance = 1000000   // A large number for representing objects to far from the agent.
)

func init() {
	rand.Seed(time.Now().Unix())
}

type agent struct {
	*player         // This agent's player.
	lock    *player // Player we're going after.
	ticks   uint64
}

// Creating agent.
func newAgent() *agent {
	self := &agent{}

	// Agent is a player that does not have a socket.
	self.player = newPlayer(nil)

	// Random name. Name collisions don't matter.
	self.control.Name = fmt.Sprintf(`AI-%d`, 1000+rand.Int31n(8999))

	// Random type.
	self.player.Entity.Kind = (1 + rand.Intn(8))
	self.player.Life = playerMaxLife / 2

	// There are different kinds of players.
	switch self.player.Entity.Kind {
	case 1, 2, 3, 4:
		self.player.Entity.Width = 90
		self.player.Entity.Height = 120
	case 5, 6, 7, 8:
		self.player.Entity.Width = 91
		self.player.Entity.Height = 91
	}

	self.player.Entity.Model = fmt.Sprintf("ship-ai-%d", self.player.Entity.Kind)

	go self.autoPilot()

	return self
}

// Quick distance calculation, without the square root.
func (self *agent) distance(p *player) int64 {
	x := int64(self.Position[0] - p.Position[0])
	y := int64(self.Position[1] - p.Position[1])
	z := x*x + y*y
	if z > largeDistance {
		return largeDistance
	}
	return z
}

// Basic navigate and attack routine.
func (self *agent) autoPilot() {
	z := 0

	defer func() {
		log.Printf("Terminating agent.")
	}()

	for {

		if self.sector == nil {
			return
		}

		if self.Life <= 0 {
			return
		}

		// Which is the nearest ship?
		if self.lock == nil {
			var nearest *player
			var dist int64

			dist = -1

			for p, _ := range self.sector.players {
				if self.sameAs(p) == false {
					test := self.distance(p)
					if dist < 0 || test < dist {
						dist = test
						nearest = p
					}
				}
			}

			self.lock = nearest
		}

		if self.lock != nil {

			if self.lock.Life == 0 {
				// We've already killed this one.
				self.lock = nil
				continue
			}

			// TODO: make movements more natural.
			dx := self.lock.Position[0] - self.Position[0]
			dy := self.lock.Position[1] - self.Position[1]

			dw := self.distance(self.lock)

			// Angle distance between our ship and the ship we've got our lock on.
			t := math.Atan2(dy, dx) - math.Atan2(self.Direction[0], self.Direction[1])

			if math.Abs(t) > 0.1 {
				// Rotate
				if t > 0 {
					self.control.X = -1
				} else if t < 0 {
					self.control.X = 1
				}
				self.control.Y = 0
			} else {
				// Keep forward!
				self.control.X = 0
				self.control.Y = -1
			}

			if dw < 4e4 {
				// If we're near the ship, shoot.
				self.control.S = 1
			} else {
				// If not, just keep forward.
				self.control.Y = -1
				self.control.S = 0
			}

		}

		z = z + 1

		if z > 10 {
			// Keep the lock for a while.
			z = 0
			self.lock = nil
		}

		// Advancing tick counter.
		self.ticks = self.ticks + 1

		if self.ticks > agentMaxTicks {
			// Self-destroying the ship after having reached tick limit.
			self.destroy()
			return
		}

		time.Sleep(time.Millisecond * agentSleep)
	}
}
