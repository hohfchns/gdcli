package core

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type GodotVersion struct {
	DisplayName string // User-friendly name for selection
	Version     string // Base version number
	DotNet      bool   // Whether this is a Mono/.NET version
	URL         string // Download URL
	OS          string // Operating System (Architecture assumed as x64)
}

var VersionManifest = []GodotVersion{
	{
		DisplayName: "4.3.0 (Standard)",
		Version:     "4.3.0",
		DotNet:      false,
		URL:         "https://github.com/godotengine/godot-builds/releases/download/4.3-stable/Godot_v4.3-stable_win64.exe.zip",
		OS:          "windows",
	},
	{
		DisplayName: "4.3.0 (Mono)",
		Version:     "4.3.0",
		DotNet:      true,
		URL:         "https://github.com/godotengine/godot-builds/releases/download/4.3-stable/Godot_v4.3-stable_mono_win64.zip",
		OS:          "windows",
	},
	{
		DisplayName: "4.3.0 (Standard)",
		Version:     "4.3.0",
		DotNet:      false,
		URL:         "https://github.com/godotengine/godot-builds/releases/download/4.3-stable/Godot_v4.3-stable_linux.x86_64.zip",
		OS:          "linux",
	},
	{
		DisplayName: "4.3.0 (Mono)",
		Version:     "4.3.0",
		DotNet:      true,
		URL:         "https://github.com/godotengine/godot-builds/releases/download/4.3-stable/Godot_v4.3-stable_mono_linux_x86_64.zip",
		OS:          "linux",
	},
	{
		DisplayName: "4.4.0 (Standard)",
		Version:     "4.4.0",
		DotNet:      false,
		URL:         "https://github.com/godotengine/godot-builds/releases/download/4.4-stable/Godot_v4.4-stable_linux.x86_64.zip",
		OS:          "linux",
	},
	// Add more versions as needed
}

func GetVersionByIdentifier(identifier string) (GodotVersion, error) {
	var matches []GodotVersion

	currentOS := runtime.GOOS
	for _, v := range VersionManifest {
		if v.OS == currentOS && strings.EqualFold(v.DisplayName, identifier) || v.Version == identifier {
			return v, nil
		}
	}

	for _, v := range VersionManifest {
		if v.OS == currentOS && strings.Contains(strings.ToLower(v.DisplayName), strings.ToLower(identifier)) {
			matches = append(matches, v)
		}
	}

	switch len(matches) {
	case 1:
		return matches[0], nil
	case 0:
		return GodotVersion{}, fmt.Errorf("no versions found matching '%s'", identifier)
	default:
		var options []string
		for _, m := range matches {
			options = append(options, m.DisplayName)
		}
		return GodotVersion{}, fmt.Errorf("multiple matches found:\n%s", strings.Join(options, "\n"))
	}
}

func InstallGodotVersion(version GodotVersion) error {
	if version.URL == "" {
		return fmt.Errorf("no URL found for version %s", version.DisplayName)
	}

	if err := os.MkdirAll("dependencies", 0755); err != nil {
		return err
	}

	gdIgnorePath := filepath.Join("dependencies", ".gdignore")
	if _, err := os.Create(gdIgnorePath); err != nil {
		return err
	}

	zipName := filepath.Base(version.URL)
	zipPath := filepath.Join("dependencies", zipName)

	fmt.Printf("Downloading %s...\n", zipName)
	if err := downloadFile(zipPath, version.URL); err != nil {
		return err
	}

	tempDir := filepath.Join("dependencies", "temp_extract")
	defer os.RemoveAll(tempDir)

	fmt.Printf("Extracting %s...\n", zipName)
	if err := extractZip(zipPath, tempDir); err != nil {
		return err
	}

	exeDir, _, err := findExePath(tempDir)
	if err != nil {
		return fmt.Errorf("error locating executables: %v", err)
	}

	if err := moveFilesFromSubdir(exeDir, "dependencies"); err != nil {
		return fmt.Errorf("error moving files: %v", err)
	}

	if err := os.Remove(zipPath); err != nil {
		return fmt.Errorf("failed to remove zip: %v", err)
	}

	if err := renameExecutables("dependencies"); err != nil {
		return err
	}

	fmt.Printf("Successfully installed Godot %s\n", version.DisplayName)
	return nil
}

