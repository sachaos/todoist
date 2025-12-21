package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"time"

	todoist "github.com/sachaos/todoist/lib"
)

type PipelineItem struct {
	Item      todoist.Item     `json:"item,omitempty"`
	Command   todoist.Command  `json:"command,omitempty"`
	QuickText string           `json:"quick_text,omitempty"`
	IsQuick   bool             `json:"is_quick"`
	IsClose   bool             `json:"is_close"`
	CloseIDs  []string         `json:"close_ids,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
}

type PipelineCache struct {
	Items []PipelineItem `json:"items"`
	mu    sync.RWMutex   `json:"-"`
}

var (
	globalPipelineCache *PipelineCache
	pipelineCacheMutex  sync.Mutex
)

func GetPipelineCache(pipelineCachePath string) (*PipelineCache, error) {
	pipelineCacheMutex.Lock()
	defer pipelineCacheMutex.Unlock()

	if globalPipelineCache == nil {
		globalPipelineCache = &PipelineCache{
			Items: []PipelineItem{},
		}
		err := LoadPipelineCache(pipelineCachePath, globalPipelineCache)
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}
	return globalPipelineCache, nil
}

func LoadPipelineCache(filename string, pc *PipelineCache) error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	jsonString, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonString, pc)
}

func WritePipelineCache(filename string, pc *PipelineCache) error {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	buf, err := json.MarshalIndent(pc, "", "  ")
	if err != nil {
		return err
	}
	err = AssureExists(filename)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, buf, os.ModePerm)
}

func (pc *PipelineCache) AddItem(item PipelineItem) error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	pc.Items = append(pc.Items, item)
	return nil
}

func (pc *PipelineCache) GetItems() []PipelineItem {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	items := make([]PipelineItem, len(pc.Items))
	copy(items, pc.Items)
	return items
}

func (pc *PipelineCache) RemoveItems(uuids []string) error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	uuidMap := make(map[string]bool)
	for _, uuid := range uuids {
		uuidMap[uuid] = true
	}

	newItems := []PipelineItem{}
	for _, item := range pc.Items {
		if !uuidMap[item.Command.UUID] {
			newItems = append(newItems, item)
		}
	}
	pc.Items = newItems
	return nil
}

func (pc *PipelineCache) Clear() error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	pc.Items = []PipelineItem{}
	return nil
}

func (pc *PipelineCache) IsEmpty() bool {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	return len(pc.Items) == 0
}
