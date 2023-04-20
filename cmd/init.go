package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize ptpr.yaml interactively",
	Long: `Create a ptpr.yaml configuration file interactively.
This command will create the file if it doesn't exist, and add or update projects and Pivotal API tokens.`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath := os.Getenv("HOME") + "/.config/ptpr.yaml"

		_, err := os.Stat(configPath)
		if os.IsNotExist(err) {
			err = os.MkdirAll(filepath.Dir(configPath), os.ModePerm)
			if err != nil {
				panic(err)
			}

			// 空のファイルを作成
			file, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				panic(err)
			}
			defer file.Close()

			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Enter your Pivotal API Token (leave blank if already set): ")
			apiToken, _ := reader.ReadString('\n')
			apiToken = strings.TrimSpace(apiToken)

			if apiToken != "" {
				config.PIVOTAL_API_TOKEN = apiToken
			}

			currentDir, err := os.Getwd()
			if err != nil {
				panic(err)
			}

			projectRoot := filepath.Base(currentDir)

			fmt.Printf("Enter the Pivotal Project ID for the current project (%s): ", projectRoot)
			var projectID int
			_, err = fmt.Scanf("%d", &projectID)
			if err != nil {
				panic(err)
			}

			fmt.Printf("Do you want to use the default Pivotal API Token for this project? (yes/no): ")
			choice, _ := reader.ReadString('\n')
			choice = strings.TrimSpace(choice)

			project := Project{PIVOTAL_PROJECT_ID: projectID}

			if choice == "no" {
				fmt.Print("Enter a custom Pivotal API Token for this project: ")
				customToken, _ := reader.ReadString('\n')
				customToken = strings.TrimSpace(customToken)
				project.PIVOTAL_API_TOKEN = customToken
			}

			if config.Projects == nil {
				config.Projects = make(map[string]Project)
			}
			config.Projects[projectRoot] = project

			data, err := yaml.Marshal(config)
			if err != nil {
				panic(err)
			}

			_, err = file.Write(data)

			if err != nil {
				panic(err)
			}

			fmt.Println("ptpr.yaml has been updated.")
		} else if err != nil {
			panic(err)
		} else {
			fmt.Println("ptpr.yaml already exists. Please update it with your editor.")
		}
	},
}
