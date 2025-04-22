Users
- system 
- org_admin
- project_admin
- environment_admin
- users

Resources
- organization
- project
- iac artifact
- cloud credential 
- git token 

Actions
- create
- get 
- delete 
- edit 
- invite 
- join 

Hierarchy
- system -> organization -> project -> environment 

function (user, resource, action):
    store(user, resource, action)

function can(user, resource, action):
    if get_store(user, resource, action):
        return true 
    else:
        return false 