package alexa

type Request struct {
	Version string  `json:"version"`
	Session Session `json:"session"`
	Body    ReqBody `json:"request"`
	Context Context `json:"context"`
}

type Response struct {
	Version           string                 `json:"version"`
	SessionAttributes map[string]interface{} `json:"sessionAttributes,omitempty"`
	Body              ResBody                `json:"response"`
}

type ResBody struct {
	OutputSpeech     *OutputSpeech `json:"outputSpeech,omitempty"`
	Card             *Card         `json:"card,omitempty"`
	ShouldEndSession bool          `json:"shouldEndSession"`
}

type Card struct {
	Type    string `json:"type,omitempty"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	Text    string `json:"text,omitempty"`
}

type OutputSpeech struct {
	Type         string `json:"type,omitempty"`
	Text         string `json:"text,omitempty"`
	SSML         string `json:"ssml,omitempty"`
	PlayBehavior string `json:"playBehavior,omitempty"`
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
		AccessToken string `json:"accessToken,omitempty"`
	} `json:"user"`
}

type ReqBody struct {
	Type        string `json:"type"`
	RequestID   string `json:"requestId"`
	Timestamp   string `json:"timestamp"`
	Locale      string `json:"locale"`
	Intent      Intent `json:"intent,omitempty"`
	Reason      string `json:"reason,omitempty"`
	DialogState string `json:"dialogState,omitempty"`
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
			DeviceID string `json:"deviceId,omitempty"`
		} `json:"device,omitempty"`
		Application struct {
			ApplicationID string `json:"applicationId,omitempty"`
		} `json:"application,omitempty"`
	} `json:"System,omitempty"`
}

type Intent struct {
	Name               string          `json:"name,omitempty"`
	ConfirmationStatus string          `json:"confirmationStatus,omitempty"`
	Slots              map[string]Slot `json:"slots,omitempty"`
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
