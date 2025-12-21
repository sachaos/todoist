package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rkoesters/xdg/basedir"
	todoist "github.com/sachaos/todoist/lib"
)

func logBackgroundSync(format string, v ...interface{}) {
	logPath := filepath.Join(basedir.CacheHome, "todoist", "background-sync.log")
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags)
	logger.Printf(format, v...)
}

func BackgroundSyncWorker(client *todoist.Client, pipelineCachePath string, cachePath string) {
	logBackgroundSync("=== Background sync started ===")

	pc, err := GetPipelineCache(pipelineCachePath)
	if err != nil {
		client.Log("Background sync: failed to get pipeline cache: %v", err)
		logBackgroundSync("ERROR: Failed to get pipeline cache: %v", err)
		return
	}

	if pc.IsEmpty() {
		client.Log("Background sync: pipeline cache is empty, nothing to sync")
		logBackgroundSync("Pipeline cache is empty, nothing to sync")
		return
	}

	items := pc.GetItems()
	if len(items) == 0 {
		logBackgroundSync("No items to sync")
		return
	}

	client.Log("Background sync: syncing %d items", len(items))
	logBackgroundSync("Syncing %d items", len(items))

	commands := todoist.Commands{}
	uuids := []string{}

	for _, pipelineItem := range items {
		if pipelineItem.IsQuick {
			err := client.QuickCommand(context.Background(), pipelineItem.QuickText)
			if err != nil {
				client.Log("Background sync: failed to sync quick command: %v", err)
				continue
			}
			uuids = append(uuids, pipelineItem.Command.UUID)
		} else {
			commands = append(commands, pipelineItem.Command)
			uuids = append(uuids, pipelineItem.Command.UUID)
		}
	}

	if len(commands) > 0 {
		logBackgroundSync("Executing %d commands...", len(commands))
		err = client.ExecCommands(context.Background(), commands)
		if err != nil {
			client.Log("Background sync: failed to sync commands: %v", err)
			logBackgroundSync("ERROR: Failed to sync commands: %v", err)
			fmt.Fprintf(os.Stderr, "\nBackground sync failed: %v\n", err)
			fmt.Fprintf(os.Stderr, "Your task has been saved locally and will sync on next 'todoist sync' or add/quick command.\n")
			fmt.Fprintf(os.Stderr, "Check log: ~/.cache/todoist/background-sync.log\n")
			return
		}
		logBackgroundSync("Commands executed successfully")
	}

	logBackgroundSync("Syncing with server...")
	err = client.Sync(context.Background())
	if err != nil {
		client.Log("Background sync: failed to sync with server: %v", err)
		logBackgroundSync("ERROR: Failed to sync with server: %v", err)
		fmt.Fprintf(os.Stderr, "\nBackground sync failed: %v\n", err)
		fmt.Fprintf(os.Stderr, "Your task has been saved locally and will sync on next 'todoist sync' or add/quick command.\n")
		fmt.Fprintf(os.Stderr, "Check log: ~/.cache/todoist/background-sync.log\n")
		return
	}
	logBackgroundSync("Server sync completed")

	logBackgroundSync("Writing cache...")
	err = WriteCache(cachePath, client.Store)
	if err != nil {
		client.Log("Background sync: failed to write cache: %v", err)
		logBackgroundSync("ERROR: Failed to write cache: %v", err)
	}

	logBackgroundSync("Removing %d synced items from pipeline cache...", len(uuids))
	err = pc.RemoveItems(uuids)
	if err != nil {
		client.Log("Background sync: failed to remove items from pipeline cache: %v", err)
		logBackgroundSync("ERROR: Failed to remove items: %v", err)
		return
	}

	err = WritePipelineCache(pipelineCachePath, pc)
	if err != nil {
		client.Log("Background sync: failed to write pipeline cache: %v", err)
		logBackgroundSync("ERROR: Failed to write pipeline cache: %v", err)
		fmt.Fprintf(os.Stderr, "\nWarning: failed to update pipeline cache: %v\n", err)
	}

	client.Log("Background sync: successfully synced %d items", len(uuids))
	logBackgroundSync("SUCCESS: Synced %d items", len(uuids))
	if len(uuids) > 0 {
		fmt.Fprintf(os.Stderr, "Background sync completed: %d task(s) synced to Todoist.\n", len(uuids))
	}
	logBackgroundSync("=== Background sync finished ===")
}

func StartBackgroundSync(client *todoist.Client, pipelineCachePath string, cachePath string) {
	executable, err := os.Executable()
	if err != nil {
		logBackgroundSync("ERROR: Failed to get executable path: %v", err)
		go BackgroundSyncWorker(client, pipelineCachePath, cachePath)
		return
	}

	cmd := exec.Command(executable, "__background_sync__")
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil

	err = cmd.Start()
	if err != nil {
		logBackgroundSync("ERROR: Failed to start background process: %v", err)
		go BackgroundSyncWorker(client, pipelineCachePath, cachePath)
		return
	}
}
