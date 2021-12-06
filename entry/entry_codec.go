package entry

import (
	"github.com/golang/protobuf/proto"
	"object-storage-server/entry/entry_pb"
)

func (entry *Entry) EncodeProto() ([]byte, error) {
	message := entry.ToPb()
	return proto.Marshal(message)
}

func (entry *Entry) DecodeProto(b []byte) error {
	message := &entry_pb.Entry{}
	if err := proto.Unmarshal(b, message); err != nil {
		return err
	}
	entry = EntryPbToInstance(message)
	return nil
}
