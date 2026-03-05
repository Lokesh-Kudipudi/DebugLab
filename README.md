# 🔬 DebugLab (`dblab`)

DebugLab is a CLI tool built in Go for practicing real-world debugging skills. It provides a local, timed, and scored environment where you can pick a broken project, fix the bugs, and verify your solution against a test suite.

Everything runs locally on your machine—no Docker or internet connection required after setup. Session data (times, scores, and status) is stored locally in an SQLite database.

## Features

- **Local Workspaces**: Copies problems to a safe, isolated `dblab-workspace` directory.
- **Timed & Scored**: Tracks how long you take to solve a problem and deducts points for overtime, extra files changed, or using hints.
- **Multiple Languages**: Supports React (JavaScript/TypeScript) and Python problems out of the box.
- **Built-in Verification**: Automatically runs Jest or Pytest to verify your fixes.
- **Git Integration**: Automatically initializes a git repository in your workspace to track the files you change.

## Installation

### Prerequisites

- [Go](https://go.dev/dl/) (1.20+ recommended)
- Node.js & npm (for React problems)
- Python 3 & pytest (for Python problems)

### Build the CLI

Clone the repository and build the binary:

```bash
git clone <repository-url>
cd DebugLab
go build -o dblab .
```

You can move the `dblab` binary to somewhere in your `$PATH` (e.g., `/usr/local/bin`) for easier access.

## Usage Workflow

The typical workflow involves choosing a problem, starting a session, fixing the code, and verifying it.

### 1. List Available Problems

```bash
./dblab list
```
Displays a table of all available problems, their difficulty, language (e.g., `react-javascript`, `python`), and your current status (`unsolved` or `✓ solved`).

### 2. Start a Debugging Session

```bash
./dblab start <problem-name>
# Example: ./dblab start react-syntax-error
```
This command:
- Copies the problem files into `./dblab-workspace/<problem-name>/` in your current directory.
- Initializes a fresh git repository in the workspace.
- Starts a timer for your session.
- Prints the "Bug Ticket" explaining the issue you need to fix.

### 3. Fix the Bugs

Open the workspace in your favorite editor and investigate the bugs described in the ticket.

*Note: For React/Node problems, remember to install dependencies first!*
```bash
cd dblab-workspace/<problem-name>
npm install
```

### 4. Verify Your Solution

Once you think you've fixed the issue, run the verification command:

```bash
./dblab verify <problem-name>
```
This runs the internal test suite (Jest or Pytest).
- **If tests pass:** The timer stops, your score is calculated (out of 100), and the session is marked as completed.
- **If tests fail:** You'll see which tests failed. Fix them and run `verify` again.

### 5. Check Your Stats

```bash
./dblab stats
```
Displays your aggregate statistics: problems attempted, problems solved, average/best scores, total time spent, and a log of recent sessions.

## Scoring System

Your starting score is **100**. Points are deducted based on:

- **Overtime limitations**: `-1 point` per minute spent over the "par time" for the problem.
- **Hints used**: `-10 points` per hint requested (future feature limitation).
- **Extra files changed**: `-5 points` for modifying files unnecessarily (beyond the expected fix scope).

## Available Problems

Currently, DebugLab includes 6 problems categorized by language:

### JavaScript (`react-javascript`)
- `react-syntax-error` (Easy) - Fix broken component syntax and rendering issues.
- `react-event-handler` (Easy) - Fix broken form submissions and state mutations in a Todo app.

### TypeScript (`react-typescript`)
- `react-cart-total` (Easy) - Fix state mutations returning stale data in a shopping cart.
- `react-async-fetch` (Medium) - Debug a custom React hook with infinite loading and stale closures.

### Python (`python`)
- `python-sort-filter` (Easy) - Fix broken list sorting, filtering, and pagination logic.
- `python-api-parser` (Medium) - Debug an API parser with mutable default arguments and swallowed exceptions.

## Project Structure

```text
DebugLab/
├── main.go                     # Entry point
├── cmd/                        # CLI Commands (Cobra)
│   ├── root, list, start, verify, stats, help
├── internal/
│   ├── store/                  # SQLite DB interactions for sessions
│   ├── scoring/                # Scoring logic calculator
│   └── runner/                 # Test runner (npm/pytest) and problem loader
├── problems/                   # Challenge definitions
│   ├── javascript/
│   ├── typescript/
│   └── python/
└── ...
```

## Data Storage

All your session data, timings, and scores are stored locally in an SQLite database located right next to your active workspaces:
`./dblab-workspace/dblab.db`
