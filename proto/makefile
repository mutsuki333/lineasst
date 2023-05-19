
SERVICE_PATH=../service
APIDOC_PATH=../src


all: generate


$(SERVICE_PATH):
	mkdir -p $(SERVICE_PATH)
$(APIDOC_PATH):
	mkdir -p $(APIDOC_PATH)

setup: $(wildcard *.proto) setup.go
	go run setup.go

generate: $(wildcard *.proto) $(SERVICE_PATH) $(APIDOC_PATH)
	protoc -I=. \
	--go_out=$(SERVICE_PATH) \
	--go_opt=paths=source_relative \
	--go-grpc_out=$(SERVICE_PATH) \
	--go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=$(SERVICE_PATH) \
	--grpc-gateway_opt=paths=source_relative \
	--grpc-gateway_opt=allow_delete_body=true \
	--openapiv2_out=$(APIDOC_PATH) \
	--openapiv2_opt=allow_merge=true \
	--openapiv2_opt=disable_service_tags=true \
	--openapiv2_opt=json_names_for_fields=false \
	--openapiv2_opt=allow_delete_body=true \
	--openapiv2_opt=use_go_templates=true \
	--openapiv2_opt=recursive-depth=3 \
	--openapiv2_opt=logtostderr=true \
	$(shell find . -maxdepth 1 -name "*.proto")


clean:
	find $(SERVICE_PATH) -type f \( -name "*pb.go" -o -name "entpb*" -o -name "*pb.gw.go" \) -delete
	find $(APIDOC_PATH) -type f -name "apidocs.swagger.json" -delete