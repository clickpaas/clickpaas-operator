all: build, deploy

build:
	docker build --pull -t registry.bizsaas.net/clickpaas-controller:v1alpha1 -f deploy/docker/Dockerfile

deploy:
	curnodeip=`hostname -i`
	curhostname=`kubectl get node -o wide |grep $curnodeip |awk '{print $1}'`
	kubectl label node $curhostname custom-controller=clickpaas-controller
	kubectl apply -f deploy/deployment.yml

.PHONY: clean
	kubectl delete -f deploy/deployment.yml
	for crd in `kubectl get crd|grep -v NAME|awk '{print $1}'`;do kubectl delete crd $crd;done
	echo "wait 10 second to stop controller container"
	sleep 10
	docker rmi
	registry.bizsaas.net/clickpaas-controller