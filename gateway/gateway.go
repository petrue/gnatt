package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	G "github.com/alsm/gnatt/gateway/gate"
)

func main() {
	var gateway G.Gateway
	stopsig := registerSignals()
	gatewayconf := setup()

	//initLogger(sys.Stdout, Sys.Stderr) // todo: configurable

	if gatewayconf.IsAggregating() {
		fmt.Println("GNATT Gateway starting in aggregating mode")
		gateway = initAggregating(gatewayconf, stopsig)
	} else {
		fmt.Println("GNATT Gateway starting in transparent mode")
		gateway = initTransparent(gatewayconf, stopsig)
	}

	gateway.Start()
}

func setup() *G.GatewayConfig {
	var configFile string
	var port int

	flag.StringVar(&configFile, "c", "", "Configuration File")
	flag.IntVar(&port, "port", 0, "MQTT-G UDP Listening Port")
	flag.Parse()

	if configFile != "" {
		if gc, err := G.ParseConfigFile(configFile); err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			return gc
		}
	}

	fmt.Println("-configuration <file> must be specified")
	os.Exit(1)
	return nil
}

func initAggregating(gc *G.GatewayConfig, stopsig chan os.Signal) *G.AggGate {
	ag := G.NewAggGate(gc, stopsig)
	return ag
}

func initTransparent(gc *G.GatewayConfig, stopsig chan os.Signal) *G.TransGate {
	tg := G.NewTransGate(gc, stopsig)
	return tg
}

func registerSignals() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	return c
}
