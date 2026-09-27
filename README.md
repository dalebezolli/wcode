# wcode - Unf*ck project navigation

> For side-project collectors of all ages and genders.
wcode (which code) provides a simple way to find and navigate to the correct project directory.

## Features
✅ Fullscreen TUI display \
✅ Searching with Linear Search (fallback) \
✅ Searching with RipGrep \
✅ Project details view \
✅ Tmux integration to be ready once navigating in a new session

![wcode Showcase](./wcode_showcase.gif)

## 🌱 How to install
1. Clone the repo.
2. Run `make build` and put `bin/wcode` on your `PATH`.
3. Set `WCODE_PATHS` to the project root directories, separated by semicolons.
4. Add `eval "$(wcode init bash)"` to your `.bashrc`.

Make sure you have git and tmux installed for the best experience

## 🌷 How to use
For the current shell, run:
```sh
eval "$(wcode init bash)"
export WCODE_PATHS="/home/user/path/to/projects_root_directory" # if you have more than one projects dir split them with a semicolon (;)
```

Then run from anywhere:
```sh
wcode
```

While in the TUI:
- Arrow Up/Down or CTRL N/P: Move through the list
- Type anything: Search through the projects
- Enter: Opens the currently selected project

## 🧑‍🌾 How to contribute
Feel free to suggest any additions or changes by opening a pull request || an issue.
