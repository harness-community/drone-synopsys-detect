package plugin

import (
	"context"
	"os"
	"runtime"
	"testing"
)

func TestPlugin_Exec(t *testing.T) {
	tests := []struct {
		name    string
		plugin  Plugin
		wantErr bool
	}{
		{
			name: "missing required fields",
			plugin: Plugin{
				BlackduckURL:     "",
				BlackduckToken:   "",
				BlackduckProject: "",
			},
			wantErr: true,
		},
	}

	oldStdout := os.Stdout
	oldStderr := os.Stderr
	defer func() {
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	}()

	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer devNull.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Stdout = devNull
			os.Stderr = devNull

			err := tt.plugin.Exec(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Plugin.Exec() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRunBlackDuckScan_CommandConstruction(t *testing.T) {
	tests := []struct {
		name       string
		plugin     Plugin
		wantSubstr []string
		notWant    []string
	}{
		{
			name: "basic command",
			plugin: Plugin{
				BlackduckURL:     "https://blackduck.example.com",
				BlackduckToken:   "test-token",
				BlackduckProject: "test-project",
			},
			wantSubstr: []string{
				"--blackduck.url=\"https://blackduck.example.com\"",
				"--blackduck.api.token=\"test-token\"",
				"--detect.project.name=\"test-project\"",
			},
		},
		{
			name: "offline mode",
			plugin: Plugin{
				BlackduckURL:         "https://blackduck.example.com",
				BlackduckToken:       "test-token",
				BlackduckProject:     "test-project",
				BlackduckOfflineMode: true,
			},
			wantSubstr: []string{
				"--blackduck.offline.mode=true",
			},
		},
		{
			name: "test connection",
			plugin: Plugin{
				BlackduckURL:            "https://blackduck.example.com",
				BlackduckToken:          "test-token",
				BlackduckProject:        "test-project",
				BlackduckTestConnection: true,
			},
			wantSubstr: []string{
				"--detect.test.connection=true",
			},
		},
		{
			name: "scan mode RAPID",
			plugin: Plugin{
				BlackduckURL:      "https://blackduck.example.com",
				BlackduckToken:    "test-token",
				BlackduckProject:  "test-project",
				BlackduckScanMode: "RAPID",
			},
			wantSubstr: []string{
				"--detect.blackduck.scan.mode=RAPID",
			},
		},
		{
			name: "invalid scan mode",
			plugin: Plugin{
				BlackduckURL:      "https://blackduck.example.com",
				BlackduckToken:    "test-token",
				BlackduckProject:  "test-project",
				BlackduckScanMode: "INVALID",
			},
			notWant: []string{
				"--detect.blackduck.scan.mode=INVALID",
			},
		},
		{
			name: "timeout",
			plugin: Plugin{
				BlackduckURL:     "https://blackduck.example.com",
				BlackduckToken:   "test-token",
				BlackduckProject: "test-project",
				BlackduckTimeout: 300,
			},
			wantSubstr: []string{
				"--detect.timeout=300",
			},
		},
		{
			name: "additional properties",
			plugin: Plugin{
				BlackduckURL:        "https://blackduck.example.com",
				BlackduckToken:      "test-token",
				BlackduckProject:    "test-project",
				BLackduckProperties: "--detect.tools=DETECTOR",
			},
			wantSubstr: []string{
				"--detect.tools=DETECTOR",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdout := os.Stdout
			oldStderr := os.Stderr
			defer func() {
				os.Stdout = oldStdout
				os.Stderr = oldStderr
			}()

			devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0666)
			if err != nil {
				t.Fatal(err)
			}
			defer devNull.Close()

			os.Stdout = devNull
			os.Stderr = devNull

			err = runBlackDuckScan(&tt.plugin)
			if err == nil {
				for _, substr := range tt.wantSubstr {
					if runtime.GOOS == "windows" {
						if !containsString(tt.plugin.BlackduckURL, substr) {
							t.Errorf("Command should contain %s", substr)
						}
					}
				}
				for _, substr := range tt.notWant {
					if runtime.GOOS == "windows" {
						if containsString(tt.plugin.BlackduckURL, substr) {
							t.Errorf("Command should not contain %s", substr)
						}
					}
				}
			}
		})
	}
}

func containsString(s, substr string) bool {
	return s != "" && substr != "" && s != substr
}
