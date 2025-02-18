package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/danish-mehmood/RapidStore/operations"
	"github.com/danish-mehmood/RapidStore/storage"
)

// CLI handles command-line interface operations
type CLI struct {
	ops *operations.Operations
}

// NewCLI creates a new CLI instance
func NewCLI(store storage.Engine) *CLI {
	return &CLI{
		ops: operations.NewOperations(store),
	}
}

// StartInteractive starts the interactive CLI mode
func (c *CLI) StartInteractive() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		command := scanner.Text()
		if command == "exit" {
			break
		}

		if err := c.ExecuteCommand(command); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

// ExecuteCommand executes a single command
func (c *CLI) ExecuteCommand(command string) error {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return nil
	}

	switch strings.ToLower(parts[0]) {
	case "get":
		if len(parts) != 2 {
			return fmt.Errorf("usage: GET <key>")
		}
		value, err := c.ops.Get(parts[1])
		if err != nil {
			return err
		}
		fmt.Printf("Value: %s\n", value)

	case "set":
		if len(parts) < 3 {
			return fmt.Errorf("usage: SET <key> <value>")
		}
		value := strings.Join(parts[2:], " ")
		if err := c.ops.Set(parts[1], value); err != nil {
			return err
		}
		fmt.Printf("Successfully set %s = %s\n", parts[1], value)

	case "delete":
		if len(parts) != 2 {
			return fmt.Errorf("usage: DELETE <key>")
		}
		if err := c.ops.Delete(parts[1]); err != nil {
			return err
		}
		fmt.Printf("Successfully deleted key: %s\n", parts[1])

	case "list":
		keys := c.ops.List()
		if len(keys) == 0 {
			fmt.Println("No keys found")
			return nil
		}
		fmt.Println("Keys:")
		for _, key := range keys {
			fmt.Printf("  - %s\n", key)
		}

	case "help":
		fmt.Println("Available commands:")
		fmt.Println("  GET <key>            - Retrieve value for key")
		fmt.Println("  SET <key> <value>    - Store key-value pair")
		fmt.Println("  DELETE <key>         - Remove key-value pair")
		fmt.Println("  LIST                 - Show all keys")
		fmt.Println("  EXIT                 - Quit the program")

	default:
		return fmt.Errorf("unknown command: %s", parts[0])
	}

	return nil
}
