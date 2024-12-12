package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		current_dir, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error getting current directory:", err)
			current_dir = "unknown"
		}
		fmt.Printf("%s> ", current_dir)
		// Read the keyboad input.
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		// Handle the execution of the input.
		if err = execInput(input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func translateCommand(input string) string {
	input = strings.TrimSpace(input)

	switch {
	case strings.HasPrefix(input, "ls -l"):
		return "Get-ChildItem | Format-List"
	case strings.HasPrefix(input, "ls"):
		return "Get-ChildItem"
	case strings.HasPrefix(input, "pwd"):
		return "Get-Location"
	case strings.HasPrefix(input, "cd"):
		return input
	case strings.HasPrefix(input, "exit"):
		os.Exit(0)
	}
	return input
}

func execInput(input string) error {
	// Translate the command.
	input = translateCommand(input)

	if strings.HasPrefix(input, "cd") {
		// Extract the directory from the input.
		parts := strings.Fields(input)
		if len(parts) < 2 {
			return fmt.Errorf("cd: missing argument")
		}
		err := os.Chdir(parts[1])
		if err != nil {
			return fmt.Errorf("cd: %v", err)
		}
		return nil
	}

	// Prepare the PowerShell command.
	cmd := exec.Command("powershell", "-Command", input)

	// Set the correct output device.
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// Execute the command and return the error.
	return cmd.Run()
}
