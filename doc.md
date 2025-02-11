# Build Ermes Platform

## DEPLOY FUNCTIONS
within faasd/functions:
- Create a template folder with the custom language within `/templates/template`
- Create a folder for each custom function within `/functions`. 
- Add the function to the stack.yml to deploy functions
  
## BUILD THE PLATFORM
- Create the json of the infrastructure
- Run the deploy.go script: `./deploy path/to/infrastructure.json`

## BEFORE BUILDING THE PLATFORM A SECOND TIME
Remove active faasd containers in each node:

    ctr -n openfaas-fn task kill --s sigkill function-name
    ctr -n openfaas-fn task rm function-name
  Todo for each deployed function.

Remove Redis container in each node:

    ctr task kill -s sigkill redis
    ctr container delete redis




## Configure Nodes
Before running the platform, some configuration steps may be required

### Configure SSH Access in each VM

    sudo apt-get update

    sudo apt-get install openssh-server

    sudo systemctl enable ssh

    sudo systemctl start ssh

On Host, Generate SSH Keys
`ssh-keygen -t rsa -b 2048`

### Copy SSH keys to VMs. 
In the Host terminal: 
    `ssh-copy-id user@1<VM_ip>`  

To verify SSH access: `ssh user@VM_ip`


### Remove sudo password
- `sudo visudo`
- In the bottom of the file, add 
    `$USER ALL=(ALL) NOPASSWD: ALL`

### Faasd authorization
`sudo cat /var/lib/faasd/secrets/basic-auth-password | faas-cli login --username admin --password-stdin --gateway http://10.62.0.5:8080/`

