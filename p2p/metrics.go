package p2p

import "github.com/empow-blockchain/go-empow/metrics"

var (
	neighborCountGauge = metrics.NewGauge("iost_p2p_neighbor_count", nil)
	routingCountGauge  = metrics.NewGauge("iost_p2p_routing_count", nil)
	byteOutCounter     = metrics.NewCounter("iost_p2p_bytes_out", []string{"mtype"})
	packetOutCounter   = metrics.NewCounter("iost_p2p_packet_out", []string{"mtype"})
	byteInCounter      = metrics.NewCounter("iost_p2p_bytes_in", []string{"mtype"})
	packetInCounter    = metrics.NewCounter("iost_p2p_packet_in", []string{"mtype"})
)
