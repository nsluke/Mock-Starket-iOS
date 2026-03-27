package worker

import (
	"encoding/json"

	"github.com/google/uuid"
)

// mustMarshal marshals v to JSON or returns nil on error.
func mustMarshal(v interface{}) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	return data
}

// parseUUID parses a UUID string, returning an error if invalid.
func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
