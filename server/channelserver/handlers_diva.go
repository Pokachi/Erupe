package channelserver

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"erupe-ce/common/stringsupport"
	_config "erupe-ce/config"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
	"math/rand"
	"time"

	"erupe-ce/common/byteframe"
	"erupe-ce/network/mhfpacket"
)

func cleanupDiva(s *Session) {
	s.server.db.Exec(`DELETE FROM events WHERE event_type='diva'`)
	s.server.db.Exec(`DELETE FROM diva_beads`)
	s.server.db.Exec(`DELETE FROM diva_beads_assignment`)
	s.server.db.Exec(`DELETE FROM diva_beads_points`)
	s.server.db.Exec(`DELETE FROM diva_buffs`)
}

func generateDivaTimestamps(s *Session, start uint32, debug bool) []uint32 {
	timestamps := make([]uint32, 6)
	midnight := TimeMidnight()
	if debug && start <= 3 {
		midnight := uint32(midnight.Unix())
		switch start {
		case 1:
			timestamps[0] = midnight
			timestamps[1] = timestamps[0] + 601200
			timestamps[2] = timestamps[1] + 3900
			timestamps[3] = timestamps[1] + 604800
			timestamps[4] = timestamps[3] + 3900
			timestamps[5] = timestamps[3] + 604800
		case 2:
			timestamps[0] = midnight - 605100
			timestamps[1] = midnight - 3900
			timestamps[2] = midnight
			timestamps[3] = timestamps[1] + 604800
			timestamps[4] = timestamps[3] + 3900
			timestamps[5] = timestamps[3] + 604800
		case 3:
			timestamps[0] = midnight - 1213800
			timestamps[1] = midnight - 608700
			timestamps[2] = midnight - 604800
			timestamps[3] = midnight - 3900
			timestamps[4] = midnight
			timestamps[5] = timestamps[3] + 604800
		}
		return timestamps
	}
	if start == 0 || TimeAdjusted().Unix() > int64(start)+2977200 {
		cleanupDiva(s)
		// Generate a new diva defense, starting midnight tomorrow
		start = uint32(midnight.Add(24 * time.Hour).Unix())
		s.server.db.Exec("INSERT INTO events (event_type, start_time) VALUES ('diva', to_timestamp($1)::timestamp without time zone)", start)
		// Generate 4 random beads
		beads := []uint8{1, 3, 4, 8, 9, 10, 11, 14, 15, 17, 18, 19, 20, 21, 22, 23, 24, 25}
		for {
			if len(beads) == 4 {
				break
			}
			result := rand.Intn(len(beads))
			beads[result] = beads[len(beads)-1]
			beads = beads[:len(beads)-1]
		}
		s.server.db.Exec(`INSERT INTO diva_beads (type) VALUES ($1), ($2), ($3), ($4)`, beads[0], beads[1], beads[2], beads[3])
	}
	timestamps[0] = start
	timestamps[1] = timestamps[0] + 601200
	timestamps[2] = timestamps[1] + 3900
	timestamps[3] = timestamps[1] + 604800
	timestamps[4] = timestamps[3] + 3900
	timestamps[5] = timestamps[3] + 604800
	return timestamps
}

func handleMsgMhfGetUdSchedule(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetUdSchedule)
	bf := byteframe.NewByteFrame()

	id, start := uint32(0xCAFEBEEF), uint32(0)
	_ = s.server.db.QueryRow("SELECT id, (EXTRACT(epoch FROM start_time)::int) as start_time FROM events WHERE event_type='diva'").Scan(&id, &start)

	var timestamps []uint32
	if s.server.erupeConfig.DevMode && s.server.erupeConfig.DevModeOptions.DivaEvent >= 0 {
		if s.server.erupeConfig.DevModeOptions.DivaEvent == 0 {
			if s.server.erupeConfig.RealClientMode >= _config.Z2 {
				doAckBufSucceed(s, pkt.AckHandle, make([]byte, 32))
			} else {
				doAckBufSucceed(s, pkt.AckHandle, make([]byte, 36))
			}
			return
		}
		timestamps = generateDivaTimestamps(s, uint32(s.server.erupeConfig.DevModeOptions.DivaEvent), true)
	} else {
		timestamps = generateDivaTimestamps(s, start, false)
	}
  
	if s.server.erupeConfig.RealClientMode >= _config.Z2 {
		bf.WriteUint32(id)
	}
	for i := range timestamps {
		bf.WriteUint32(timestamps[i])
	}

	// Apparently buff multipliers
	bf.WriteUint16(22)
	bf.WriteUint16(45)
	bf.WriteUint16(2)
	bf.WriteUint16(2)

	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}

func handleMsgMhfGetUdInfo(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetUdInfo)
	// Message that appears on the Diva Defense NPC and triggers the green exclamation mark
	udInfos := []struct {
		Text      string
		StartTime time.Time
		EndTime   time.Time
	}{}

	resp := byteframe.NewByteFrame()
	resp.WriteUint8(uint8(len(udInfos)))
	for _, udInfo := range udInfos {
		resp.WriteBytes(stringsupport.PaddedString(udInfo.Text, 1024, true))
		resp.WriteUint32(uint32(udInfo.StartTime.Unix()))
		resp.WriteUint32(uint32(udInfo.EndTime.Unix()))
	}

	doAckBufSucceed(s, pkt.AckHandle, resp.Data())
}

