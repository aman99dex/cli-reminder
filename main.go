package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
)

const (
	markName  = "GOLANG_CLI_REMINDER"
	markValue = "1"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <hh:mm or natural language time> <text message>\n", os.Args[0])
		fmt.Println("Example: ./reminder '14:30' 'Call mom' or ./reminder 'in 5 minutes' 'Meeting'")
		os.Exit(1)
	}

	now := time.Now()
	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)

	t, err := w.Parse(os.Args[1], now)
	if err != nil {
		fmt.Printf("Error parsing time: %v\n", err)
		os.Exit(2)
	}
	if t == nil {
		fmt.Println("Unable to parse time!")
		os.Exit(2)
	}

	// Check if the parsed time is in the future
	if now.After(t.Time) || now.Equal(t.Time) {
		fmt.Println("Please set a future time!")
		os.Exit(3)
	}

	diff := t.Time.Sub(now)

	if os.Getenv(markName) == markValue {
		// This is the spawned process: sleep and notify
		time.Sleep(diff)
		err := beeep.Alert("Reminder", strings.Join(os.Args[2:], " "), "assets/information.png")
		if err != nil {
			fmt.Printf("Error showing notification: %v\n", err)
			os.Exit(4)
		}
	} else {
		// This is the initial process: spawn a new instance
		cmd := exec.Command(os.Args[0], os.Args[1:]...)
		cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%s", markName, markValue))
		if err := cmd.Start(); err != nil {
			fmt.Printf("Error starting reminder process: %v\n", err)
			os.Exit(5)
		}
		fmt.Printf("Reminder set for %s after %v\n", t.Time.Format("15:04"), diff.Round(time.Second))
		os.Exit(0)
	}
}