func getGodotExe(inDir string) (string, string, error) {
	var exeName string
	var consoleExeName string

	entries, err := os.ReadDir(inDir)
	if err != nil {
		return "", "", err
	}

	for _, info := range entries {
		switch runtime.GOOS {
		case "linux":
			// TODO handle other architectures?
			if strings.HasPrefix(info.Name(), "Godot_") && strings.EqualFold(filepath.Ext(info.Name()), ".x86_64") {
				exeName = info.Name()
				break
			}
		case "windows":
			name := info.Name()
			lowerName := strings.ToLower(name)
			if strings.Contains(lowerName, "_console") && strings.HasSuffix(lowerName, ".exe") {
				consoleExeName = name
			} else if strings.HasSuffix(lowerName, ".exe") && !strings.Contains(lowerName, "_console") {
				exeName = name
			}
		}
	}

	return exeName, consoleExeName, err
}

var errFound = fmt.Errorf("found")

func findExePath(searchDir string) (string, string, error) {
    var foundExeDir, foundExe string

    err := filepath.Walk(searchDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() {
            return nil
        }

        exe, _, err := getGodotExe(path)
        if err != nil {
            return err
        }
        if exe != "" {
            foundExeDir = path
            foundExe = exe
            // Return a sentinel error to stop the walk.
            return errFound
        }
        return nil
    })

    // If we broke out because we found the exe, ignore the sentinel error.
    if err != nil && err != errFound {
        return "", "", err
    }
    if foundExeDir == "" || foundExe == "" {
        return "", "", fmt.Errorf("no executables found in %s", searchDir)
    }
    return foundExeDir, foundExe, nil
}


func moveFilesFromSubdir(src, dest string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %v", err)
	}

	if err := os.MkdirAll(dest, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %v", err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if err := os.Rename(srcPath, destPath); err != nil {
			return fmt.Errorf("failed to move %s: %v", entry.Name(), err)
		}
	}

	return nil
}

func renameExecutables(targetDir string) error {
	newMainPath := filepath.Join(targetDir, "godot.exe")
	newConsolePath := filepath.Join(targetDir, "godot_console.exe")

	if _, err := os.Stat(newMainPath); err == nil {
		if err := os.Remove(newMainPath); err != nil {
			return fmt.Errorf("failed to remove existing godot.exe: %v", err)
		}
	}
	if _, err := os.Stat(newConsolePath); err == nil {
		if err := os.Remove(newConsolePath); err != nil {
			return fmt.Errorf("failed to remove existing godot_console.exe: %v", err)
		}
	}

	// Wait for the filesystem to settle.
	time.Sleep(1 * time.Second)

	var mainExe, consoleExe string

	// Retry mechanism in case of delays in file availability.
	retryCount := 3
	for i := 0; i < retryCount; i++ {
		// Search for the Godot executables
		var err error
		mainExe, consoleExe, err = getGodotExe(targetDir)

		if err != nil {
			return fmt.Errorf("failed to read directory: %v", err)
		}

		// If at least the main exe is found, break out of the loop.
		if mainExe != "" {
			break
		}

		fmt.Println("Retrying file detection...")
		time.Sleep(1 * time.Second)
	}

	// Copy main executable.
	if mainExe != "" {
		mainPath := filepath.Join(targetDir, mainExe)
		if err := copyFile(mainPath, newMainPath); err != nil {
			return fmt.Errorf("failed to copy main executable: %v", err)
		}
		if err := os.Remove(mainPath); err != nil {
			fmt.Printf("Warning: failed to remove original main executable %s: %v\n", mainExe, err)
		}
		fmt.Printf("Copied %s -> godot.exe\n", mainExe)

		// Probably the same with "darwin" AKA MacOS
		if runtime.GOOS == "linux" {
			fmt.Printf("Attempting to set executable permission...")
			os.Chmod(newMainPath, os.FileMode(0755))
		}
	} else {
		return fmt.Errorf("main executable not found after extraction")
	}

	// Copy console executable if found.
	if consoleExe != "" {
		consolePath := filepath.Join(targetDir, consoleExe)
		if err := copyFile(consolePath, newConsolePath); err != nil {
			return fmt.Errorf("failed to copy console executable: %v", err)
		}
		if err := os.Remove(consolePath); err != nil {
			fmt.Printf("Warning: failed to remove original console executable %s: %v\n", consoleExe, err)
		}
		fmt.Printf("Copied %s -> godot_console.exe\n", consoleExe)
	} else {
		fmt.Println("Warning: Console executable not found after extraction")
	}

	return nil
}

func downloadFile(path string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func extractZip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		// Skip directories
		if f.FileInfo().IsDir() {
			continue
		}

		// Create file path
		path := filepath.Join(dest, f.Name)

		// Create parent directories
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		// Extract file
		out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer out.Close()

		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		_, err = io.Copy(out, rc)
		if err != nil {
			return err
		}
	}
	return nil
}

// copyFile copies a single file from src to dst.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}