func handleMsgMhfGetKijuInfo(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetKijuInfo)
	beads := []uint8{1, 3, 4, 8}
	rows, err := s.server.db.Query(`SELECT * FROM diva_beads`)
	if err == nil {
		var i int
		for rows.Next() {
			var bead uint8
			rows.Scan(&bead)
			beads[i] = bead
			i++
		}
	}
	kijuInfo := []struct {
		Color  uint8
		Effect uint8
	}{
		{1, beads[0]},
		{2, beads[1]},
		{3, beads[2]},
		{4, beads[3]},
	}
	bf := byteframe.NewByteFrame()
	bf.WriteUint8(uint8(len(kijuInfo)))
	for _, kiju := range kijuInfo {
		bf.WriteBytes(stringsupport.PaddedString(s.server.dict[fmt.Sprintf(`prayerBead%dName`, kiju.Effect)], 32, true))
		bf.WriteBytes(stringsupport.PaddedString(s.server.dict[fmt.Sprintf(`prayerBead%dDescription`, kiju.Effect)], 512, true))
		bf.WriteUint8(kiju.Color)
		bf.WriteUint8(kiju.Effect)
	}
	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}

func handleMsgMhfSetKiju(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfSetKiju)
	midday := TimeMidnight().Add(12 * time.Hour)
	if TimeAdjusted().After(midday) {
		midday = midday.Add(12 * time.Hour)
	}
	s.server.db.Exec(`INSERT INTO diva_beads_assignment VALUES ($1, $2, $3) ON CONFLICT (character_id, expiry) DO UPDATE SET bead_index = excluded.bead_index`, s.charID, pkt.BeadIndex, midday)
	doAckBufSucceed(s, pkt.AckHandle, []byte{0})
}

func handleMsgMhfAddUdPoint(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfAddUdPoint)
	var beadIndex = 0
	s.server.db.QueryRow(`SELECT bead_index FROM diva_beads_assignment WHERE character_id=$1 AND expiry>$3`, s.charID, TimeAdjusted()).Scan(&beadIndex)

	if (beadIndex == 0) {
		s.server.db.QueryRow(`SELECT bead_index FROM diva_beads_assignment WHERE character_id=$1 ORDER BY expiry DESC`, s.charID,).Scan(&beadIndex)
		if (beadIndex == 0) {
			doAckBufSucceed(s, pkt.AckHandle, []byte{0})	// bead not selected, do nothing
			return;
		} else {						// carry over last selected bead
			midday := TimeMidnight().Add(12 * time.Hour)
			if TimeAdjusted().After(midday) {
				midday = midday.Add(12 * time.Hour)
			}
			s.server.db.Exec(`INSERT INTO diva_beads_assignment VALUES ($1, $2, $3) ON CONFLICT (character_id, expiry) DO UPDATE SET bead_index = excluded.bead_index`, s.charID, beadIndex, midday)
		}
	}
	
	s.server.db.Exec(`INSERT INTO diva_beads_points VALUES ($1, $2, now(), $3)`, s.charID, pkt.Points * 1000, beadIndex)
	doAckBufSucceed(s, pkt.AckHandle, []byte{0})
}

func handleMsgMhfGetUdMyPoint(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetUdMyPoint)
	bf := byteframe.NewByteFrame()
	bf.WriteUint8(0) // No error
	var startTime time.Time
	s.server.db.QueryRow(`SELECT start_time FROM events WHERE event_type='diva'`).Scan(&startTime)
	for i := 0; i < 8; i++ {
		for j := 0; j < 2; j++ {
			var bead uint8
			var points uint32
			s.server.db.QueryRow(`SELECT bead_index FROM diva_beads_assignment WHERE expiry=$1 AND character_id=$2`, startTime, s.charID).Scan(&bead)
			s.server.db.QueryRow(`SELECT COALESCE(SUM(points), 0) FROM diva_beads_points WHERE character_id=$4 AND bead_index=$1 AND timestamp BETWEEN $2 AND $3`, bead, startTime.Add(time.Hour*-12), startTime, s.charID).Scan(&points)
			bf.WriteUint8(bead)
			bf.WriteUint32(points)
			bf.WriteUint32(points)
			startTime = startTime.Add(time.Hour * 12)
		}
	}
	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}

type UdPointTargets struct {
	Type  uint8
	Value uint64
}

