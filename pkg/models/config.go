package models

/*
API_PORT=3000
DB_URI="host=localhost user=iacmaster password=password dbname=iacmaster port=5430 sslmode=disable"
NODE_NAME=main
NODE_TYPE=primary
CLUSTER="main=localhost,zollikofen=127.0.0.1,zurich=0.0.0.0"
SERVICE_URL=http://localhost:2020
SECRET_KEY=c1876cdd21d43cd460fe1a3cb4bc5d0847018cee3a33744bfdf0b57ecff2f53f059a6f8b1312790db2bb451a8cd5f368d2c9f69bc36a5cfbd144006ca1f57d1dc2ac1285163fbfb6ea5302ecf8f5fe2f2b00628e7f85b8ade813a3738251dd919a0c29b6bef280f39054f0bbe4eafd9627572cb6f6c7e16c1def3497a9bf39ec28d573f04d25d9d1eff1780220b091fc2aea532b5e39770b3e9d0b4f686db6510935f885e2e1abdf999f9debf69a1d7da994abf38e4263cc5ff4ef8e4f68b8e48fdc81387e2a4572b8aaf2b08f909b9112a927811d70e97938067772703089cf18f32ef6ceeee42d3e5a79b494ebe8f9858955b5b9bd3568d529b44d57f1be33
EXECUTION_PLATFORM=docker
IACMASTER_SYSTEM_ADDRESS=192.168.1.128
IACMASTER_SYSTEM_PORT=3434
RUNNER_WORKING_DIR=/runner
*/

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
