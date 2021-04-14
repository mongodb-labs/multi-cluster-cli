package main

const svc = `
apiVersion: v1
kind: Service
metadata:
  # labels:
  #   app: my-replica-set-0-svc
  # name: my-replica-set-0-svc
spec:
  clusterIP: None
  ports:
  - name: mongodb
    port: 27017
    protocol: TCP
    targetPort: 27017
  publishNotReadyAddresses: true
  selector:
    # app: my-replica-set-0
`

const pod = `
apiVersion: v1
kind: Pod
metadata:
  labels:
    app: my-replica-set-0
    statefulset.kubernetes.io/pod-name: my-replica-set-0
  name: my-replica-set-0
spec:
  containers:
  - command:
    - /opt/scripts/agent-launcher.sh
    env:
    - name: AGENT_API_KEY
      valueFrom:
        secretKeyRef:
          key: agentApiKey
          name: 605c39d90882a963fc3734f2-group-secret
    - name: AGENT_FLAGS
      value: -logFile,/var/log/mongodb-mms-automation/automation-agent.log,
    - name: BASE_URL
      value: https://cloud-qa.mongodb.com
    - name: GROUP_ID
      value: 605c39d90882a963fc3734f2
    - name: LOG_LEVEL
      value: DEBUG
    - name: SSL_REQUIRE_VALID_MMS_CERTIFICATES
      value: "true"
    image: quay.io/mongodb/mongodb-enterprise-database:2.0.0
    imagePullPolicy: Always
    livenessProbe:
      exec:
        command:
        - /opt/scripts/probe.sh
      failureThreshold: 6
      initialDelaySeconds: 60
      periodSeconds: 30
      successThreshold: 1
      timeoutSeconds: 30
    name: mongodb-enterprise-database
    ports:
    - containerPort: 27017
      protocol: TCP
    readinessProbe:
      exec:
        command:
        - /opt/scripts/readinessprobe
      failureThreshold: 240
      initialDelaySeconds: 5
      periodSeconds: 5
      successThreshold: 1
      timeoutSeconds: 1
    resources:
      requests:
        cpu: 200m
        memory: 300M
    securityContext:
      runAsNonRoot: true
      runAsUser: 2000
    terminationMessagePath: /dev/termination-log
    terminationMessagePolicy: File
    volumeMounts:
    - mountPath: /opt/scripts
      name: database-scripts
      readOnly: true
  dnsPolicy: ClusterFirst
  enableServiceLinks: true
  hostname: my-replica-set-2
  imagePullSecrets:
  - name: image-registries-secret
  initContainers:
  - image: 268558157000.dkr.ecr.eu-west-1.amazonaws.com/raj/ubuntu/mongodb-enterprise-init-database:latest
    imagePullPolicy: Always
    name: mongodb-enterprise-init-database
    resources: {}
    securityContext:
      runAsNonRoot: true
      runAsUser: 2000
    terminationMessagePath: /dev/termination-log
    terminationMessagePolicy: File
    volumeMounts:
    - mountPath: /opt/scripts
      name: database-scripts
  priority: 0
  restartPolicy: Always
  schedulerName: default-scheduler
  securityContext:
    fsGroup: 2000
  subdomain: my-replica-set-svc
  terminationGracePeriodSeconds: 600
  volumes:
  - emptyDir: {}
    name: database-scripts

`

