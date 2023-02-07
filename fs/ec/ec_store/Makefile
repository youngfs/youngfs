all: ec_store_pb/suite.pb.go ec_store_pb/ec_plan.pb.go

ec_store_pb/suite.pb.go:
	protoc --go_out=. suite.proto

ec_store_pb/ec_plan.pb.go:
	protoc --go_out=. ec_plan.proto