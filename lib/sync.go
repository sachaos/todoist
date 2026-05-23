package todoist

type Store struct {
	CollaboratorStates []interface{} `json:"collaborator_states"`
	Collaborators      []interface{} `json:"collaborators"`
	DayOrders          interface{}   `json:"day_orders"`
	DayOrdersTimestamp string        `json:"day_orders_timestamp"`
	Filters            []struct {
		ID        string `json:"id"`        // The ID of the filter.
		Name      string `json:"name"`      // The name of the filter.
		Query     string `json:"query"`     // The query to search for. Examples of searches can be found in the Todoist help page.
		Color     string `json:"color"`    // The color of the filter icon. Refer to the name column in the Colors guide for more info.
		ItemOrder int    `json:"item_order"`     // Filter's order in the filter list (where the smallest value should place the filter at the top).
		IsDeleted bool   `json:"is_deleted"`    // Whether the filter is marked as deleted (a true or false value).
		IsFavorite bool  `json:"is_favorite"`   // Whether the filter is a favorite (a true or false value).
		IsFrozen  bool   `json:"is_frozen"`     // Filters from a canceled subscription cannot be changed. This is a read-only attribute (a true or false value).
	} `json:"filters"`
	FullSync          bool   `json:"full_sync"`
	Items             Items  `json:"items"`
	Labels            Labels `json:"labels"`
	LiveNotifications []struct {
		CompletedTasks   int     `json:"completed_tasks"`
		Created          string  `json:"created_at"`
		DateReached      string  `json:"date_reached"`
		ID               string  `json:"id"`
		IsDeleted        bool    `json:"is_deleted"`
		KarmaLevel       int     `json:"karma_level"`
		NotificationKey  string  `json:"notification_key"`
		NotificationType string  `json:"notification_type"`
		PromoImg         string  `json:"promo_img"`
		SeqNo            int64   `json:"seq_no"`
		TopProcent       float32 `json:"top_procent"`
	} `json:"live_notifications"`
	LiveNotificationsLastReadID string        `json:"live_notifications_last_read_id"`
	Locations                   []interface{} `json:"locations"`
	Notes                       []struct {
		Content        string      `json:"content"`
		FileAttachment interface{} `json:"file_attachment"`
		ID             string      `json:"id"`
		IsArchived     int         `json:"is_archived"`
		IsDeleted      bool        `json:"is_deleted"`
		ItemID         string      `json:"item_id"`
		Posted         string      `json:"posted_at"`
		PostedUID      string      `json:"posted_uid"`
		ProjectID      *string     `json:"project_id"`
		UidsToNotify   interface{} `json:"uids_to_notify"`
	} `json:"notes"`
	ProjectNotes []interface{} `json:"project_notes"`
	Projects     Projects      `json:"projects"`
	Sections     Sections      `json:"sections"`
	Reminders    []struct {
		DateLang     string      `json:"date_lang"`
		Due          *Due        `json:"due"`
		ID           string      `json:"id"`
		IsDeleted    bool        `json:"is_deleted"`
		ItemID       string      `json:"item_id"`
		MinuteOffset int         `json:"minute_offset"`
		NotifyUID    string      `json:"notify_uid"`
		Service      interface{} `json:"service"`
		Type         string      `json:"type"`
	} `json:"reminders"`
	SyncToken     string              `json:"sync_token"`
	SchemaVersion int                 `json:"schema_version"`
	TempIDMapping struct{}            `json:"temp_id_mapping"`
	User          User                `json:"user"`
	RootItem      *Item               `json:"-"`
	RootProject   *Project            `json:"-"`
	ItemMap       map[string]*Item    `json:"-"`
	ProjectMap    map[string]*Project `json:"-"`
	LabelMap      map[string]*Label   `json:"-"`
	SectionMap    map[string]*Section `json:"-"`
}

func (s *Store) FindItem(id string) *Item {
	return s.ItemMap[id]
}

func (s *Store) FindProject(id string) *Project {
	return s.ProjectMap[id]
}

func (s *Store) FindSection(id string) *Section {
	return s.SectionMap[id]
}

func (s *Store) FindLabel(id string) *Label {
	return s.LabelMap[id]
}

func addToBrotherItem(item *Item, b *Item) {
	i := item
	for {
		if i.BrotherItem == nil {
			i.BrotherItem = b
			return
		}
		i = i.BrotherItem
	}
}

func addToChildItem(item *Item, b *Item) {
	if item.ChildItem == nil {
		item.ChildItem = b
		return
	}
	addToBrotherItem(item.ChildItem, b)
}

func addToBrotherProject(project *Project, b *Project) {
	i := project
	for {
		if i.BrotherProject == nil {
			i.BrotherProject = b
			return
		}
		i = i.BrotherProject
	}
}

func addToChildProject(project *Project, b *Project) {
	if project.ChildProject == nil {
		project.ChildProject = b
		return
	}
	addToBrotherProject(project.ChildProject, b)
}

func (s *Store) ConstructItemTree() {
	s.LabelMap = map[string]*Label{}
	s.ProjectMap = map[string]*Project{}
	s.ItemMap = map[string]*Item{}
	s.SectionMap = map[string]*Section{}

	for i, label := range s.Labels {
		s.LabelMap[label.ID] = &s.Labels[i]
	}

	for i, item := range s.Items {
		s.ItemMap[item.ID] = &s.Items[i]
		s.Items[i].ChildItem = nil
		s.Items[i].BrotherItem = nil
	}

	for i, project := range s.Projects {
		s.ProjectMap[project.ID] = &s.Projects[i]
		s.Projects[i].ChildProject = nil
		s.Projects[i].BrotherProject = nil
	}

	for i, section := range s.Sections {
		s.SectionMap[section.ID] = &s.Sections[i]
	}

	for _, item := range s.Items {
		if item.ParentID == nil {
			s.RootItem = &item
			break
		}
	}

	for _, project := range s.Projects {
		if project.ParentID == nil {
			s.RootProject = &project
			break
		}
	}

	for i := range s.Items {
		if s.Items[i].ID == s.RootItem.ID {
			continue
		}

		if s.Items[i].ParentID == nil {
			addToBrotherItem(s.RootItem, &s.Items[i])
			continue
		}
		id, _ := s.Items[i].GetParentID()
		parent := s.FindItem(id)
		addToChildItem(parent, &s.Items[i])
	}

	for i := range s.Projects {
		if s.Projects[i].ID == s.RootProject.ID {
			continue
		}

		if s.Projects[i].ParentID == nil {
			addToBrotherProject(s.RootProject, &s.Projects[i])
			continue
		}
		id, _ := s.Projects[i].GetParentID()
		parent := s.FindProject(id)
		addToChildProject(parent, &s.Projects[i])
	}
}
