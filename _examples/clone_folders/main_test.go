package main

import (
	"fmt"
	"testing"
)

type wtf struct {
	n int
}

func readChan(in chan *wtf) {
	for {
		val := <-in
		fmt.Println(val.n)
	}
}

func TestChannels(t *testing.T) {

	theChan := make(chan *wtf)
	for i := 0; i < 10; i++ {
		fmt.Println(i)
		theChan <- &wtf{n: i}
	}
	go readChan(theChan)

}
