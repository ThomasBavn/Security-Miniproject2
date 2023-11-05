package main

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	gRPC "github.com/ThomasBavn/Security-Miniproject2/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var prime = 87178291199
var hospitalId = 5000

type peer struct {
	gRPC.UnimplementedNodeServer
	id             interface{}
	clients        map[int]gRPC.NodeClient
	receivedShares []int
	ctx            context.Context
}

func main() {
	log.SetFlags(log.Lshortfile)
	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32)
	ownPort := int(arg1) + hospitalId //hospital is 5000, clients are 5001 5002 5003

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := &peer{
		id:             ownPort,
		clients:        make(map[int]gRPC.NodeClient),
		receivedShares: []int{},
		ctx:            ctx,
	}

	// set up server
	list, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", ownPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	serverCert, err := credentials.NewServerTLSFromFile("certificate/server.crt", "certificate/priv.key")
	if err != nil {
		log.Fatalln("failed to create certificate", err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(serverCert))
	gRPC.RegisterNodeServer(grpcServer, p)

	// start the server
	go func() {
		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("failed to serve %v", err)
		}
	}()

	// Credit to Thor Andersen (TCLA) for helping setting up the TLS secure connection
	// Dial the other peers
	for i := 0; i < 4; i++ {
		port := hospitalId + i

		if port == ownPort {
			continue
		}

		// Set up client connections
		clientCert, err := credentials.NewClientTLSFromFile("certificate/server.crt", "")
		if err != nil {
			log.Fatalln("failed to create cert\n", err)
		}

		log.Printf("Trying to dial: %v\n", port)
		conn, err := grpc.Dial(fmt.Sprintf("localhost:%v", port), grpc.WithTransportCredentials(clientCert), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %s", err)
		}
		defer conn.Close()

		c := gRPC.NewNodeClient(conn)
		p.clients[port] = c
	}

	if ownPort == hospitalId {
		log.Println("Waiting for participants to send data")
	} else {
		//random := rand.New(rand.NewSource(int64(ownPort)))
		//secret := random.Intn(ownPort)

		// send its own port as secret for easy verification
		secret := ownPort
		p.DistributeData(secret)
	}

	// Keep main function alive until the program is terminated
	for true {
	}

}

func (p *peer) SendTo(clientId, data int) {
	client := p.clients[clientId]
	request := &gRPC.ExchangeRequest{Share: int32(data)}
	log.Printf("Client %v received %v", clientId, data)

	_, err := client.Exchange(p.ctx, request)
	if err != nil {
		log.Fatalf("Error!: %v", err)
	}
	log.Printf("Client %v succesfully received share", clientId)
}

func (p *peer) DistributeData(secret int) {
	shares := createShares(secret)

	// ensure that the peer itself gets the first share
	p.receivedShares = append(p.receivedShares, shares[0])
	log.Printf("port %v gave itself %v", p.id, shares[0])
	shareI := 1
	// send the rest of the shares to the other peers
	for id := range p.clients {
		if id == p.id || id == hospitalId {
			continue
		}
		p.SendTo(id, shares[shareI])
		shareI++
	}
	if len(p.receivedShares) == 3 {
		log.Printf("port %v has the following shares: %v", p.id, p.receivedShares)
		if p.id == hospitalId {
			log.Printf("Hospital has the final value %v", sum(p.receivedShares))
		} else {
			// send to hospital
			log.Printf("port %v sent %v to hospital", p.id, sum(p.receivedShares))
			p.SendTo(hospitalId, sum(p.receivedShares))
		}
	}
}

func (p *peer) Exchange(_ context.Context, request *gRPC.ExchangeRequest) (*emptypb.Empty, error) {

	log.Printf("port %v received %v", p.id, request.Share)

	p.receivedShares = append(p.receivedShares, int(request.Share))

	log.Printf("port %v has %v shares", p.id, len(p.receivedShares))

	if len(p.receivedShares) == 3 {

		if p.id == 5000 { // if peer is hospital
			log.Printf("Hospital has the final value %v", sum(p.receivedShares))

		} else {
			// send to hospital
			log.Printf("port %v sent %v to hospital", p.id, sum(p.receivedShares))
			p.SendTo(hospitalId, sum(p.receivedShares))
		}
	}
	return &emptypb.Empty{}, nil
}

func sum(shares []int) int {
	sum := 0
	for _, share := range shares {
		sum += share
	}
	return sum
}

func createShares(secret int) []int {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Ensure that a and b is between -prime and prime
	a := random.Intn(2*prime+1) - prime
	b := random.Intn(2*prime+1) - prime
	c := secret - a - b
	return []int{a, b, c}
}
