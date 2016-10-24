package main

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func main() {
	var config Config
	var token string
	var sync Sync
	config, err := ParseConfig(os.Getenv("HOME") + "/.todoist.config.json")
	if err != nil {
		fmt.Scan(&token)
		config = Config{Token: token}
		err = CreateConfig(os.Getenv("HOME")+"/.todoist.config.json", config)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	token = config.Token

	sync, err = LoadCache(os.Getenv("HOME") + "/.todoist.cache.json")
	if err != nil {
		sync, err = FetchCache(token)
		if err != nil {
			return
		}
		err = SaveCache(os.Getenv("HOME")+"/.todoist.cache.json", sync)
		if err != nil {
			return
		}
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 4, 1, ' ', 0)

	for _, item := range sync.Items {
		fmt.Fprintf(w, "%d\tp%d\t%s\t%s\n", item.ID, item.Priority, LabelsString(item, sync.Labels), item.Content)
		// for _, label_id := range item.LabelIDs {
		// 	label, err := FindByID(sync.Labels, label_id)
		// 	if err != nil {
		// 		return
		// 	}
		// 	fmt.Printf("@%s", label.Name)
		// }
		// fmt.Printf("\n")
	}
	w.Flush()
}
