package main

import (
	"fmt"
	"time"
)

// TODO Try to simulate windowsize and dropping data / reordering.

type header struct {
	syn        bool
	ack        int
	seq        int
	windowSize int
}

type clientSimulation struct {
	header         *header
	bytes          int
	channelRecieve chan header
	channelSend    chan header
}

type serverSimulation struct {
	header         *header
	channelRecieve chan header
	channelSend    chan header
	expectedSeq    int
}

func main() {
	// Ceater Header
	head := header{syn: false, ack: 0, seq: 0, windowSize: 0}
	// Create client
	client := clientSimulation{
		header:         &head,
		bytes:          10,
		channelRecieve: make(chan header, 1),
		channelSend:    make(chan header, 1),
	}
	// Create server
	server := serverSimulation{
		header:         &head,
		expectedSeq:    0,
		channelRecieve: make(chan header, 1),
		channelSend:    make(chan header, 1),
	}

	// Go routine for client and server
	go clientRoutine(&client, server.channelRecieve)
	go serverRoutine(&server, client.channelRecieve)

	// Remember not to exit main loop until all goroutines are done.
	time.Sleep(time.Second * 10000)
}

func clientRoutine(client *clientSimulation, serverChannel chan header) {
	client.header.syn = true
	client.header.seq = 0
	// First handshake. Syn flag on, and we chose to start seq at 0
	fmt.Println("Client sent first handshake, syn: ", client.header.syn, " and seq:", client.header.seq)
	serverChannel <- *client.header
	// Simulate we want to send X bytes (move to param)
	client.bytes = 20
	client.header.windowSize = 1000
	// Third handshake, and now connection is established and data can be sent.
	thirdHandshake := <-client.channelRecieve
	fmt.Println("Third handshake, server established connection, with Acknowledge", thirdHandshake.ack)
	fmt.Println("Client is now sending data of total length:", client.bytes)
	client.header.seq = thirdHandshake.ack
	serverChannel <- *client.header
	time.Sleep(time.Second * 3)
	for {
		select {
		case msg := <-client.channelRecieve:
			// Rediger her vores header inden vi sender til serveren.
			client.header.seq = msg.ack
			fmt.Println("Client fik besked fra server. med ACK: ", msg.ack)
			if msg.ack == client.bytes {
				// Close vores connection.
				fmt.Println("Client fik besked fra server at vi er færdige med at sende data.")
				break
			}
			serverChannel <- *client.header
		default:
			continue
		}
	}
}

func serverRoutine(server *serverSimulation, clientChannel chan header) {
	// Wait until host wants to establish connection
	firstHandshake := <-server.channelRecieve // First handshake from client.
	server.header.syn = firstHandshake.syn
	server.header.ack = firstHandshake.seq + 1
	fmt.Println("Server sent second handshake, syn:", server.header.syn, " and ack:", server.header.ack)
	clientChannel <- *server.header
	for {
		select {
		case msg := <-server.channelRecieve:
			server.expectedSeq = server.header.ack
			server.header.ack = msg.seq + 1
			fmt.Println("Server fik besked fra client. Med SEQ:", msg.seq, "and expected:", server.expectedSeq)
			clientChannel <- *server.header
		default:
			continue
		}
	}
}