func handleMsgMhfGetUdTotalPointInfo(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetUdTotalPointInfo)
	bf := byteframe.NewByteFrame()
	bf.WriteUint8(0) // No error
	targets := make([]UdPointTargets, 64)
	defaultTargets := []UdPointTargets{
		{0, 500000},
		{0, 1000000},
		{0, 2000000},
		{0, 3000000},
		{0, 4000000},
		{0, 5000000},
		{0, 6000000},
		{0, 7000000},
		{0, 8000000},
		{0, 9000000},
		{0, 10000000},
		{0, 15000000},
		{0, 20000000},
		{0, 25000000},
		{0, 30000000},
		{0, 35000000},
		{0, 40000000},
		{0, 45000000},
		{0, 50000000},
		{0, 55000000},
		{0, 60000000},
		{0, 70000000},
		{0, 80000000},
		{0, 90000000},
		{0, 100000000},
		{1, 9000000},
		{2, 30000000},
		{3, 55000000},
	}

	for i, target := range defaultTargets {
		targets[i].Type = target.Type
		targets[i].Value = target.Value
	}

	for _, target := range targets {
		bf.WriteUint64(target.Value)
	}
	for _, target := range targets {
		bf.WriteUint8(target.Type)
	}

	var totalSouls uint64
	s.server.db.QueryRow(`SELECT SUM(points) FROM diva_beads_points WHERE bead_index IS NOT NULL`).Scan(&totalSouls)
	bf.WriteUint64(totalSouls)
	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}

func handleMsgMhfGetUdSelectedColorInfo(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetUdSelectedColorInfo)
	bf := byteframe.NewByteFrame()
	bf.WriteUint8(0) // No error
	var startTime time.Time
	s.server.db.QueryRow(`SELECT start_time FROM events WHERE event_type='diva'`).Scan(&startTime)
	for x := 0; x < 8; x++ {
		var j, k, l int
		for i := 0; i < 4; i++ {
			s.server.db.QueryRow(`SELECT COALESCE(SUM(points), 0) FROM diva_beads_points WHERE bead_index=$1 AND timestamp BETWEEN $2 AND $3`, i+1, startTime, startTime.Add(time.Hour*24)).Scan(&j)
			if j > k {
				k = j
				l = i + 1
			}
		}
		startTime = startTime.Add(time.Hour * 24)
		bf.WriteUint8(uint8(l))
	}
	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}

