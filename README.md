> [!CAUTION]
> This software is not verified, will have breaking changes and might have bugs.

# Ror

<p align="center">
  <img src="./docs/images/ror-logo.webp" alt="Ror Logo">
</p>

Ror is a [PFE](https://en.wikipedia.org/wiki/Not_invented_here) task runner.

[See Examples](#examples)

# Format Description

Ror uses [KDL v2](https://kdl.dev/) for its configuration file `ror.kdl`. The file defines a set of tasks that can be executed.

## Structure

### Tasks

Tasks are defined using the `task` node followed by the task name.

```kdl
task my-task {
    ...
}
```

You generally don't need quotes for task names unless they contain special characters (like spaces, brackets `()[]{}` or symbols like `=`, `,`, `;`, `\`, `/`, `<`, `>`) or look like a number, boolean, or null.

Common safe characters include letters, numbers, `-`, and `_`.

### Properties

Inside a task, you can define the following properties:

-   `description`: A string describing what the task does.
-   `cmd`: The command to run.
-   `depends`: A block listing dependencies using `on`.

### Dependencies

Dependencies define other tasks that must run successfully before the current task. They are specified within a `depends` block using `on`.

```kdl
task install {
  depends {
    on "build"
  }
  
  cmd "cp ./app /usr/local/bin/app"
}
```

### Command Execution

The `cmd` node specifies the shell command to execute. It supports variable substitution.

### Variables

Variables are defined using `where` nodes attached to the `cmd` node or nested within other `where` nodes. They are substituted using `%%variable_name%%`.

Variable name can only contain symbols `A-Z`, `a-z`, `0-9`, `_` and `-`.

#### Static Variables

```kdl
where my-var="value"
```

#### Dynamic Variables

Variables can be the result of a shell command execution by adding an `{ execute }` block.

```kdl
where current-date="date +%Y-%m-%d" { execute }
```

We can go even deeper:

```kdl
where current-date="date %%format%%" { 
  where format="+%Y-%m-%d"

  execute
}
```

### Examples

#### Basic Usage

A simple task to run a command.

```kdl
task run {
  description "Run the project"

  cmd "go run ./main.go"
}
```

#### Complex Build (Variables)

A build task that injects git commit hash and compilation time into the binary.

```kdl
task build {
  description "Build the application"

  cmd "go build -o out/app -ldflags='-X main.commit=%%commit%% -X main.date=%%date%%' ./cmd/app" {
    where commit="git describe --tags --always --dirty" { execute }
    where date="date +%Y-%m-%dT%H:%M:%SZ" { execute }
  }
}
```

#### Dependencies

A deploy task that ensures the application is tested and built first.

```kdl
task deploy {
  description "Deploy the application"
  depends {
    on "test"
    on "build"
  }
  
  cmd "./scripts/deploy.sh"
}
```
