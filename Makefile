PROTO=$(shell find api/esi -name \*.proto)

proto:
	protoc -I. --go_out=,paths=source_relative:. ${PROTO}

clean:
	rm api/esi/*.pb.go