func handleMsgMhfGetUdMonsterPoint(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetUdMonsterPoint)

	monsterPoints := []struct {
		MID    uint8
		Points uint16
	}{
		{MID: 0x01, Points: 0x3C}, // em1 Rathian
		{MID: 0x02, Points: 0x5A}, // em2 Fatalis
		{MID: 0x06, Points: 0x14}, // em6 Yian Kut-Ku
		{MID: 0x07, Points: 0x50}, // em7 Lao-Shan Lung
		{MID: 0x08, Points: 0x28}, // em8 Cephadrome
		{MID: 0x0B, Points: 0x3C}, // em11 Rathalos
		{MID: 0x0E, Points: 0x3C}, // em14 Diablos
		{MID: 0x0F, Points: 0x46}, // em15 Khezu
		{MID: 0x11, Points: 0x46}, // em17 Gravios
		{MID: 0x14, Points: 0x28}, // em20 Gypceros
		{MID: 0x15, Points: 0x3C}, // em21 Plesioth
		{MID: 0x16, Points: 0x32}, // em22 Basarios
		{MID: 0x1A, Points: 0x32}, // em26 Monoblos
		{MID: 0x1B, Points: 0x0A}, // em27 Velocidrome
		{MID: 0x1C, Points: 0x0A}, // em28 Gendrome
		{MID: 0x1F, Points: 0x0A}, // em31 Iodrome
		{MID: 0x21, Points: 0x50}, // em33 Kirin
		{MID: 0x24, Points: 0x64}, // em36 Crimson Fatalis
		{MID: 0x25, Points: 0x3C}, // em37 Pink Rathian
		{MID: 0x26, Points: 0x1E}, // em38 Blue Yian Kut-Ku
		{MID: 0x27, Points: 0x28}, // em39 Purple Gypceros
		{MID: 0x28, Points: 0x50}, // em40 Yian Garuga
		{MID: 0x29, Points: 0x5A}, // em41 Silver Rathalos
		{MID: 0x2A, Points: 0x50}, // em42 Gold Rathian
		{MID: 0x2B, Points: 0x3C}, // em43 Black Diablos
		{MID: 0x2C, Points: 0x3C}, // em44 White Monoblos
		{MID: 0x2D, Points: 0x46}, // em45 Red Khezu
		{MID: 0x2E, Points: 0x3C}, // em46 Green Plesioth
		{MID: 0x2F, Points: 0x50}, // em47 Black Gravios
		{MID: 0x30, Points: 0x1E}, // em48 Daimyo Hermitaur
		{MID: 0x31, Points: 0x3C}, // em49 Azure Rathalos
		{MID: 0x32, Points: 0x50}, // em50 Ashen Lao-Shan Lung
		{MID: 0x33, Points: 0x3C}, // em51 Blangonga
		{MID: 0x34, Points: 0x28}, // em52 Congalala
		{MID: 0x35, Points: 0x50}, // em53 Rajang
		{MID: 0x36, Points: 0x6E}, // em54 Kushala Daora
		{MID: 0x37, Points: 0x50}, // em55 Shen Gaoren
		{MID: 0x3A, Points: 0x50}, // em58 Yama Tsukami
		{MID: 0x3B, Points: 0x6E}, // em59 Chameleos
		{MID: 0x40, Points: 0x64}, // em64 Lunastra
		{MID: 0x41, Points: 0x6E}, // em65 Teostra
		{MID: 0x43, Points: 0x28}, // em67 Shogun Ceanataur
		{MID: 0x44, Points: 0x0A}, // em68 Bulldrome
		{MID: 0x47, Points: 0x6E}, // em71 White Fatalis
		{MID: 0x4A, Points: 0xFA}, // em74 Hypnocatrice
		{MID: 0x4B, Points: 0xFA}, // em75 Lavasioth
		{MID: 0x4C, Points: 0x46}, // em76 Tigrex
		{MID: 0x4D, Points: 0x64}, // em77 Akantor
		{MID: 0x4E, Points: 0xFA}, // em78 Bright Hypnoc
		{MID: 0x4F, Points: 0xFA}, // em79 Lavasioth Subspecies
		{MID: 0x50, Points: 0xFA}, // em80 Espinas
		{MID: 0x51, Points: 0xFA}, // em81 Orange Espinas
		{MID: 0x52, Points: 0xFA}, // em82 White Hypnoc
		{MID: 0x53, Points: 0xFA}, // em83 Akura Vashimu
		{MID: 0x54, Points: 0xFA}, // em84 Akura Jebia
		{MID: 0x55, Points: 0xFA}, // em85 Berukyurosu
		{MID: 0x59, Points: 0xFA}, // em89 Pariapuria
		{MID: 0x5A, Points: 0xFA}, // em90 White Espinas
		{MID: 0x5B, Points: 0xFA}, // em91 Kamu Orugaron
		{MID: 0x5C, Points: 0xFA}, // em92 Nono Orugaron
		{MID: 0x5E, Points: 0xFA}, // em94 Dyuragaua
		{MID: 0x5F, Points: 0xFA}, // em95 Doragyurosu
		{MID: 0x60, Points: 0xFA}, // em96 Gurenzeburu
		{MID: 0x63, Points: 0xFA}, // em99 Rukodiora
		{MID: 0x65, Points: 0xFA}, // em101 Gogomoa
		{MID: 0x67, Points: 0xFA}, // em103 Taikun Zamuza
		{MID: 0x68, Points: 0xFA}, // em104 Abiorugu
		{MID: 0x69, Points: 0xFA}, // em105 Kuarusepusu
		{MID: 0x6A, Points: 0xFA}, // em106 Odibatorasu
		{MID: 0x6B, Points: 0xFA}, // em107 Disufiroa
		{MID: 0x6C, Points: 0xFA}, // em108 Rebidiora
		{MID: 0x6D, Points: 0xFA}, // em109 Anorupatisu
		{MID: 0x6E, Points: 0xFA}, // em110 Hyujikiki
		{MID: 0x6F, Points: 0xFA}, // em111 Midogaron
		{MID: 0x70, Points: 0xFA}, // em112 Giaorugu
		{MID: 0x72, Points: 0xFA}, // em114 Farunokku
		{MID: 0x73, Points: 0xFA}, // em115 Pokaradon
		{MID: 0x74, Points: 0xFA}, // em116 Shantien
		{MID: 0x77, Points: 0xFA}, // em119 Goruganosu
		{MID: 0x78, Points: 0xFA}, // em120 Aruganosu
		{MID: 0x79, Points: 0xFA}, // em121 Baruragaru
		{MID: 0x7A, Points: 0xFA}, // em122 Zerureusu
		{MID: 0x7B, Points: 0xFA}, // em123 Gougarf
		{MID: 0x7D, Points: 0xFA}, // em125 Forokururu
		{MID: 0x7E, Points: 0xFA}, // em126 Meraginasu
		{MID: 0x7F, Points: 0xFA}, // em127 Diorekkusu
		{MID: 0x80, Points: 0xFA}, // em128 Garuba Daora
		{MID: 0x81, Points: 0xFA}, // em129 Inagami
		{MID: 0x82, Points: 0xFA}, // em130 Varusaburosu
		{MID: 0x83, Points: 0xFA}, // em131 Poborubarumu
		{MID: 0x8B, Points: 0xFA}, // em139 Gureadomosu
		{MID: 0x8C, Points: 0xFA}, // em140 Harudomerugu
		{MID: 0x8D, Points: 0xFA}, // em141 Toridcless
		{MID: 0x8E, Points: 0xFA}, // em142 Gasurabazura
		{MID: 0x90, Points: 0xFA}, // em144 Yama Kurai
		{MID: 0x92, Points: 0x78}, // em146 Zinogre
		{MID: 0x93, Points: 0x78}, // em147 Deviljho
		{MID: 0x94, Points: 0x78}, // em148 Brachydios
		{MID: 0x96, Points: 0xFA}, // em150 Toa Tesukatora
		{MID: 0x97, Points: 0x78}, // em151 Barioth
		{MID: 0x98, Points: 0x78}, // em152 Uragaan
		{MID: 0x99, Points: 0x78}, // em153 Stygian Zinogre
		{MID: 0x9A, Points: 0xFA}, // em154 Guanzorumu
		{MID: 0x9E, Points: 0xFA}, // em158 Voljang
		{MID: 0x9F, Points: 0x78}, // em159 Nargacuga
		{MID: 0xA0, Points: 0xFA}, // em160 Keoaruboru
		{MID: 0xA1, Points: 0xFA}, // em161 Zenaserisu
		{MID: 0xA2, Points: 0x78}, // em162 Gore Magala
		{MID: 0xA4, Points: 0x78}, // em164 Shagaru Magala
		{MID: 0xA5, Points: 0x78}, // em165 Amatsu
		{MID: 0xA6, Points: 0xFA}, // em166 Elzelion
		{MID: 0xA9, Points: 0x78}, // em169 Seregios
		{MID: 0xAA, Points: 0xFA}, // em170 Bogabadorumu
	}

	resp := byteframe.NewByteFrame()
	resp.WriteUint8(uint8(len(monsterPoints)))
	for _, mp := range monsterPoints {
		resp.WriteUint8(mp.MID)
		resp.WriteUint16(mp.Points)
	}

	doAckBufSucceed(s, pkt.AckHandle, resp.Data())
}

