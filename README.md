# wcode - Unf*ck your project folder

> For side-project collectors of all ages and genders.
wcode (which code) provides a simple way to find and navigate to the correct project directory.

## Features
✅ Fullscreen TUI display \
✅ Searching with Linear Search (fallback) \
✅ Searching with RipGrep \
✅ Project details view \
✅ Tmux integration to be ready once navigating in a new session \
🟥 Handle edge-cases:
- no git
- add blinking cursor???

🟥 Improve rendering \
🟥 Remove ripgrep dependency

![wcode Showcase](./wcode_showcase.gif)

## 🌱 How to install
1. Clone the repo.
2. Set variable WCODE_PATHS with all the paths (space separated) the tool will look for projects
3. Profit?!?

Make sure you have tmux, ripgrep & git installed for the best experience

## 🌷 How to use
Simply run the following command
```sh
source ./wcode.sh
export WCODE_PATHS="/home/user/path/to/projects_root_directory" # if you have more than one projects dir split them with a semicolon (;)
```
*maybe alias the above in your `.bashrc` or anyplace they'd be run on shell startup.*

and then run from anywhere
```sh
wcode
```

While in the TUI:
- Arrow Up/Down or CTRL N/P: Move through the list
- Type anything: Search through the projects
- Enter: Opens the currently selected project

## 🧑‍🌾 How to contribute
Feel free to suggest any additions or changes by opening a pull request || an issue.
