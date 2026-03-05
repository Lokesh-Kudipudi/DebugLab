package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/lokeshkudipudi/dblab/internal/runner"
	"github.com/lokeshkudipudi/dblab/internal/scoring"
	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify <problem-name>",
	Short: "Verify your solution by running the test suite",
	Long:  "Runs the test suite for the problem and calculates your score if all tests pass.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		problemName := args[0]

		// Get the active session
		session, err := db.GetActiveSessionByProblem(problemName)
		if err != nil {
			return fmt.Errorf("no active session for %q. Run 'dblab start %s' first", problemName, problemName)
		}

		if cfg == nil {
			return fmt.Errorf("no workspace configured. Please run 'dblab init <directory>' first")
		}
		
		workspaceDir := filepath.Join(cfg.WorkspaceDir, problemName)

		// Check workspace exists
		if _, err := os.Stat(workspaceDir); os.IsNotExist(err) {
			return fmt.Errorf("workspace not found at %s. Run 'dblab start %s' first", workspaceDir, problemName)
		}

		// Load the problem definition from the workspace
		problem, err := runner.LoadProblem(workspaceDir)
		if err != nil {
			return fmt.Errorf("could not load problem definition: %w", err)
		}

		fmt.Println()
		color.New(color.FgCyan, color.Bold).Printf("🔬 Verifying: %s\n", problem.Title)
		fmt.Println()
		color.White("  Running test suite...")
		fmt.Println()

		// Run tests
		result, err := runner.RunTests(workspaceDir, problem.Language)
		if err != nil {
			return fmt.Errorf("failed to run tests: %w", err)
		}

		// Calculate time taken
		timeTaken := time.Since(session.StartedAt)
		timeTakenSeconds := int(timeTaken.Seconds())

		// Count changed files
		filesChanged, err := runner.CountChangedFiles(workspaceDir)
		if err != nil {
			filesChanged = 0
		}

		// Print test results
		fmt.Println("  " + color.New(color.FgWhite, color.Bold).Sprint("─────────────────────────────────────────────"))
		fmt.Println()

		if result.AllPassed {
			color.New(color.FgGreen, color.Bold).Println("  ✅ All tests passed!")
		} else {
			color.New(color.FgRed, color.Bold).Println("  ❌ Some tests failed!")
		}

		fmt.Println()
		color.New(color.FgWhite).Printf("  Tests passed:  ")
		color.New(color.FgGreen).Printf("%d", result.Passed)
		color.New(color.FgWhite).Printf("/%d\n", result.Total)

		if result.Failed > 0 {
			color.New(color.FgWhite).Printf("  Tests failed:  ")
			color.New(color.FgRed).Printf("%d\n", result.Failed)
		}

		fmt.Println()
		color.New(color.FgWhite).Printf("  ⏱  Time taken:    %s\n", formatDuration(timeTaken))
		color.New(color.FgWhite).Printf("  📝 Files changed: %d\n", filesChanged)

		if result.AllPassed {
			// Calculate score
			parTime := time.Duration(problem.ParTimeMinutes) * time.Minute
			score := scoring.Calculate(timeTaken, parTime, session.HintsUsed, filesChanged, problem.ExpectedFiles)

			// Update session as completed
			if err := db.CompleteSession(session.ID, timeTakenSeconds, score, result.Passed, result.Failed, filesChanged); err != nil {
				return fmt.Errorf("failed to save session: %w", err)
			}

			fmt.Println()
			scoreColor := scoreGradeColor(score)
			color.New(color.FgWhite, color.Bold).Print("  🏆 Score: ")
			scoreColor.Printf("%d/100\n", score)
			fmt.Println()

			printScoreBreakdown(timeTaken, time.Duration(problem.ParTimeMinutes)*time.Minute, session.HintsUsed, filesChanged, problem.ExpectedFiles)
		} else {
			// Update session as failed
			if err := db.FailSession(session.ID, result.Passed, result.Failed); err != nil {
				return fmt.Errorf("failed to save session: %w", err)
			}

			fmt.Println()
			color.Yellow("  Fix the failing tests and run 'dblab verify %s' again.", problemName)
		}

		fmt.Println()
		return nil
	},
}

func formatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

func scoreGradeColor(score int) *color.Color {
	switch {
	case score >= 90:
		return color.New(color.FgGreen, color.Bold)
	case score >= 70:
		return color.New(color.FgYellow, color.Bold)
	case score >= 50:
		return color.New(color.FgYellow)
	default:
		return color.New(color.FgRed, color.Bold)
	}
}

func printScoreBreakdown(timeTaken, parTime time.Duration, hintsUsed, filesChanged, expectedFiles int) {
	color.New(color.FgWhite, color.Faint).Println("  Score breakdown:")

	// Base
	color.New(color.FgWhite, color.Faint).Println("    Base score:         100")

	// Hints
	if hintsUsed > 0 {
		color.New(color.FgRed, color.Faint).Printf("    Hints used (%d):     -%d\n", hintsUsed, hintsUsed*10)
	}

	// Extra files
	extraFiles := filesChanged - expectedFiles
	if extraFiles > 0 {
		color.New(color.FgRed, color.Faint).Printf("    Extra files (%d):    -%d\n", extraFiles, extraFiles*5)
	}

	// Overtime
	overtimeMinutes := int(timeTaken.Minutes()) - int(parTime.Minutes())
	if overtimeMinutes > 0 {
		color.New(color.FgRed, color.Faint).Printf("    Overtime (%dm):      -%d\n", overtimeMinutes, overtimeMinutes)
	}
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
