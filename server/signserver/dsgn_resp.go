package signserver

import (
	"erupe-ce/common/byteframe"
	ps "erupe-ce/common/pascalstring"
	"erupe-ce/common/stringsupport"
	_config "erupe-ce/config"
	"erupe-ce/server/channelserver"
	"fmt"
	"go.uber.org/zap"
	"strings"
	"time"
)

func (s *Session) makeSignResponse(uid uint32) []byte {
	// Get the characters from the DB.
	chars, err := s.server.getCharactersForUser(uid)
	if len(chars) == 0 && uid != 0 {
		err = s.server.newUserChara(uid)
		if err == nil {
			chars, err = s.server.getCharactersForUser(uid)
		}
	}
	if err != nil {
		s.logger.Warn("Error getting characters from DB", zap.Error(err))
	}

	bf := byteframe.NewByteFrame()
	var tokenID uint32
	var sessToken string
	if uid == 0 && s.psn != "" {
		tokenID, sessToken, err = s.server.registerPsnToken(s.psn)
	} else {
		tokenID, sessToken, err = s.server.registerUidToken(uid)
	}
	if err != nil {
		bf.WriteUint8(uint8(SIGN_EABORT))
		return bf.Data()
	}

	bf.WriteUint8(uint8(SIGN_SUCCESS)) // resp_code
	if (s.server.erupeConfig.PatchServerManifest != "" && s.server.erupeConfig.PatchServerFile != "") || s.client == PS3 {
		bf.WriteUint8(2)
	} else {
		bf.WriteUint8(0)
	}
	bf.WriteUint8(1) // entrance server count
	bf.WriteUint8(uint8(len(chars)))
	bf.WriteUint32(tokenID)
	bf.WriteBytes([]byte(sessToken))
	bf.WriteUint32(uint32(channelserver.TimeAdjusted().Unix()))
	if s.client == PS3 {
		ps.Uint8(bf, fmt.Sprintf(`ps3-%s.zerulight.cc`, s.server.erupeConfig.Language), false)
		ps.Uint8(bf, fmt.Sprintf(`ps3-%s.zerulight.cc`, s.server.erupeConfig.Language), false)
	} else {
		if s.server.erupeConfig.PatchServerManifest != "" && s.server.erupeConfig.PatchServerFile != "" {
			ps.Uint8(bf, s.server.erupeConfig.PatchServerManifest, false)
			ps.Uint8(bf, s.server.erupeConfig.PatchServerFile, false)
		}
	}
	if strings.Split(s.rawConn.RemoteAddr().String(), ":")[0] == "127.0.0.1" {
		ps.Uint8(bf, fmt.Sprintf("127.0.0.1:%d", s.server.erupeConfig.Entrance.Port), false)
	} else {
		ps.Uint8(bf, fmt.Sprintf("%s:%d", s.server.erupeConfig.Host, s.server.erupeConfig.Entrance.Port), false)
	}

	lastPlayed := uint32(0)
	for _, char := range chars {
		if lastPlayed == 0 {
			lastPlayed = char.ID
		}
		bf.WriteUint32(char.ID)

		// Exp, HR[x] is split by 0, 1, 30, 50, 99, 299, 998, 999
		if s.server.erupeConfig.DevMode && s.server.erupeConfig.DevModeOptions.MaxLauncherHR {
			bf.WriteUint16(999)
		} else {
			bf.WriteUint16(char.HRP)
		}

		bf.WriteUint16(char.WeaponType)                                          // Weapon, 0-13.
		bf.WriteUint32(char.LastLogin)                                           // Last login date, unix timestamp in seconds.
		bf.WriteBool(char.IsFemale)                                              // Sex, 0=male, 1=female.
		bf.WriteBool(char.IsNewCharacter)                                        // Is new character, 1 replaces character name with ?????.
		bf.WriteUint8(0)                                                         // Old GR
		bf.WriteBool(true)                                                       // Use uint16 GR, no reason not to
		bf.WriteBytes(stringsupport.PaddedString(char.Name, 16, true))           // Character name
		bf.WriteBytes(stringsupport.PaddedString(char.UnkDescString, 32, false)) // unk str
		if s.server.erupeConfig.RealClientMode >= _config.G7 {
			bf.WriteUint16(char.GR)
			bf.WriteUint8(0) // Unk
			bf.WriteUint8(0) // Unk
		}
	}

	friends := s.server.getFriendsForCharacters(chars)
	if len(friends) == 0 {
		bf.WriteUint8(0)
	} else {
		if len(friends) > 255 {
			bf.WriteUint8(255)
			bf.WriteUint16(uint16(len(friends)))
		} else {
			bf.WriteUint8(uint8(len(friends)))
		}
		for _, friend := range friends {
			bf.WriteUint32(friend.CID)
			bf.WriteUint32(friend.ID)
			ps.Uint8(bf, friend.Name, true)
		}
	}

	guildmates := s.server.getGuildmatesForCharacters(chars)
	if len(guildmates) == 0 {
		bf.WriteUint8(0)
	} else {
		if len(guildmates) > 255 {
			bf.WriteUint8(255)
			bf.WriteUint16(uint16(len(guildmates)))
		} else {
			bf.WriteUint8(uint8(len(guildmates)))
		}
		for _, guildmate := range guildmates {
			bf.WriteUint32(guildmate.CID)
			bf.WriteUint32(guildmate.ID)
			ps.Uint8(bf, guildmate.Name, true)
		}
	}

	if s.server.erupeConfig.HideLoginNotice {
		bf.WriteBool(false)
	} else {
		bf.WriteBool(true)
		bf.WriteUint8(0)
		bf.WriteUint8(0)
		ps.Uint16(bf, strings.Join(s.server.erupeConfig.LoginNotices[:], "<PAGE>"), true)
	}

	bf.WriteUint32(s.server.getLastCID(uid))
	bf.WriteUint32(s.server.getUserRights(uid))
	ps.Uint16(bf, "", false) // filters
	if s.client == VITA || s.client == PS3 {
		var psnUser string
		s.server.db.QueryRow("SELECT psn_id FROM users WHERE id = $1", uid).Scan(&psnUser)
		bf.WriteBytes(stringsupport.PaddedString(psnUser, 20, true))
	}
	bf.WriteUint16(0xCA10)
	bf.WriteUint16(0x4E20)
	ps.Uint16(bf, "", false) // unk key
	bf.WriteUint8(0x00)
	bf.WriteUint16(0xCA11)
	bf.WriteUint16(0x0001)
	bf.WriteUint16(0x4E20)
	ps.Uint16(bf, "", false) // unk ipv4
	bf.WriteUint32(uint32(s.server.getReturnExpiry(uid).Unix()))
	bf.WriteUint32(0)

	tickets := []uint32{
		s.server.erupeConfig.GameplayOptions.MezfesSoloTickets,
		s.server.erupeConfig.GameplayOptions.MezfesGroupTickets,
	}
	stalls := []uint8{
		10, 3, 6, 9, 4, 8, 5, 7,
	}
	if s.server.erupeConfig.GameplayOptions.MezFesSwitchMinigame {
		stalls[4] = 2
	}

	// We can just use the start timestamp as the event ID
	bf.WriteUint32(uint32(channelserver.TimeWeekStart().Unix()))
	// Start time
	bf.WriteUint32(uint32(channelserver.TimeWeekNext().Add(-time.Duration(s.server.erupeConfig.GameplayOptions.MezFesDuration) * time.Second).Unix()))
	// End time
	bf.WriteUint32(uint32(channelserver.TimeWeekNext().Unix()))
	bf.WriteUint8(uint8(len(tickets)))
	for i := range tickets {
		bf.WriteUint32(tickets[i])
	}
	bf.WriteUint8(uint8(len(stalls)))
	for i := range stalls {
		bf.WriteUint8(stalls[i])
	}
	return bf.Data()
}
