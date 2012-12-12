package main

import (
    cache "code.google.com/p/vitess/go/cache"
	"fmt"
	zmq "github.com/alecthomas/gozmq"
	"hash"
	"hash/fnv"
	"os"
	"time"
	"unsafe"
)

type CacheValue struct {
	expires time.Time
}

const cache_value_size = uint64(unsafe.Sizeof(CacheValue{}))

func (self *CacheValue) Size() int {
	return int(cache_value_size)
}

func main() {
	cache := cache.NewLRUCache(cache_value_size * 100000)

	context, _ := zmq.NewContext()

	receiver, _ := context.NewSocket(zmq.SUB)
	receiver.SetSockOptString(zmq.SUBSCRIBE, "")
	receiver.Connect("tcp://master.eve-emdr.com:8050")
	receiver.Connect("tcp://secondary.eve-emdr.com:8050")
	//receiver.Connect("tcp://relay-us-central-1.eve-emdr.com:8050")

	sender, _ := context.NewSocket(zmq.PUB)
	sender.Bind("tcp://0.0.0.0:8050")

	cache_duration, err := time.ParseDuration("5m")
    if err != nil {
        println(err.Error())
        os.Exit(1)
    }

	println("Listening on port 8050...")

	for {
		msg, zmq_err := receiver.Recv(0)
		now := time.Now()

		if zmq_err != nil {
			println("RECV ERROR:", zmq_err.Error())
		}

		var h hash.Hash = fnv.New32()
		h.Write(msg)

		checksum := h.Sum([]byte{})
		cache_key := fmt.Sprintf("%x", checksum)

		cache_item, cache_hit := cache.Get(cache_key)
		item_expired := false
		if cache_item != nil {
			item_expired = now.After(cache_item.(*CacheValue).expires)
		}
		cache_item = &CacheValue{expires: now.Add(cache_duration)}
		// Insert (or reset) the cache entry to prevent future re-sends of this message.
		cache.Set(cache_key, cache_item)

		if !cache_hit {
			// A cache miss means that the incoming message is not a dupe.
			// Send the message to subscribers.
			sender.Send(msg, 0)
		} else {
			// Cache hit, dupe if not expired.
			if !item_expired {
				// Not expired - dupe.
			} else {
				// Expired - valid. Send to subscribers.
				sender.Send(msg, 0)
			}
		}

	}
}
