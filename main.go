package main

import (
	"fmt"
)

type Item struct {
	AllDay         bool        `json:"all_day"`
	AssignedByUID  int         `json:"assigned_by_uid"`
	Checked        int         `json:"checked"`
	Collapsed      int         `json:"collapsed"`
	Content        string      `json:"content"`
	DateAdded      string      `json:"date_added"`
	DateLang       string      `json:"date_lang"`
	DateString     string      `json:"date_string"`
	DayOrder       int         `json:"day_order"`
	DueDateUtc     interface{} `json:"due_date_utc"`
	HasMoreNotes   bool        `json:"has_more_notes"`
	ID             int         `json:"id"`
	InHistory      int         `json:"in_history"`
	Indent         int         `json:"indent"`
	IsArchived     int         `json:"is_archived"`
	IsDeleted      int         `json:"is_deleted"`
	ItemOrder      int         `json:"item_order"`
	LabelIDs       []int       `json:"labels"`
	ParentID       interface{} `json:"parent_id"`
	Priority       int         `json:"priority"`
	ProjectID      int         `json:"project_id"`
	ResponsibleUID interface{} `json:"responsible_uid"`
	SyncID         interface{} `json:"sync_id"`
	UserID         int         `json:"user_id"`
}

type Project struct {
	Collapsed    int         `json:"collapsed"`
	Color        int         `json:"color"`
	HasMoreNotes bool        `json:"has_more_notes"`
	ID           int         `json:"id"`
	InboxProject bool        `json:"inbox_project"`
	Indent       int         `json:"indent"`
	IsArchived   int         `json:"is_archived"`
	IsDeleted    int         `json:"is_deleted"`
	ItemOrder    int         `json:"item_order"`
	Name         string      `json:"name"`
	ParentID     interface{} `json:"parent_id"`
	Shared       bool        `json:"shared"`
}

type Label struct {
	Color     int    `json:"color"`
	ID        int    `json:"id"`
	IsDeleted int    `json:"is_deleted"`
	ItemOrder int    `json:"item_order"`
	Name      string `json:"name"`
}

func FindByID(lables []Label, id int) (Label, interface{}) {
	for _, label := range lables {
		if label.ID == id {
			return label, nil
		}
	}
	return Label{}, "NotFound"
}

type User struct {
	AutoReminder      int         `json:"auto_reminder"`
	AvatarBig         string      `json:"avatar_big"`
	AvatarMedium      string      `json:"avatar_medium"`
	AvatarS640        string      `json:"avatar_s640"`
	AvatarSmall       string      `json:"avatar_small"`
	BusinessAccountID interface{} `json:"business_account_id"`
	CompletedCount    int         `json:"completed_count"`
	CompletedToday    int         `json:"completed_today"`
	DailyGoal         int         `json:"daily_goal"`
	DateFormat        int         `json:"date_format"`
	DefaultReminder   string      `json:"default_reminder"`
	Email             string      `json:"email"`
	Features          struct {
		Beta             int  `json:"beta"`
		GoldTheme        bool `json:"gold_theme"`
		HasPushReminders bool `json:"has_push_reminders"`
		Restriction      int  `json:"restriction"`
	} `json:"features"`
	FullName     string      `json:"full_name"`
	ID           int         `json:"id"`
	ImageID      string      `json:"image_id"`
	InboxProject int         `json:"inbox_project"`
	IsBizAdmin   bool        `json:"is_biz_admin"`
	IsPremium    bool        `json:"is_premium"`
	JoinDate     string      `json:"join_date"`
	Karma        float32     `json:"karma"`
	KarmaTrend   string      `json:"karma_trend"`
	MobileHost   interface{} `json:"mobile_host"`
	MobileNumber interface{} `json:"mobile_number"`
	NextWeek     int         `json:"next_week"`
	PremiumUntil string      `json:"premium_until"`
	SortOrder    int         `json:"sort_order"`
	StartDay     int         `json:"start_day"`
	StartPage    string      `json:"start_page"`
	Theme        int         `json:"theme"`
	TimeFormat   int         `json:"time_format"`
	Token        string      `json:"token"`
	TzInfo       struct {
		GmtString string `json:"gmt_string"`
		Hours     int    `json:"hours"`
		IsDst     int    `json:"is_dst"`
		Minutes   int    `json:"minutes"`
		Timezone  string `json:"timezone"`
	} `json:"tz_info"`
}

func main() {
	var config Config
	var token string
	var sync Sync
	config, err := ParseConfig("./.todoist.config.json")
	if err != nil {
		fmt.Scan(&token)
		config = Config{Token: token}
		err = CreateConfig("./.todoist.config.json", config)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	token = config.Token

	sync, err = LoadCache("./.todoist.cache.json")
	if err != nil {
		sync, err = FetchCache(token)
		if err != nil {
			return
		}
		err = SaveCache("./.todoist.cache.json", sync)
		if err != nil {
			return
		}
	}

	for _, item := range sync.Items {
		fmt.Printf("%d,p%d,%s,", item.ID, item.Priority, item.Content)
		for _, label_id := range item.LabelIDs {
			label, err := FindByID(sync.Labels, label_id)
			if err != nil {
				return
			}
			fmt.Printf("@%s", label.Name)
		}
		fmt.Printf("\n")
	}
}
