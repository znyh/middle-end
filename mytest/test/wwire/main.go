package main

import (
    "fmt"
)

type (
    Event struct {
        Greeter Greeter // <- adding a Greeter field
    }

    Greeter struct {
        Message Message
    }
    Message string
)

func NewEvent(g Greeter) Event {
    return Event{Greeter: g}
}
func NewGreeter(m Message) Greeter {
    return Greeter{Message: m}
}
func NewMessage() Message {
    return Message("Hi there!")
}

func main() {
    event := InitializeEvent()
    event.Start()
}

func (g Greeter) Greet() Message {
    return g.Message
}

func (e Event) Start() {
    msg := e.Greeter.Greet()
    fmt.Println(msg)
}
