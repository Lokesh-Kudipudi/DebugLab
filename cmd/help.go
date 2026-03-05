package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Show help for all commands",
	Long:  "Displays a list of all available commands with descriptions.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()
		color.New(color.FgCyan, color.Bold).Println("  🔬 DebugLab — Practice debugging real-world bugs")
		fmt.Println()
		color.New(color.FgWhite, color.Bold).Println("  USAGE:")
		color.White("    dblab <command> [arguments]")
		fmt.Println()
		color.New(color.FgWhite, color.Bold).Println("  COMMANDS:")
		fmt.Println()
		printCommand("list", "List all available debugging problems")
		printCommand("start <name>", "Start a debugging session for a problem")
		printCommand("verify <name>", "Run tests and score your solution")
		printCommand("stats", "Show your debugging statistics")
		printCommand("help", "Show this help message")
		fmt.Println()
		color.New(color.FgWhite, color.Bold).Println("  WORKFLOW:")
		fmt.Println()
		color.White("    1. Run " + color.CyanString("dblab list") + " to see available problems")
		color.White("    2. Run " + color.CyanString("dblab start <name>") + " to begin a session")
		color.White("    3. Fix the bugs in your workspace")
		color.White("    4. Run " + color.CyanString("dblab verify <name>") + " to check your solution")
		color.White("    5. Run " + color.CyanString("dblab stats") + " to track your progress")
		fmt.Println()
	},
}

func printCommand(name, description string) {
	color.New(color.FgCyan).Printf("    %-20s ", name)
	color.New(color.FgWhite).Println(description)
}

func init() {
	rootCmd.AddCommand(helpCmd)
}
