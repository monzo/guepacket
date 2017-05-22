package guepacket

import (
	"github.com/google/gopacket"
	gplayers "github.com/google/gopacket/layers"
)

// GUELayerType should be populated when GUE has been registered with gopacket
var GUELayerType gopacket.LayerType

// GUE represents a packet encoded with Generic UDP Encapsulation
// It should sit "under" a UDP layer
//
// For more information about the meaning of the fields, see
// https://tools.ietf.org/html/draft-ietf-intarea-gue-04#section-3.1
type GUE struct {
	Version    uint8
	C          bool
	Protocol   gplayers.IPProtocol
	Flags      uint16
	Extensions []byte
	Data       []byte
}

func (l GUE) LayerType() gopacket.LayerType {
	return GUELayerType
}

func (l GUE) LayerContents() []byte {
	b := make([]byte, 4, 4+len(l.Extensions))
	hlen := uint8(len(l.Extensions))
	b[0] = l.Version<<6 | hlen
	if l.C {
		b[0] |= 0x20
	}
	b[0] |= hlen
	b[1] = byte(l.Protocol)
	b[2] = byte(l.Flags >> 8)
	b[3] = byte(l.Flags & 0xff)
	b = append(b, l.Extensions...)
	return b
}

func (l GUE) LayerPayload() []byte {
	return l.Data
}

func (l GUE) SerializeTo(buf gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	b := l.LayerContents()
	writeTo, err := buf.PrependBytes(len(b))
	if err != nil {
		return err
	}
	copy(writeTo, b)
	return nil
}

func (l GUE) CanDecode() gopacket.LayerClass {
	return GUELayerType
}

func (l *GUE) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	l.Version = data[0] >> 6
	l.C = data[0]&0x20 != 0
	l.Protocol = gplayers.IPProtocol(data[1])
	l.Flags = (uint16(data[2]) << 8) | uint16(data[3])
	hlen := data[0] & 0x1f
	l.Extensions = data[4 : 4+hlen]
	l.Data = data[4+hlen:]
	return nil
}

func (l GUE) NextLayerType() gopacket.LayerType {
	return l.Protocol.LayerType()
}

func DecodeGUE(data []byte, p gopacket.PacketBuilder) error {
	l := GUE{}
	if err := l.DecodeFromBytes(data, gopacket.NilDecodeFeedback); err != nil {
		return err
	}
	p.AddLayer(l)
	return p.NextDecoder(l.NextLayerType())
}
