# ror - Fast, KDL-Configured Task Runner

## Overview
**ror** is a high-performance task runner designed for developers, DevOps engineers, and CI environments. Built with speed and configurability in mind, it utilizes **KDL (Key Document Language)** for its configuration, offering a cleaner, more human-readable alternative to YAML or JSON.

The project aims to provide granular control over task execution (fine-tuning) while remaining simple to set up. It is currently being built in Go and is self-hosting (dogfooding) its own configuration.

## Core Features (MVP)
*   **KDL Configuration:** Tasks are defined in a `ror.kdl` file using a user-friendly syntax.
*   **Speed:** Minimized overhead for task execution.
*   **Task Execution:** Simple command running capabilities.
*   **Dependencies:** (Planned) Ability to define execution order between tasks.

## Future Roadmap / TBD
The following features are identified as important but deferred for the initial MVP:
*   **Concurrency:** Parallel execution of independent tasks.
*   **Caching Strategy:** Intelligent avoidance of redundant work (strategy TBD).
*   **Watch Mode:** Re-running tasks on file changes.
*   **Argument Forwarding:** Passing CLI args to defined commands (e.g., `ror test -- -v`).
*   **Environment Management:** Loading `.env` files and defining per-task variables.
*   **Dry Run / Explain:** Visualizing the execution plan without running commands.
*   **Platform Specifics:** OS-dependent commands or skipping.
*   **Shell Selection:** Customizing the shell or interpreter (bash, python, etc.).
*   **Namespaces:** Grouping tasks (e.g., `db:migrate`).
*   **Hooks:** Pre/post execution scripts.
*   **Interactive Prompts:** User confirmation or input during execution.

## Target Audience
*   Software Developers
*   DevOps Engineers
*   CI/CD Pipelines
