package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/shirou/gopsutil/process"
)

func main() {
	fmt.Println("=============================================================")
	fmt.Println("   Dynamic Syncer v1.0")
	fmt.Println("   made w/ hate by Izumi")
	fmt.Println("   you can change the delay with -delay in case you need to")
	fmt.Println("=============================================================")

	delayPtr := flag.Int("delay", 6, "Delay in seconds between detecting the process and launching Dynamic")
	flag.Parse()

	processToMonitor := "SSOClient.exe"

	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("Error determining executable path: %v\n", err)
		os.Exit(1)
	}

	execDir := filepath.Dir(execPath)
	dynamicPath := filepath.Join(execDir, "dynamic_loader.exe")

	if _, err := os.Stat(dynamicPath); os.IsNotExist(err) {
		fmt.Println("ERROR: dynamic_loader.exe not found in the same directory!")
		fmt.Println("Please make sure Dynamic is in the same folder as this program.")
		fmt.Println("Press any key to exit...")
		fmt.Scanln()
		os.Exit(1)
	}

	fmt.Printf("Starting monitor for Star Stable...\n")
	fmt.Printf("Using %d second delay before launching Dynamic\n", *delayPtr)
	fmt.Printf("When Star Stable is detected, will launch Dynamic and then exit\n")

	processWasRunning := false

	for {
		isRunning := isProcessRunning(processToMonitor)

		if isRunning && !processWasRunning {
			fmt.Printf("%s - Star Stable detected! Waiting %d seconds before launching Dynamic...\n",
				time.Now().Format("2006-01-02 15:04:05"),
				*delayPtr)

			time.Sleep(time.Duration(*delayPtr) * time.Second)

			cmd := exec.Command("./dynamic_loader.exe")
			err := cmd.Start()
			if err != nil {
				fmt.Printf("%s - ERROR: Failed to launch Dynamic: %v\n", time.Now().Format("2006-01-02 15:04:05"), err)
			} else {
				fmt.Printf("%s - Dynamic launched successfully.\n", time.Now().Format("2006-01-02 15:04:05"))
				fmt.Println("Mission complete! Exiting Dynamic Syncer now.")
				os.Exit(0)
			}

			processWasRunning = true
		}

		if !isRunning {
			processWasRunning = false
		}

		time.Sleep(2 * time.Second)
	}
}

func isProcessRunning(processName string) bool {
	processes, err := process.Processes()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting process list: %v\n", err)
		return false
	}

	for _, p := range processes {
		name, err := p.Name()
		if err == nil && name == processName {
			return true
		}
	}

	return false
}
