# shooter.io

![shooter](./shooter.png)

Source code of the [shooter.io][1] server.

You may want to checkout the [shooter-html5][2] repo for a HTML5 client of this
server.

This is a work in progress.

```
cd ~/projects
git clone https://github.com/xiam/shooter-server.git
cd shooter-server
make
cd src
make
MONGO_HOST="10.0.0.123" ./shooter-server -listen 127.0.0.1:3223
```

## License

> Copyright 2014 José Carlos Nieto
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

[1]: http://shooter.io
[2]: https://github.com/xiam/shooter-html5
