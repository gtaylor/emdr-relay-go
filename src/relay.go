package main

import (
	"fmt"
	gozmq "github.com/alecthomas/gozmq"
	"github.com/bradfitz/gomemcache/memcache"
	"hash"
	"hash/fnv"
	"os"
)

// Since our relays are supervised by process monitors, and our architecture
// is extremely fault-tolerant, un-caught errors result in the relay exiting
// with an error code. The relay will be re-started immediately.
func errorHandler(errorPrefix string, errorMessage string) {
    println(errorPrefix, errorMessage)
    os.Exit(1)
}

// Handles the setup of the various variables, and the startup of the
// main relay loop.
func RunRelay() {
    mcClient := memcache.New("127.0.0.1:11211")

    context, zmqContextError := gozmq.NewContext()
    if zmqContextError != nil {
        errorHandler("zmqContextError:", zmqContextError.Error())
    }

    listener, listenSocketError := context.NewSocket(gozmq.SUB)
    if listenSocketError != nil {
        errorHandler("listenSocketError", listenSocketError.Error())
    }

    // This ZeroMQ socket is what we receive incoming messages on.
    listener.SetSockOptString(gozmq.SUBSCRIBE, "")

    // Connect to all of the requested listeners.
	for _, listener_uri := range RelayConfig.listeners {
	    listener.Connect(listener_uri)
    }

    // This ZeroMQ socket is where we relay the messages back out over.
    sender, senderSocketError := context.NewSocket(gozmq.PUB)
    if senderSocketError != nil {
        errorHandler("senderSocketError", senderSocketError.Error())
    }

    //"tcp://master.eve-emdr.com:8050"
	//"tcp://secondary.eve-emdr.com:8050"
	//"tcp://relay-us-central-1.eve-emdr.com:8050"

	// Bind all of the requested senders.
    for _, sender_uri := range RelayConfig.senders {
        sender.Bind(sender_uri)
    }

    // Let's get this party started.
    relayLoop(listener, sender, mcClient)
}

// The main relay loop. Receives incoming messages, spits them back out to
// any connected subscribers.
func relayLoop(listener gozmq.Socket, sender gozmq.Socket, mcClient *memcache.Client) {

	for {
	    // This blocks until something comes down the pipe.
		msg, listenRecvError := listener.Recv(0)
		if listenRecvError != nil {
		    errorHandler("listenRecvError", listenRecvError.Error())
		}

		var h hash.Hash = fnv.New32()
		h.Write(msg)

		checksum := h.Sum([]byte{})
		var checksum_str string = fmt.Sprintf("%x", checksum)

		_, mcError := mcClient.Get(checksum_str)

        // This needs to get set in the case of a cache miss error, or no error.
		mcItem := &memcache.Item{Key: checksum_str, Value: []byte{1}, Expiration: 300}

		if mcError == memcache.ErrCacheMiss {
			// A cache miss means that the incoming message is not a dupe.
			// Set the cache entry to prevent future re-sends of this message.
			mcClient.Set(mcItem)
			// Send the message to subscribers.
			sender.Send(msg, 0)
			continue
		}

		if mcError == nil {
			// No cache miss error means this is a dupe. Just re-set the
			// cache entry to bump the key expiration time.
			mcClient.Set(mcItem)
			//fmt.Printf("Dupe %x\n", checksum)
		} else {
			// Something else unexpected happened.
			errorHandler("mcError", mcError.Error())
		}

	}
}
