build:
	go build -C ./src/ -o ../render

run-test01: build
	./render --template test/test01/nginx-ingress-controller.yaml -data test/test01/cluster.json > test/test01/nginx-ingress-controller.yaml.out
