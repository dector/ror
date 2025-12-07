# Project Index: ror

**Generated:** 2025-12-07
**Language:** Go 1.25.4
**License:** Not specified
**Status:** ⚠️ Alpha (not verified, breaking changes expected)

---

## 📋 Quick Summary

**ror** is a high-performance task runner built in Go that uses **KDL (Key Document Language)** for configuration. It provides a clean, human-readable alternative to YAML/JSON-based task runners like Make or Task. The project emphasizes speed, simplicity, and fine-grained control over task execution.

**Lines of Code:** ~1,065 Go LOC
**Configuration:** `ror.kdl` (KDL v2)
**Self-hosting:** Yes (dogfooding - project uses itself for build/run tasks)

---

## 📁 Project Structure

```
ror/
├── cmd/ror/              # CLI entry point
│   ├── main.go          # Main application logic, argument parsing
│   └── runner.go        # Task execution with dependency resolution
├── internal/            # Internal packages
│   ├── commands/        # Built-in commands (list, etc.)
│   ├── compat/          # Taskfile.yml compatibility layer
│   ├── config/          # KDL config parser
│   ├── env/             # Version/environment info
│   ├── io/              # I/O abstractions (files, shell, std)
│   ├── task/            # Task models and expander
│   │   ├── models.go    # Task, Project, WhereVariable definitions
│   │   └── expander.go  # Variable expansion engine
│   └── utils/           # Utility functions (slices, etc.)
├── docs/                # Documentation
│   └── images/          # Logo and assets
├── out/                 # Build output directory
├── ror.kdl              # Project's own task definitions
├── go.mod               # Go module definition
├── README.md            # User documentation
└── DESCRIPTION.md       # Project overview and roadmap
```

---

## 🚀 Entry Points

### CLI Entry
- **Path:** `cmd/ror/main.go`
- **Function:** `main()` → `execMain(io)`
- **Purpose:** Command-line interface for task execution

### Core Workflows
1. **Argument Parsing:** `parseArguments()` splits ror flags, task name, and task args
2. **Project Loading:** `buildProject()` reads and parses `ror.kdl`
3. **Task Execution:** `execute()` → `runTaskWithDependencies()` (in runner.go)

---

## 📦 Core Modules

### `cmd/ror/main.go`
- **Exports:** `main()`, `execMain()`
- **Purpose:** CLI entry point, argument parsing, command routing
- **Key Features:**
  - Verbose output flags (`-v`, `-vvv`)
  - Built-in commands: `version`, `help`
  - Taskfile.yml compatibility fallback
  - `--export-taskfile` flag for migration

### `cmd/ror/runner.go`
- **Purpose:** Task execution engine with dependency resolution
- **Key Features:**
  - Topological sorting of task dependencies
  - Circular dependency detection
  - Variable expansion for commands
  - Environment variable injection

### `internal/task/models.go`
- **Exports:** `Task`, `Project`, `WhereVariable`, `CommandTemplate`
- **Purpose:** Core data structures
- **Key Types:**
  - `Task`: Name, Command, Dependencies, EnvVars, Description
  - `Project`: Ordered map of tasks
  - `WhereVariable`: Variable definitions (static or command execution)
  - `CommandTemplate`: Command with `%%variable%%` placeholders

### `internal/task/expander.go`
- **Purpose:** Variable expansion engine
- **Features:**
  - Recursive variable substitution
  - Command execution for dynamic variables
  - Nested `where` block support

### `internal/config/config.go`
- **Exports:** `ParseProject()`, `CheckTaskfileExists()`
- **Purpose:** KDL configuration parser
- **Parses:**
  - Task definitions
  - Dependencies (`depends { on "task-name" }`)
  - Commands with `where` variables
  - Environment variables (`env KEY="value"`)

### `internal/io/`
- **Purpose:** I/O abstraction layer (testable, mockable)
- **Components:**
  - `io.go`: Interface definitions
  - `files.go`: File system operations
  - `std.go`: Standard I/O (stdin/stdout/stderr)
  - `shell.go`: Shell command execution

### `internal/compat/taskfile.go`
- **Purpose:** Taskfile.yml compatibility
- **Features:**
  - Export `ror.kdl` → `Taskfile.yml`
  - Fallback execution of existing Taskfile.yml

### `internal/commands/list.go`
- **Purpose:** `ror` (no args) lists all available tasks

---

## 🔧 Configuration Format

### File: `ror.kdl`

