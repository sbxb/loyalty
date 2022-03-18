GM_PATH=cmd/gophermart
GM_APP=gophermart
GM_ADDRESS=localhost:8080
ACCRUAL_ADDRESS=http://localhost:8888

.PHONY: run
run:
	go run cmd/gophermart/main.go

.PHONY: build
build:
	cd ./${GM_PATH} && rm -f ${GM_APP} && go build -o ${GM_APP} .

.PHONY: full
full: build
	./${GM_PATH}/${GM_APP} -a ${GM_ADDRESS} -d postgres://loyalty:loyalty@localhost/loyalty -r ${ACCRUAL_ADDRESS}

.PHONY: accrual
accrual:
	$(info ************************************)
	$(info *  Do not forget to run make init  *)
	$(info ************************************)
	cmd/accrual/accrual_linux_amd64 -a ${ACCRUAL_ADDRESS}

.PHONY: init
init:
	@./accrual_init.sh

.PHONY: test
test:
	go test -v -count=1 ./...