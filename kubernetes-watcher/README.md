# Extra Manifests

Usage:

```bash
k exec -it deploy/pepr-k8s-eks-watcher -n pepr-system -- 
cat /tmp/k8s-watcher.log
```


```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: pepr-k8s-watcher
  namespace: pepr-system
  annotations:
    pepr.dev/description: Watch Kubernetes EndpointSlice and Service
  labels:
    app: pepr-k8s-watcher
    pepr.dev/controller: admission
    pepr.dev/uuid: k8s-watcher
spec:
  replicas: 2
  selector:
    matchLabels:
      app: pepr-k8s-watcher
      pepr.dev/controller: admission
  template:
    metadata:
      annotations:
        buildTimestamp: '1733326822584'
      labels:
        app: pepr-k8s-watcher
        pepr.dev/controller: admission
    spec:
      terminationGracePeriodSeconds: 5
      priorityClassName: system-node-critical
      serviceAccountName: pepr-k8s-watcher
      securityContext:
        runAsUser: 65532
        runAsGroup: 65532
        runAsNonRoot: true
        fsGroup: 65532
      containers:
        - name: server
          image: ghcr.io/defenseunicorns/pepr/controller:v0.40.1
          imagePullPolicy: IfNotPresent
          command:
            - node
            - /app/node_modules/pepr/dist/controller.js
            - 2bf62d994d332564c1ac3de56042ee5b02efc998a1bd30fcd011c6b355ad22e3
          readinessProbe:
            httpGet:
              path: /healthz
              port: 3000
              scheme: HTTPS
            initialDelaySeconds: 10
          livenessProbe:
            httpGet:
              path: /healthz
              port: 3000
              scheme: HTTPS
            initialDelaySeconds: 10
          ports:
            - containerPort: 3000
          resources:
            requests:
              memory: 256Mi
              cpu: 200m
            limits:
              memory: 512Mi
              cpu: 500m
          env:
            - name: PEPR_WATCH_MODE
              value: 'false'
            - name: PEPR_PRETTY_LOG
              value: 'false'
            - name: LOG_LEVEL
              value: info
          securityContext:
            runAsUser: 65532
            runAsGroup: 65532
            runAsNonRoot: true
            allowPrivilegeEscalation: false
            capabilities:
              drop:
                - ALL
          volumeMounts:
            - name: logs
              mountPath: /tmp
            - name: tls-certs
              mountPath: /etc/certs
              readOnly: true
            - name: api-token
              mountPath: /app/api-token
              readOnly: true
            - name: module
              mountPath: /app/load
              readOnly: true
      volumes:
        - name: logs
          emptyDir: {} 
        - name: tls-certs
          secret:
            secretName: pepr-k8s-watcher-tls
        - name: api-token
          secret:
            secretName: pepr-k8s-watcher-api-token
        - name: module
          secret:
            secretName: pepr-k8s-watcher-module
```
