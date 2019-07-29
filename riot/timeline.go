package riot

import "encoding/json"

type (
	// MatchTimeline is the set of events for a match
	MatchTimeline struct {
		FrameInterval int             `json:"frameInterval"`
		Frames        []TimelineFrame `json:"frames"`

		// GameID is NOT provided by Riot API, must be set - mostly for use with
		// BigQuery
		GameID int64 `json:"gameId"`
	}

	// TimelineFrame records the state of a game at a point in time
	TimelineFrame struct {
		Timestamp         int                `json:"timestamp"`
		ParticipantFrames *ParticipantFrames `json:"participantFrames"`
		// Events            []interface{}               `json:"events"`
	}

	// ParticipantFrames is a container for all participants - no maps in BigQuery
	ParticipantFrames struct {
		Participant1  ParticipantFrame `json:"1"`
		Participant2  ParticipantFrame `json:"2"`
		Participant3  ParticipantFrame `json:"3"`
		Participant4  ParticipantFrame `json:"4"`
		Participant5  ParticipantFrame `json:"5"`
		Participant6  ParticipantFrame `json:"6"`
		Participant7  ParticipantFrame `json:"7"`
		Participant8  ParticipantFrame `json:"8"`
		Participant9  ParticipantFrame `json:"9"`
		Participant10 ParticipantFrame `json:"10"`
	}

	// ParticipantFrame records the state of a participant at a point in time
	ParticipantFrame struct {
		ParticipantID       int `json:"participantId"`
		Level               int `json:"level"`
		Xp                  int `json:"xp"`
		TotalGold           int `json:"totalGold"`
		CurrentGold         int `json:"currentGold"`
		MinionsKilled       int `json:"minionsKilled"`
		JungleMinionsKilled int `json:"jungleMinionsKilled"`
		Position            struct {
			Y int `json:"y"`
			X int `json:"x"`
		} `json:"position"`
	}
)

// MarshalJSON is for BigQuery, where we can't have column names starting with
// numbers, so pop these into an array for UNNEST
func (f *ParticipantFrames) MarshalJSON() ([]byte, error) {
	if f == nil {
		return json.Marshal([]ParticipantFrame{})
	}
	return json.Marshal([]ParticipantFrame{
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
	})
}
