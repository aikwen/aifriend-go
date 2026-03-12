package chat

type EventType string

const (
    EventDelta EventType = "delta"
    EventDone  EventType = "done"
    EventError EventType = "error"
)

type Usage struct {
    InputTokens  int
    OutputTokens int
    TotalTokens  int
}

type StreamEvent struct {
    Type  EventType
    Text  string
    Usage *Usage
}