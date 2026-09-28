# wcode - Unf*ck project navigation

> For side-project collectors of all ages and genders.
wcode (which code) provides a simple way to find and navigate to the correct project directory.

## Features
✅ Fullscreen TUI display \
✅ Searching with regex and linear matching \
✅ Project details view \
✅ Tmux integration to be ready once navigating in a new session

![wcode Showcase](./wcode_showcase.gif)

## Requirements
- A Linux, macOS, or WSL environment with an interactive terminal.
- Supported shells are currently: Bash, Zsh, and Fish.
- Go 1.23 or newer and `make` to build wcode from source.
- `WCODE_PATHS` set to one or more existing directories, separated by semicolons. wcode treats each immediate child directory as a project.

And optionally:
- **Git**: to get project details directly from Git.
- **Tmux**: to immediately jump in a tmux session with possible configuration.

## 🌱 How to install

1. Either build the binary **or** download the built binary:

To clone the repository and build the binary run:
   ```sh
   git clone https://github.com/dalebezolli/wcode.git
   cd wcode
   make build
   ```

Otherwise, download the appropriate binary from the [releases section](https://github.com/dalebezolli/wcode/releases).

2. Set up your shell.

### Bash
Add these lines to `~/.bashrc`:
```sh
export PATH="/path/to/wcode/bin:$PATH"
export WCODE_PATHS="/path/to/projects1;/path/to/projects2"
eval "$(wcode init bash)"
```

Open a new Bash session or run `source ~/.bashrc`.

### Zsh
Add these lines to `~/.zshrc`:
```sh
export PATH="/path/to/wcode/bin:$PATH"
export WCODE_PATHS="/path/to/projects1;/path/to/projects2"
eval "$(wcode init zsh)"
```

Open a new Zsh session or run `source ~/.zshrc`.

### Fish
Add these lines to `~/.config/fish/config.fish`:
```fish
set -gx PATH /path/to/wcode/bin $PATH
set -gx WCODE_PATHS "/path/to/projects1;/path/to/projects2"
wcode init fish | source
```

Open a new Fish session or run `source ~/.config/fish/config.fish`.

3. And like this, you're ready to make your project browsing enjoyable again.

## 🌷 How to use
Now you can run from anywhere you want:
```sh
wcode
```

Controls:
- Arrow Up/Down or CTRL N/P: Move through the list
- Type anything: Search through the projects
- Enter: Opens the currently selected project

## 🔌 Integrations
### Git
When Git is installed, wcode adds repository information to each project's details: the current branch and latest commit. If the repository has an upstream branch configured, wcode also shows how many commits the local branch is ahead of or behind it.

wcode reads the local Git data and does not fetch from the remote. The ahead and behind counts therefore reflect the last time the repository's remote tracking information was updated.
If Git is not installed, wcode falls back to basic project details and leaves out the Git information.

### Tmux
After you select a project, wcode changes the shell's current directory to that project.
When tmux is available and wcode is not already running inside a tmux session, it looks for a session named after the project. wcode attaches to that session if it exists or creates a new one.
For a new session, wcode loads the project's `.tmux.conf` when that file exists.

If tmux is not installed, or wcode is already inside tmux, selecting a project still changes the directory and no new tmux session is opened.

Put a `.tmux.conf` file in a project's root to prepare your workspace when wcode creates a new tmux session. For example, this is the configuration I while developing wcode:

```tmux
rename-window main
send-keys "vim ." C-m

new-window -n terminal
send-keys "git fetch --verbose" C-m
split-window -h

next-window
```

This names the first window `main` and opens Vim in it, then creates a `terminal` window, runs `git fetch --verbose`, splits the window, and returns to `main`.

You don't have to do all your work inside the terminal to benefit from tmux. A project config can open VS Code with `code .` and use tmux to automate the boring, repetitive parts of your life as a developer, such as starting Docker containers, development servers, or running any other commands the project needs.

If you want to learn more feel free to start by reading the [tmux Getting Started guide](https://github.com/tmux/tmux/wiki/Getting-Started).

## 🧑‍🌾 How to contribute
Feel free to suggest any additions or changes by opening a pull request || an issue.

## 📜 License
Wcode is licensed under the [MIT License](LICENSE).
