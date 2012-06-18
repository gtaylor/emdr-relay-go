package main

import (
	"fmt"
	zmq "github.com/alecthomas/gozmq"
	"github.com/bradfitz/gomemcache/memcache"
	"hash"
	"hash/fnv"
	"os"
)

func main() {
	mc := memcache.New("127.0.0.1:11211")

	context, _ := zmq.NewContext()

	receiver, _ := context.NewSocket(zmq.SUB)
	receiver.SetSockOptString(zmq.SUBSCRIBE, "")
	receiver.Connect("tcp://master.eve-emdr.com:8050")
	receiver.Connect("tcp://secondary.eve-emdr.com:8050")
	//receiver.Connect("tcp://relay-us-central-1.eve-emdr.com:8050")

	sender, _ := context.NewSocket(zmq.PUB)
	sender.Bind("tcp://0.0.0.0:8050")

	println("Listening on port 8050...")

	for {
		msg, zmq_err := receiver.Recv(0)

		if zmq_err != nil {
			println("RECV ERROR:", zmq_err.Error())
		}

		var h hash.Hash = fnv.New32()
		h.Write(msg)

		checksum := h.Sum([]byte{})
		var checksum_str string = fmt.Sprintf("%x", checksum)

		_, mc_err := mc.Get(checksum_str)

		mc_item := &memcache.Item{Key: checksum_str, Value: []byte{1}, Expiration: 300}

		if mc_err == memcache.ErrCacheMiss {
			// A cache miss means that the incoming message is not a dupe.
			// Set the cache entry to prevent future re-sends of this message.
			mc.Set(mc_item)
			// Send the message to subscribers.
			sender.Send(msg, 0)
			continue
		}

		if mc_err == nil {
			// No cache miss error means this is a dupe. Just re-set the
			// cache entry to bump the key expiration time.
			mc.Set(mc_item)
			//fmt.Printf("Dupe %x\n", checksum)
		} else {
			// Something else unexpected happened.
			println("MC ERROR:", mc_err.Error())
			println(" Checksum:", checksum)
			os.Exit(1)
		}

	}
}
