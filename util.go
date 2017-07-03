package kv

import "encoding/json"

// ReadJSON returns the value decoded from json of the key.
func ReadJSON(s Reader, key string, value interface{}) error {
	data, err := s.Read(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}

// WriteJSON sets the value encoded as json for the key.
func WriteJSON(s Store, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.Write(key, data)
}
