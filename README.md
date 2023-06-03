# README Enhancer
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

README Enhancer is a Go project that aims to enhance the quality and readability of README files in GitHub repositories. It utilizes the OpenAI GPT-3.5 model to generate improved content for the README files.

The application will cycle over your public repositories, retrieve the README.md file for each repository, and enhance it using the OpenAI GPT-3.5 model.
The enhanced content will be used to create a new branch and update the README.md file. A pull request will be automatically created with the improved changes.

Please note that you need to have appropriate permissions and valid GitHub and OpenAI API keys with the necessary scopes to read and modify repositories and use the OpenAI GPT-3.5 model.

## Installation

To use README Enhancer, you need to have Go installed on your system. You can install Go by following the official installation guide: [https://golang.org/doc/install](https://golang.org/doc/install).

Once Go is installed, you can clone the repository:

```shell
git clone https://github.com/franpandol/readme_enhancer.git
````
Change to the project directory:

```shell
cd readme_enhancer
```

Build the project:

```shell
go build
```

## Usage

Set up the necessary environment variables. Create a .env file in the project root directory and provide the following values:

```code
GITHUB_API_KEY=your_github_api_key
GITHUB_REPOSITORY_OWNER=your_github_username
GITHUB_BASE_BRANCH=your_base_branch_name
OPENAI_API_KEY=your_openai_api_key
```

* Replace `your_github_api_key` with your GitHub API key. If you don't have an API key, you can generate one in your GitHub account settings.
* Replace `your_github_username` with your GitHub username.
* Replace `your_base_branch_name` with the name of the base branch in your repository (e.g., "main" or "master").
* Replace `your_openai_api_key` with your OpenAI API key. If you don't have an API key, you can obtain one from the OpenAI website.

Run the application:
```shell
./readme_enhancer
```

