package controllers

import (
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

func (i *IaCArtifactController) GetRepo(url, token, tokenUsername, revision, proxyUrl, proxyUsername, proxyPassword, environment string) error {
	log.Println("Cloning -> ", url)
	_, err := os.Stat(i.TmpFolderPath + "/" + environment)
	if err != nil {
		log.Println("The folder exist, it will be remove")
		os.RemoveAll(i.TmpFolderPath + "/" + environment + "/repo")
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

	_, err = git.PlainClone(i.TmpFolderPath+"/"+environment, false, &git.CloneOptions{
		URL:           url,
		Progress:      os.Stdout,
		ReferenceName: plumbing.ReferenceName(revision),
		ProxyOptions:  proxyOptions,
		Auth:          &auth,
	})
	return err
}
