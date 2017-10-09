package guepacket

import (
	"encoding/hex"
	"os"
	"testing"

	"github.com/google/gopacket"
	gplayers "github.com/google/gopacket/layers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	GUELayerType = gopacket.RegisterLayerType(120001, gopacket.LayerTypeMetadata{
		Name:    "GUE",
		Decoder: gopacket.DecodeFunc(DecodeGUE)})
	gplayers.RegisterUDPPortLayerType(7777, GUELayerType)
	os.Exit(m.Run())
}

func TestDecoding(t *testing.T) {
	// This is a packet generated by the Linux kernel GUE implementation, captured
	// by pcap. It includes:
	// - Ethernet
	// - IPv4
	// - UDP
	// - GUE (port 7777)
	// - IPv4
	// - ICMP (ping)
	ph := `02427b2522f502420ae0d90608004500007451ea4000ff119d050ae0d9060afa9da88c0f1e6100608cfa000400004500005459f240004001e2cd0ae0d9060afd0f0608000a7e005f000cea811f59000000005daa080000000000000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f2021222324252627`
	pr, err := hex.DecodeString(ph)
	if err != nil {
		t.Errorf("Error decoding hex packet: %v", err)
	}

	p := gopacket.NewPacket(pr, gplayers.LayerTypeEthernet, gopacket.Default)
	require.Nil(t, p.ErrorLayer())
	t.Logf("%v", p)

	gue := p.Layer(GUELayerType).(*GUE)
	require.NotNil(t, gue)
	assert.Equal(t, uint8(0), gue.Version)
}
