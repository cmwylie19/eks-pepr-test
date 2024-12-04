k exec -it deploy/k8s-reporter -n k8s-reporter -- cat /tmp/k8s-watcher.log
kubectl cp k8s-reporter/k8s-reporter-bb4569d8-p79gl:/tmp/k8s-watcher.log ./k8s-watcher.log
cat ./k8s-watcher.log
