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
	"log"
	"os"
	"strings"
	"time"
	"upper.io/db.v2"
	"upper.io/db.v2/postgresql"
)

type mark struct {
	Name    string    `db:"name" json:"name"`
	Points  uint64    `db:"points" json:"points"`
	Created time.Time `db:"created" json:"-"`
}

var (
	sess   db.Database
	scores db.Collection
)

const (
	defaultDatabase = "shooter"
	defaultHost     = "127.0.0.1"
)

func init() {
	var err error

	dbAddr := os.Getenv("POSTGRESQL_PORT_5432_TCP_ADDR")
	if dbAddr == "" {
		dbAddr = defaultHost
	}

	settings := postgresql.ConnectionURL{
		User:     os.Getenv("POSTGRESQL_USER"),
		Password: os.Getenv("POSTGRESQL_PASSWORD"),
		Host:     dbAddr,
		Database: defaultDatabase,
	}

	if sess, err = postgresql.Open(settings); err != nil {
		log.Fatal("db.Open: ", err)
	}

	scores = sess.Collection("scores")
}

const maxScores = 5

type score struct {
}

func newScore() *score {
	return &score{}
}

func (self *score) Add(name string, points uint64) {
	name = strings.TrimSpace(name)

	if name != "" && points > 0 {
		_, err := scores.Insert(mark{name, points, time.Now()})
		if err != nil {
			log.Printf("Append: %v", err)
		}
	}
}

func (self *score) GetTop() (list []mark) {
	list = make([]mark, 0, maxScores)

	res := scores.Find().OrderBy("-points").Limit(maxScores)

	err := res.All(&list)
	if err != nil {
		log.Printf("res.All: %v", err)
	}

	return list
}
