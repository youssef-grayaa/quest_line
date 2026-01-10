package domain

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const dataFile = "quests.json"

// SaveProjects saves the projects to a JSON file
func SaveProjects(projects []Project) error {
	data, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Dir(dataFile)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return os.WriteFile(dataFile, data, 0644)
}

// LoadProjects loads the projects from a JSON file
func LoadProjects() ([]Project, error) {
	data, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Project{}, nil
		}
		return nil, err
	}
	var projects []Project
	err = json.Unmarshal(data, &projects)
	return projects, err
}
