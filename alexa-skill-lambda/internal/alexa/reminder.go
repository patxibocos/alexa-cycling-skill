package alexa

type reminderRequest struct {
	RequestTime      string           `json:"requestTime"`
	Trigger          trigger          `json:"trigger"`
	AlertInfo        alertInfo        `json:"alertInfo"`
	PushNotification pushNotification `json:"pushNotification"`
}

type trigger struct {
	Type          string `json:"type"`
	ScheduledTime string `json:"scheduledTime"`
	TimeZoneID    string `json:"timeZoneId"`
}

type alertInfo struct {
	SpokenInfo spokenInfo `json:"spokenInfo"`
}

type spokenInfo struct {
	Content []content `json:"content"`
}

type content struct {
	Locale string `json:"locale"`
	Text   string `json:"text"`
	Ssml   string `json:"ssml,omitempty"`
}

type pushNotification struct {
	Status string `json:"status"`
}

type remindersResponse struct {
	TotalCount string  `json:"totalCount"`
	Alerts     []alert `json:"alerts"`
}

type alert struct {
	Trigger trigger `json:"trigger"`
}
