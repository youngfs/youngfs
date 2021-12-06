package iam

type Action string

type Identity struct {
	Name      string
	AccessKey string
	Actions   []Action
}

func (action Action) overBucket(bucket string) bool {
	return true
}