func handleMsgMhfGetUdDailyPresentList(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetUdDailyPresentList)
	// Temporary canned response
	data, _ := hex.DecodeString("0100001600000A5397DF00000000000000000000000000000000")
	doAckBufSucceed(s, pkt.AckHandle, data)
}

func handleMsgMhfGetUdNormaPresentList(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetUdNormaPresentList)
	// Temporary canned response
	data, _ := hex.DecodeString("0100001600000A5397DF00000000000000000000000000000000")
	doAckBufSucceed(s, pkt.AckHandle, data)
}

func handleMsgMhfAcquireUdItem(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfAcquireUdItem)
	// from gal
	// daily = 0
	// personal = 1
	// personal rank = 2
	// guild rank = 3
	// gcp = 4
	// from cat
	// treasure achievement = 5
	// personal achievement = 6
	// guild achievement = 7
	bf := byteframe.NewByteFrame()
	bf.WriteUint16(0) // NumRewards
	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}

func handleMsgMhfGetUdRanking(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetUdRanking)
	bf := byteframe.NewByteFrame()
	// Temporary
	for i := 0; i < 100; i++ {
		bf.WriteUint16(uint16(i + 1))
		bf.WriteBytes(stringsupport.PaddedString("", 25, true))
		if pkt.RankType == 1 || pkt.RankType == 3 { // "Total" type
			bf.WriteBytes(stringsupport.PaddedString("", 16, true))
		}
		bf.WriteUint32(0)
	}
	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}

func handleMsgMhfGetUdMyRanking(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetUdMyRanking)
	bf := byteframe.NewByteFrame()
	bf.WriteUint32(0) // ranking
	bf.WriteUint32(0) // rankingDupe?
	bf.WriteUint32(0) // guildPoints
	bf.WriteUint32(0) // unk
	bf.WriteUint32(0) // unkDupe?
	bf.WriteUint32(0) // guildPointsDupe?
	bf.WriteBytes(stringsupport.PaddedString("", 25, true))
	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}

type Tile struct {
	ID          uint16
	NextID      uint16
	BranchID    uint16
	QuestFile1  uint16
	QuestFile2  uint16
	QuestFile3  uint16
	BranchIndex uint8
	Type        uint8
	PointsReq   int32
	Claimed     bool
	Unk1        uint8
	Unk2        uint32
}

type InterceptionMaps struct {
	Maps     []MapData
	Branches []MapBranch
}

func (im *InterceptionMaps) Scan(val interface{}) (err error) {
	switch v := val.(type) {
	case []byte:
		err = json.Unmarshal(v, &im)
	case string:
		err = json.Unmarshal([]byte(v), &im)
	}

	return
}

func (im *InterceptionMaps) Value() (valuer driver.Value, err error) {
	return json.Marshal(im)
}

func (md *MapData) GetClaimed() uint32 {
	var claimed uint32
	for _, tile := range md.Tiles {
		if md.Points[tile.QuestFile1]-tile.PointsReq > 0 {
			tile.Claimed = true
			if tile.PointsReq > 0 {
				claimed++
			}
			md.Points[tile.QuestFile1] -= tile.PointsReq
		}
	}
	return claimed
}

func (md *MapData) TotalPoints() int32 {
	var points int32
	for i := range md.Tiles {
		if md.Tiles[i].Type > 2 {
			continue
		}
		points += md.Tiles[i].PointsReq
	}
	return points
}

func (md *MapData) Completed() bool {
	if md.Points[0] > md.TotalPoints() {
		return true
	}
	return false
}

