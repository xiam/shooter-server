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
	"upper.io/db"
	_ "upper.io/db/mongo"
)

type mark struct {
	Name    string    `json:"name" bson:"name"`
	Points  uint64    `json:"points" bson:"points"`
	Created time.Time `json:"-" bson:"created"`
}

var settings db.Settings
var sess db.Database
var scores db.Collection

const (
	defaultDatabase = "shooter"
	defaultHost     = "127.0.0.1"
)

func init() {
	var err error

	host := os.Getenv("MONGO_HOST")

	if host == "" {
		host = defaultHost
	}

	settings = db.Settings{
		Host:     host,
		Database: defaultDatabase,
	}

	if sess, err = db.Open("mongo", settings); err != nil {
		log.Fatal("db.Open: ", err)
	}

	log.Printf("Connected to mongo://%s/%s.\n", host, defaultDatabase)

	scores, err = sess.Collection("scores")
	if err != nil {
		if err != db.ErrCollectionDoesNotExists {
			log.Fatal("db.Collection: ", err)
		}
	}
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
		_, err := scores.Append(mark{name, points, time.Now()})

		if err != nil {
			log.Printf("Append: %v", err)
		}
	}
}

func (self *score) GetTop() (list []mark) {
	list = make([]mark, 0, maxScores)

	res := scores.Find().Sort("-points").Limit(maxScores)

	err := res.All(&list)

	if err != nil {
		log.Printf("res.All: %v", err)
	}

	return list
}
