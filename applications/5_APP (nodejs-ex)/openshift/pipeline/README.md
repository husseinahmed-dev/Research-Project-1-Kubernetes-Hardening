This directory contains a Jenkinsfile which can be used to build
nodejs-ex using an OpenShift build pipeline.

To do this, run:

```bash
# create the nodejs example as usual
oc new-app https://github.com/husseinahmed-dev/Research-Project-1-Kubernetes-Hardening/tree/main/applications/5_APP%20(nodejs-ex)

# now create the pipeline build controller from the openshift/pipeline
# subdirectory
oc new-app [https://github.com/sclorg/nodejs-ex](https://github.com/husseinahmed-dev/Research-Project-1-Kubernetes-Hardening/tree/main/applications/5_APP%20(nodejs-ex)) \
  --context-dir=openshift/pipeline --name nodejs-ex-pipeline
```