**Example:**
```kdl
task build {
  description "Compiles the application"

  cmd "go build -o out/ror %%flags%% ./cmd/ror" {
    where flags="-ldflags='%%ldflags%%'" {
      where ldflags="-X '%%pkg%%.commit=%%commit%%'" {
        where pkg="github.com/dector/ror/internal/env"
        where commit="git describe --tags --always --long --dirty" { execute }
      }
    }
  }
}

task install {
  description "Installs the in ~/.local/bin"
  depends {
    on "build"
  }

  cmd "cp ./out/ror ~/.local/bin/ror"
}

task test:env {
  description "Setup env for tasks"

  env GREET="Hello"
  env NAME="User"

  cmd "echo $GREET $NAME"
}
```

**Key Features:**
- **Variables:** `%%variable%%` syntax with static or dynamic (`{ execute }`) values
- **Nested Variables:** Recursive `where` blocks for complex substitutions
- **Dependencies:** `depends { on "task" }` for execution ordering
- **Environment:** `env KEY="value"` for task-specific env vars

---

## 🔗 Key Dependencies

| Dependency | Version | Purpose |
|------------|---------|---------|
| `github.com/dector/kdly` | v0.0.0-20251130231056 | KDL v2 parser |
| `github.com/fatih/color` | v1.18.0 | Terminal colors |
| `github.com/wk8/go-ordered-map/v2` | v2.1.8 | Ordered task map |
| `gopkg.in/yaml.v3` | v3.0.1 | Taskfile.yml export |

---

## 📚 Documentation

| File | Purpose |
|------|---------|
| `README.md` | User guide, CLI usage, KDL format reference, examples |
| `DESCRIPTION.md` | Project overview, MVP features, future roadmap |
| `docs/images/ror-logo.webp` | Project logo |

---

## 📝 Quick Start

### Installation
```bash
# Build from source
ror build

# Install to ~/.local/bin
ror install
```

### Basic Usage
```bash
# List all tasks
ror

# Run a task
ror build

# Run with verbose output
ror -v build

# Show version
ror version
```

### Create Tasks
Edit `ror.kdl`:
```kdl
task mytask {
  description "My custom task"
  cmd "echo Hello World"
}
```

---

## 🧪 Testing & Quality

- **Tests:** Not currently implemented (no `tests/` directory found)
- **Test Coverage:** 0% (TODO)
- **Quality Tools:** None configured yet

---

## 🎯 Current Features (MVP)

✅ **Implemented:**
- KDL v2 configuration parsing
- Task execution with dependencies
- Variable substitution (static & dynamic)
- Nested `where` blocks
- Environment variable support
- Taskfile.yml export/compatibility
- Colorful CLI output
- Ordered task execution

🚧 **Planned (from DESCRIPTION.md):**
- Concurrency (parallel task execution)
- Caching strategy
- Watch mode
- Argument forwarding (`ror test -- -v`)
- Dry run / explain mode
- Platform-specific commands
- Shell selection
- Namespaces
- Hooks (pre/post execution)
- Interactive prompts

---

## 🏗️ Architecture Notes

### Design Patterns
- **I/O Abstraction:** All file/shell/stdio operations abstracted for testability
- **Ordered Execution:** Uses ordered map to preserve task definition order
- **Dependency Resolution:** Topological sort with cycle detection
- **Variable Expansion:** Recursive template engine

### Key Algorithms
- **Dependency Resolution:** `cmd/ror/runner.go` - topological sort with visited/visiting sets
- **Variable Expansion:** `internal/task/expander.go` - recursive substitution with command execution

---

## 🔍 Notable Implementation Details

1. **Self-hosting:** The project uses its own `ror.kdl` for build/install tasks
2. **Compatibility Layer:** Falls back to Taskfile.yml if `ror.kdl` not found
3. **Version Injection:** Build process injects git commit and timestamp via `-ldflags`
4. **Ordered Tasks:** Preserves task definition order for deterministic listing
5. **Nested Variables:** Supports arbitrary depth of `where` block nesting

---

## 📊 Project Stats

- **Total Go Files:** 14
- **Total Lines:** ~1,065 LOC
- **Packages:** 8 internal packages
- **Dependencies:** 9 (4 direct, 5 indirect)
- **Git Commits:** 10+ recent commits
- **Latest Features:** Colorful output, stable task ordering

---

## 🚦 Project Status

**Maturity:** Alpha / Proof of Concept
**API Stability:** Unstable (breaking changes expected)
**Production Ready:** No

**From README:**
> ⚠️ This software is not verified, will have breaking changes and might have bugs.

---

## 💡 Use Cases

- **Local Development:** Replace Makefiles with readable KDL syntax
- **CI/CD Pipelines:** Fast task execution with dependency management
- **DevOps Workflows:** Scriptable build/deploy/test pipelines
- **Multi-step Builds:** Complex builds with variable injection

---

**Index Size:** ~3.5 KB
**Token Estimate:** ~3,000 tokens
**Token Savings vs Full Read:** ~94% (58K → 3K)
