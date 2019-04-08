.PHONY: generate
generate:
	rm -rf generated/*
	protoc --proto_path=definitions signature.proto --go_out=generated