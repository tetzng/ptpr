# ptpr

ptpr is a CLI tool for generating pull request options from your Pivotal Tracker project. The tool reads the API token and project details from a YAML configuration file located at `~/.config/ptpr.yaml`. The `gen` command retrieves the story information from Pivotal Tracker API using the story ID obtained from the current Git branch name. The retrieved story information (name and URL) is printed as a string in the format of `--title=[#STORY_ID]STORY_NAME --body=STORY_URL`.

## Installation

Clone this repository:

```sh
git clone https://github.com/tetzng/ptpr.git
```

Change into the project directory:

```sh
cd ptpr
```

Build the CLI using the provided Makefile:

```sh
make build
```

Add the resulting binary to your PATH, or move it to a directory that is already in your PATH:

```sh
mv ptpr /usr/local/bin/

```

## Usage

Create a YAML configuration file at `~/.config/ptpr.yaml` with the following format:

```yaml
PIVOTAL_API_TOKEN: YOUR_PIVOTAL_API_TOKEN # This is default PIVOTAL_API_TOKEN
Projects:
  path/to/project1:
    PIVOTAL_PROJECT_ID: 12345678
  path/to/project2:
    PIVOTAL_PROJECT_ID: 23456789
    PIVOTAL_API_TOKEN: ANOTHER_PIVOTAL_API_TOKEN # this is optional: if you want to use another PIVOTAL_API_TOKEN
```

Replace the placeholders with your actual Pivotal API tokens and project details.

Navigate to your project directory and run the `ptpr gen` command:

```sh
ptpr gen
```

This will output the pull request options based on the story ID found in the current Git branch name.

## Setup

To create a configuration file using the init command, run:

```sh
ptpr init
```

This will interactively create or update the configuration file at ~/.config/ptpr.yaml. The init command will prompt you for your Pivotal API Token and the Pivotal Project ID for the current project, and it will create a configuration file with this information. If the configuration file already exists, the init command will add or update the project and API token for the current project. To overwrite existing configuration file without prompt, use -f option.

After running the init command, you can use the gen command to generate pull request options based on the story ID found in the current Git branch name.

## Integration with GitHub CLI

The output of the `gen` command is formatted to be used directly as arguments for the `gh pr create` command from the [GitHub CLI](https://cli.github.com/). You can use the `xargs` command to pass the output of `ptpr gen` directly to `gh pr create` like this:

```sh
ptpr gen | xargs gh pr create
```

This command will create a new pull request on GitHub using the generated title and body from the Pivotal Tracker story.

## Contributing

1. Fork this repository and create a new branch for your feature or bugfix.
1. Make your changes and commit them to your branch.
1. Open a pull request to merge your changes back into the main repository.
