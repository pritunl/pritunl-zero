package audit

const (
	AdminLogin                = "admin_login"
	AdminLoginFailed          = "admin_login_failed"
	Login                     = "login"
	LoginFailed               = "login_failed"
	DuoApprove                = "duo_approve"
	DuoDeny                   = "duo_deny"
	OneLoginApprove           = "one_login_approve"
	OneLoginDeny              = "one_login_deny"
	SshApprove                = "ssh_approve"
	SshDeny                   = "ssh_deny"
	KeybaseAssociationApprove = "keybase_association_approve"
	KeybaseAssociationDeny    = "keybase_association_deny"
)
