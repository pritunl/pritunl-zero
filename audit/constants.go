package audit

const (
	AdminLogin                = "admin_login"
	AdminLoginFailed          = "admin_login_failed"
	ProxyLogin                = "proxy_login"
	ProxyLoginFailed          = "proxy_login_failed"
	UserLogin                 = "user_login"
	UserLoginFailed           = "user_login_failed"
	AdminDeviceRegister       = "admin_device_register"
	DeviceRegister            = "device_register"
	DeviceRegisterFailed      = "device_register_failed"
	DuoApprove                = "duo_approve"
	DuoDeny                   = "duo_deny"
	OneLoginApprove           = "one_login_approve"
	OneLoginDeny              = "one_login_deny"
	OktaApprove               = "okta_approve"
	OktaDeny                  = "okta_deny"
	SshApprove                = "ssh_approve"
	SshDeny                   = "ssh_deny"
	KeybaseAssociationApprove = "keybase_association_approve"
	KeybaseAssociationDeny    = "keybase_association_deny"
)
