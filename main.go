package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
	"golang.org/x/oauth2"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	requestGithub()
}

// Function to generate a random branch name
func generateRandomBranchName() string {
	rand.Seed(time.Now().UnixNano())

	timestamp := time.Now().Format("20050102150405") // Current timestamp
	randomString := make([]byte, 8)

	for i := range randomString {
		randomString[i] = byte(rand.Intn(26) + 97) //nolint:gosec // we don't need cryptographically secure random numbers here
	}

	return fmt.Sprintf("branch-%s-%s", string(randomString), timestamp)
}

func updateFile(ctx context.Context, repositoryName string, client *github.Client, improvedContent string) {
	// Repository information
	owner := os.Getenv("GITHUB_REPOSITORY_OWNER")
	repo := repositoryName
	baseBranch := os.Getenv("GITHUB_BASE_BRANCH")
	newBranch := generateRandomBranchName()
	filePath := "README.md"

	// Check if the new branch already exists
	_, _, err := client.Repositories.GetBranch(ctx, owner, repo, newBranch)
	if err == nil {
		log.Fatal("Branch already exists")
	}

	// Get the reference of the base branch
	baseBranchRef, _, err := client.Git.GetRef(ctx, owner, repo, "heads/"+baseBranch)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new branch based on the base branch
	newBranchRef := &github.Reference{
		Ref: github.String("refs/heads/" + newBranch),
		Object: &github.GitObject{
			SHA: baseBranchRef.Object.SHA,
		},
	}

	_, _, err = client.Git.CreateRef(ctx, owner, repo, newBranchRef)
	if err != nil {
		log.Fatal(err)
	}

	// Get the latest commit of the README file in the new branch
	readmeFile, _, _, err := client.Repositories.GetContents(ctx, owner, repo, filePath, &github.RepositoryContentGetOptions{Ref: newBranch})
	if err != nil {
		log.Fatal(err)
	}

	// Update the content of the README in the new branch
	updateContent := &github.RepositoryContentFileOptions{
		Message: github.String("Update README.md"),
		Content: []byte(improvedContent),
		SHA:     readmeFile.SHA,
		Branch:  github.String(newBranch),
	}

	_, _, err = client.Repositories.UpdateFile(ctx, owner, repo, filePath, updateContent)
	if err != nil {
		log.Fatal(err)
	}

	// create a pull request
	pr := &github.NewPullRequest{
		Title:               github.String("Update README.md"),
		Head:                github.String(newBranch),
		Base:                github.String(baseBranch),
		Body:                github.String("This is an automated pull request to update the README.md file."),
		MaintainerCanModify: github.Bool(true),
	}
	pullRequest, _, err := client.PullRequests.Create(ctx, owner, repo, pr)
	if err != nil {
		log.Fatal(err)
	}

	// print pr url
	fmt.Println(*pullRequest.HTMLURL)
}

func requestGithub() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_API_KEY")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// list all repositories for the authenticated user
	repositories, _, _ := client.Repositories.List(ctx, "", nil)
	for _, repository := range repositories {
		// get README.md from every repo 
		readme, _, _ := client.Repositories.GetReadme(ctx, *repository.Owner.Login, *repository.Name, nil)
		if readme == nil {
			continue
		}

		// download README.md and read it
		resp, _ := client.Repositories.DownloadContents(ctx, *repository.Owner.Login, *repository.Name, "README.md", nil)
		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(resp)

		if err != nil {
			fmt.Println(err)
		}

		readmeContent := buf.String()
		prompt := fmt.Sprintf("Improve the README.md for the %s repository.\n\n%s", *repository.Name, readmeContent)
		newReadme := requestOpenAI(prompt)

		updateFile(ctx, *repository.Name, client, newReadme)
	}
}

func requestOpenAI(prompt string) string {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return ""
	}

	return resp.Choices[0].Message.Content
}
