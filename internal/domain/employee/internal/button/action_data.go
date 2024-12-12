package button

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/action"
)

type ActionType int

const (
	ActionTypeCancel ActionType = iota + 1
)

type ActionData struct {
	ID   action.ActionID
	Type ActionType
}

func DecodeFromString(s string) (ActionData, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return ActionData{}, fmt.Errorf("base64 decode: %w", err)
	}

	var data ActionData
	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&data); err != nil {
		return ActionData{}, fmt.Errorf("gob decode: %w", err)
	}

	return data, nil
}

func (a ActionData) EncodeToString() (string, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(a); err != nil {
		return "", fmt.Errorf("gob endode: %w", err)
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
