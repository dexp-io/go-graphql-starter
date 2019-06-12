package dexp

import "time"

type Activity struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	EntityID   int64     `json:"entity_id"`
	EntityType string    `json:"entity_type"`
	Action     string    `json:"action"`
	Payload    string    `json:"payload"`
	CreatedAt  time.Time `json:"created_at"`
}

func (a *Activity) Create() error {
	res, err := DB.Insert("activity").Fields(map[string]interface{}{
		"user_id":     a.UserID,
		"entity_id":   a.EntityID,
		"entity_type": a.EntityType,
		"action":      a.Action,
		"payload":     a.Payload,
		"created_at":  time.Now(),
	}).Execute()

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return err
	}

	a.ID = id

	return nil
}
