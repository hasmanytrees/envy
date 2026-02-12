# envy

`envy` is a lightweight environment variable manager that automatically loads and unloads environment variables as you move between directories. It's designed to be simple, shell-native, and unobtrusive.

## How it Works

`envy` works by looking for `envy.sh` files in your current directory and all its parent directories. When you enter a directory:
1. It unloads any environment variables previously managed by `envy` (returning them to their original state).
2. It finds all `envy.sh` files from the root down to your current directory.
3. It sources these files in order, allowing nested configurations to override or append to parent ones.
4. It captures the changes and prepares an "undo" script for when you leave the directory.

This ensures that your environment is always tailored to the project you are currently working on.

## Installation

### 1. Build the binary

Ensure you have Go installed, then build and install the binary:

```bash
go install ./cmd/envy
```

### 2. Initialize your shell

Add the following to your `.zshrc` (currently Zsh is the only supported shell):

```bash
eval "$(envy init zsh)"
```

This command injects the necessary hooks into your shell to trigger `envy` whenever you change directories (`chpwd` hook).

## Usage

Create an `envy.sh` file in any directory where you want to manage environment variables:

```bash
# ~/projects/my-app/envy.sh
export API_KEY="secret-key-123"
export DEBUG=true
```

When you `cd` into `~/projects/my-app`, these variables will be automatically set. When you `cd` out, they will be unset or restored to their previous values.

### Shared configurations

Because `envy` searches parent directories, you can have shared configurations:

```text
~/projects/
  envy.sh          <-- Shared variables for all projects
  project-a/
    envy.sh        <-- Project-specific overrides
  project-b/
    envy.sh
```

## Commands

### `init SHELL`
Outputs the shell initialization script. Use this in your shell's configuration file (e.g., `.zshrc`) via `eval "$(envy init zsh)"`. It sets up the session and hooks.

### `gen`
The core logic of `envy`. It:
- Locates relevant `envy.sh` files.
- Calculates the difference between the current environment and the desired environment.
- Generates `load` and `undo` shell scripts in `~/.cache/envy/`.
- This command is usually called automatically by the shell hooks.

### `export`
Dumps the current environment variables to standard output. This is a helper command used internally by `gen` to capture the environment of a subshell.

## Configuration

`envy` uses the following environment variables (set automatically by `init`):
- `ENVY_SHELL`: The type of shell being used.
- `ENVY_SESSION_KEY`: A unique ID for the current shell session, used to manage temporary scripts in `~/.cache/envy/`.

