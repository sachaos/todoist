package lib

import (
	"encoding/json"
	"net/url"
)

type Sync struct {
	CollaboratorStates []interface{} `json:"collaborator_states"`
	Collaborators      []interface{} `json:"collaborators"`
	DayOrders          interface{}   `json:"day_orders"`
	DayOrdersTimestamp string        `json:"day_orders_timestamp"`
	Filters            []struct {
		Color     int    `json:"color"`
		ID        int    `json:"id"`
		IsDeleted int    `json:"is_deleted"`
		ItemOrder int    `json:"item_order"`
		Name      string `json:"name"`
		Query     string `json:"query"`
	} `json:"filters"`
	FullSync          bool   `json:"full_sync"`
	Items             Items  `json:"items"`
	Labels            Labels `json:"labels"`
	LiveNotifications []struct {
		CompletedTasks   int    `json:"completed_tasks"`
		Created          int    `json:"created"`
		DateReached      int    `json:"date_reached"`
		ID               int    `json:"id"`
		IsDeleted        int    `json:"is_deleted"`
		KarmaLevel       int    `json:"karma_level"`
		NotificationKey  string `json:"notification_key"`
		NotificationType string `json:"notification_type"`
		PromoImg         string `json:"promo_img"`
		SeqNo            int    `json:"seq_no"`
		TopProcent       int    `json:"top_procent"`
	} `json:"live_notifications"`
	LiveNotificationsLastReadID int           `json:"live_notifications_last_read_id"`
	Locations                   []interface{} `json:"locations"`
	Notes                       []struct {
		Content        string      `json:"content"`
		FileAttachment interface{} `json:"file_attachment"`
		ID             int         `json:"id"`
		IsArchived     int         `json:"is_archived"`
		IsDeleted      int         `json:"is_deleted"`
		ItemID         int         `json:"item_id"`
		Posted         string      `json:"posted"`
		PostedUID      int         `json:"posted_uid"`
		ProjectID      int         `json:"project_id"`
		UidsToNotify   interface{} `json:"uids_to_notify"`
	} `json:"notes"`
	ProjectNotes []interface{} `json:"project_notes"`
	Projects     Projects      `json:"projects"`
	Reminders    []struct {
		DateLang     string `json:"date_lang"`
		DueDateUtc   string `json:"due_date_utc"`
		ID           int    `json:"id"`
		IsDeleted    int    `json:"is_deleted"`
		ItemID       int    `json:"item_id"`
		MinuteOffset int    `json:"minute_offset"`
		NotifyUID    int    `json:"notify_uid"`
		Service      string `json:"service"`
		Type         string `json:"type"`
	} `json:"reminders"`
	SyncToken     string   `json:"sync_token"`
	TempIDMapping struct{} `json:"temp_id_mapping"`
	User          User     `json:"user"`
}

func SyncAll(token string) (Sync, error) {
	var sync Sync
	body, err := SyncRequest(
		url.Values{"token": {token}, "sync_token": {"*"}, "resource_types": {"[\"all\"]"}},
	)
	err = json.Unmarshal(body, &sync)
	if err != nil {
		return Sync{}, err
	}
	return sync, nil
}
