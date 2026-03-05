package runner

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Problem represents a debugging problem loaded from problem.json.
type Problem struct {
	Name             string   `json:"name"`
	Title            string   `json:"title"`
	Difficulty       string   `json:"difficulty"`
	Language         string   `json:"language"`
	ParTimeMinutes   int      `json:"par_time_minutes"`
	ExpectedFiles    int      `json:"expected_files_to_fix"`
	Ticket           string   `json:"ticket"`
	Hints            []string `json:"hints"`
}

// TestResult holds the result of running a test suite.
type TestResult struct {
	Passed    int
	Failed    int
	Total     int
	Output    string
	AllPassed bool
}

// LoadProblem reads and parses problem.json from the given directory.
func LoadProblem(problemDir string) (*Problem, error) {
	data, err := os.ReadFile(filepath.Join(problemDir, "problem.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to read problem.json: %w", err)
	}

	var p Problem
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("failed to parse problem.json: %w", err)
	}

	return &p, nil
}

// RunTests executes the test suite in the given workspace directory.
// Supports React (npm test), Python (pytest), and more.
func RunTests(workspaceDir string, language string) (*TestResult, error) {
	var cmd *exec.Cmd

	switch {
	case strings.Contains(language, "python"):
		cmd = exec.Command("uv", "run", "pytest", "tests/", "-v")
	case strings.Contains(language, "react"), strings.Contains(language, "typescript"), strings.Contains(language, "javascript"):
		cmd = exec.Command("npm", "test", "--", "--watchAll=false")
	default:
		return nil, fmt.Errorf("unsupported language: %s", language)
	}

	cmd.Dir = workspaceDir
	cmd.Env = append(os.Environ(), "CI=true")

	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	result := parseTestOutput(outputStr, language)

	// If the command failed but we got test results, it's a test failure not an execution error
	if err != nil && result.Total == 0 {
		return nil, fmt.Errorf("failed to run tests: %w\nOutput:\n%s", err, outputStr)
	}

	return result, nil
}

// parseTestOutput parses test output for pass/fail counts.
func parseTestOutput(output string, language string) *TestResult {
	result := &TestResult{
		Output: output,
	}

	if strings.Contains(language, "python") {
		// Pytest output pattern: "3 passed, 1 failed" or "3 passed"
		passedRe := regexp.MustCompile(`(\d+)\s+passed`)
		failedRe := regexp.MustCompile(`(\d+)\s+failed`)

		if m := passedRe.FindStringSubmatch(output); len(m) >= 2 {
			result.Passed, _ = strconv.Atoi(m[1])
		}
		if m := failedRe.FindStringSubmatch(output); len(m) >= 2 {
			result.Failed, _ = strconv.Atoi(m[1])
		}
		result.Total = result.Passed + result.Failed
		result.AllPassed = result.Failed == 0 && result.Passed > 0
		return result
	}

	// Jest output patterns:
	// Tests:       3 passed, 3 total
	// Tests:       1 failed, 2 passed, 3 total
	testsRe := regexp.MustCompile(`Tests:\s+(?:(\d+)\s+failed,\s+)?(\d+)\s+passed,\s+(\d+)\s+total`)
	matches := testsRe.FindStringSubmatch(output)

	if len(matches) >= 4 {
		if matches[1] != "" {
			result.Failed, _ = strconv.Atoi(matches[1])
		}
		result.Passed, _ = strconv.Atoi(matches[2])
		result.Total, _ = strconv.Atoi(matches[3])
	} else {
		// Fallback: count PASS/FAIL lines
		scanner := bufio.NewScanner(strings.NewReader(output))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "✓") || strings.Contains(line, "✕") || strings.Contains(line, "●") {
				if strings.Contains(line, "✓") || strings.Contains(line, "pass") {
					result.Passed++
				}
				if strings.Contains(line, "✕") || strings.Contains(line, "fail") {
					result.Failed++
				}
			}
		}
		result.Total = result.Passed + result.Failed
	}

	result.AllPassed = result.Failed == 0 && result.Passed > 0
	return result
}

// CountChangedFiles counts how many files differ between the original problem
// and the workspace. It compares against the git initial commit if available,
// otherwise does a simple directory comparison.
func CountChangedFiles(workspaceDir string) (int, error) {
	// Try git diff first (we init a git repo when starting)
	cmd := exec.Command("git", "diff", "--name-only", "HEAD")
	cmd.Dir = workspaceDir
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		count := 0
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				count++
			}
		}
		return count, nil
	}

	// Fallback: count all source files (not ideal, but safe)
	count := 0
	err = filepath.Walk(filepath.Join(workspaceDir, "src"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})

	return count, err
}