func (im *InterceptionMaps) CurrPrevID() (uint32, uint32) {
	var currID, prevID uint32
	for i := range im.Maps {
		prevID = currID
		currID = im.Maps[i].ID
		if im.Maps[i].Points[0] < im.Maps[i].TotalPoints() {
			break
		}
	}
	return currID, prevID
}

type MapData struct {
	ID     uint32
	NextID uint32
	Points map[uint16]int32
	Tiles  []Tile
}

type MapProg struct {
	ID    uint32
	Unk   uint16
	Tiles []Tile
	Bytes *byteframe.ByteFrame
}

type MapBranch struct {
	MapIndex   uint32
	ItemType   uint8
	ItemID     uint16
	Quantity   uint16
	TileIndex1 uint16 // Sequential
	TileIndex2 uint16 // Sequential, last = 99
	ChestType  uint8
}

func handleMsgMhfGetUdGuildMapInfo(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetUdGuildMapInfo)

	guild, err := GetGuildInfoByCharacterId(s, s.charID)
	if err != nil || guild == nil {
		doAckBufSucceed(s, pkt.AckHandle, []byte{0xFF})
		return
	}
	isApplicant, _ := guild.HasApplicationForCharID(s, s.charID)
	if err != nil || isApplicant {
		doAckBufSucceed(s, pkt.AckHandle, []byte{0xFF})
		return
	}

	var interceptionMaps InterceptionMaps
	err = s.server.db.QueryRow(`SELECT interception_maps FROM guilds WHERE id=$1`, guild.ID).Scan(&interceptionMaps)
	if err != nil {
		s.server.logger.Error("Failed to load interception map data", zap.Error(err))
		doAckBufSucceed(s, pkt.AckHandle, []byte{0xFF})
		return
	}

	bf := byteframe.NewByteFrame()
	bf.WriteUint8(0) // No error
	var tilesClaimed uint32
	currentMapID, prevMapID := interceptionMaps.CurrPrevID()
	currProg := byteframe.NewByteFrame()
	prevProg := byteframe.NewByteFrame()
	bf.WriteUint8(uint8(len(interceptionMaps.Maps)))
	for _, _map := range interceptionMaps.Maps {
		bf.WriteUint32(_map.ID)
		bf.WriteUint32(_map.NextID)
		for _, tile := range _map.Tiles {
			bf.WriteUint16(tile.ID)
			bf.WriteUint16(tile.NextID)
			bf.WriteUint16(tile.BranchID)
			bf.WriteUint16(tile.QuestFile1)
			bf.WriteUint16(tile.QuestFile2)
			bf.WriteUint16(tile.QuestFile3)
			bf.WriteUint8(tile.BranchIndex)
			bf.WriteUint8(tile.Type)
			bf.WriteInt32(tile.PointsReq)

			bf.WriteUint8(tile.Unk1)
			bf.WriteUint32(tile.Unk2)
		}
		bf.WriteBytes(make([]byte, 23*(64-len(_map.Tiles)))) // Fill out 64 tiles

		if _map.Completed() && _map.ID != prevMapID {
			tilesClaimed += _map.GetClaimed()
		}

		if _map.ID == currentMapID {
			currProg.WriteUint32(_map.ID)
			currProg.WriteUint16(1)
			currProg.WriteUint8(uint8(len(_map.Tiles)))
			for _, tile := range _map.Tiles {
				if tile.Type != 1 {
					if _map.Points[tile.QuestFile1]-tile.PointsReq > 0 {
						tile.Claimed = true
						tilesClaimed++
						_map.Points[tile.QuestFile1] -= tile.PointsReq
						currProg.WriteInt32(tile.PointsReq)
					} else {
						currProg.WriteInt32(_map.Points[tile.QuestFile1])
						_map.Points[tile.QuestFile1] = 0
					}
				} else {
					currProg.WriteInt32(0)
				}
				currProg.WriteInt32(tile.PointsReq)
				currProg.WriteUint16(tile.ID)
				currProg.WriteUint16(tile.NextID)
				currProg.WriteUint16(tile.BranchID)
				currProg.WriteUint16(tile.QuestFile1)
				currProg.WriteUint16(tile.QuestFile2)
				currProg.WriteUint16(tile.QuestFile3)
				currProg.WriteUint8(tile.BranchIndex)
				currProg.WriteUint8(tile.Type)
				if tile.Claimed || tile.Type == 1 {
					currProg.WriteBool(true)
				} else {
					currProg.WriteBool(false)
				}
			}
		}
		if _map.ID == prevMapID {
			prevProg.WriteUint32(_map.ID)
			prevProg.WriteUint16(1)
			prevProg.WriteUint8(uint8(len(_map.Tiles)))
			for _, tile := range _map.Tiles {
				if tile.Type != 1 {
					if _map.Points[tile.QuestFile1]-tile.PointsReq > 0 {
						tile.Claimed = true
						tilesClaimed++
						_map.Points[tile.QuestFile1] -= tile.PointsReq
						prevProg.WriteInt32(tile.PointsReq)
					} else {
						prevProg.WriteInt32(_map.Points[tile.QuestFile1])
						_map.Points[tile.QuestFile1] = 0
					}
				} else {
					prevProg.WriteInt32(0)
				}
				prevProg.WriteInt32(tile.PointsReq)
				prevProg.WriteUint16(tile.ID)
				prevProg.WriteUint16(tile.NextID)
				prevProg.WriteUint16(tile.BranchID)
				prevProg.WriteUint16(tile.QuestFile1)
				prevProg.WriteUint16(tile.QuestFile2)
				prevProg.WriteUint16(tile.QuestFile3)
				prevProg.WriteUint8(tile.BranchIndex)
				prevProg.WriteUint8(tile.Type)
				if tile.Claimed || tile.Type == 1 {
					prevProg.WriteBool(true)
				} else {
					prevProg.WriteBool(false)
				}
			}
		}
	}

	bf.WriteUint16(uint16(len(interceptionMaps.Branches)))
	for _, branch := range interceptionMaps.Branches {
		bf.WriteUint32(branch.MapIndex)
		bf.WriteUint8(branch.ItemType)
		bf.WriteUint16(branch.ItemID)
		bf.WriteUint16(branch.Quantity)
		bf.WriteUint16(branch.TileIndex1)
		bf.WriteUint16(branch.TileIndex2)
		bf.WriteUint8(branch.ChestType)
	}

	if prevMapID > 0 {
		bf.WriteUint8(2)
	} else {
		bf.WriteUint8(1)
	}
	bf.WriteBytes(currProg.Data())
	if prevMapID > 0 {
		bf.WriteBytes(prevProg.Data())
	}
	bf.WriteUint32(tilesClaimed)
	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}

