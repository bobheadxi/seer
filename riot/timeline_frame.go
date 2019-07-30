package riot

import (
	"encoding/json"
	"sort"
)

// Position denotes a set of coordinates
type Position struct {
	Y int `json:"y"`
	X int `json:"x"`
}

// FrameEventType is the type of a frame event
type FrameEventType string

const (
	// FrameEventChampionKill = CHAMPION_KILL
	FrameEventChampionKill FrameEventType = "CHAMPION_KILL"
	// FrameEventBuildingKill = BUILDING_KILL
	FrameEventBuildingKill FrameEventType = "BUILDING_KILL"
	// FrameEventMonsterKill = ELITE_MONSTER_KILL
	FrameEventMonsterKill FrameEventType = "ELITE_MONSTER_KILL"

	// FrameEventWardPlaced = WARD_PLACED
	FrameEventWardPlaced FrameEventType = "WARD_PLACED"
	// FrameEventWardKill = WARD_KILL
	FrameEventWardKill FrameEventType = "WARD_KILL"

	// FrameEventItemPurchased = ITEM_PURCHASED
	FrameEventItemPurchased FrameEventType = "ITEM_PURCHASED"
	// FrameEventItemSold = ITEM_SOLD
	FrameEventItemSold FrameEventType = "ITEM_SOLD"
	// FrameEventItemDestroyed = ITEM_DESTROYED
	FrameEventItemDestroyed FrameEventType = "ITEM_DESTROYED"
	// FrameEventItemUndo = ITEM_UNDO
	FrameEventItemUndo FrameEventType = "ITEM_UNDO"

	// FrameEventSkillUp = SKILL_LEVEL_UP
	FrameEventSkillUp FrameEventType = "SKILL_LEVEL_UP"
)

type (
	// FrameEvent denotes a game event that occured in a frame
	FrameEvent struct {
		Type      FrameEventType `json:"type"`
		Timestamp int            `json:"timestamp"`

		// data will populate one of the following
		*FrameEventChampionKillData
		*FrameEventBuildingKillData
		*FrameEventMonsterKillData
		*FrameEventWardData
		*FrameEventSkillData
		*FrameEventItemData
	}

	// FrameEventChampionKillData contains data for CHAMPION_KILL
	FrameEventChampionKillData struct {
		KillerID     int   `json:"killerId"`
		AssistingIDs []int `json:"assistingParticipantIds"`

		VictimID int `json:"victimId"`

		Position *Position `json:"position"`
	}

	// FrameEventBuildingKillData contains data for BUILDING_KILL
	FrameEventBuildingKillData struct {
		KillerID     int   `json:"killerId"`
		AssistingIDs []int `json:"assistingParticipantIds"`

		TeamID       int     `json:"teamId"`
		BuildingType string  `json:"buildingType"`
		TowerType    *string `json:"towerType"`
		LaneType     *string `json:"laneType"`

		Position *Position `json:"position"`
	}

	// FrameEventMonsterKillData contains data for ELITE_MONSTER_KILL
	FrameEventMonsterKillData struct {
		KillerID    int    `json:"killerId"`
		MonsterType string `json:"monsterType"`

		Position *Position `json:"position"`
	}

	// FrameEventWardData contains data for WARD_PLACED, WARD_KILL
	FrameEventWardData struct {
		CreatorID *int    `json:"creatorId"`
		KillerID  *int    `json:"killerId"`
		WardType  *string `json:"wardType"`

		Position *Position `json:"position"`
	}

	// FrameEventSkillData contains data for SKILL_LEVEL_UP
	FrameEventSkillData struct {
		ParticipantID int    `json:"participantId"`
		SkillSlot     int    `json:"skillSlot"`
		LevelUpType   string `json:"levelUpType"`
	}

	// FrameEventItemData contains data for ITEM_PURCHASED, ITEM_SOLD, ITEM_DESTROYED, ITEM_UNDO
	FrameEventItemData struct {
		ParticipantID int `json:"participantId"`
		ItemID        int `json:"itemId"`
	}
)

// ParticipantFrame records the state of a participant at a point in time
type ParticipantFrame struct {
	ParticipantID       int      `json:"participantId"`
	Level               int      `json:"level"`
	Xp                  int      `json:"xp"`
	TotalGold           int      `json:"totalGold"`
	CurrentGold         int      `json:"currentGold"`
	MinionsKilled       int      `json:"minionsKilled"`
	JungleMinionsKilled int      `json:"jungleMinionsKilled"`
	Position            Position `json:"position"`
}

// MarshalJSON is for BigQuery, where we can't have column names starting with
// numbers, so pop these into an array for UNNEST
func (f *ParticipantFrames) MarshalJSON() ([]byte, error) {
	if f == nil {
		return json.Marshal(participantFrames{})
	}
	parts := participantFrames{
		f.Participant1,
		f.Participant2,
		f.Participant3,
		f.Participant4,
		f.Participant5,
		f.Participant6,
		f.Participant7,
		f.Participant8,
		f.Participant9,
		f.Participant10,
	}
	sort.Sort(parts)
	return json.Marshal(parts)
}

// participantFrames sorts by participant ID
type participantFrames []ParticipantFrame

func (f participantFrames) Len() int           { return len(f) }
func (f participantFrames) Less(i, j int) bool { return f[i].ParticipantID < f[j].ParticipantID }
func (f participantFrames) Swap(i, j int) {
	tmp := f[i]
	f[i] = f[j]
	f[j] = tmp
}
