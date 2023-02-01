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
	rootCmd.AddCommand(genCmd)

	data, err := ioutil.ReadFile(os.Getenv("HOME") + "/.config/pivo_conf.yml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
}

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate PR options from your Pivotal Tracker Project",
  Long: `Generate Pull Request (PR) options from a Pivotal Tracker project.
The tool reads the API token and project details from a YAML configuration file
located at ~/.config/pivo_conf.yml. The gen command retrieves the story
information from Pivotal Tracker API using the story ID obtained from the
current Git branch name. The retrieved story information (name and URL)
is printed as a string in the format of "--title=[#STORY_ID]STORY_NAME --body=STORY_URL".`,
	Run: func(cmd *cobra.Command, args []string) {
		projectRoot := getProjectRoot()
		storyID, err := getStoryID()

    if len(storyID) == 0 {
      fmt.Println("No story ID")
      return
    }

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
      fmt.Printf("%v",resp)
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
    if story["name"] == nil {
      fmt.Println("No story name.")
      return
    }
    if story["url"] == nil {
      fmt.Println("No story url.")
      return
    }
    fmt.Printf("--title=[#%s]%s --body=%s\n", storyID, story["name"], story["url"])
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
	if err != nil {
		panic(err)
	}

	for k, v := range config.Projects {
		expandedTildeRoot := expandTilde(v.Root)
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

	ref := strings.TrimSpace(string(branch))
	if !strings.HasPrefix(ref, "ref: refs/heads/") {
		panic("Not on a branch")
	}

	branchName := strings.TrimPrefix(ref, "ref: refs/heads/")
	return extractStoryID(branchName)
}
