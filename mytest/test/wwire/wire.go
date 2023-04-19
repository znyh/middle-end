//go:build wireinject
// +build wireinject

package main

import "github.com/google/wire"

func InitializeEvent() Event {
    panic(wire.Build(NewEvent, NewGreeter, NewMessage))
}
