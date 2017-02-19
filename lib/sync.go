package todoist

import (
	"encoding/json"
	"net/url"
	"sort"
)

type Sync struct {
	CollaboratorStates []interface{} `json:"collaborator_states"`
	Collaborators      []interface{} `json:"collaborators"`
	DayOrders          interface{}   `json:"day_orders"`
	DayOrdersTimestamp string        `json:"day_orders_timestamp"`
	ItemOrders         ItemOrders    `json:"-"`
	ProjectOrders      Orders        `json:"-"`
	LabelOrders        Orders        `json:"-"`
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
		CompletedTasks   int     `json:"completed_tasks"`
		Created          int     `json:"created"`
		DateReached      int     `json:"date_reached"`
		ID               int     `json:"id"`
		IsDeleted        int     `json:"is_deleted"`
		KarmaLevel       int     `json:"karma_level"`
		NotificationKey  string  `json:"notification_key"`
		NotificationType string  `json:"notification_type"`
		PromoImg         string  `json:"promo_img"`
		SeqNo            int64   `json:"seq_no"`
		TopProcent       float32 `json:"top_procent"`
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

func (sync *Sync) ConstructItemOrder() {
	sort.Sort(sync.Projects)
	sort.Sort(sync.Items)
	sort.Sort(sync.Labels)

	sync.ProjectOrders = make(Orders, len(sync.Projects))
	for i := 0; i < len(sync.Projects); i++ {
		project := sync.Projects[i]
		sync.ProjectOrders[i] = Order{Num: project.ItemOrder, ID: project.ID, Data: project}
	}
	sort.Sort(sync.ProjectOrders)

	sync.LabelOrders = make(Orders, len(sync.Labels))
	for i := 0; i < len(sync.Labels); i++ {
		label := sync.Labels[i]
		sync.LabelOrders[i] = Order{Num: label.ItemOrder, ID: label.ID, Data: label}
	}
	sort.Sort(sync.LabelOrders)

	sync.ItemOrders = make(ItemOrders, len(sync.Items))
	for i := 0; i < len(sync.Items); i++ {
		item := sync.Items[i]
		project, err := SearchByID(sync.Projects, item.ProjectID)
		if err != nil {
			panic(err)
		}
		sync.ItemOrders[i] = ItemOrder{Order: Order{Num: item.ItemOrder, ID: item.ID, Data: item}, ProjectOrder: project.(Project).ItemOrder}
	}
	sort.Sort(sync.ItemOrders)
}

func SyncAll(token string) (Sync, error) {
	var sync Sync
	body, err := APIRequest("sync",
		url.Values{"token": {token}, "sync_token": {"*"}, "resource_types": {"[\"all\"]"}},
	)
	err = json.Unmarshal(body, &sync)
	if err != nil {
		return Sync{}, err
	}
	sync.ConstructItemOrder()
	return sync, nil
}
