package main

import (
	"bufio"
	"fmt"
	"strings"

	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var version = "0.1.0"

type Project struct {
	Name         string `yaml:"name"`
	Path         string `yaml:"path"`
	Repo         string `yaml:"repo"`
	Type         string `yaml:"type"`
	BuildCommand string `yaml:"buildCommand"`
}

type Config struct {
	IntervalMinutes int       `yaml:"intervalMinutes"`
	Projects        []Project `yaml:"projects"`
}

func main() {
	rootCmd := &cobra.Command{
		Use:     "updatectl",
		Version: version,
	}
	rootCmd.AddCommand(initCmd, watchCmd)
	rootCmd.Execute()
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize updatectl configuration and daemon",
	Run: func(cmd *cobra.Command, args []string) {
		var configDir, configPath string

		if runtime.GOOS == "windows" {
			configDir = filepath.Join(os.Getenv("ProgramData"), "updatectl")
		} else {
			configDir = "/etc/updatectl"
		}
		configPath = filepath.Join(configDir, "updatectl.yaml")

		os.MkdirAll(configDir, 0755)

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			defaultConfig := []byte(`intervalMinutes: 10
projects:
  - name: example
    path: /srv/example
    repo: https://github.com/user/example.git
    type: docker
    buildCommand: docker compose up -d --build
`)
			os.WriteFile(configPath, defaultConfig, 0644)
			fmt.Println("Created config at", configPath)
		} else {
			fmt.Println("Config already exists at", configPath)
		}

		if runtime.GOOS == "windows" {
			taskCmd := `schtasks /Create /TN "updatectl" /TR "updatectl watch" /SC ONSTART /RL HIGHEST /F`
			exec.Command("cmd", "/C", taskCmd).Run()
			fmt.Println("Created Windows Task Scheduler job for updatectl.")
		} else {
			fmt.Print("Enter the user for the systemd service (default: root): ")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			user := scanner.Text()
			if user == "" {
				user = "root"
			}
			servicePath := "/etc/systemd/system/updatectl.service"
			service := fmt.Sprintf(`[Unit]
Description=Updatectl Daemon - Auto-update your projects
After=network.target

[Service]
ExecStart=/usr/local/bin/updatectl watch
WorkingDirectory=/etc/updatectl
Restart=always
User=%s

[Install]
WantedBy=multi-user.target
`, user)
			os.WriteFile(servicePath, []byte(service), 0644)
			exec.Command("systemctl", "daemon-reload").Run()
			exec.Command("systemctl", "enable", "--now", "updatectl").Run()
			fmt.Println("Systemd service installed and started.")
		}
	},
}

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Run updatectl daemon to auto-update projects",
	Run: func(cmd *cobra.Command, args []string) {
		config := loadConfig()
		fmt.Printf("Running updatectl every %d minutes...\n", config.IntervalMinutes)

		for {
			for _, p := range config.Projects {
				fmt.Println("Checking", p.Name)
				updateProject(p)
			}
			time.Sleep(time.Duration(config.IntervalMinutes) * time.Minute)
		}
	},
}

func loadConfig() Config {
	var configPath string
	if runtime.GOOS == "windows" {
		configPath = filepath.Join(os.Getenv("ProgramData"), "updatectl", "updatectl.yaml")
	} else {
		configPath = "/etc/updatectl/updatectl.yaml"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println("Failed to read config:", err)
		os.Exit(1)
	}

	var c Config
	yaml.Unmarshal(data, &c)
	return c
}

func updateProject(p Project) {
	if _, err := os.Stat(p.Path); os.IsNotExist(err) {
		fmt.Println("Path not found:", p.Path)
		return
	}

	fmt.Println("→ Pulling latest changes for", p.Name)
	gitPull := exec.Command("git", "-C", p.Path, "pull")
	output, err := gitPull.CombinedOutput()
	if err != nil {
		fmt.Println("Git pull failed:", err)
		return
	}
	fmt.Print(string(output))

	if strings.Contains(string(output), "Already up to date.") {
		fmt.Println("No new commits for", p.Name)
		return
	}

	switch p.Type {
	case "pm2":
		fmt.Println("→ Restarting PM2 process:", p.Name)
		cmd := exec.Command("pm2", "restart", p.Name)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	case "docker":
		fmt.Println("→ Rebuilding Docker container for", p.Name)
		cmd := exec.Command("bash", "-c", p.BuildCommand)
		cmd.Dir = p.Path
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	case "static":
		// No additional action needed, git pull is sufficient
	default:
		fmt.Println("Unknown type:", p.Type)
	}
}
