package channelserver

func getLangStrings(s *Server) map[string]string {
	strings := make(map[string]string)
	switch s.erupeConfig.Language {
	case "jp":
		strings["language"] = "日本語"
		strings["cafeReset"] = "%d/%dにリセット"

		strings["prayerBead1Name"] = "暴風の祈珠"
		strings["prayerBead1Description"] = "ーあらしまかぜのきじゅー\n暴風とは猛る思い。\n聞く者に勇気を与える。"
		strings["prayerBead3Name"] = "断力の祈珠"
		strings["prayerBead3Description"] = "ーだんりきのきじゅー\n断力とは断ち切る思い。\n聴く者に新たな利からを授ける。"
		strings["prayerBead4Name"] = "風韻の祈珠"
		strings["prayerBead4Description"] = "ーふういんのきじゅー\n風韻とは歌姫の艶。\n時々で異なる趣を醸し出す。"
		strings["prayerBead8Name"] = "斬刃の祈珠"
		strings["prayerBead8Description"] = "ーざんばのきじゅー\n斬刃とはすべてを切り裂く力。\n集めるほどに声の透明感は増す。"
		strings["prayerBead9Name"] = "打明の祈珠"
		strings["prayerBead9Description"] = "ーうちあかりのきじゅー\n打明とは熱い力。\n聴く者に活力を与える。"
		strings["prayerBead10Name"] = "弾起の祈珠"
		strings["prayerBead10Description"] = "ーたまおこしのきじゅー\n弾起とは悠遠の記憶。\n聴く者に更なる力を授ける。"
		strings["prayerBead11Name"] = "変続の祈珠"
		strings["prayerBead11Description"] = "ーへんぞくのきじゅー\n変続とは永久の言葉。\n聴く者に継続力を授ける。"
		strings["prayerBead14Name"] = "万雷の祈珠"
		strings["prayerBead14Description"] = "ーばんらいのきじゅー\n万雷とは歌姫に集う民の意識。\n歌姫の声を伝播させる。"
		strings["prayerBead15Name"] = "不動の祈珠"
		strings["prayerBead15Description"] = "ーうごかずのきじゅー\n不動とは圧力。聞く者に圧倒する力を与える。"
		strings["prayerBead17Name"] = "結集の祈珠"
		strings["prayerBead17Description"] = "ーけっしゅうのきじゅー\n結集とは確固たる信頼。\n集めるほどに狩人たちの精神力となる。"
		strings["prayerBead18Name"] = "歌護の祈珠"
		strings["prayerBead18Description"] = "ーうたまもりのきじゅー\n歌護とは歌姫の護り。\n集めるほどに狩人たちの支えとなる。"
		strings["prayerBead19Name"] = "強撃の祈珠"
		strings["prayerBead19Description"] = "ーきょうげきのきじゅー\n強撃とは強い声色。\n聞く者の力を研ぎ澄ます。"
		strings["prayerBead20Name"] = "封火の祈珠"
		strings["prayerBead20Description"] = "ーふうかのきじゅー"
		strings["prayerBead21Name"] = "封水の祈珠"
		strings["prayerBead21Description"] = "ーふうすいのきじゅー"
		strings["prayerBead22Name"] = "封氷の祈珠"
		strings["prayerBead22Description"] = "ーふうひょうのきじゅー"
		strings["prayerBead23Name"] = "封龍の祈珠"
		strings["prayerBead23Description"] = "ーふうりゅうのきじゅー"
		strings["prayerBead24Name"] = "封雷の祈珠"
		strings["prayerBead24Description"] = "ーふうらいのきじゅー"
		strings["prayerBead25Name"] = "封属の祈珠"
		strings["prayerBead25Description"] = "ーふうぞくのきじゅー"

		strings["commandDisabled"] = "%sのコマンドは無効です"
		strings["commandReload"] = "リロードします"
		strings["commandKqfGet"] = "現在のキークエストフラグ：%x"
		strings["commandKqfSetError"] = "キークエコマンドエラー　例：%s set xxxxxxxxxxxxxxxx"
		strings["commandKqfSetSuccess"] = "キークエストのフラグが更新されました。ワールド／ランドを移動してください"
		strings["commandKqfVersion"] = "This command is disabled prior to MHFG10"
		strings["commandRightsError"] = "コース更新コマンドエラー　例：%s x"
		strings["commandRightsSuccess"] = "コース情報を更新しました：%d"
		strings["commandCourseError"] = "コース確認コマンドエラー　例：%s <name>"
		strings["commandCourseDisabled"] = "%sコースは無効です"
		strings["commandCourseEnabled"] = "%sコースは有効です"
		strings["commandCourseLocked"] = "%sコースはロックされています"
		strings["commandTeleportError"] = "テレポートコマンドエラー　構文：%s x y"
		strings["commandTeleportSuccess"] = "%d %dにテレポート"
		strings["commandPSNError"] = "PSN連携コマンドエラー　例：%s <psn id>"
		strings["commandPSNSuccess"] = "PSN「%s」が連携されています"
		strings["commandPSNExists"] = "PSNは既存のユーザに接続されています"

		strings["commandRaviNoCommand"] = "ラヴィコマンドが指定されていません"
		strings["commandRaviStartSuccess"] = "大討伐を開始します"
		strings["commandRaviStartError"] = "大討伐は既に開催されています"
		strings["commandRaviMultiplier"] = "ラヴィダメージ倍率：ｘ%.2f"
		strings["commandRaviResSuccess"] = "復活支援を実行します"
		strings["commandRaviResError"] = "復活支援は実行されませんでした"
		strings["commandRaviSedSuccess"] = "鎮静支援を実行します"
		strings["commandRaviRequest"] = "鎮静支援を要請します"
		strings["commandRaviError"] = "ラヴィコマンドが認識されません"
		strings["commandRaviNoPlayers"] = "誰も大討伐に参加していません"
		strings["commandRaviVersion"] = "This command is disabled outside of MHFZZ"
		strings["commandRoadError"] = "Error in command. Format: %s <road type>"
		strings["commandRoadSuccess"] = "Road is updated to: %s"
		strings["commandRoadCurrent"] = "Road is currently set to: %s"

		strings["ravienteBerserk"] = "<大討伐：猛狂期>が開催されました！"
		strings["ravienteExtreme"] = "<大討伐：猛狂期【極】>が開催されました！"
		strings["ravienteExtremeLimited"] = "<大討伐：猛狂期【極】(制限付)>が開催されました！"
		strings["ravienteBerserkSmall"] = "<大討伐：猛狂期(小数)>が開催されました！"

		strings["guildInviteName"] = "猟団勧誘のご案内"
		strings["guildInvite"] = "猟団「%s」からの勧誘通知です。\n「勧誘に返答」より、返答を行ってください。"

		strings["guildInviteSuccessName"] = "成功"
		strings["guildInviteSuccess"] = "あなたは「%s」に参加できました。"

		strings["guildInviteAcceptedName"] = "承諾されました"
		strings["guildInviteAccepted"] = "招待した狩人が「%s」への招待を承諾しました。"

		strings["guildInviteRejectName"] = "却下しました"
		strings["guildInviteReject"] = "あなたは「%s」への参加を却下しました。"

		strings["guildInviteDeclinedName"] = "辞退しました"
		strings["guildInviteDeclined"] = "招待した狩人が「%s」への招待を辞退しました。"
	default:
		strings["language"] = "English"
		strings["cafeReset"] = "Resets on %d/%d"

		strings["prayerBead1Name"] = "Bead of Storms"
		strings["prayerBead1Description"] = "ーあらしまかぜのきじゅー\n暴風とは猛る思い。\n聞く者に勇気を与える。"
		strings["prayerBead3Name"] = "Bead of Severing"
		strings["prayerBead3Description"] = "All damage types can sever tails\nPower to sever, inspire with might.\nEmpower those who hear, in new light."
		strings["prayerBead4Name"] = "Bead of Vitality"
		strings["prayerBead4Description"] = "Increased red health recovery speed\nDiva's allure, a soothing balm.\nRenews one's vigor, with vitality and calm."
		strings["prayerBead8Name"] = "Bead of Slashing"
		strings["prayerBead8Description"] = "Damage up for slashing weapons\nWith every slash, its voice rings out.\nGrowing ever sharper, without a doubt."
		strings["prayerBead9Name"] = "Bead of Striking"
		strings["prayerBead9Description"] = "Damage up for striking weapons\nWith every blow, you strike with force.\nLet the power guide your course."
		strings["prayerBead10Name"] = "Bead of Firing"
		strings["prayerBead10Description"] = "Damage up for shooting weapons\nA memory of might, empowering those who hear.\nBullet and body, soaring without fear."
		strings["prayerBead11Name"] = "Bead of Tenacity"
		strings["prayerBead11Description"] = "ーへんぞくのきじゅー\n変続とは永久の言葉。\n聴く者に継続力を授ける。"
		strings["prayerBead14Name"] = "Bead of Elements"
		strings["prayerBead14Description"] = "ーばんらいのきじゅー\n万雷とは歌姫に集う民の意識。\n歌姫の声を伝播させる。"
		strings["prayerBead15Name"] = "Bead of Restraint"
		strings["prayerBead15Description"] = "ーうごかずのきじゅー\n不動とは圧力。聞く者に圧倒する力を与える。"
		strings["prayerBead17Name"] = "Bead of Unity"
		strings["prayerBead17Description"] = "ーけっしゅうのきじゅー\n結集とは確固たる信頼。\n集めるほどに狩人たちの精神力となる。"
		strings["prayerBead18Name"] = "Bead of Warding"
		strings["prayerBead18Description"] = "ーうたまもりのきじゅー\n歌護とは歌姫の護り。\n集めるほどに狩人たちの支えとなる。"
		strings["prayerBead19Name"] = "Bead of Fury"
		strings["prayerBead19Description"] = "ーきょうげきのきじゅー\n強撃とは強い声色。\n聞く者の力を研ぎ澄ます。"
		strings["prayerBead20Name"] = "Bead of Fireproof"
		strings["prayerBead20Description"] = "ーふうかのきじゅー"
		strings["prayerBead21Name"] = "Bead of Waterproof"
		strings["prayerBead21Description"] = "ーふうすいのきじゅー"
		strings["prayerBead22Name"] = "Bead of Iceproof"
		strings["prayerBead22Description"] = "ーふうひょうのきじゅー"
		strings["prayerBead23Name"] = "Bead of Dragonproof"
		strings["prayerBead23Description"] = "ーふうりゅうのきじゅー"
		strings["prayerBead24Name"] = "Bead of Thunderproof"
		strings["prayerBead24Description"] = "ーふうらいのきじゅー"
		strings["prayerBead25Name"] = "Bead of Immunity"
		strings["prayerBead25Description"] = "ーふうぞくのきじゅー"

		strings["commandDisabled"] = "%s command is disabled"
		strings["commandReload"] = "Reloading players..."
		strings["commandKqfGet"] = "KQF: %x"
		strings["commandKqfSetError"] = "Error in command. Format: %s set xxxxxxxxxxxxxxxx"
		strings["commandKqfSetSuccess"] = "KQF set, please switch Land/World"
		strings["commandKqfVersion"] = "This command is disabled prior to MHFG10"
		strings["commandRightsError"] = "Error in command. Format: %s x"
		strings["commandRightsSuccess"] = "Set rights integer: %d"
		strings["commandCourseError"] = "Error in command. Format: %s <name>"
		strings["commandCourseDisabled"] = "%s Course disabled"
		strings["commandCourseEnabled"] = "%s Course enabled"
		strings["commandCourseLocked"] = "%s Course is locked"
		strings["commandTeleportError"] = "Error in command. Format: %s x y"
		strings["commandTeleportSuccess"] = "Teleporting to %d %d"
		strings["commandPSNError"] = "Error in command. Format: %s <psn id>"
		strings["commandPSNSuccess"] = "Connected PSN ID: %s"
		strings["commandPSNExists"] = "PSN ID is connected to another account!"

		strings["commandRaviNoCommand"] = "No Raviente command specified!"
		strings["commandRaviStartSuccess"] = "The Great Slaying will begin in a moment"
		strings["commandRaviStartError"] = "The Great Slaying has already begun!"
		strings["commandRaviMultiplier"] = "Raviente multiplier is currently %.2fx"
		strings["commandRaviResSuccess"] = "Sending resurrection support!"
		strings["commandRaviResError"] = "Resurrection support has not been requested!"
		strings["commandRaviSedSuccess"] = "Sending sedation support if requested!"
		strings["commandRaviRequest"] = "Requesting sedation support!"
		strings["commandRaviError"] = "Raviente command not recognised!"
		strings["commandRaviNoPlayers"] = "No one has joined the Great Slaying!"
		strings["commandRaviVersion"] = "This command is disabled outside of MHFZZ"
		strings["commandRoadError"] = "Error in command. Format: %s <road type>"
		strings["commandRoadSuccess"] = "Road is updated to: %s"
		strings["commandRoadCurrent"] = "Road is currently set to: %s"

		strings["ravienteBerserk"] = "<Great Slaying: Berserk> is being held!"
		strings["ravienteExtreme"] = "<Great Slaying: Extreme> is being held!"
		strings["ravienteExtremeLimited"] = "<Great Slaying: Extreme (Limited)> is being held!"
		strings["ravienteBerserkSmall"] = "<Great Slaying: Berserk (Small)> is being held!"

		strings["guildInviteName"] = "Invitation!"
		strings["guildInvite"] = "You have been invited to join\n「%s」\nDo you want to accept?"

		strings["guildInviteSuccessName"] = "Success!"
		strings["guildInviteSuccess"] = "You have successfully joined\n「%s」."

		strings["guildInviteAcceptedName"] = "Accepted"
		strings["guildInviteAccepted"] = "The recipient accepted your invitation to join\n「%s」."

		strings["guildInviteRejectName"] = "Rejected"
		strings["guildInviteReject"] = "You rejected the invitation to join\n「%s」."

		strings["guildInviteDeclinedName"] = "Declined"
		strings["guildInviteDeclined"] = "The recipient declined your invitation to join\n「%s」."
	}
	return strings
}
