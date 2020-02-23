package mhfpacket

import (
	"github.com/Andoryuuta/Erupe/network"
	"github.com/Andoryuuta/byteframe"
)

// MsgMhfSaveHunterNavi represents the MSG_MHF_SAVE_HUNTER_NAVI
type MsgMhfSaveHunterNavi struct {
	AckHandle      uint32
	DataSize       uint32
	Unk0           bool
	RawDataPayload []byte
}

// Opcode returns the ID associated with this packet type.
func (m *MsgMhfSaveHunterNavi) Opcode() network.PacketID {
	return network.MSG_MHF_SAVE_HUNTER_NAVI
}

// Parse parses the packet from binary
func (m *MsgMhfSaveHunterNavi) Parse(bf *byteframe.ByteFrame) error {
	m.AckHandle = bf.ReadUint32()
	m.DataSize = bf.ReadUint32()
	m.Unk0 = bf.ReadBool()
	m.RawDataPayload = bf.ReadBytes(uint(m.DataSize))
	return nil
}

// Build builds a binary packet from the current data.
func (m *MsgMhfSaveHunterNavi) Build(bf *byteframe.ByteFrame) error {
	panic("Not implemented")
}
