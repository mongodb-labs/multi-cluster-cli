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
  generateName: my-replica-set-
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
      limits:
        cpu: 500m
        memory: 700M
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