package protocol

//go:generate protoc -I ./ --go_out=plugins=grpc:./ ./raft.proto
