package core

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type GodotVersion struct {
    DisplayName string // User-friendly name for selection
    Version     string // Base version number
    DotNet      bool   // Whether this is a Mono/.NET version
    URL         string // Download URL
}

var VersionManifest = []GodotVersion{
    {
        DisplayName: "4.3.0 (Standard)",
        Version:     "4.3.0",
        DotNet:      false,
        URL:         "https://github.com/godotengine/godot-builds/releases/download/4.3-stable/Godot_v4.3-stable_win64.exe.zip",
    },
    {
        DisplayName: "4.3.0 (Mono)",
        Version:     "4.3.0",
        DotNet:      true,
        URL:         "https://github.com/godotengine/godot-builds/releases/download/4.3-stable/Godot_v4.3-stable_mono_win64.zip",
    },
    // Add more versions as needed
}

func GetVersionByIdentifier(identifier string) (GodotVersion, error) {
    var matches []GodotVersion
    
    // First try exact matches
    for _, v := range VersionManifest {
        if strings.EqualFold(v.DisplayName, identifier) || v.Version == identifier {
            return v, nil
        }
    }
    
    // Then try partial matches
    for _, v := range VersionManifest {
        if strings.Contains(strings.ToLower(v.DisplayName), strings.ToLower(identifier)) {
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
    // Directly use the version struct fields
    if version.URL == "" {
        return fmt.Errorf("no URL found for version %s", version.DisplayName)
    }

    // Create dependencies directory
    if err := os.MkdirAll("dependencies", 0755); err != nil {
        return err
    }

    // Rest of the installation logic using version.URL directly
    zipName := filepath.Base(version.URL)
    zipPath := filepath.Join("dependencies", zipName)
    
    fmt.Printf("Downloading %s...\n", zipName)
    if err := downloadFile(zipPath, version.URL); err != nil {
        return err
    }

    tempDir := filepath.Join("dependencies", "temp_extract")
    defer os.RemoveAll(tempDir) // Clean up temp directory

    fmt.Printf("Extracting %s...\n", zipName)
    if err := extractZip(zipPath, tempDir); err != nil {
        return err
    }

    if err := moveFilesFromSubdir(tempDir, "dependencies"); err != nil {
        return err
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

func moveFilesFromSubdir(src, dest string) error {
    entries, err := os.ReadDir(src)
    if err != nil {
        return fmt.Errorf("failed to read temp directory: %v", err)
    }

    // Handle nested directory structure
    for _, entry := range entries {
        srcPath := filepath.Join(src, entry.Name())
        destPath := filepath.Join(dest, entry.Name())

        if entry.IsDir() {
            // Recursively move directory contents
            err := filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
                if err != nil {
                    return err
                }

                relPath, _ := filepath.Rel(srcPath, path)
                newPath := filepath.Join(dest, entry.Name(), relPath)

                if info.IsDir() {
                    return os.MkdirAll(newPath, info.Mode())
                }
                
                return os.Rename(path, newPath)
            })
            
            if err != nil {
                return fmt.Errorf("failed to move directory: %v", err)
            }
        } else {
            // Move file directly
            if err := os.Rename(srcPath, destPath); err != nil {
                return fmt.Errorf("failed to move file: %v", err)
            }
        }
    }
    return nil
}


func renameExecutables(targetDir string) error {
	files, err := os.ReadDir(targetDir)
	if err != nil {
		return err
	}

	var mainExe, consoleExe string

	// Identify files (case-insensitive check)
	for _, f := range files {
		name := f.Name()
		lowerName := strings.ToLower(name)

		if strings.Contains(lowerName, "_console") && strings.HasSuffix(lowerName, ".exe") {
			consoleExe = name
		} else if strings.HasSuffix(lowerName, ".exe") && !strings.Contains(lowerName, "_console") {
			mainExe = name
		}
	}

	// Define new file paths
	newMainPath := filepath.Join(targetDir, "godot.exe")
	newConsolePath := filepath.Join(targetDir, "godot_console.exe")

	// Remove existing files if they exist
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

	// Wait for the filesystem to release deleted files
	time.Sleep(1 * time.Second)

	// Re-scan the directory after deletion
	files, err = os.ReadDir(targetDir)
	if err != nil {
		return fmt.Errorf("failed to read directory after deletion: %v", err)
	}

	mainExe, consoleExe = "", ""
	for _, f := range files {
		name := f.Name()
		lowerName := strings.ToLower(name)

		if strings.Contains(lowerName, "_console") && strings.HasSuffix(lowerName, ".exe") {
			consoleExe = name
		} else if strings.HasSuffix(lowerName, ".exe") && !strings.Contains(lowerName, "_console") {
			mainExe = name
		}
	}

	// Retry mechanism in case of delays in file availability
	retryCount := 3
	for i := 0; i < retryCount; i++ {
		if mainExe == "" || consoleExe == "" {
			fmt.Println("Retrying file detection...")
			time.Sleep(1 * time.Second)

			// Re-scan the directory
			files, err = os.ReadDir(targetDir)
			if err != nil {
				return fmt.Errorf("failed to read directory: %v", err)
			}

			mainExe, consoleExe = "", ""
			for _, f := range files {
				name := f.Name()
				lowerName := strings.ToLower(name)

				if strings.Contains(lowerName, "_console") && strings.HasSuffix(lowerName, ".exe") {
					consoleExe = name
				} else if strings.HasSuffix(lowerName, ".exe") && !strings.Contains(lowerName, "_console") {
					mainExe = name
				}
			}
		} else {
			break
		}
	}

	// Rename main executable
	if mainExe != "" {
		mainPath := filepath.Join(targetDir, mainExe)
		if err := os.Rename(mainPath, newMainPath); err != nil {
			return fmt.Errorf("failed to rename main executable: %v", err)
		}
		fmt.Printf("Renamed %s -> godot.exe\n", mainExe)
	} else {
		return fmt.Errorf("main executable not found after extraction")
	}

	// Rename console executable
	if consoleExe != "" {
		consolePath := filepath.Join(targetDir, consoleExe)
		if err := os.Rename(consolePath, newConsolePath); err != nil {
			return fmt.Errorf("failed to rename console executable: %v", err)
		}
		fmt.Printf("Renamed %s -> godot_console.exe\n", consoleExe)
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