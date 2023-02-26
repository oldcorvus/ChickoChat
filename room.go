package main

type room struct {
	// channel that holds incoming messages
	forward chan []byte
}
