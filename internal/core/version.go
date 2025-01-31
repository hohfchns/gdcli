package core

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type GodotVersion struct {
	Version  string
	DotNet   bool
	URL      string
	Checksum string
}

var VersionManifest = []GodotVersion{
	{
		Version: "4.3.0",
		DotNet:  false,
		URL:     "https://github.com/godotengine/godot-builds/releases/download/4.3-stable/Godot_v4.3-stable_win64.exe.zip",
	},
	{
		Version: "4.3.0",
		DotNet:  true,
		URL:     "https://github.com/godotengine/godot-builds/releases/download/4.3-stable/Godot_v4.3-stable_mono_win64.zip",
	},
	// Add other versions and platforms as needed
}

func InstallGodotVersion(version string, dotnet bool) error {
	// Find version in manifest
	var dlUrl string
	for _, v := range VersionManifest {
		if v.Version == version && v.DotNet == dotnet {
			dlUrl = v.URL
			break
		}
	}

	if dlUrl == "" {
		return fmt.Errorf("version %s (%s) not found", version, dotnetStr(dotnet))
	}

	if err := os.MkdirAll("dependencies", 0755); err != nil {
		return err
	}

	zipName := filepath.Base(dlUrl)
	zipPath := filepath.Join("dependencies", zipName)
	
	fmt.Printf("Downloading %s...\n", zipName)
	if err := downloadFile(zipPath, dlUrl); err != nil {
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

    fmt.Printf("Successfully installed Godot %s (%s)\n", version, dotnetStr(dotnet))
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

    // Rename main executable
    if mainExe != "" {
        mainPath := filepath.Join(targetDir, mainExe)
        newMainPath := filepath.Join(targetDir, "godot.exe")
        if err := os.Rename(mainPath, newMainPath); err != nil {
            return fmt.Errorf("failed to rename main executable: %v", err)
        }
        fmt.Printf("Renamed %s -> godot.exe\n", mainExe)
    }

    // Rename console executable
    if consoleExe != "" {
        consolePath := filepath.Join(targetDir, consoleExe)
        newConsolePath := filepath.Join(targetDir, "godot_console.exe")
        if err := os.Rename(consolePath, newConsolePath); err != nil {
            return fmt.Errorf("failed to rename console executable: %v", err)
        }
        fmt.Printf("Renamed %s -> godot_console.exe\n", consoleExe)
    }

    // Enhanced file walking for cleanup
    return filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return nil
        }

        // Skip root directory
        if path == targetDir {
            return nil
        }

        // Delete if not an exe file
        if !info.IsDir() && !strings.EqualFold(filepath.Ext(info.Name()), ".exe") {
            if err := os.Remove(path); err != nil {
                fmt.Printf("Warning: Failed to remove file %s: %v\n", path, err)
            }
        }
        
        // Remove empty directories
        if info.IsDir() {
            if entries, _ := os.ReadDir(path); len(entries) == 0 {
                if err := os.Remove(path); err != nil {
                    fmt.Printf("Warning: Failed to remove directory %s: %v\n", path, err)
                }
            }
        }
        
        return nil
    })
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

func dotnetStr(dotnet bool) string {
	if dotnet {
		return "Mono"
	}
	return "Standard"
}