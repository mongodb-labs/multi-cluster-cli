### mccli
multi-cluster-cli (aka mccli) allows you to run MongoDB deployments accross multiple [Kind](https://kind.sigs.k8s.io/) clusters on your local machine. 

#### Prerequisite

Create or Use an exiting cloud-qa organization and put the credentials in a `config.json` file:

```json
{
  "orgId": "${ORG_ID}",
  "public": "${PUBLIC_KEY}",
  "private": "${PRIVATE_KEY}"
}
```

### CLI commands

* create the kind clusters locally -- creates 2 clusters with name `cluster-a` and `cluster-b`.

  `mccli --op cluster`

* create the Cloud-QA project and automation config
   `mccli --op project`

* deploy MongoDB pods/nodes accross `cluster-a` and `cluster-b`:
   deploys one pod in `cluster-a` and 2 pods in `cluster-b`

   `mccli --op mongo`

Wait for the deployment to show up on cloud-qa UI.

**Note: This project is only for experimental deployments and is not supported officially by MongoDB. DO NOT use it in production**.