func getNeighbourTiles(tiles [][]uint16, tile uint16) []uint16 {
	var vals []uint16
	var temp []uint16
	if tile%2 == 0 {
		temp = []uint16{tile - 100, tile - 1, tile + 1, tile + 99, tile + 100, tile + 101}
	} else {
		temp = []uint16{tile - 101, tile - 100, tile - 99, tile - 1, tile + 1, tile + 100}
	}

	for _, val := range temp {
		for x := range tiles {
			for y := range tiles[x] {
				if tiles[x][y] == val {
					vals = append(vals, val)
				}
			}
		}
	}
	return vals
}

func getBranchTile(tiles [][]uint16, excluded []uint16, tile uint16) uint16 {
	neighbours := getNeighbourTiles(tiles, tile)
	var validNeighbours, validBranchTiles []uint16
	for i := range neighbours {
		if !slices.Contains(excluded, neighbours[i]) {
			// Neighbour tiles that are not in the path
			validNeighbours = append(validNeighbours, neighbours[i])
		}
	}
	if len(validNeighbours) == 0 {
		return 0
	}
	for i := range validNeighbours {
		subNeighbours := getNeighbourTiles(tiles, validNeighbours[i])
		var invalid bool
		var cleanSubNeighbours []uint16
		for _, subNeighbour := range subNeighbours {
			if subNeighbour != validNeighbours[i] && subNeighbour != tile {
				cleanSubNeighbours = append(cleanSubNeighbours, subNeighbour)
			}
		}
		for _, subNeighbour := range cleanSubNeighbours {
			if slices.Contains(excluded, subNeighbour) {
				invalid = true
				break
			}
		}
		if !invalid {
			validBranchTiles = append(validBranchTiles, validNeighbours[i])
		}
	}
	if len(validBranchTiles) == 0 {
		return 0
	}
	rand.Seed(time.Now().UnixNano())
	return validBranchTiles[rand.Intn(len(validBranchTiles))]
}

