package chat

type EventType string

const (
    EventDelta EventType = "delta"
    EventDone  EventType = "done"
    EventError EventType = "error"
)

type StreamEvent struct {
    Type  EventType
    Text  string
}