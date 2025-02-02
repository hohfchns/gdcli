
**Description:**

Displays the version information of gdcli, including the build version, commit hash, and build time.

**Usage:**

```bash
gdcli version [--full]
```

**Parameters:**

- `--full` (optional): Displays detailed version information, including the commit hash and build time.

**Behavior:**

- Without the `--full` flag, displays the gdcli version.

- With the `--full` flag, displays the gdcli version along with the commit hash and build time.

**Example:**

```bash
$ gdcli version
1.0.0

$ gdcli version --full
gdcli version 1.0.0
Commit: abcdef123456
Build time: 2025-02-02T17:38:03Z
```

