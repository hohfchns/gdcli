package config

import (
	"encoding/json"
	"os"
)

type GodotConfig struct {
    EngineVersion string `json:"engine_version"`
    ProjectName   string `json:"project_name"`
    IsDotNet      bool   `json:"is_dotnet"`
}

func CreateConfig(version, name string, dotnet bool) error {
    cfg := GodotConfig{
        EngineVersion: version,
        ProjectName:   name,
        IsDotNet:      dotnet,
    }

    data, err := json.MarshalIndent(cfg, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile("godot.json", data, 0644)
}

func LoadConfig() (*GodotConfig, error) {
    data, err := os.ReadFile("godot.json")
    if err != nil {
        return nil, err
    }

    var cfg GodotConfig
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }
    
    return &cfg, nil
}