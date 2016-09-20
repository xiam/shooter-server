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
	"flag"
	"fmt"
	"github.com/davecheney/profile"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"runtime"
	"runtime/debug"
	"time"
)

const (
	fps  = 24               // Frames per second.
	fpsl = 1000 / fps       // Duration of a single (milliseconds)
	fpsn = 1000000000 / fps // Duration of a single frame (nanoseconds)
)

var (
	listenAddr      = flag.String("listen", "127.0.0.1:3223", "HTTP service address.")
	enableProfiling = flag.Bool("profile", false, "Enable application profiling.")
)

var (
	highScores *score
	mainSector *sector
)

func init() {
	mainSector = newSector()
	highScores = newScore()
}

func wsHandler(w http.ResponseWriter, r *http.Request) {

	// Websocket handshake.
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		return
	}

	addr := ws.RemoteAddr()

	log.Printf("Websocket accepted: %s\n", addr)

	// Creating player.
	one := newPlayer(ws)

	// Adding player to sector.
	mainSector.addPlayer(one)

	// Spawing writing goroutine.
	go one.writer()

	// Sending player identification.
	one.ident()

	// Blocking this function on a reader.
	one.reader()

	// Reader has stopped, good bye!
	log.Printf("Websocket finalized: %s", addr)
}

func memProfile() {
	var gcstats *debug.GCStats
	var stats *runtime.MemStats

	stats = &runtime.MemStats{}
	gcstats = &debug.GCStats{}

	for {
		fmt.Println("STATS")

		runtime.ReadMemStats(stats)
		fmt.Printf("EnableGC: %v.\n", stats.EnableGC)
		fmt.Printf("LastGC: %d.\n", stats.LastGC)
		fmt.Printf("Mallocs: %d.\n", stats.Mallocs)
		fmt.Printf("Frees: %d.\n", stats.Frees)
		fmt.Printf("Mallocs - Frees: %d.\n", stats.Mallocs-stats.Frees)

		debug.ReadGCStats(gcstats)
		fmt.Printf("LastGC: %v.\n", gcstats.LastGC)
		fmt.Printf("NumGC: %d.\n", gcstats.NumGC)

		time.Sleep(time.Second * 2)
		fmt.Println("")
		fmt.Println("")
	}
}

func main() {

	flag.Parse()

	if *enableProfiling == true {
		defer profile.Start(profile.CPUProfile).Stop()
	}

	//http.Handle("/", http.FileServer(http.Dir("../html/")))
	http.HandleFunc("/w/", wsHandler)

	log.Printf("Listening on %s.\n", *listenAddr)
	//go memProfile()

	if err := http.ListenAndServe(*listenAddr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
