package dto

type ReqCaptureEvent struct {
	IdEvent      int    `json:"id_event"`
	EventType    string `json:"event_type"`
	EventContent string `json:"event_content"`
}

type ExternalAPIRequest struct {
	Method string            `json:"method"` // "GET" or "POST"
	URL    string            `json:"url"`
	Params map[string]string `json:"params,omitempty"` // For GET
	Body   map[string]any    `json:"body,omitempty"`   // For POST
}
