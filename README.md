go-macaron websocket example
============================

> An example app using [Macaron](https://github.com/go-macaron/macaron) with websockets for Golang

This package uses the [go-macaron/sockets](https://github.com/go-macaron/sockets) middleware to provide sockets. The middleware is built on [gorilla/websocket](https://github.com/gorilla/websocket), so all the kids love it.

### Running

Run with `go run main.go`, then visit `http://localhost:4000/`

### Includes

* send messages back and forth
* using/displaying username
* connect/disconnect messages
* session-based username
* server-based disconnect
* additional route for displaying messages via POST (`/webhook`)

### Example

![demo](https://cl.ly/452U04203S28/Screen%20Recording%202017-07-19%20at%2004.20%20PM.gif)