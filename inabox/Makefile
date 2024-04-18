

dt := $(shell date '+%YY-%mM-%dD-%HH-%MM-%SS')

setup-geth:
	cd geth && rm -rf data && ./run.sh setup 10 ../secrets/ecdsa_keys/private_key_hex.txt ./secret/geth-account-password

start-geth:
	cd geth && docker compose up -d 
	./wait-for http://0.0.0.0:8545 -- echo "Geth up"

stop-geth:
	cd geth && docker compose down


.PHONY: new-anvil new-docker-anvil chain localstack exp deploy-all stop-infra run-e2e run-e2e-nochurner run-e2e-nograph clean

new-anvil:
	mkdir -p "testdata/$(dt)"
	cp ./templates/testconfig-anvil.yaml "testdata/$(dt)/config.yaml"

new-anvil-nochurner:
	mkdir -p "testdata/$(dt)"
	cp ./templates/testconfig-anvil-nochurner.yaml "testdata/$(dt)/config.yaml"


new-docker-anvil:
	mkdir -p "testdata/$(dt)"
	cp ./templates/testconfig-anvil-docker.yaml "testdata/$(dt)/config.yaml"

chain:
	go run ./deploy/cmd -localstack-port 4570 chain

localstack:
	GATEWAY_LISTEN=0.0.0.0:4570 LOCALSTACK_HOST=localhost.localstack.cloud:4570 localstack start

localstack-resources:
	AWS_URL=http://localhost:4570 ./create-s3-bucket.sh
	go run ./deploy/cmd -localstack-port 4570 resources

exp: 
	go run ./deploy/cmd exp

deploy-all:
	go run ./deploy/cmd  -localstack-port 4570 -deploy-resources true  all

stop-infra:
	./bin.sh stop-anvil
	./bin.sh stop-graph
	docker stop localstack-test

run-e2e:
	go test ./tests -v -config=../templates/testconfig-anvil.yaml

run-e2e-nochurner:
	go test ./tests -v -config=../templates/testconfig-anvil-nochurner.yaml

run-e2e-nograph:
	go test ./tests -v -config=../templates/testconfig-anvil-nograph.yaml

clean:
	rm -rf testdata/*

run-inabox:
	rm -rf testdata/*

	mkdir -p "testdata/$(dt)"
	cp ./templates/testconfig-anvil-nochurner.yaml "testdata/$(dt)/config.yaml"

	mkdir -p testdata/$(dt)/logs

	anvil --host 0.0.0.0 > testdata/$(dt)/logs/anvil.log 2>&1 &

	sleep 1

	cd ./thegraph 
	docker-compose up > testdata/$(dt)/logs/graph-node.log 2>&1 &
	cd ../

	sleep 5

	GATEWAY_LISTEN=0.0.0.0:4570 LOCALSTACK_HOST=localhost.localstack.cloud:4570 localstack start > testdata/$(dt)/logs/localstack.log 2>&1 &

	AWS_URL=http://localhost:4570 ./create-s3-bucket.sh
	go run ./deploy/cmd -localstack-port 4570 resources

	sleep 1

	make exp

	sleep 5

	./bin.sh start-detached

stop-services:
	./bin.sh stop-anvil
	./bin.sh stop-graph
	docker stop localstack-test
	./bin.sh stop-detached





	





	



	

