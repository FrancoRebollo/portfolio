package domain

type Event struct {
	IdEvent      int
	EventType    string
	EventContent string
}

type ExternalAPIRequest struct {
	Method string
	URL    string
	Params map[string]string
	Body   map[string]any
}

type ExternalAPIResponse struct {
	Status     string
	StatusCode int
	Data       any
}
