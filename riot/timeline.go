package riot

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
		Events            []FrameEvent       `json:"events"`
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
)
