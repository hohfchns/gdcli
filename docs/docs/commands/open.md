
**Description:**

Opens the current Godot project in the Godot editor. If the project is not initialized, it will first set up a new Godot project.

**Usage:**

```bash
gdcli open
```

**Behavior:**

- Checks for the existence of the Godot executable in the `dependencies` directory. If not found, prompts the user to run `gdcli install`.

- If a `project.godot` file does not exist, initializes a new Godot project.

- Launches the Godot editor with the current project.

**Example:**

```bash
$ gdcli open
Launching Godot editor...
```

