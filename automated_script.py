import subprocess

# Install Ansible using pip
subprocess.run(["pip", "install", "ansible"])

# Download the playbook from the URL
subprocess.run(["curl", "-O", "https://raw.githubusercontent.com/husseinahmed-dev/Research-Project-1-Kubernetes-Hardening/main/01-deploy-loadbalancer.yml"])
subprocess.run(["curl", "-O", "https://raw.githubusercontent.com/husseinahmed-dev/Research-Project-1-Kubernetes-Hardening/main/02-prepare-master-worker.yml"])
subprocess.run(["curl", "-O", "https://raw.githubusercontent.com/husseinahmed-dev/Research-Project-1-Kubernetes-Hardening/main/03-install-kubernetes-master-worker.yml"])

# Apply the playbook using ansible-playbook
subprocess.run(["ansible-playbook", "01-deploy-loadbalancer.yml"])
subprocess.run(["ansible-playbook", "02-prepare-master-worker.yml"])
subprocess.run(["ansible-playbook", "03-install-kubernetes-master-worker.yml"])
