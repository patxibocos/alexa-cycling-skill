package alexa

const version = "1.0"
const plainText = "PlainText"

type Response struct {
	Version           string                 `json:"version"`
	SessionAttributes map[string]interface{} `json:"sessionAttributes,omitempty"`
	Body              ResBody                `json:"response"`
}

type ResBody struct {
	OutputSpeech     *OutputSpeech `json:"outputSpeech,omitempty"`
	Card             *Card         `json:"card,omitempty"`
	ShouldEndSession bool          `json:"shouldEndSession"`
	Directives       []interface{} `json:"directives,omitempty"`
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

type SendRequestDirective struct {
	Type    string           `json:"type"`
	Name    string           `json:"name"`
	Payload DirectivePayload `json:"payload"`
	Token   string           `json:"token"`
}

type DirectivePayload struct {
	Type             string                     `json:"@type"`
	Version          string                     `json:"@version"`
	PermissionScopes []DirectivePermissionScope `json:"permissionScopes"`
}

type DirectivePermissionScope struct {
	PermissionScope string `json:"permissionScope"`
	ConsentLevel    string `json:"consentLevel"`
}

func newResponse() Response {
	return Response{
		Version: version,
		Body: ResBody{
			OutputSpeech: &OutputSpeech{
				Type: plainText,
			},
		},
	}
}

func (r Response) shouldEndSession(shouldEndSession bool) Response {
	r.Body.ShouldEndSession = shouldEndSession
	return r
}

func (r Response) text(text string) Response {
	r.Body.OutputSpeech.Text = text
	return r
}

func (r Response) sessionAttributes(sessionAttributes map[string]interface{}) Response {
	r.SessionAttributes = sessionAttributes
	return r
}

func (r Response) directives(directives []interface{}) Response {
	r.Body.Directives = directives
	return r
}
