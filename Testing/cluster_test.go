package cluster

import (
	"fmt"
	//zmq "github.com/pebbe/zmq4"
	cluster "github.com/gargprateek26/zmq_cluster"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

var mutex = &sync.Mutex{}

var serverCount = 6

var serverCluster = make([]cluster.ThisServer, serverCount)

var msgRecvdCount int
var msgSentCount int

func TestBroadcast(t *testing.T) {
	msgSentCount = 0

	msgRecvdCount = 0

	for id := 0; id < serverCount; id++ {
		serverCluster[id] = cluster.ConfigServer(id*11, "config.json")

		//go serverCluster[id].Msg_Rcv() //calling this function would flush the channel twice
		go countRecv(&serverCluster[id], &msgRecvdCount)
	}
	for i := 1; i <= 10; i++ {
		data_msg1 := cluster.Envelope{}
		data_msg1.Pid = cluster.BROADCAST
		data_msg1.MsgId = i

		data_msg1.Msg = strconv.Itoa(i) + " somebody there...?? :-)"
		serverCluster[2].Outbox <- &data_msg1 //server 2 broadcasts the message to all other servers
		msgSentCount++

	}
	time.Sleep(2 * time.Second)
	var totalMsgSent = msgSentCount * (serverCount - 1)

	fmt.Println("Messages Sent: ", totalMsgSent, "Messages Received: ", msgRecvdCount)
	fmt.Println("Messages Dropped: ", totalMsgSent-msgRecvdCount)
	if (totalMsgSent - msgRecvdCount) != 0 {
		t.Errorf("Error in transmission")
	} else {
		fmt.Println("Broadcast Test Done\n")
	}

}

func TestUnicast(t *testing.T) {
	msgSentCount = 0
	msgRecvdCount = 0
	/*var serverCluster = make([]cluster.ThisServer, serverCount)

	for id := 0; id < serverCount; id++ {
		serverCluster[id] = cluster.ConfigServer(id*11, "config.json")

		//go serverCluster[id].Msg_Rcv() //calling this function would flush the channel twice
		go countRecv(&serverCluster[id], &msgRecvdCount)
	}
	*/

	for i := 1; i <= 100; i++ {

		data_msg1 := cluster.Envelope{}
		data_msg1.Pid = 11
		data_msg1.MsgId = i

		data_msg1.Msg = strconv.Itoa(i) + " somebody there...?? :-)"

		serverCluster[5].Outbox <- &data_msg1
		msgSentCount++
		//fmt.Println("sent")
	}
	//fmt.Println("Uploaded into channel ")
	time.Sleep(2 * time.Second)

	fmt.Println("Messages Sent: ", msgSentCount, "Messages Received: ", msgRecvdCount)
	fmt.Println("Messages Dropped: ", msgSentCount-msgRecvdCount)
	if (msgSentCount - msgRecvdCount) != 0 {
		t.Errorf("Error in transmission")
	} else {
		fmt.Println("Unicast Test Done\n")

	}

}

func TestRandomBroadcast(t *testing.T) {
	msgSentCount = 0
	msgRecvdCount = 0
	//var serverCluster = make([]cluster.ThisServer, serverCount)

	/*for id := 0; id < serverCount; id++ {
		serverCluster[id] = cluster.ConfigServer(id*11, "config.json")

		//go serverCluster[id].Msg_Rcv() //calling this function would flush the channel twice
	go countRecv(serverCluster[id], &msgRecvdCount)
	}*/

	var selServer int
	for i := 0; i < 100; i++ {
		selServer = rand.Intn(100) % (serverCount)
		data_msg1 := cluster.Envelope{}
		data_msg1.Pid = cluster.BROADCAST
		data_msg1.MsgId = i
		data_msg1.Msg = strconv.Itoa(i) + " somebody there...?? :-)"
		serverCluster[selServer].Outbox <- &data_msg1
		msgSentCount++
	}
	time.Sleep(2 * time.Second)
	var totalMsgSent = msgSentCount * (serverCount - 1)

	fmt.Println("Messages Sent: ", totalMsgSent, "Messages Received: ", msgRecvdCount)
	
	fmt.Println("Messages dropped: ", totalMsgSent-msgRecvdCount)

	if (totalMsgSent - msgRecvdCount) != 0 {
		t.Errorf("Error in transmission")
	} else {
		fmt.Println("Random Broadcast Done\n")

	}

}

func TestLargeMsgBroadcast(t *testing.T) {
	msgSentCount = 0
	msgRecvdCount = 0
	/*var serverCluster = make([]cluster.ThisServer, serverCount)

	for id := 0; id < serverCount; id++ {
		serverCluster[id] = cluster.ConfigServer(id*11, "config.json")

		//go serverCluster[id].Msg_Rcv() //calling this function would flush the channel twice
		go countRecv(serverCluster[id], &msgRecvdCount)
	}*/
	for i := 1; i <= 10; i++ {
		data_msg1 := cluster.Envelope{}
		data_msg1.Pid = cluster.BROADCAST
		data_msg1.MsgId = i

		data_msg1.Msg = strconv.Itoa(i) + strings.Repeat(" somebody there...?? :-)", 1024)
		serverCluster[2].Outbox <- &data_msg1 //server 2 broadcasts the message to all other servers
		msgSentCount++

		//fmt.Println("Uploaded into channel ")
	}
	time.Sleep(2 * time.Second)
	var totalMsgSent = msgSentCount * (serverCount - 1)
	fmt.Println("Messages Sent: ", totalMsgSent, "Messages Received: ", msgRecvdCount)
	

	fmt.Println("Messages Dropped: ", totalMsgSent-msgRecvdCount)
	if (totalMsgSent - msgRecvdCount) != 0 {
		t.Errorf("Error in transmission")
	} else {
		fmt.Println("Large Message Broadcast Test Done\n")

	}
}

func countRecv(s *cluster.ThisServer, recvdCount *int) {
	for {

		_ = <-s.In_box()

		mutex.Lock()

		*recvdCount += 1

		mutex.Unlock()

		//fmt.Println("still counting", *recvdCount)

	}
}
