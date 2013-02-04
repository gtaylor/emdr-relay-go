package main

import (
	cache "code.google.com/p/vitess/go/cache"
	"fmt"
	zmq "github.com/alecthomas/gozmq"
	"hash"
	"hash/fnv"
	"unsafe"
)

// The presence of the cache value is all we need, so keep this super simple.
type CacheValue struct {
	found bool
}

// Calculate the size (in bytes) of our struct.
const cache_value_size = uint64(unsafe.Sizeof(CacheValue{}))

// Determines the max cache size, in bytes.
const cache_size_limit = cache_value_size * 1000

// Satisfies the Value interface.
func (self *CacheValue) Size() int {
	return int(cache_value_size)
}

func main() {
	cache := cache.NewLRUCache(cache_size_limit)

	context, _ := zmq.NewContext()
	receiver, _ := context.NewSocket(zmq.SUB)
	receiver.SetSockOptString(zmq.SUBSCRIBE, "")
	receiver.Connect("tcp://master.eve-emdr.com:8050")
	receiver.Connect("tcp://secondary.eve-emdr.com:8050")
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
