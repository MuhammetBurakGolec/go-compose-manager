package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "projmanager",
	Short: "Project Manager CLI",
	Long:  `A simple CLI to manage Docker Compose projects`,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List projects",
	Run: func(cmd *cobra.Command, args []string) {
		projects, err := listProjects()
		if err != nil {
			fmt.Println("Error listing projects:", err)
			return
		}
		fmt.Println("Available projects:")
		for i, project := range projects {
			fmt.Printf("%d: %s\n", i+1, project)
		}
	},
}

func listProjects() ([]string, error) {
	files, err := filepath.Glob("docker-compose-files/*.yml")
	if err != nil {
		return nil, err
	}
	return files, nil
}

func selectProject() (string, error) {
	projects, err := listProjects()
	if err != nil {
		return "", err
	}

	prompt := promptui.Select{
		Label: "Select Project",
		Items: projects,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

func runDockerComposeCommand(project, action string) {
	command := fmt.Sprintf("docker-compose -f %s %s", project, action)
	fmt.Println("Running command:", command)

	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running %s on %s: %v\n", action, project, err)
	}
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop and remove containers, networks, images, and volumes",
	Run: func(cmd *cobra.Command, args []string) {
		project, err := selectProject()
		if err != nil {
			fmt.Println("Error selecting project:", err)
			return
		}
		runDockerComposeCommand(project, "down")
	},
}

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull service images",
	Run: func(cmd *cobra.Command, args []string) {
		project, err := selectProject()
		if err != nil {
			fmt.Println("Error selecting project:", err)
			return
		}
		runDockerComposeCommand(project, "pull")
	},
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Create and start containers",
	Run: func(cmd *cobra.Command, args []string) {
		project, err := selectProject()
		if err != nil {
			fmt.Println("Error selecting project:", err)
			return
		}
		runDockerComposeCommand(project, "up -d")
	},
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy containers",
	Run: func(cmd *cobra.Command, args []string) {
		project, err := selectProject()
		if err != nil {
			fmt.Println("Error selecting project:", err)
			return
		}
		runDockerComposeCommand(project, "down")
		time.Sleep(500 * time.Millisecond)
		runDockerComposeCommand(project, "pull")
		time.Sleep(500 * time.Millisecond)
		runDockerComposeCommand(project, "up -d")
	},
}

func main() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(downCmd)
	rootCmd.AddCommand(pullCmd)
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(deployCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
