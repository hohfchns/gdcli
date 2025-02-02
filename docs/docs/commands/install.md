
**Description:**

Installs a specified Godot engine version or uses the version defined in the `gdproj.json` configuration file. This is particularly useful when setting up a project cloned from a repository.

**Usage:**

```bash
gdcli install [version]
```

**Parameters:**

- `version` (optional): The specific Godot version to install (e.g., `4.3.0-mono`). If omitted, the version specified in `gdproj.json` will be used.

**Behavior:**

- If a version is provided as an argument, gdcli attempts to install that specific version.

- If no version is provided, gdcli checks the `gdproj.json` file for the required version.

- Downloads and installs the specified Godot version into the `dependencies` directory.

**Example:**

```bash
# Install a specific version
$ gdcli install 4.3.0-mono
Installing Godot 4.3.0-mono...
Successfully installed Godot 4.3.0-mono

# Install version from configuration
$ gdcli install
Installing Godot 4.3.0...
Successfully installed Godot 4.3.0
```
