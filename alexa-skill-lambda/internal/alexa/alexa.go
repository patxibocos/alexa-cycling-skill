package alexa

type Request struct {
	Version string  `json:"version"`
	Session Session `json:"session"`
	Body    ReqBody `json:"request"`
	Context Context `json:"context"`
}

type Response struct {
	Version           string                 `json:"version"`
	SessionAttributes map[string]interface{} `json:"sessionAttributes"`
	Body              ResBody                `json:"response"`
}

type ResBody struct {
	OutputSpeech     *OutputSpeech `json:"outputSpeech"`
	Card             *Card         `json:"card"`
	ShouldEndSession bool          `json:"shouldEndSession"`
}

type Card struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Text    string `json:"text"`
}

type OutputSpeech struct {
	Type         string `json:"type"`
	Text         string `json:"text"`
	SSML         string `json:"ssml"`
	PlayBehavior string `json:"playBehavior"`
}

type Session struct {
	New         bool                   `json:"new"`
	SessionID   string                 `json:"sessionId"`
	Attributes  map[string]interface{} `json:"attributes"`
	Application struct {
		ApplicationID string `json:"applicationId"`
	} `json:"application"`
	User struct {
		UserID      string `json:"userId"`
		AccessToken string `json:"accessToken"`
	} `json:"user"`
}

type ReqBody struct {
	Type        string `json:"type"`
	RequestID   string `json:"requestId"`
	Timestamp   string `json:"timestamp"`
	Locale      string `json:"locale"`
	Intent      Intent `json:"intent"`
	Reason      string `json:"reason"`
	DialogState string `json:"dialogState"`
	Error       Error  `json:"error"`
}

type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type Context struct {
	System struct {
		ApiAccessToken string `json:"apiAccessToken"`
		ApiEndpoint    string `json:"apiEndpoint"`
		Device         struct {
			DeviceID string `json:"deviceId"`
		} `json:"device"`
		Application struct {
			ApplicationID string `json:"applicationId"`
		} `json:"application"`
	} `json:"system"`
}

type Intent struct {
	Name               string          `json:"name"`
	ConfirmationStatus string          `json:"confirmationStatus"`
	Slots              map[string]Slot `json:"slots"`
}

type Slot struct {
	Name               string      `json:"name"`
	Value              string      `json:"value"`
	Resolutions        Resolutions `json:"resolutions"`
	SlotValue          SlotValue   `json:"slotValue"`
	ConfirmationStatus string      `json:"confirmationStatus"`
	Source             string      `json:"source"`
}

type SlotValue struct {
	Resolutions Resolutions `json:"resolutions"`
	Type        string      `json:"type"`
	Value       string      `json:"value"`
	Values      []SlotValue `json:"values"`
}

type Resolutions struct {
	ResolutionsPerAuthority []struct {
		Authority string `json:"authority"`
		Status    struct {
			Code string `json:"code"`
		} `json:"status"`
		Values []struct {
			Value struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"value"`
		} `json:"values"`
	} `json:"resolutionsPerAuthority"`
}
