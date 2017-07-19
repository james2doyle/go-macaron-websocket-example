package main

import (
  "log"
  "net/http"
  "time"

  "github.com/go-macaron/gzip"
  "github.com/go-macaron/sockets"
  "github.com/gorilla/websocket"
  "gopkg.in/macaron.v1"
)

type (
  eventMessage struct {
    User    string `json:"user" binding:"Required"`
    Message string `json:"message" binding:"Required"`
  }
  errorResponse struct {
    Message string `json:"message"`
  }
)

func main() {
  m := macaron.Classic()

  // macaron.Classic() is a wrapper for:
  // m := macaron.New()
  // m.Use(Logger())
  // m.Use(Recovery())
  // m.Use(Static("public"))

  // support HEAD
  m.SetAutoHead(true)
  // render ctx responses
  m.Use(macaron.Renderer())
  // gzip responses
  m.Use(gzip.Gziper())

  // collect all the channels that need to be notified
  senders := make(map[string]chan<- *eventMessage)

  m.Get("/ws", sockets.JSON(eventMessage{}), func(receiver <-chan *eventMessage, sender chan<- *eventMessage, done <-chan bool, disconnect chan<- int, errorChannel <-chan error, ctx *macaron.Context) {
    // count down 30 minutes for disconnect
    ticker := time.After(30 * time.Minute)
    for {
      select {
      case msg := <-receiver:
        // here we simply echo the received message to the sender for demonstration purposes
        // In your app, collect the senders of different clients and do something useful with them
        // sender <- msg
        if senders[msg.User] == nil {
          senders[msg.User] = sender
        }
        // range over the connections and send the message out to each one
        for k := range senders {
          senders[k] <- msg
        }
      case <-ticker:
        // This will close the connection after 30 minutes no matter what
        // To demonstrate use of the disconnect channel
        // You can use close codes according to RFC 6455
        disconnect <- websocket.CloseNormalClosure
      case <-done:
        // the client disconnected, so you should return / break if the done channel gets sent a message
        return
      case err := <-errorChannel:
        // Uh oh, we received an error. This will happen before a close if the client did not disconnect regularly.
        // Maybe useful if you want to store statistics
        ctx.Error(500, "an error occured")
        log.Println(err)
      }
    }
  })

  // handle Internal Server Error issues
  m.InternalServerError(func(ctx *macaron.Context) {
    response := errorResponse{
      Message: "A server error has occurred",
    }

    ctx.JSON(http.StatusInternalServerError, &response)
  })

  // handle Not Found issues
  m.NotFound(func(ctx *macaron.Context) {
    response := errorResponse{
      Message: "The route could not be found",
    }

    ctx.JSON(http.StatusNotFound, &response)
  })

  // run the application
  log.Println("Server is running...")
  log.Println(http.ListenAndServe("0.0.0.0:4000", m))
}
