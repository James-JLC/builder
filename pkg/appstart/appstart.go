// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package appstart contains code related to creating the app_start.json config file.
package appstart

import (
	"encoding/json"
	"fmt"

	gcp "github.com/GoogleCloudPlatform/buildpacks/pkg/gcpbuildpack"
)

const (
	// ConfigDir is the location relative to the user's application where the configFile lives.
	ConfigDir  = ".googleconfig"
	configFile = ConfigDir + "/app_start.json"
)

// Config holds the parameters to pass into app_start.json
type Config struct {
	Runtime        string     `json:"runtime,omitempty"`
	Entrypoint     Entrypoint `json:"entrypoint"`
	MainExecutable string     `json:"main,omitempty"`
}

// Entrypoint contains the command to start the application
type Entrypoint struct {
	Type    string `json:"type"`
	Command string `json:"command"`
	WorkDir string `json:"workdir"`
}

// EntrypointType represents how the entrypoint command was created.
type EntrypointType int

// This enum and its String() method need to be updated together.
// Generated represents a command generated by a language specific command generator.
// Default represents the default start command.
const (
	EntrypointDefault EntrypointType = iota
	EntrypointGenerated
	EntrypointUser
)

func (et EntrypointType) String() string {
	return [...]string{"Default", "Generated", "User"}[et]
}

// EntrypointGenerator is a function that returns an entrypoint.
type EntrypointGenerator func(*gcp.Context) (*Entrypoint, error)

// Write writes the config to the default config file in a new layer.
func (c Config) Write(ctx *gcp.Context) error {
	l := ctx.Layer("config", gcp.LaunchLayer)

	cb, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("marshalling JSON: %w", err)
	}

	ctx.RemoveAll(ConfigDir)
	ctx.Symlink(l.Path, ConfigDir)
	ctx.WriteFile(configFile, cb, 0444)

	return nil
}
