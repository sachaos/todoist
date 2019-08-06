package todoist

import (
	"sort"
)

type Store struct {
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
		Created          string  `json:"created"`
		DateReached      string  `json:"date_reached"`
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
		Due          *Due   `json:"due"`
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

func (s *Store) ConstructItemOrder() {
	sort.Sort(s.Projects)
	sort.Sort(s.Items)
	sort.Sort(s.Labels)

	s.ProjectOrders = make(Orders, len(s.Projects))
	for i := 0; i < len(s.Projects); i++ {
		project := s.Projects[i]
		s.ProjectOrders[i] = Order{Num: project.ItemOrder, ID: project.ID, Data: project}
	}
	sort.Sort(s.ProjectOrders)

	s.LabelOrders = make(Orders, len(s.Labels))
	for i := 0; i < len(s.Labels); i++ {
		label := s.Labels[i]
		s.LabelOrders[i] = Order{Num: label.ItemOrder, ID: label.ID, Data: label}
	}
	sort.Sort(s.LabelOrders)

	s.ItemOrders = make(ItemOrders, len(s.Items))
	for i := 0; i < len(s.Items); i++ {
		item := s.Items[i]
		project, err := SearchByID(s.Projects, item.ProjectID)
		var pjtOrder int
		if err != nil {
			// Set unknown project order to 0
			pjtOrder = 0
		} else {
			pjtOrder = project.(Project).ItemOrder
		}
		s.ItemOrders[i] = ItemOrder{Order: Order{Num: item.ItemOrder, ID: item.ID, Data: item}, ProjectOrder: pjtOrder}
	}
	sort.Sort(s.ItemOrders)
}
