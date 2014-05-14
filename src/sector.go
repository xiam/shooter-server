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
	"math/rand"
	"time"
)

const hyperWatchMaxTime = 1000.0

type sector struct {
	players  map[*player]bool
	fires    map[*fire]bool
	powerups map[*powerup]bool
	ids      map[string]interface{}
	offset   [2]int64
	bounds   [2]int64

	gbgPlayer  chan *player
	gbgFire    chan *fire
	gbgPowerup chan *powerup
}

func newSector() *sector {
	self := &sector{}
	self.players = map[*player]bool{}
	self.fires = map[*fire]bool{}
	self.powerups = map[*powerup]bool{}
	self.ids = map[string]interface{}{}

	self.offset = [2]int64{0, 0}
	self.bounds = [2]int64{10e3, 10e3}

	// Channels for removing discarded elements.
	self.gbgPlayer = make(chan *player, 512)
	self.gbgPowerup = make(chan *powerup, 512)
	self.gbgFire = make(chan *fire, 512)

	go self.run()
	go self.hyperWatch()
	return self
}

func (self *sector) run() {
	var p *player
	var u *powerup
	var f *fire

	var start, diff, sleep int64

	for {
		start = time.Now().UnixNano()

		for u = range self.powerups {
			u.Tick()
		}

		for f = range self.fires {
			f.Tick()
			f.update()
		}

		for p = range self.players {
			p.Tick()
			p.update()
		}

		// Removing discarded elements.
		removing := true
		for removing {
			select {
			case player := <-self.gbgPlayer:
				self.removePlayer(player)
			case powerup := <-self.gbgPowerup:
				self.removePowerup(powerup)
			case fire := <-self.gbgFire:
				self.removeFire(fire)
			default:
				removing = false
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
}

func (self *sector) removePowerup(u *powerup) {
	if _, ok := self.powerups[u]; ok == true {
		self.powerups[u] = false
		delete(self.powerups, u)
		delete(self.ids, u.Id)
		u.sector = nil
	}
}

func (self *sector) removeFire(f *fire) {
	if _, ok := self.fires[f]; ok == true {
		self.fires[f] = false
		delete(self.fires, f)
		delete(self.ids, f.Id)
		f.sector = nil
	}
}

func (self *sector) removePlayer(p *player) {
	if _, ok := self.players[p]; ok == true {
		self.players[p] = false
		delete(self.players, p)
		delete(self.ids, p.Id)
		p.sector = nil
	}
}

func (self *sector) takeId(prefix string, el interface{}) string {
	var id string
	for i := 0; ; i++ {
		id = fmt.Sprintf("%s-%d", prefix, rand.Int31n(9999))
		if _, ok := self.ids[id]; ok == false {
			self.ids[id] = el
			return id
		}
	}
	panic("reached")
}

func (self *sector) addPlayer(p *player) {
	var chunk []byte

	p.sector = self

	p.Id = self.takeId("ship", p)
	self.players[p] = true
	p.correct()

	if p.ws == nil {
		chunk = createFn("ship-ai", p.Id, p.Serialize())
	} else {
		chunk = createFn("ship", p.Id, p.Serialize())
	}

	// Announcing new player
	self.broadcast(chunk)

	// Announcing existing elements.
	for other, _ := range self.players {
		if p.sameAs(other) == false {
			p.notice(other)
		}
	}

}

func (self *sector) addPowerup(u *powerup) {
	u.sector = self

	u.Id = self.takeId("powerup", u)
	self.powerups[u] = true

	chunk := createFn("powerup", u.Id, u.Serialize())
	self.broadcast(chunk)
}

func (self *sector) addFire(b *fire) {
	b.sector = self

	b.Id = self.takeId("fire", b)
	self.fires[b] = true

	chunk := createFn("fire", b.Id, b.Serialize())
	self.broadcast(chunk)
}

func (self *sector) broadcast(chunk []byte) {
	for p, ok := range self.players {
		if ok {
			if p.ws != nil {
				p.write(chunk)
			}
		}
	}
}

func (self *sector) hyperWatch() {
	ratio := int64(1e3)
	for {
		// Keeping number of players in field constant.
		if len(self.players) < 20 {
			ai := newAgent()
			ai.player.Position = [2]float64{
				float64(rand.Int63n(ratio) - ratio/2),
				float64(rand.Int63n(ratio) - ratio/2),
			}
			self.addPlayer(ai.player)
		}
		// Keeping number of items constant.
		if len(self.powerups) < 30 {
			pu := newPowerup()
			// Random position.
			pu.Position = [2]float64{
				float64(rand.Int63n(self.bounds[0]) - self.bounds[0]/2),
				float64(rand.Int63n(self.bounds[1]) - self.bounds[1]/2),
			}
			self.addPowerup(pu)
		}
		log.Printf("players: %v, items: %v\n", len(self.players), len(self.powerups))
		time.Sleep(time.Millisecond * time.Duration(rand.Float32()*hyperWatchMaxTime))
	}
}
