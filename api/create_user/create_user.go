package create_user

type CreateUserInfo struct {
	BucketName string `form:"BucketName" json:"BucketName" uri:"BucketName" binding:"required"`
	AccessKey  string `form:"AccessKey" json:"AccessKey" uri:"AccessKey" binding:"required"`
	ObjectName string `form:"ObjectName" json:"ObjectName" uri:"ObjectName" binding:"required"`
	DataTime   string `form:"DataTime" json:"DataTime" uri:"DataTime" binding:"required"`
}
