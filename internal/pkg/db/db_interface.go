package db

import (
	"time"
)

// DBInterface is the interface that wraps the functions chronicler uses to communicate with the database.
type DBInterface interface {
	// Write a single event to a given stream.
	WriteEvent(streamName string, event string, uuidString string, version string) (int64, error)

	// Conditionally write an event to a given stream.  If the current last event in the stream is at previousPosition,
	// then the write will succeed; otherwise, it will fail and an error will be returned.
	ConditionallyWriteEvent(streamName string, event string, version string, previousPosition uint64) error

	// Read a single event from a given stream at a given position.
	ReadEvent(streamName string, position uint64) (string, error)

	// Reads all of the events from a given stream.
	ReadAllEvents(streamName string) (string, error)

	// Reads all of the events from a given stream, starting at a given position.
	ReadEventsFromPosition(streamName string, position uint64) (string, error)

	// Reads all of the events from a given stream, starting at a given datetime.
	ReadEventsFromTime(streamName string, datetime time.Time) (string, error)
}
