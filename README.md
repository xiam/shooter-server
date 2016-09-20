# shooter.io

![shooter](./shooter.png)

[shooter.io][1] is a network based multiplayer space shooter game.

## shooter-server

`shooter-server` is a Go program that runs a universe for space ships that can
shoot and destroy each other.

It also accepts incoming websockets connections, assigns each connection its
own ship, and lets the user control her ship.

## shooter-html5

The [shooter-html5][2] repo is a HTML5 client that uses canvas, JavaScript and
Websockets to display whatever is happening on `shooter-server`.

## How to run the server with vagrant?

```
cd ~/projects
git clone https://github.com/xiam/shooter-vagrant.git
vagrant up
# Go grab some coffee.
```

Open `10.2.2.10` in your browser!

See [shooter-vagrant][3].

## How to run the server manually?

```
cd ~/projects
git clone https://github.com/xiam/shooter-server.git
cd shooter-server
make
cd src
go get -d
make
MONGO_HOST="10.0.0.123" ./shooter-server -listen 127.0.0.1:3223
```

Now you have a running `shooter-server` that creates a virtual universe for
space ships, see [shooter-html5][2] for a client to interact with this
universe.

Note that, in order to connect to the `shooter-server`, you may change a variable
within the `shooter-html5/src/js/main.js` file. For instance, this line

```
var WEBSOCKET_SERVICE = 'ws://127.0.0.1:3223/w/';
```

will instruct the client to connect to the `shooter-server` that is listening
on `127.0.0.1:3223`.

## Current state

* This project was never finished, code is messy, racy and undocumented. I
  don't plan to continue working on this again anytime soon.
* I'd like to create a native client for mobile too.
* I'd like to experiment with UDP messages instead of TCP websockets.

## License

> Copyright 2014-today José Carlos Nieto
>
> Licensed under the Apache License, Version 2.0 (the "License");
> you may not use this file except in compliance with the License.
> You may obtain a copy of the License at
>
>     http:>www.apache.org/licenses/LICENSE-2.0
>
> Unless required by applicable law or agreed to in writing, software
> distributed under the License is distributed on an "AS IS" BASIS,
> WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
> See the License for the specific language governing permissions and
> limitations under the License.

[1]: https://shooter.io
[2]: https://github.com/xiam/shooter-html5
[3]: https://github.com/xiam/shooter-vagrant
