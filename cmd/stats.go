package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show your debugging statistics",
	Long:  "Displays a summary of all your debugging sessions including problems attempted, solved, scores, and times.",
	RunE: func(cmd *cobra.Command, args []string) error {
		sessions, err := db.ListSessions()
		if err != nil {
			return fmt.Errorf("failed to load sessions: %w", err)
		}

		if len(sessions) == 0 {
			fmt.Println()
			color.Yellow("  No sessions yet. Run 'dblab start <problem>' to begin!")
			fmt.Println()
			return nil
		}

		// Calculate stats
		attempted := make(map[string]bool)
		solved := make(map[string]bool)
		var totalScore, bestScore, totalTimeSec int
		var solvedCount int

		for _, s := range sessions {
			attempted[s.ProblemName] = true

			if s.Status == "completed" && s.Score != nil {
				solved[s.ProblemName] = true
				solvedCount++
				totalScore += *s.Score
				if *s.Score > bestScore {
					bestScore = *s.Score
				}
				if s.TimeTakenSeconds != nil {
					totalTimeSec += *s.TimeTakenSeconds
				}
			}
		}

		avgScore := 0
		avgTimeSec := 0
		if solvedCount > 0 {
			avgScore = totalScore / solvedCount
			avgTimeSec = totalTimeSec / solvedCount
		}

		// Print stats
		fmt.Println()
		color.New(color.FgCyan, color.Bold).Println("  📊 DebugLab Statistics")
		fmt.Println("  " + strings.Repeat("─", 40))
		fmt.Println()

		printStat("Problems attempted", fmt.Sprintf("%d", len(attempted)))
		printStat("Problems solved", fmt.Sprintf("%d", len(solved)))
		fmt.Println()
		printStat("Average score", fmt.Sprintf("%d/100", avgScore))
		printStat("Best score", fmt.Sprintf("%d/100", bestScore))
		fmt.Println()
		printStat("Average time", formatSeconds(avgTimeSec))
		printStat("Total time spent", formatSeconds(totalTimeSec))
		fmt.Println()

		// Recent sessions
		color.New(color.FgWhite, color.Bold).Println("  Recent Sessions:")
		fmt.Println("  " + strings.Repeat("─", 55))

		limit := len(sessions)
		if limit > 5 {
			limit = 5
		}

		for _, s := range sessions[:limit] {
			status := s.Status
			var statusColor *color.Color

			switch status {
			case "completed":
				statusColor = color.New(color.FgGreen)
				status = fmt.Sprintf("✓ %d/100", *s.Score)
			case "started":
				statusColor = color.New(color.FgYellow)
				status = "⏳ in progress"
			case "failed":
				statusColor = color.New(color.FgRed)
				status = "✕ failed"
			default:
				statusColor = color.New(color.FgWhite)
			}

			fmt.Print("  ")
			color.New(color.FgWhite).Printf("%-25s ", s.ProblemName)
			statusColor.Printf("%-15s ", status)
			color.New(color.FgWhite, color.Faint).Printf("%s\n", s.StartedAt.Format("2006-01-02 15:04"))
		}

		fmt.Println()
		return nil
	},
}

func printStat(label, value string) {
	color.New(color.FgWhite).Printf("  %-20s ", label)
	color.New(color.FgCyan, color.Bold).Println(value)
}

func formatSeconds(totalSec int) string {
	if totalSec == 0 {
		return "—"
	}
	hours := totalSec / 3600
	minutes := (totalSec % 3600) / 60
	seconds := totalSec % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
