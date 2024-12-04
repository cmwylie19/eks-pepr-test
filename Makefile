IMAGE_NAME = k8s-reporter
K8S_DIR = k8s
DOCKERFILE = Dockerfile
TAG = 0.0.1

build-image:
	k3d cluster delete --all;
	k3d cluster create;
	docker build -t $(IMAGE_NAME):$(TAG) -f kubernetes-reporter/$(DOCKERFILE) kubernetes-reporter/
	k3d image import k8s-reporter:0.0.1 -c k3s-default 

apply-k8s:
	kubectl apply -f kubernetes-reporter/$(K8S_DIR)

deploy: build-image apply-k8s

k8s-deploy: 
	kubectl apply -f kubernetes-reporter/$(K8S_DIR)
	kubectl apply -f kubernetes-watcher/dist
