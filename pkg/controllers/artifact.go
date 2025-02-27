package controllers

import (
	"fmt"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type IaCArtifactController struct {
	TmpFolderPath string
}

func CreateIaCArtifactController(tmp string) *IaCArtifactController {
	return &IaCArtifactController{
		TmpFolderPath: tmp,
	}
}

func (i *IaCArtifactController) UpdateRepo(url, token, tokenUsername, revision, proxyUrl, proxyUsername, proxyPassword, environment string) error {
	localPath := i.TmpFolderPath + "/" + environment
	repo, err := git.PlainOpen(localPath)
	if err != nil {
		log.Printf("Error opening repository: %s\n", err)
		return err
	}

	// Get the working directory for the repository
	worktree, err := repo.Worktree()
	if err != nil {
		log.Printf("Error getting worktree: %s\n", err)
		return err
	}
	var proxyOptions transport.ProxyOptions
	var auth http.BasicAuth

	if proxyUrl != "" {
		proxyOptions.URL = proxyUrl
		proxyOptions.Username = proxyUsername
		proxyOptions.Password = proxyPassword
	}

	if token != "" {
		auth.Username = tokenUsername
		auth.Password = token
	}

	// Pull the latest changes from the remote
	log.Println("Pulling latest changes...")
	err = worktree.Pull(&git.PullOptions{
		Auth:         &auth,
		ProxyOptions: proxyOptions,
	})
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			fmt.Println("Repository is already up-to-date.")
		} else {
			fmt.Printf("Error pulling repository: %s\n", err)
		}
		return err
	}
	fmt.Println("Repository updated successfully.")
	return nil
}

func (i *IaCArtifactController) GetRepo(url, token, tokenUsername, revision, proxyUrl, proxyUsername, proxyPassword, environment string) error {
	log.Println("Cloning repository -> ", url)
	var proxyOptions transport.ProxyOptions
	var auth http.BasicAuth

	if proxyUrl != "" {
		proxyOptions.URL = proxyUrl
		proxyOptions.Username = proxyUsername
		proxyOptions.Password = proxyPassword
	}

	if token != "" {
		auth.Username = tokenUsername
		auth.Password = token
	}

	_, err := git.PlainClone(i.TmpFolderPath+"/"+environment, false, &git.CloneOptions{
		URL:           url,
		Progress:      os.Stdout,
		ReferenceName: plumbing.ReferenceName(revision),
		ProxyOptions:  proxyOptions,
		Auth:          &auth,
	})
	return err
}
