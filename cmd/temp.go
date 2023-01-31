/*
Copyright © 2023 Teppei Taguchi tetzng.tt@gmail.com
*/
package cmd

import (
	"fmt"

	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type Project struct {
	PIVOTAL_PROJECT_ID int    `yaml:"PIVOTAL_PROJECT_ID"`
	Root               string `yaml:"Root"`
}

var config struct {
	PIVOTAL_API_TOKEN string             `yaml:"PIVOTAL_API_TOKEN"`
	Projects          map[string]Project `yaml:"Projects"`
}

func init() {
	rootCmd.AddCommand(tempCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tempCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tempCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	data, err := ioutil.ReadFile(os.Getenv("HOME") + "/.config/pivo_conf.yml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf("error: %v", err)
		panic(err)
	}
}

// tempCmd represents the temp command
var tempCmd = &cobra.Command{
	Use:   "temp",
	Short: "Print PR template",
	Long: `Print Pull Request template from your ~/.config/pivo_conf.yml
and your current branch name.`,
	Run: func(cmd *cobra.Command, args []string) {
		projectRoot := getProjectRoot()
		storyID, err := getStoryID()
		if err != nil {
			panic(err)
		}
		client := &http.Client{}

		req, err := http.NewRequest("GET", fmt.Sprintf("https://www.pivotaltracker.com/services/v5/projects/%d/stories/%s",
			config.Projects[projectRoot].PIVOTAL_PROJECT_ID, storyID), nil)
		if err != nil {
			panic(err)
		}
		req.Header.Set("X-TrackerToken", config.PIVOTAL_API_TOKEN)

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		var story map[string]interface{}
		err = yaml.Unmarshal(body, &story)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Title: [#%s]%s\n", storyID, story["name"])
		fmt.Printf("Body: %s\n", story["url"])
	},
}

func expandTilde(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}

	usr, err := user.Current()
	if err != nil {
		return ""
	}

	return strings.Replace(path, "~", usr.HomeDir, 1)
}

func getProjectRoot() string {
	dir, err := os.Getwd()
	fmt.Println("dir: ", dir)
	if err != nil {
		panic(err)
	}

	for k, v := range config.Projects {
		expandedTildeRoot := expandTilde(v.Root)
		fmt.Println("expandedTildeRoot: ", expandedTildeRoot)
		path := os.ExpandEnv(expandedTildeRoot)
		if strings.HasPrefix(dir, path) {
			return k
		}
	}

	panic("Not in any project")
}

func extractStoryID(branch string) (string, error) {
	re := regexp.MustCompile(`(?:[a-z]+/)*#?(\d+)(?:-\w+)?`)
	matches := re.FindStringSubmatch(branch)
	if len(matches) < 2 {
		return "", fmt.Errorf("No story ID found in %s", branch)
	}
	return matches[1], nil
}

func getStoryID() (string, error) {
	branch, err := ioutil.ReadFile(".git/HEAD")
	if err != nil {
		panic(err)
	}
	fmt.Println("branch: ", string(branch))

	ref := strings.TrimSpace(string(branch))
	if !strings.HasPrefix(ref, "ref: refs/heads/") {
		panic("Not on a branch")
	}

	branchName := strings.TrimPrefix(ref, "ref: refs/heads/")
	return extractStoryID(branchName)
}
