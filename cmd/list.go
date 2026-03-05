package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/lokeshkudipudi/dblab/internal/runner"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available debugging problems",
	Long:  "Reads all folders in the problems directory and displays their name, difficulty, language, and status.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Fetch problems catalog
		color.Cyan("\nFetching latest problems catalog from GitHub...\n")
		catalog, err := runner.FetchCatalog()
		if err != nil {
			return fmt.Errorf("failed to fetch problem catalog: %w", err)
		}

		if len(catalog) == 0 {
			color.Yellow("No problems found in the catalog.")
			return nil
		}

		// Get solved problems from database
		solved, err := db.GetSolvedProblems()
		if err != nil {
			solved = make(map[string]bool) // continue without status
		}

		// Print header
		header := color.New(color.FgWhite, color.Bold)
		fmt.Println()
		header.Printf("  %-25s %-12s %-20s %s\n", "NAME", "DIFFICULTY", "LANGUAGE", "STATUS")
		fmt.Println("  " + strings.Repeat("─", 70))

		// Print problems
		for _, p := range catalog {
			status := "unsolved"
			statusColor := color.New(color.FgYellow)
			if solved[p.Name] {
				status = "✓ solved"
				statusColor = color.New(color.FgGreen, color.Bold)
			}

			diffColor := difficultyColor(p.Difficulty)

			fmt.Print("  ")
			color.New(color.FgWhite).Printf("%-25s ", p.Name)
			diffColor.Printf("%-12s ", p.Difficulty)
			color.New(color.FgCyan).Printf("%-20s ", p.Language)
			statusColor.Printf("%s\n", status)
		}
		fmt.Println()

		return nil
	},
}

func difficultyColor(difficulty string) *color.Color {
	switch strings.ToLower(difficulty) {
	case "easy":
		return color.New(color.FgGreen)
	case "medium":
		return color.New(color.FgYellow)
	case "hard":
		return color.New(color.FgRed)
	default:
		return color.New(color.FgWhite)
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
}
