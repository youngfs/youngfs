package iam

type setIAM struct {
	User // user_iam name
	Set  // set_iam name
}

func (setIam *setIAM) Key() string {
	return string(setIam.User) + setIAMKey
}