const processesJSON = `
{

	"auth": {
		"usersWanted": [],
		"usersDeleted": [],
		"disabled": true,
		"authoritativeSet": false,
		"autoAuthMechanisms": [],
		"autoAuthRestrictions": []
	},
	"processes": [{
			"name": "my-replica-set-0",
			"processType": "mongod",
			"version": "4.4.0",
			"authSchemaVersion": 5,
			"featureCompatibilityVersion": "4.4",
      "disabled": false,
      "manualMode": false,
			"hostname": "my-replica-set-0-svc.mdb.svc.cluster.local",
			"args2_6": {
				"net": {
					"port": 27017,
					"tls": {
						"mode": "disabled"
					}
				},
				"replication": {
					"replSetName": "my-replica-set"
				},
				"storage": {
					"dbPath": "/data"
				},
				"systemLog": {
					"destination": "file",
					"path": "/var/log/mongodb-mms-automation/mongodb.log"
				}
			},
			"horizons": {},
			"logRotate": {
				"sizeThresholdMB": 1000,
				"timeThresholdHrs": 24
			}
		},
		{
			"name": "my-replica-set-1",
			"processType": "mongod",
			"version": "4.4.0",
			"authSchemaVersion": 5,
			"featureCompatibilityVersion": "4.4",
      "disabled": false,
      "manualMode": false,
			"hostname": "my-replica-set-1-svc.mdb.svc.cluster.local",
			"args2_6": {
				"net": {
					"port": 27017,
					"tls": {
						"mode": "disabled"
					}
				},
				"replication": {
					"replSetName": "my-replica-set"
				},
				"storage": {
					"dbPath": "/data"
				},
				"systemLog": {
					"destination": "file",
					"path": "/var/log/mongodb-mms-automation/mongodb.log"
				}
			},
			"horizons": {},
			"logRotate": {
				"sizeThresholdMB": 1000,
				"timeThresholdHrs": 24
			}
		},
		{
			"name": "my-replica-set-2",
			"processType": "mongod",
			"version": "4.4.0",
			"authSchemaVersion": 5,
			"featureCompatibilityVersion": "4.4",
      "disabled": false,
      "manualMode": false,
			"hostname": "my-replica-set-2-svc.mdb.svc.cluster.local",
			"args2_6": {
				"net": {
					"port": 27017,
					"tls": {
						"mode": "disabled"
					}
				},
				"replication": {
					"replSetName": "my-replica-set"
				},
				"storage": {
					"dbPath": "/data"
				},
				"systemLog": {
					"destination": "file",
					"path": "/var/log/mongodb-mms-automation/mongodb.log"
				}
			},
			"horizons": {},
			"logRotate": {
				"sizeThresholdMB": 1000,
				"timeThresholdHrs": 24
			}
		}
	],
	"replicaSets": [{
		"_id": "my-replica-set",
		"members": [{
				"_id": 0,
				"arbiterOnly": false,
				"hidden": false,
				"priority": 1,
				"slaveDelay": 0,
				"votes": 1,
				"buildIndexes": true,
				"tags": {},
				"host": "my-replica-set-0"
			},
			{
				"_id": 1,
				"arbiterOnly": false,
				"hidden": false,
				"priority": 1,
				"slaveDelay": 0,
				"votes": 1,
				"buildIndexes": true,
				"tags": {},
				"host": "my-replica-set-1"
			},
			{
				"_id": 2,
				"arbiterOnly": false,
				"hidden": false,
				"priority": 1,
				"slaveDelay": 0,
				"votes": 1,
				"buildIndexes": true,
				"tags": {},
				"host": "my-replica-set-2"
			}
		],
		"protocolVersion": "1",
		"settings": {}
	}],
	"monitoringVersions": [{
			"name": "6.4.0.433-1",
			"hostname": "my-replica-set-0-svc.mdb.svc.cluster.local"
		},
		{
			"name": "6.4.0.433-1",
			"hostname": "my-replica-set-1-svc.mdb.svc.cluster.local"
		},
		{
			"name": "6.4.0.433-1",
			"hostname": "my-replica-set-2-svc.mdb.svc.cluster.local"
		}
	],
	"backupVersions": [{
			"name": "6.6.0.959-1",
			"hostname": "my-replica-set-0-svc.mdb.svc.cluster.local"
		},
		{
			"name": "6.6.0.959-1",
			"hostname": "my-replica-set-1-svc.mdb.svc.cluster.local"
		},
		{
			"name": "6.6.0.959-1",
			"hostname": "my-replica-set-2-svc.mdb.svc.cluster.local"
		}
	],
	"agentVersion": {
		"name": "10.27.1.6801-1",
		"directoryUrl": "https://s3.amazonaws.com/mciuploads/mms-automation/mongodb-mms-build-agent/builds/automation-agent/qa/"
	},
	"options": {
		"downloadBase": "/var/lib/mongodb-mms-automation",
		"downloadBaseWindows": "%SystemDrive%\\MMSAutomation\\versions"
	}
}
`
