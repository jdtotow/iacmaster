package models

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
	CloudDestination      string            `json:"cloud_destination"`
	EnvironmentParameters map[string]string `json:"environment_parameters`
	Status                string            `json:"status`
	GitData               GitData           `json:"git_data"`
	EnvironmentID         string            `json:"environment_id"`
}
