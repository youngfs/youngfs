package iam

type setIAM struct {
	User // user_iam name
	Set  // set_iam name
}

// read permission
func (iam *setIAM) ReadKey() string {
	return string(iam.User) + setReadIAMKv
}

// write permission
func (iam *setIAM) WriteKey() string {
	return string(iam.User) + setWriteIAMKv
}
