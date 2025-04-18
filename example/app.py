import requests, json 

#creating organization 
iacmaster_url = "http://localhost:3000"
org_data = {"name":"swisscom"}
headers = {"Content-Type":"application/json"}

"""

result = requests.post(iacmaster_url+"/organization", json.dumps(org_data), headers=headers)
print(result.text)
org_id = json.loads(result.text)["id"]

#creating project 
project_data = {"name":"aurora", "organization_id": org_id}
result = requests.post(iacmaster_url+"/project", json.dumps(project_data), headers=headers)
project_id = json.loads(result.text)["id"]

#creation of artifact 

artifact_data = {
    "name": "vpn-to-vpn",
    "type": "terraform",
    "scm_url": "https://github.com/futurice/terraform-examples",
    "revision": "master",
    "home_folder": "azure/azure_linux_docker_app_service",
    "project_id": project_id
}

result = requests.post(iacmaster_url+"/iacartifact", json.dumps(artifact_data), headers=headers)
artifact_id = json.loads(result.text)["id"]

#creating token 

token_data = {
    "name": "test",
    "type": "git",
    "username": "jd",
    "token": "",
    "project_id": project_id
}

result = requests.post(iacmaster_url+"/token", json.dumps(token_data), headers=headers)
token_id = json.loads(result.text)["id"]


credential_data = {
    "name": "access_azure",
    "destination_cloud": "azure",
    "project_id": project_id,
    "variables":{
        "AZURE_SUBSCRIPTION_ID":"your-subscription-id",
        "AZURE_CLIENT_ID":"your-client-id",
        "AZURE_CLIENT_SECRET":"your-client-secret",
        "AZURE_TENANT_ID":"your-tenant-id",
    }
}

result = requests.post(iacmaster_url+"/cloudcredential", json.dumps(credential_data), headers=headers)
print(result.text)
credential_id = json.loads(result.text)["id"]

setting_data = {
    "terraform_version": "1.9.4",
    "backend_type": "local",
    "state_file_storage": "local",
    "destination_cloud": "azure",
    "cloudcredential_id": credential_id,
    "token_id": token_id
}

result = requests.post(iacmaster_url+"/settings", json.dumps(setting_data), headers=headers)
print(result.text)
setting_id = json.loads(result.text)["id"]

env_data = {
    "name": "vpn-env",
    "project_id": project_id,
    "iacartifact_id": artifact_id,
    "iac_execution_settings_id": setting_id,
    "status": "init"
}

result = requests.post(iacmaster_url+"/environment", json.dumps(env_data), headers=headers)
print(result.text)
environment_id = json.loads(result.text)["id"]

"""

files = {'file': open('variables.tfvars','rb')}
values = {'artifact': 'terraform', 'environment_id': 'f8b0b561-417e-41c1-aabd-960319704a6f'}

r = requests.post(iacmaster_url+"/environment/f8b0b561-417e-41c1-aabd-960319704a6f/variables", files=files, data=values)
