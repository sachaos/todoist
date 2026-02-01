package todoist

import (
	"context"
)

// Reminder represents a Todoist reminder
// API Reference: https://developer.todoist.com/sync/v9/#reminders
type Reminder struct {
	ID           string `json:"id,omitempty"`
	ItemID       string `json:"item_id"`
	Type         string `json:"type,omitempty"`          // "relative" or "absolute"
	Due          *Due   `json:"due,omitempty"`           // For absolute reminders
	MinuteOffset int    `json:"minute_offset,omitempty"` // For relative reminders (minutes before due)
	IsDeleted    bool   `json:"is_deleted,omitempty"`
}

type Reminders []Reminder

// AddReminder creates a new reminder for a task using the Sync API
// Uses reminder_add command: https://developer.todoist.com/sync/v9/#add-a-reminder
func (c *Client) AddReminder(ctx context.Context, itemID string, due *Due) error {
	c.Log("AddReminder: called for item %s", itemID)

	// Build reminder_add command parameters
	params := map[string]interface{}{
		"item_id": itemID,
		"type":    "absolute",
	}

	// Add due date info if provided
	if due != nil {
		dueParam := map[string]interface{}{}
		if due.Date != "" {
			dueParam["date"] = due.Date
		} else if due.String != "" {
			dueParam["string"] = due.String
		}
		if due.TimeZone != "" {
			dueParam["timezone"] = due.TimeZone
		}
		if len(dueParam) > 0 {
			params["due"] = dueParam
		}
	}

	commands := Commands{
		NewCommand("reminder_add", params),
	}

	return c.ExecCommands(ctx, commands)
}

// AddReminderRelative creates a relative reminder (X minutes before due date)
func (c *Client) AddReminderRelative(ctx context.Context, itemID string, minutesBefore int) error {
	c.Log("AddReminderRelative: called for item %s, %d minutes before", itemID, minutesBefore)

	params := map[string]interface{}{
		"item_id":       itemID,
		"type":          "relative",
		"minute_offset": minutesBefore,
	}

	commands := Commands{
		NewCommand("reminder_add", params),
	}

	return c.ExecCommands(ctx, commands)
}