func GenerateUdGuildMaps() ([]MapData, []MapBranch) {
	tiles := make([][]uint16, 5)
	for i := range tiles {
		tiles[i] = make([]uint16, 12)
		for j := range tiles[i] {
			tiles[i][j] = uint16(((i + 1) * 100) + j + 1)
		}
	}

	var mapData []MapData
	var mapBranches []MapBranch

	for i := 0; i < 5; i++ {
		var startTile, endTile uint16
		var randTemp []uint16
		rand.Seed(time.Now().UnixNano())
		randTemp = tiles[rand.Intn(len(tiles))]
		startTile = randTemp[rand.Intn(len(randTemp))]
		for {
			rand.Seed(time.Now().UnixNano())
			randTemp = tiles[rand.Intn(len(tiles))]
			endTile = randTemp[rand.Intn(len(randTemp))]
			invalidTiles := append(getNeighbourTiles(tiles, startTile), startTile)
			if !slices.Contains(invalidTiles, endTile) {
				break
			}
		}

		var tilePath []uint16
		var iterations int
		var tooDifficult bool
		for {
			var pathFailed bool
			var evictedTiles []uint16
			tilePath = []uint16{startTile}
			for {
				var possibleTiles []uint16
				tempTiles := getNeighbourTiles(tiles, tilePath[len(tilePath)-1])
				for _, tile := range tempTiles {
					if !slices.Contains(evictedTiles, tile) {
						possibleTiles = append(possibleTiles, tile)
					}
				}
				if len(possibleTiles) == 0 {
					pathFailed = true
					break
				}
				for _, tile := range possibleTiles {
					evictedTiles = append(evictedTiles, tile)
				}
				newTile := possibleTiles[rand.Intn(len(possibleTiles))]
				tilePath = append(tilePath, newTile)
				if tilePath[len(tilePath)-1] == endTile {
					if len(tilePath) < 20 {
						pathFailed = true
					}
					break
				}
			}
			if !pathFailed {
				break
			}
			if pathFailed {
				iterations = iterations + 1
			}
			if iterations > 1000 {
				tooDifficult = true
				break
			}
		}

		if tooDifficult {
			i--
			continue
		}

		var mapTiles []Tile
		for j, tile := range tilePath {
			mapTile := Tile{}
			mapTile.ID = tile
			mapTile.BranchIndex = uint8(j + 1)
			switch j {
			case 0:
				mapTile.Type = 1
				mapTile.NextID = tilePath[j+1]
			case len(tilePath) - 1:
				mapTile.Type = 2
			default:
				mapTile.NextID = tilePath[j+1]
			}
			switch i {
			case 0:
				mapTile.PointsReq = int32(2500 + 150*(j-1))
			case 1:
				mapTile.PointsReq = int32(5500 + 600*(j-1))
			case 2:
				mapTile.PointsReq = int32(6500 + 800*(j-1))
			case 3:
				mapTile.PointsReq = int32(7500 + 1000*(j-1))
			case 4:
				mapTile.PointsReq = int32(8500 + 1000*(j-1))
			}
			if mapTile.Type == 1 {
				mapTile.PointsReq = 0
			}
			mapTiles = append(mapTiles, mapTile)
		}

		var evictedTiles []uint16
		for _, tile := range tilePath {
			evictedTiles = append(evictedTiles, tile)
		}

		var branchTiles []Tile
		for j := range mapTiles {
			if mapTiles[j].Type != 0 {
				continue
			}
			var newBranchTile uint16
			var branchIndex int
			currentBranchTile := mapTiles[j]
			for {
				newBranchTile = getBranchTile(tiles, evictedTiles, currentBranchTile.ID)
				if newBranchTile == 0 {
					if currentBranchTile != mapTiles[j] {
						branchTiles[len(branchTiles)-1].Type = 4
						branchTiles[len(branchTiles)-1].Unk1 = 1
						branchTiles[len(branchTiles)-1].Unk2 = 2
						// Make treasure more interesting, 2000GCP for now
						mapBranches = append(mapBranches, MapBranch{
							MapIndex:   uint32(i + 1),
							ItemType:   26,
							ItemID:     0,
							Quantity:   2000,
							TileIndex1: uint16(branchIndex),
							TileIndex2: 99,
							ChestType:  2,
						})
					}
					break
				} else {
					if currentBranchTile.ID == mapTiles[j].ID {
						mapTiles[j].BranchID = newBranchTile
						mapTiles[j].Type = 3
					} else {
						branchTiles[len(branchTiles)-1].NextID = newBranchTile
					}
					branchIndex++
					newTile := Tile{
						ID:          newBranchTile,
						QuestFile1:  uint16(j%5 + 58079),
						BranchIndex: uint8(branchIndex),
						Type:        0,
						PointsReq:   100,
						Unk1:        0,
						Unk2:        0,
					}
					branchTiles = append(branchTiles, newTile)
					for _, k := range getNeighbourTiles(tiles, currentBranchTile.ID) {
						evictedTiles = append(evictedTiles, k)
					}
					currentBranchTile = newTile
				}
			}
		}
		for j := range branchTiles {
			mapTiles = append(mapTiles, branchTiles[j])
		}
		if i >= 4 {
			mapData = append(mapData, MapData{uint32(i + 1), 4, make(map[uint16]int32), mapTiles})
		} else {
			mapData = append(mapData, MapData{uint32(i + 1), uint32(i + 2), make(map[uint16]int32), mapTiles})
		}
	}
	return mapData, mapBranches
}

func handleMsgMhfGenerateUdGuildMap(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGenerateUdGuildMap)
	guild, err := GetGuildInfoByCharacterId(s, s.charID)
	if err != nil || guild == nil {
		doAckBufSucceed(s, pkt.AckHandle, []byte{0xFF})
		return
	}
	isApplicant, _ := guild.HasApplicationForCharID(s, s.charID)
	if err != nil || isApplicant {
		doAckBufSucceed(s, pkt.AckHandle, []byte{0xFF})
		return
	}
	interceptionMaps := &InterceptionMaps{}
	interceptionMaps.Maps, interceptionMaps.Branches = GenerateUdGuildMaps()
	_, err = s.server.db.Exec(`UPDATE guilds SET interception_maps = $1 WHERE id = $2`, interceptionMaps, guild.ID)
	if err != nil {
		s.server.logger.Debug("err", zap.Error(err))
	}
	doAckBufSucceed(s, pkt.AckHandle, []byte{0})
}
