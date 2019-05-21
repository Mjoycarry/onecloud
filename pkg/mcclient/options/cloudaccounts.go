package options

type SUserPasswordCredential struct {
	Username string `help:"Username" positional:"true"`
	Password string `help:"Password" positional:"true"`
}

type SVMwareCredentialWithEnvironment struct {
	SUserPasswordCredential

	Host string `help:"VMware VCenter/ESXi host" positional:"true"`
	Port string `help:"VMware VCenter/ESXi host port" default:"443"`
}

type SAzureCredential struct {
	ClientID     string `help:"Azure client_id" positional:"true"`
	ClientSecret string `help:"Azure clinet_secret" positional:"true"`
}

type SAzureCredentialWithEnvironment struct {
	DirectoryID string `help:"Azure directory_id" positional:"true"`

	SAzureCredential

	Environment string `help:"Cloud environment" choices:"AzureGermanCloud|AzureChinaCloud|AzurePublicCloud" default:"AzureChinaCloud"`
}

type SQcloudCredential struct {
	AppID     string `help:"Qcloud appid" positional:"true"`
	SecretID  string `help:"Qcloud secret_id" positional:"true"`
	SecretKey string `help:"Qcloud secret_key" positional:"true"`
}

type SOpenStackCredential struct {
	ProjectName string `help:"OpenStack project_name" positional:"true"`

	SUserPasswordCredential

	DomainName string `help:"OpenStack domain name"`
}

type SOpenStackCredentialWithAuthURL struct {
	SOpenStackCredential

	AuthURL string `help:"OpenStack auth_url" positional:"true" json:"auth_url"`
}

type SAccessKeyCredential struct {
	AccessKeyID     string `help:"Access_key_id" positional:"true"`
	AccessKeySecret string `help:"Access_key_secret" positional:"true"`
}

type SAccessKeyCredentialWithEnvironment struct {
	SAccessKeyCredential
	Environment string `help:"Cloud environment" choices:"InternationalCloud|ChinaCloud" default:"ChinaCloud"`
}

/// create options

type SCloudAccountCreateBaseOptions struct {
	Name string `help:"Name of cloud account" positional:"true"`
	// PROVIDER string `help:"Driver for cloud account" choices:"VMware|Aliyun|Azure|Qcloud|OpenStack|Huawei|Aws"`
	Desc string `help:"Description" token:"desc" json:"description"`

	AutoCreateProject bool `help:"Enable the account with same name project"`
	EnableAutoSync    bool `help:"Enable automatically synchronize resources of this account"`

	SyncIntervalSeconds int `help:"Interval to synchronize if auto sync is enable" metavar:"SECONDS"`
}

type SVMwareCloudAccountCreateOptions struct {
	SCloudAccountCreateBaseOptions
	SVMwareCredentialWithEnvironment
}

type SAliyunCloudAccountCreateOptions struct {
	SCloudAccountCreateBaseOptions
	SAccessKeyCredential
}

type SAzureCloudAccountCreateOptions struct {
	SCloudAccountCreateBaseOptions
	SAzureCredentialWithEnvironment
}

type SQcloudCloudAccountCreateOptions struct {
	SCloudAccountCreateBaseOptions
	SQcloudCredential
}

type SAWSCloudAccountCreateOptions struct {
	SCloudAccountCreateBaseOptions
	SAccessKeyCredentialWithEnvironment
}

type SOpenStackCloudAccountCreateOptions struct {
	SCloudAccountCreateBaseOptions
	SOpenStackCredentialWithAuthURL
}

type SHuaweiCloudAccountCreateOptions struct {
	SCloudAccountCreateBaseOptions
	SAccessKeyCredentialWithEnvironment
}

// update credential options

type SCloudAccountUpdateCredentialBaseOptions struct {
	ID string `help:"ID or Name of cloud account" json:"-"`
}

type SVMwareCloudAccountUpdateCredentialOptions struct {
	SCloudAccountUpdateCredentialBaseOptions
	SUserPasswordCredential
}

type SAliyunCloudAccountUpdateCredentialOptions struct {
	SCloudAccountUpdateCredentialBaseOptions
	SAccessKeyCredential
}

type SAzureCloudAccountUpdateCredentialOptions struct {
	SCloudAccountUpdateCredentialBaseOptions
	SAzureCredential
}

type SQcloudCloudAccountUpdateCredentialOptions struct {
	SCloudAccountUpdateCredentialBaseOptions
	SQcloudCredential
}

type SAWSCloudAccountUpdateCredentialOptions struct {
	SCloudAccountUpdateCredentialBaseOptions
	SAccessKeyCredential
}

type SOpenStackCloudAccountUpdateCredentialOptions struct {
	SCloudAccountUpdateCredentialBaseOptions
	SOpenStackCredential
}

type SHuaweiCloudAccountUpdateCredentialOptions struct {
	SCloudAccountUpdateCredentialBaseOptions
	SAccessKeyCredential
}
