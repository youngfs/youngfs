all:
	make master
	make volume

master: master_pb/master.pb.go master_pb/master_grpc.pb.go

master_pb/master.pb.go:
master_pb/master_grpc.pb.go:
	protoc --go_out=. --go-grpc_out=. master.proto

volume: volume_pb/volume.pb.go volume_pb/volume_grpc.pb.go

volume_pb/volume.pb.go:
volume_pb/volume_grpc.pb.go:
	protoc --go_out=. --go-grpc_out=. volume.proto
