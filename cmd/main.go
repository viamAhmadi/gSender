package main

import (
	"flag"
	"fmt"
	gConn "github.com/viamAhmadi/gReceiver2/pkg/conn"
	"github.com/viamAhmadi/gSender/pkg/conn"
	"io/ioutil"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func main() {
	addrBroker := flag.String("b", "tcp://127.0.0.1:5555", "Broker address")
	addrDes := flag.String("d", "tcp://127.0.0.3:5555", "Destination of messages")
	msgCount := flag.Int("c", 10000, "Number of messages, MAX=50K")
	connId := flag.String("i", "generate", "Connection id, REQUIRE_LEN=20")
	flag.Parse()

	if *connId == "generate" {
		*connId = randomString(20)
	}

	if len(*connId) != 20 {
		panic("Connection id, REQUIRE_LEN=20")
	}

	data, err := ioutil.ReadFile("./txt")
	if err != nil {
		panic(err)
	}

	fmt.Println("start sender")
	c := conn.NewSendConn(*addrDes, *connId, *msgCount)
	if err := c.OpenConnection(*addrBroker); err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case s := <-c.SendCh:
				fmt.Println("start sending, connId: ", s)
				sendMsg(c, data)
			case f := <-c.FactorCh:
				c.Factor = f
				result := "NO"
				if f.Successful == gConn.YES {
					result = "YES"
				}
				fmt.Printf("FACTOR\t\tsuccessful: %s\tlost: %d\n", result, len(*f.List))
				//c.Close()
			case <-time.After(4 * time.Second):
				if c.Factor == nil {
					fmt.Println("failed")
				}
				fmt.Println("\n\ntimeout 4 sec")
				c.Close()
				return
			}
		}
	}()

	if err := c.Receiving(); err != nil {
		fmt.Println(err)
	}
}

func sendMsg(c *conn.SendConn, data []byte) {
	sent := 0
	for i := 1; i <= c.Count; i++ {
		if i == 10000 || i == 20000 || i == 30000 {
			time.Sleep(1 * time.Second)
		}
		size := rand.Intn((8192 - 50 + 1) + 50)
		fmt.Printf("\nsend message id: %d \t size: %d byte", i, size)
		if err := c.Send(&gConn.Message{Id: i, ConnId: c.Id, Content: string(data[size])}); err != nil {
			fmt.Println(err.Error())
			if err == gConn.ErrDealer {
				return
			}
		}
		sent += 1
	}
	fmt.Println("\n\nsuccessfully sent ", sent)
}
