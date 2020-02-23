package mhfpacket

import (
	"github.com/Andoryuuta/Erupe/network"
	"github.com/Andoryuuta/byteframe"
)

// MsgMhfSaveScenarioData represents the MSG_MHF_SAVE_SCENARIO_DATA
type MsgMhfSaveScenarioData struct {
	AckHandle      uint32
	DataSize       uint32
	RawDataPayload []byte
}

// Opcode returns the ID associated with this packet type.
func (m *MsgMhfSaveScenarioData) Opcode() network.PacketID {
	return network.MSG_MHF_SAVE_SCENARIO_DATA
}

// Parse parses the packet from binary
func (m *MsgMhfSaveScenarioData) Parse(bf *byteframe.ByteFrame) error {
	m.AckHandle = bf.ReadUint32()
	m.DataSize = bf.ReadUint32()
	m.RawDataPayload = bf.ReadBytes(uint(m.DataSize))
	return nil
}

// Build builds a binary packet from the current data.
func (m *MsgMhfSaveScenarioData) Build(bf *byteframe.ByteFrame) error {
	panic("Not implemented")
}
