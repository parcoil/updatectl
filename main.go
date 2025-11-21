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
	Name         string            `yaml:"name"`
	Path         string            `yaml:"path"`
	Repo         string            `yaml:"repo"`
	Type         string            `yaml:"type"`
	BuildCommand string            `yaml:"buildCommand"`
	Image        string            `yaml:"image"`        // Docker image to pull (e.g., "ghcr.io/user/vite-app:main")
	Port         string            `yaml:"port"`         // Port mapping (e.g., "80:80" or "3000:80")
	Env          map[string]string `yaml:"env"`          // Environment variables
	ContainerName string           `yaml:"containerName"` // Optional custom container name
}

type Config struct {
	IntervalMinutes int       `yaml:"intervalMinutes"`
	Projects        []Project `yaml:"projects"`
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured projects",
	Run: func(cmd *cobra.Command, args []string) {
		config := loadConfig()
		if len(config.Projects) == 0 {
			fmt.Println("No projects configured.")
			return
		}
		fmt.Println("Configured projects:")
		for _, p := range config.Projects {
			if p.Type == "image" && p.Image != "" {
				portInfo := ""
				if p.Port != "" {
					portInfo = fmt.Sprintf(", port=%s", p.Port)
				}
				fmt.Printf("- %s (%s): image=%s%s\n", p.Name, p.Type, p.Image, portInfo)
			} else {
				fmt.Printf("- %s (%s): %s\n", p.Name, p.Type, p.Path)
			}
		}
	},
}

