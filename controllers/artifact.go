package controllers

import (
	"log"
	"os"

	"github.com/go-git/go-git/v5"
)

type IaCArtifactController struct {
	TmpFolderPath string
}

func CreateIaCArtifactController(tmp string) *IaCArtifactController {
	return &IaCArtifactController{
		TmpFolderPath: tmp,
	}
}

func (i *IaCArtifactController) GetRepo(url, token, environment string) error {
	log.Println("Cloning -> ", url)
	_, err := os.Stat(i.TmpFolderPath + "/" + environment)
	if err != nil {
		log.Println("The folder exist, it will be remove")
		os.RemoveAll(i.TmpFolderPath + "/" + environment)
	}

	_, err = git.PlainClone(i.TmpFolderPath+"/"+environment, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	return err
}
