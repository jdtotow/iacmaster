package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BackendType string

const LocalBackend BackendType = "local"
const RemoteBackend BackendType = "remote"
const PrivateBackend BackendType = "private"

type StateFileStorage string

const LocalStateFile StateFileStorage = "local"
const S3StateFile StateFileStorage = "s3"
const AzureStateFile StateFileStorage = "azure"
const HTTPServerStateFile StateFileStorage = "http"

type DestinationCloud string

const Azure DestinationCloud = "azure"
const AWS DestinationCloud = "aws"
const GCP DestinationCloud = "gcp"

type IaCExecutionSettings struct {
	gorm.Model
	ID               uint
	TerraformVersion string           `json:"terraform_version"`
	BackendType      BackendType      `json:"backend_type"`
	StateFileStorage StateFileStorage `json:"state_file_storage"`
	DestinationCloud DestinationCloud `json:"destination_cloud"`
	Credential       CloudCredential  `json:"credential" gorm:"foreignKey:CredentialUuid"`
	CredentialUuid   uuid.UUID
	Uuid             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
}

func (i *IaCExecutionSettings) SetTerraformVersion(version string) {
	i.TerraformVersion = version

}
func (i *IaCExecutionSettings) SetBackendType(backend string) {
	i.BackendType = BackendType(backend)
}
func (i *IaCExecutionSettings) SetStateFileStorage(storage string) {
	i.StateFileStorage = StateFileStorage(storage)
}
func (i *IaCExecutionSettings) SetDestinationCloud(destination string) {
	i.DestinationCloud = DestinationCloud(destination)
}
func (i *IaCExecutionSettings) SetCredential(credential CloudCredential) {
	i.Credential = credential
}
