package main

import (
	"flag"
	"strings"
	"fmt"
)

// This struct serves as storage for all of our relay configuration.
type RelayConfigStore struct {
    // ZMQ URIs to receive messages on.
    listeners []string
    // ZMQ URIs to relay messages over.
    senders []string
}

// Instantiate a global instance for config storage.
var RelayConfig RelayConfigStore

// Global temporary variables used for storing un-parsed user input.
var (
    flagZmqListeners string
    flagZmqSenders string
)

// Given a pointer to an un-parsed string containing the user's listeners or
// senders input, split it on commas and stuff it into the appropriate varible
// on the RelayConfig instance.
func parseZmqFlagString(flagStr string, zmqConfigDirective *[]string) {
    flagSplit := strings.Split(flagStr, ",")

    for _, uri := range flagSplit {
        // Add this ZMQ URI to appropriate slice of URIs.
        *zmqConfigDirective = append(*zmqConfigDirective, uri)
        fmt.Printf("  - %v\n", uri)
    }
}

// Set up the various CLI flags.
func init() {
	flag.StringVar(&flagZmqListeners, "listeners", "ipc:///tmp/announcer-sender.sock", "Comma-separated list of listener URIs.")
	flag.StringVar(&flagZmqSenders, "senders", "ipc:///tmp/relay-sender.sock", "Comma-separated list of sender URIs.")

    flag.Parse()
    println("================================================================================")
    println("                               ## emdr-relay-go ##")
    println("--------------------------------------------------------------------------------")

    println("* Connect to Announcers via SUB:")

    // Parse, store, and print the relay's listeners.
	parseZmqFlagString(flagZmqListeners, &RelayConfig.listeners)

	println("* Accepting SUB connections on:")

    // Parse, store, and print the relay's senders.
	parseZmqFlagString(flagZmqSenders, &RelayConfig.senders)

	println("================================================================================")
}

func main() {
	RunRelay()
}
