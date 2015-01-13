package main

import (
	cache "github.com/gtaylor/emdr-relay-go/cache"
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"hash"
	"hash/fnv"
	"unsafe"
	"time"
	"os"
)

// The presence of the cache value is all we need, so keep this super simple.
type CacheValue struct {
	found bool
}

// Calculate the size (in bytes) of our struct.
const cache_value_size = int64(unsafe.Sizeof(CacheValue{}))

// Determines the max cache size, in bytes.
const cache_size_limit = cache_value_size * 1000

// Satisfies the Value interface.
func (self *CacheValue) Size() int {
	return int(cache_value_size)
}

func periodic_exiter() {
	// We exit periodically so that the process supervisor can restart us.
	// This helps recover from some edge cases where connections to the
	// announcers aren't picked back up.
	// Currently hardcoded to every 12 hours.
	ticker := time.NewTicker(12 * 3600 * time.Second)
	for {
		select {
		case <- ticker.C:
			println("Exiting so we can auto-restart.")
			os.Exit(0)
		}
	}
}

func main() {
	println("=====================[ emdr-relay-go ]=====================")
	println("Starting emdr-relay-go 1.1...")
	cache := cache.NewLRUCache(cache_size_limit)

	receiver, _ := zmq.NewSocket(zmq.SUB)
	receiver.Connect("tcp://master.eve-emdr.com:8050")
	receiver.Connect("tcp://secondary.eve-emdr.com:8050")
	receiver.SetSubscribe("")
	defer receiver.Close()

	sender, _ := zmq.NewSocket(zmq.PUB)
	sender.Bind("tcp://*:8050")
	defer sender.Close()

	println("Listening on port 8050.")
	println("===========================================================")
	//  Ensure subscriber connection has time to complete
	time.Sleep(time.Second)

	// We auto-exit every number of hours so we can recover from some
	// weird edge case conditions that disrupt the network. They're not common,
	// but we'll do this to be absolutely sure.
	go periodic_exiter()

	for {
		msg, zmq_err := receiver.Recv(0)
		if zmq_err != nil {
			println("RECV ERROR:", zmq_err.Error())
		}

		var h hash.Hash = fnv.New32()
		h.Write([]byte(msg))

		checksum := h.Sum([]byte{})
		cache_key := fmt.Sprintf("%x", checksum)

		cache_item, cache_hit := cache.Get(cache_key)
		if cache_hit {
			// We've already seen this before, ignore it.
			continue
		}

		// At this point, we know we've encountered a message we haven't
		// seen in the recent past.
		cache_item = &CacheValue{found: true}
		// Insert the cache entry to prevent future re-sends of this message.
		cache.Set(cache_key, cache_item)

		// A cache miss means that the incoming message is not a dupe.
		// Send the message to subscribers.
		sender.Send(msg, 0)
	}
}
