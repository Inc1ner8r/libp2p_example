package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	ma "github.com/multiformats/go-multiaddr"
)

func main() {
	masternode()
}

func masternode() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start a libp2p node with default settings
	node, err := libp2p.New(
	// libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/2000"),
	)
	if err != nil {
		panic(err)
	}

	// print the node's listening addresses
	fmt.Println("Listen addresses:", node.Addrs())

	// node.SetStreamHandler("/echo/1.0.0", func(s network.Stream) {
	// 	log.Println("listener received new stream")
	// 	if err := doEcho(s); err != nil {
	// 		log.Println(err)
	// 		s.Reset()
	// 	} else {
	// 		s.Close()
	// 	}
	// })

	// log.Println("listening for connections")

	StartListener(ctx, node, 2000)
	// wait for a SIGINT or SIGTERM signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("Received signal, shutting down...")

	// shut the node down
	if err := node.Close(); err != nil {
		panic(err)
	}

}

func StartListener(ctx context.Context, ha host.Host, listenPort int) {
	fullAddr := getHostAddress(ha)
	log.Printf("I am %s\n", fullAddr)

	// Set a stream handler on host A. /echo/1.0.0 is
	// a user-defined protocol name.
	ha.SetStreamHandler("/echo/1.0.0", func(s network.Stream) {
		log.Println("listener received new stream")
		if err := doEcho(s); err != nil {
			log.Println(err)
			s.Reset()
		} else {
			s.Close()
		}
	})

	log.Println("listening for connections")

}

func getHostAddress(ha host.Host) string {
	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr("/p2p/12D3KooWMHwXLsXMKto6b5yGXjHe5sXtHcJd5We6gd2HWpxhVZQy")

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addr := ha.Addrs()[0]
	return addr.Encapsulate(hostAddr).String()
}

// doEcho reads a line of data a stream and writes it back
func doEcho(s network.Stream) error {
	buf := bufio.NewReader(s)
	str, err := buf.ReadString('\n')
	if err != nil {
		return err
	}

	log.Printf("read: %s", str)
	_, err = s.Write([]byte(str))
	return err
}
