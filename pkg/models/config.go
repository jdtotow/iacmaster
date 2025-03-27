package models

type Config struct {
	ApiPort                int
	DbUri                  string
	NodeName               string
	NodeType               NodeType
	Cluster                []string
	SecretKey              string
	ExecutionPlatform      string
	IACMasterSystemAddress string
	IACMasterSystemPort    int
	RunnerWorkingDir       string
}