func main() {
	rootCmd := &cobra.Command{
		Use:     "updatectl",
		Version: version,
	}
	rootCmd.AddCommand(initCmd, watchCmd, buildCmd, listCmd, logsCmd)
	rootCmd.Execute()
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize updatectl configuration and daemon",
	Run: func(cmd *cobra.Command, args []string) {
		if runtime.GOOS != "windows" && os.Geteuid() != 0 {
			fmt.Println("Error: This command requires root privileges on Linux.")
			fmt.Println("Please run: sudo updatectl init")
			os.Exit(1)
		}

		var configDir, configPath string
		if runtime.GOOS == "windows" {
			configDir = filepath.Join(os.Getenv("USERPROFILE"), "updatectl")
		} else {
			configDir = "/etc/updatectl"
		}
		configPath = filepath.Join(configDir, "updatectl.yaml")

		if err := os.MkdirAll(configDir, 0755); err != nil {
			fmt.Printf("Failed to create config directory: %v\n", err)
			os.Exit(1)
		}

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			defaultConfig := []byte(`intervalMinutes: 10
projects:
  # Git-based project with Docker build
  - name: example-git
    path: /srv/example
    repo: https://github.com/user/example.git
    type: docker
    buildCommand: docker compose up -d --build

  # Vite app from GitHub Container Registry
  - name: vite-app
    type: image
    image: ghcr.io/user/vite-app:main
    port: "80:80"
    env:
      NODE_ENV: production
      API_URL: https://api.example.com

  # React app from Docker Hub
  - name: react-dashboard
    type: image
    image: user/react-dashboard:latest
    port: "3000:80"
    containerName: my-dashboard
`)
			if err := os.WriteFile(configPath, defaultConfig, 0644); err != nil {
				fmt.Printf("Failed to write config file: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Created config at", configPath)
		} else {
			fmt.Println("Config already exists at", configPath)
		}

		if runtime.GOOS == "windows" {
			taskName := "updatectl"
			configDir := filepath.Join(os.Getenv("USERPROFILE"), "updatectl")
			batScript := fmt.Sprintf(`@echo off
start "" /b "%s" watch
`, filepath.Join(configDir, "updatectl.exe"))
			batScriptPath := filepath.Join(configDir, "run_updatectl.bat")
			err := os.WriteFile(batScriptPath, []byte(batScript), 0644)
			if err != nil {
				fmt.Println("Failed to write batch wrapper script:", err)
				return
			}
			taskRun := batScriptPath
			createCmd := exec.Command(
				"schtasks",
				"/Create",
				"/TN", taskName,
				"/TR", taskRun,
				"/SC", "ONSTART",
				"/RL", "HIGHEST",
				"/F",
			)
			output, err := createCmd.CombinedOutput()
			if err != nil {
				fmt.Printf("Failed to create scheduled task: %v\nOutput: %s\n", err, output)
				return
			}
			fmt.Println("Created Windows Task Scheduler job for updatectl.")
			runCmd := exec.Command("schtasks", "/Run", "/TN", taskName)
			runOutput, runErr := runCmd.CombinedOutput()
			if runErr != nil {
				fmt.Printf("Failed to run scheduled task immediately: %v\nOutput: %s\n", runErr, runOutput)
			} else {
				fmt.Println("Scheduled task started immediately.")
			}
		} else {
			fmt.Print("Enter the user for the systemd service (default: root): ")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			user := strings.TrimSpace(scanner.Text())
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
			if err := os.WriteFile(servicePath, []byte(service), 0644); err != nil {
				fmt.Printf("Failed to write systemd service file: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Created systemd service file at", servicePath)

			reloadCmd := exec.Command("systemctl", "daemon-reload")
			if err := reloadCmd.Run(); err != nil {
				fmt.Printf("Failed to reload systemd daemon: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Reloaded systemd daemon")

			enableCmd := exec.Command("systemctl", "enable", "--now", "updatectl")
			if output, err := enableCmd.CombinedOutput(); err != nil {
				fmt.Printf("Failed to enable and start service: %v\nOutput: %s\n", err, output)
				os.Exit(1)
			}
			fmt.Println("Systemd service installed and started.")
			fmt.Println("\nCheck status with: sudo systemctl status updatectl")
			fmt.Println("View logs with: sudo journalctl -u updatectl -f")
		}
	},
}

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "View updatectl daemon logs",
	Long:  "View logs from the updatectl daemon service",
	Run: func(cmd *cobra.Command, args []string) {
		follow, _ := cmd.Flags().GetBool("follow")
		lines, _ := cmd.Flags().GetInt("lines")

		if runtime.GOOS == "windows" {
			fmt.Println("Viewing Windows Task Scheduler logs...")
			fmt.Println("You can view task history in Task Scheduler GUI or use:")
			fmt.Println("  Get-WinEvent -LogName Microsoft-Windows-TaskScheduler/Operational | Where-Object {$_.Message -like '*updatectl*'}")
			return
		}

		// Linux - use journalctl
		journalArgs := []string{"-u", "updatectl"}
		
		if follow {
			journalArgs = append(journalArgs, "-f")
		}
		
		if lines > 0 {
			journalArgs = append(journalArgs, "-n", fmt.Sprintf("%d", lines))
		}

		journalCmd := exec.Command("journalctl", journalArgs...)
		journalCmd.Stdout = os.Stdout
		journalCmd.Stderr = os.Stderr
		journalCmd.Stdin = os.Stdin

		if err := journalCmd.Run(); err != nil {
			fmt.Printf("Failed to view logs: %v\n", err)
			fmt.Println("Try: sudo journalctl -u updatectl -f")
			os.Exit(1)
		}
	},
}

func init() {
	logsCmd.Flags().BoolP("follow", "f", false, "Follow log output (live tail)")
	logsCmd.Flags().IntP("lines", "n", 50, "Number of log lines to show")
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

var buildCmd = &cobra.Command{
	Use:   "build [project-name]",
	Short: "Run build command for a specific project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		config := loadConfig()

		for _, p := range config.Projects {
			if p.Name == projectName {
				if p.BuildCommand == "" {
					fmt.Printf("No build command configured for project %s\n", projectName)
					return
				}

				fmt.Printf("Building project %s...\n", projectName)
				err := runBuildCommand(p.BuildCommand, p.Path)
				if err != nil {
					fmt.Printf("Build failed for %s: %v\n", projectName, err)
				} else {
					fmt.Printf("Build completed for %s\n", projectName)
				}
				return
			}
		}
		fmt.Printf("Project %s not found in configuration\n", projectName)
	},
}

func loadConfig() Config {
	var configPath string
	if runtime.GOOS == "windows" {
		configPath = filepath.Join(os.Getenv("USERPROFILE"), "updatectl", "updatectl.yaml")
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

func runBuildCommand(command, dir string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("bash", "-c", command)
	}
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func getImageDigest(image string) (string, error) {
	cmd := exec.Command("docker", "inspect", "--format={{index .RepoDigests 0}}", image)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func pullDockerImage(image string) error {
	fmt.Println("→ Pulling Docker image:", image)
	cmd := exec.Command("docker", "pull", image)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func restartDockerContainer(p Project) error {
	containerName := p.ContainerName
	if containerName == "" {
		containerName = p.Name
	}

	// Stop and remove old container if it exists
	fmt.Println("→ Stopping old container:", containerName)
	stopCmd := exec.Command("docker", "stop", containerName)
	stopCmd.Run() // Ignore error if container doesn't exist
	
	rmCmd := exec.Command("docker", "rm", containerName)
	rmCmd.Run() // Ignore error if container doesn't exist

	// Build docker run command
	args := []string{"run", "-d", "--name", containerName}
	
	// Add port mapping if specified
	if p.Port != "" {
		args = append(args, "-p", p.Port)
	}
	
	// Add environment variables
	for key, value := range p.Env {
		args = append(args, "-e", fmt.Sprintf("%s=%s", key, value))
	}
	
	// Add restart policy
	args = append(args, "--restart", "unless-stopped")
	
	// Add image
	args = append(args, p.Image)
	
	fmt.Println("→ Starting new container:", containerName)
	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
func getRemoteImageDigest(image string) (string, error) {
	cmd := exec.Command("docker", "manifest", "inspect", image)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, `"digest":`) {
			// Extract digest value: "digest": "sha256:xxxxx"
			parts := strings.Split(line, `"`)
			if len(parts) >= 4 {
				return parts[3], nil
			}
		}
	}
	return "", fmt.Errorf("could not parse digest from manifest")
}
func updateProject(p Project) {
if p.Type == "image" {
		if p.Image == "" {
			fmt.Println("✘ No image specified for project:", p.Name)
			return
		}

		containerName := p.ContainerName
		if containerName == "" {
			containerName = p.Name
		}

		checkCmd := exec.Command("docker", "inspect", "-f", "{{.State.Running}}", containerName)
		output, err := checkCmd.Output()
		containerRunning := err == nil && strings.TrimSpace(string(output)) == "true"

		currentDigest, _ := getImageDigest(p.Image)
		fmt.Println("→ Current local digest:", currentDigest)
		
		remoteDigest, err := getRemoteImageDigest(p.Image)
		if err != nil {
			fmt.Println("→ Could not check remote digest, forcing pull:", err)
			remoteDigest = ""
		} else {
			fmt.Println("→ Remote registry digest:", remoteDigest)
		}

		imageNeedsUpdate := currentDigest == "" || remoteDigest == "" || !strings.Contains(currentDigest, remoteDigest)

		if !imageNeedsUpdate && containerRunning {
			fmt.Println("● Image already up to date and container running:", p.Name)
			return
		}

		fmt.Println("→ Pulling latest image:", p.Image)
		if err := pullDockerImage(p.Image); err != nil {
			fmt.Println("✘ Failed to pull image:", err)
			return
		}

		imageUpdated := imageNeedsUpdate

		if !imageUpdated && containerRunning {
			fmt.Println("● Image already up to date and container running:", p.Name)
			return
		}

		if imageUpdated {
			fmt.Println("✓ New image version detected:", p.Name)
		} else {
			fmt.Println("→ Container not running, starting it:", p.Name)
		}

		if err := restartDockerContainer(p); err != nil {
			fmt.Println("✘ Failed to restart container:", err)
			return
		}
		fmt.Println("✓ Container started successfully")

		return
	}

	if _, err := os.Stat(p.Path); os.IsNotExist(err) {
		fmt.Println("✘ Path not found:", p.Path)
		return
	}

	fmt.Println("→ Pulling latest changes for", p.Name)
	gitPull := exec.Command("git", "-C", p.Path, "pull")
	output, err := gitPull.CombinedOutput()
	if err != nil {
		fmt.Println("✘ Git pull failed:", err)
		return
	}
	fmt.Print(string(output))

	if strings.Contains(string(output), "Already up to date.") {
		fmt.Println("● No new commits for", p.Name)
		return
	}

	if p.BuildCommand != "" {
		fmt.Println("→ Running build command for", p.Name)
		runBuildCommand(p.BuildCommand, p.Path)
	}

	switch p.Type {
	case "pm2":
		fmt.Println("→ Restarting PM2 process:", p.Name)
		cmd := exec.Command("pm2", "restart", p.Name)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	case "docker":
		// Build command already run above
	case "static":
		// No additional action needed
	default:
		fmt.Println("Unknown type:", p.Type)
	}
}