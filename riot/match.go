package riot

type (
	// MatchSummoner represents a user in a match
	MatchSummoner struct {
		AccountID        string `json:"accountId"`
		CurrentAccountID string `json:"currentAccountId,omitmepty"`
	}

	// Match represents a match overview
	Match struct {
		Lane       string `json:"lane"`
		GameID     int64  `json:"gameId"`
		Champion   int    `json:"champion"`
		PlatformID string `json:"platformId"`
		Timestamp  int64  `json:"timestamp"`
		Queue      int    `json:"queue"`
		Role       string `json:"role"`
		Season     int    `json:"season"`
	}

	// MatchDetails represents detailed statistics about a match
	MatchDetails struct {
		SeasonID     int    `json:"seasonId"`
		QueueID      int    `json:"queueId"`
		GameID       int64  `json:"gameId"`
		GameDuration int    `json:"gameDuration"`
		GameCreation int64  `json:"gameCreation"`
		GameVersion  string `json:"gameVersion"`
		PlatformID   string `json:"platformId"`
		GameMode     string `json:"gameMode"`
		MapID        int    `json:"mapId"`
		GameType     string `json:"gameType"`

		ParticipantIdentities []struct {
			Player        MatchSummoner `json:"player"`
			ParticipantID int           `json:"participantId"`
		} `json:"participantIdentities"`

		Teams []struct {
			TeamID int    `json:"teamId"`
			Win    string `json:"win"`
			Bans   []struct {
				PickTurn   int `json:"pickTurn"`
				ChampionID int `json:"championId"`
			} `json:"bans"`
			FirstDragon     bool `json:"firstDragon"`
			FirstInhibitor  bool `json:"firstInhibitor"`
			FirstRiftHerald bool `json:"firstRiftHerald"`
			FirstBaron      bool `json:"firstBaron"`
			BaronKills      int  `json:"baronKills"`
			RiftHeraldKills int  `json:"riftHeraldKills"`
			FirstBlood      bool `json:"firstBlood"`
			FirstTower      bool `json:"firstTower"`
			VilemawKills    int  `json:"vilemawKills"`
			InhibitorKills  int  `json:"inhibitorKills"`
			TowerKills      int  `json:"towerKills"`
			DragonKills     int  `json:"dragonKills"`
		} `json:"teams"`

		Participants []struct {
			ParticipantID             int    `json:"participantId"`
			TeamID                    int    `json:"teamId"`
			ChampionID                int    `json:"championId"`
			HighestAchievedSeasonTier string `json:"highestAchievedSeasonTier,omitempty"`
			Spell1ID                  int    `json:"spell1Id"`
			Spell2ID                  int    `json:"spell2Id"`
			Stats                     struct {
				ParticipantID int  `json:"participantId"`
				Win           bool `json:"win"`

				ChampLevel int `json:"champLevel"`
				Kills      int `json:"kills"`
				Deaths     int `json:"deaths"`
				Assists    int `json:"assists"`

				TurretKills                 int `json:"turretKills"`
				TotalDamageTaken            int `json:"totalDamageTaken"`
				TotalDamageDealtToChampions int `json:"totalDamageDealtToChampions"`
				TotalMinionsKilled          int `json:"totalMinionsKilled"`
				GoldEarned                  int `json:"goldEarned"`
				VisionScore                 int `json:"visionScore"`

				Item0 int `json:"item0"`
				Item1 int `json:"item1"`
				Item2 int `json:"item2"`
				Item3 int `json:"item3"`
				Item4 int `json:"item4"`
				Item6 int `json:"item6"`
				Item5 int `json:"item5"`

				NeutralMinionsKilledTeamJungle  int  `json:"neutralMinionsKilledTeamJungle"`
				MagicDamageDealtToChampions     int  `json:"magicDamageDealtToChampions"`
				NeutralMinionsKilled            int  `json:"neutralMinionsKilled"`
				DamageDealtToTurrets            int  `json:"damageDealtToTurrets"`
				PhysicalDamageDealtToChampions  int  `json:"physicalDamageDealtToChampions"`
				DamageDealtToObjectives         int  `json:"damageDealtToObjectives"`
				TotalUnitsHealed                int  `json:"totalUnitsHealed"`
				WardsKilled                     int  `json:"wardsKilled"`
				LargestCriticalStrike           int  `json:"largestCriticalStrike"`
				LargestKillingSpree             int  `json:"largestKillingSpree"`
				QuadraKills                     int  `json:"quadraKills"`
				MagicDamageDealt                int  `json:"magicDamageDealt"`
				FirstBloodAssist                bool `json:"firstBloodAssist"`
				DamageSelfMitigated             int  `json:"damageSelfMitigated"`
				MagicalDamageTaken              int  `json:"magicalDamageTaken"`
				FirstInhibitorKill              bool `json:"firstInhibitorKill"`
				TrueDamageTaken                 int  `json:"trueDamageTaken"`
				GoldSpent                       int  `json:"goldSpent"`
				TrueDamageDealt                 int  `json:"trueDamageDealt"`
				PhysicalDamageDealt             int  `json:"physicalDamageDealt"`
				SightWardsBoughtInGame          int  `json:"sightWardsBoughtInGame"`
				PhysicalDamageTaken             int  `json:"physicalDamageTaken"`
				NeutralMinionsKilledEnemyJungle int  `json:"neutralMinionsKilledEnemyJungle"`
				WardsPlaced                     int  `json:"wardsPlaced"`
				PerkPrimaryStyle                int  `json:"perkPrimaryStyle"`
				PerkSubStyle                    int  `json:"perkSubStyle"`
				FirstBloodKill                  bool `json:"firstBloodKill"`
				TrueDamageDealtToChampions      int  `json:"trueDamageDealtToChampions"`
				KillingSprees                   int  `json:"killingSprees"`
				UnrealKills                     int  `json:"unrealKills"`
				FirstTowerAssist                bool `json:"firstTowerAssist"`
				FirstTowerKill                  bool `json:"firstTowerKill"`
				DoubleKills                     int  `json:"doubleKills"`
				InhibitorKills                  int  `json:"inhibitorKills"`
				FirstInhibitorAssist            bool `json:"firstInhibitorAssist"`
				VisionWardsBoughtInGame         int  `json:"visionWardsBoughtInGame"`
				PentaKills                      int  `json:"pentaKills"`
				TotalHeal                       int  `json:"totalHeal"`
				TimeCCingOthers                 int  `json:"timeCCingOthers"`
				TripleKills                     int  `json:"tripleKills"`
				LargestMultiKill                int  `json:"largestMultiKill"`
				TotalTimeCrowdControlDealt      int  `json:"totalTimeCrowdControlDealt"`
				LongestTimeSpentLiving          int  `json:"longestTimeSpentLiving"`
			} `json:"stats"`
		} `json:"participants"`
	}
)
