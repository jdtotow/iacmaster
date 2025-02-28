package models

import "time"

type GitData struct {
	Url           string `json:"url"`
	Token         string `json:"token"`
	TokenUsername string `json:"token_username"`
	Revision      string `json:"revision"`
	ProxyUrl      string `json:"proxy_url"`
	ProxyUsername string `json:"proxy_username"`
	ProxyPassword string `json:"proxy_password"`
}

type Deployment struct {
	Name                  string            `json:"name`
	WorkingDir            string            `json:"working_dir`
	HomeFolder            string            `json:"home_folder`
	CloudDestination      string            `json:"cloud_destination"`
	EnvironmentParameters map[string]string `json:"environment_parameters`
	Status                string            `json:"status`
	GitData               GitData           `json:"git_data"`
	EnvironmentID         string            `json:"environment_id"`
	Error                 string            `json:"error"`
	Activities            []string          `json:"activities"`
}

func (d *Deployment) SetError(_error string) {
	d.Error = _error
}

func (d *Deployment) AddActivity(activity string) {
	d.Activities = append(d.Activities, time.Now().Format("01/02/2006 - 15:04:05")+" "+activity)
}
