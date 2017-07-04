package kv

import (
	"context"
	"encoding/json"
)

// ReadJSON returns the value decoded from json of the key.
func ReadJSON(ctx context.Context, s Reader, key string, value interface{}) error {
	data, err := s.Read(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}

// WriteJSON sets the value encoded as json for the key.
func WriteJSON(ctx context.Context, s Store, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.Write(ctx, key, data)
}
