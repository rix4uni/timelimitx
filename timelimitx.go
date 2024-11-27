package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/pflag"
)

// prints the version message
const version = "v0.0.1"

func PrintVersion() {
	fmt.Printf("Current timelimitx version %s\n", version)
}

// Function to parse time limits like 1s, 1m, etc.
func parseTimeLimit(timeLimit string) (time.Duration, error) {
	if len(timeLimit) < 2 {
		return 0, fmt.Errorf("invalid time format")
	}

	unit := timeLimit[len(timeLimit)-1]
	amountStr := timeLimit[:len(timeLimit)-1]
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		return 0, fmt.Errorf("invalid number in time limit")
	}

	var duration time.Duration
	switch unit {
	case 's':
		duration = time.Duration(amount) * time.Second
	case 'm':
		duration = time.Duration(amount) * time.Minute
	case 'h':
		duration = time.Duration(amount) * time.Hour
	case 'd':
		duration = time.Duration(amount) * 24 * time.Hour
	default:
		return 0, fmt.Errorf("invalid time unit, must be s, m, h, or d")
	}
	return duration, nil
}

// Function to convert signal name to syscall.Signal
func parseSignal(signalName string) (syscall.Signal, error) {
	switch strings.ToUpper(signalName) {
	case "SIGTERM":
		return syscall.SIGTERM, nil
	case "SIGINT":
		return syscall.SIGINT, nil
	case "SIGKILL":
		return syscall.SIGKILL, nil
	default:
		return 0, fmt.Errorf("unsupported signal: %s", signalName)
	}
}

// Function to terminate the process group with a specific signal
func terminateProcessGroup(cmd *exec.Cmd, signal syscall.Signal, verbose bool) {
	if cmd.Process == nil {
		return
	}

	// Send the specified signal to the process group
	if verbose {
		fmt.Printf("Time limit exceeded, sending %s to the process group...\n", signal)
	}
	_ = syscall.Kill(-cmd.Process.Pid, signal)

	time.Sleep(1 * time.Millisecond) // Grace period

	// If the process group is still running, force kill it
	_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
}

func main() {
	// Define flags
	var timeLimit string
	var signalName string
	var versionFlag bool
	var verbose bool

	pflag.StringVarP(&timeLimit, "time", "t", "", "Time limit (e.g., 1s, 1m, 1h)")
	pflag.StringVarP(&signalName, "signal", "s", "SIGTERM", "Signal to send on timeout (e.g., SIGTERM, SIGINT, SIGKILL)")
	pflag.BoolVar(&versionFlag, "version", false, "Print the version of the tool and exit.")
	pflag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
	pflag.Parse()

	if versionFlag {
		PrintVersion()
		return
	}

	// Validate flags
	if timeLimit == "" {
		fmt.Println("Error: --time flag is required")
		pflag.Usage()
		os.Exit(1)
	}

	// Parse time limit
	duration, err := parseTimeLimit(timeLimit)
	if err != nil {
		fmt.Println("Error parsing time limit:", err)
		os.Exit(1)
	}

	// Parse signal
	signal, err := parseSignal(signalName)
	if err != nil {
		fmt.Println("Error parsing signal:", err)
		os.Exit(1)
	}

	// Check for command
	args := pflag.Args()
	if len(args) < 1 {
		fmt.Println("Error: Command to execute is required")
		pflag.Usage()
		os.Exit(1)
	}

	// Combine the command arguments into a single string
	command := strings.Join(args, " ")

	// Use /bin/sh to execute the command
	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set process group ID to allow killing the entire group
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// Start the command
	err = cmd.Start()
	if err != nil {
		fmt.Println("Error starting command:", err)
		os.Exit(1)
	}

	// Create a timer for the time limit
	timer := time.NewTimer(duration)

	// Wait for the command to finish or timeout
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-timer.C:
		terminateProcessGroup(cmd, signal, verbose)
	case err := <-done:
		if err != nil {
			fmt.Println("Command finished with error:", err)
		} else if verbose {
			fmt.Println("Command finished successfully.")
		}
	}
}
