package test

import (
	"testing"
)

var chat chan string

func TestChan(t *testing.T) {
	var msgChanI <-chan string = chat
	select {
	case a := <-msgChanI:
		_ = a
	}

	var msgChanO chan<- string = chat
	msgChanO <- ""
}
