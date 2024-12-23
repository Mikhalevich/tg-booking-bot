package actionpayload

import (
	"encoding/json"
	"fmt"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/empl"
)

type EmployeeID struct {
	ID empl.EmployeeID `json:"employee_id"`
}

func BytesFromEmployeeID(id empl.EmployeeID) ([]byte, error) {
	b, err := json.Marshal(EmployeeID{
		ID: id,
	})
	if err != nil {
		return nil, fmt.Errorf("json marshal: %w", err)
	}

	return b, nil
}

func EmployeeIDFromBytes(payload []byte) (EmployeeID, error) {
	var eID EmployeeID
	if err := json.Unmarshal(payload, &eID); err != nil {
		return EmployeeID{}, fmt.Errorf("json unmarshal: %w", err)
	}

	return eID, nil
}
