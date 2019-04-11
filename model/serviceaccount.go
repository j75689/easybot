package model

// ServiceAccount api access token account
type ServiceAccount struct {
	Name     string `json:"name" bson:"id"`
	EMail    string `json:"email" bson:"email"`
	Domain   string `json:"domain" bson:"domain"`
	Provider string `json:"provider" bson:"provider"`
	Scope    string `json:"scope" bson:"scope"`
	Active   int    `json:"active" bson:"active"`

	Generate int64  `json:"generate" bson:"generate"`
	Expired  int64  `json:"expired" bson:"expired"`
	Token    string `json:"token" bson:"token"`
}
