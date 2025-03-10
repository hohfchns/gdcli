# gdcli (Godot Project CLI)

Cli tool similar to npm but for Godot projects

![Icon](<icon.png>)

## Overview

`gdcli` is a command-line interface tool for starting and managing Godot projects more efficiently. It simplifies the process of managing different Godot versions and setting up project structures.

### Problems It Solves

For people experimenting with multiple Godot projects, managing Godot versions and starting repositories with the correct file structure can be cumbersome. `gdcli` addresses these issues with simple commands.

### Features

- **Initialize Projects**: With `gdcli init`, the tool prompts you with questions and creates a Godot project with a ready-to-use folder structure.
- **Manage Godot Versions**: After cloning a project from a repository, you no longer need to manually find the correct Godot version. Simply use:
  - `gdcli install`
  - `gdcli open`

## Getting Started

1. Download the latest release from the [Releases](https://github.com/IgorBayerl/gdcli/releases) page.
2. Place the downloaded file in a directory of your choice.
3. Add the directory to your PATH environment variable.

An installer will be provided in future updates to automate the PATH addition.

## Build from Source

To build the `gdcli` tool from source, run the following command:

```sh
go build -o bin/gdcli.exe main.go
```

This will compile the code and generate the `gdcli.exe` executable in the `bin` directory.

## How to Create a New Release

To create a new release, follow these steps:

1. **Create a new version tag** in the repository:
   ```sh
   git tag -a vX.Y.Z -m "Release vX.Y.Z"
   git push origin vX.Y.Z
   ```
   Replace `X.Y.Z` with the new version number.

2. **Trigger a release workflow manually** *(if needed)* from GitHub Actions:
   - Go to the [Actions tab](https://github.com/IgorBayerl/gdcli/actions).
   - Select **Release Workflow**.
   - Click **Run Workflow** and provide a tag version (e.g., `vX.Y.Z`).

## TODO

- [ ] Add build script
- [ ] Add custom scripts similar to npm options for Node.js
- [ ] Add support for global extensions, allowing extensions to be installed globally for use in every project
- [ ] Support more versions and variants, hopefully dynamic versions
  - [ ] For now, just Godot 4.3 and 4.3 Mono
  - [x] Support for Linux
  - [ ] In the future, the objective is to support all versions dynamically
- [ ] Support templates for starting projects
  - [ ] example: menu, platformer, 2d, 3d, etc.
- [ ] Add support for custom Godot versions
  - [ ] example: custom Godot Steam version

## How to Contribute

We aim to turn `gdcli` into a useful tool for Godot developers. Contributions are welcome!

1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Commit your changes.
4. Open a pull request with a detailed description of your changes.
