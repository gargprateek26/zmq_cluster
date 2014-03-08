package main

import (
	"fmt"
	"strconv"
	//	"strings"
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

import zmq "github.com/pebbe/zmq4"

const (BROADCAST = -1)

type Envelope struct {
	// On the sender side, Pid identifies the receiving peer. If instead, Pid is
	// set to cluster.BROADCAST, the message is sent to all peers. On the receiver side, the
	// Id is always set to the original sender. If the Id is not found, the message is silently dropped
	Pid int

	// An id that globally and uniquely identifies the message, meant for duplicate detection at
	// higher levels. It is opaque to this package.
	MsgId int

	// the actual message.
	Msg interface{}
}

type jsonfileObj struct {
	Server []Server_Address
}

type Server_Address struct { //struct correspondong to json config file
	Server_Id int
	Port_Addr int
}

type IServer interface {
	S_Id() int
	S_Addr() int

	Peer_Id() []int
	Peer_Addr() []int

	Out_box() chan *Envelope
	In_box() chan *Envelope
}

func (s *ThisServer) S_Id() int {
	return s.ServerId

}

func (s *ThisServer) S_Addr() int {
	return s.ServerAddress

}

func (s *ThisServer) Peer_Id() []int {
	return s.PeerId

}

func (s *ThisServer) Peer_Addr() []int {
	return s.PeerAddr

}

func (s *ThisServer) Out_box() chan *Envelope {
	return s.Outbox

}

func (s *ThisServer) In_box() chan *Envelope {
	return s.Inbox

}

type ThisServer struct {
	ServerId      int //server's own Id
	ServerAddress int //server's full address

	PeerId   []int //serverIds of the peers connected
	PeerAddr []int //addresses of the peers conncted

	Outbox chan *Envelope
	Inbox  chan *Envelope
}

func configServer(sId int, fileName string) ThisServer { //function to confiure the server

	conFile, err := ioutil.ReadFile("./" + fileName)
	if err != nil {
		fmt.Println("Error in reading file\n")
		os.Exit(1)
	}

	var Clus jsonfileObj
	err = json.Unmarshal(conFile, &Clus)
	if err != nil {
		fmt.Print("Error:", err)
	}

	//fmt.Println(Clus)
	//return Clus[0]

	var tempServer ThisServer

	serverCount := len(Clus.Server)
	peerCount := serverCount - 1

	tempServer.PeerId = make([]int, peerCount)
	tempServer.PeerAddr = make([]int, peerCount)

	tempIndex := 0

	for _, x := range Clus.Server {

		if sId != x.Server_Id {

			tempServer.PeerId[tempIndex] = x.Server_Id
			tempServer.PeerAddr[tempIndex] = x.Port_Addr
			//fmt.Println(tempServer.PeerId[tempIndex])
			tempIndex++
		} else {

			tempServer.ServerId = x.Server_Id
			tempServer.ServerAddress = x.Port_Addr
			//fmt.Println(tempServer.ServerId)

		}

	}

	tempServer.Inbox = make(chan *Envelope)
	tempServer.Outbox = make(chan *Envelope)
	
	go tempServer.In_Msg_Ntwrk()
//	time.Sleep(2*time.Second)
	go tempServer.Out_Msg_Ntwrk()
	
	
	return tempServer
}

func (s *ThisServer) Out_Msg_Ntwrk() { // sends message from Outbox channel to network
	//onto the network
	//fmt.Println("In Sending network method")

	sockets := make([]*zmq.Socket, len(s.Peer_Id()))

	for i, addr := range s.Peer_Addr() {
		sockets[i], _ = zmq.NewSocket(zmq.PUSH)
		sockets[i].Connect("tcp://localhost:" + strconv.Itoa(addr))

	}

	//sock.Bind(conn_str)
	for {
		//fmt.Println("In Out_Msg_Ntwrk")
		msg := <-s.Outbox
		//fmt.Println("Got msg")
		if BROADCAST == msg.Pid {
			for _, sock := range sockets {

				b, _ := json.Marshal(msg)
				sock.SendBytes(b, zmq.DONTWAIT)
				fmt.Println("message sent broadcast")

			}
		} else {
			for i, pid := range s.Peer_Id() {
				if pid == msg.Pid {
					b, _ := json.Marshal(msg)
					sockets[i].SendBytes(b, zmq.DONTWAIT)
					
					fmt.Println("message sent unicast")

				}
			}
		}

	}
}

func (s *ThisServer) In_Msg_Ntwrk() { //receives message from Network and places it on the Inbox channel
	
	sock, _ := zmq.NewSocket(zmq.PULL)

	conn_str := "tcp://*:" + strconv.Itoa(s.S_Addr())
	//fmt.Println("Address :- ",s.S_Addr())
	sock.Bind(conn_str)

	for {
		//fmt.Println("In In_Msg_Ntwrk")
		//fmt.Println(s.Inbox)
		data_rcvd, _ := sock.RecvBytes(0)
		//fmt.Println("sth recvd")
		
		var env Envelope
		json.Unmarshal(data_rcvd, &env)
		s.Inbox <- &env
		//fmt.Println("message received")
	} 

}

func (s *ThisServer) Msg_Rcv() {
	//rcvd := new(Envelope)
	//fmt.Println("receiving from inbox channel.")
	for {
		//fmt.Println("In Msg_Rcv")
		//rcvd := new(Envelope)
		rcvd := <- s.Inbox //get whatever message is on the channel
		//msg_slice := strings.Split(rcvd, ":")
		fmt.Println("Pid:", rcvd.Pid, "MsgId:", rcvd.MsgId, "Msg:", rcvd.Msg)

	}
}

func main() {

	var serverId int
	if len(os.Args) == 2 {
		serverId, _ = strconv.Atoi(os.Args[1])
	} else {
		fmt.Print("Recheck Arguments Passed\n")
		os.Exit(1)
	}

	myServer := configServer(serverId, "config.json")

	
	
	go myServer.Msg_Rcv()
	// Upload the messages onto the channel to be sent

	for i := 1; i <= 3; i++ {
		data_msg1 := Envelope{}
		data_msg1.Pid = 22 //sends a message to server with Id 22
		data_msg1.MsgId = 5
		
		data_msg1.Msg = strconv.Itoa(i) + " somebody there...?? :-)"
		myServer.Outbox <- &data_msg1         
		//fmt.Println("Uploaded into channel ") 
		
	}

	time.Sleep(10 * time.Second)

}

