/// <reference path="./References.d.ts"/>
import * as UserActions from './actions/UserActions';
import * as SessionActions from './actions/SessionActions';
import * as DeviceActions from './actions/DeviceActions';
import * as AlertActions from './actions/AlertActions';
import * as CheckActions from './actions/CheckActions';
import * as AuditActions from './actions/AuditActions';
import * as SshcertificateActions from './actions/SshcertificateActions';
import * as NodeActions from './actions/NodeActions';
import * as PolicyActions from './actions/PolicyActions';
import * as AuthorityActions from './actions/AuthorityActions';
import * as CertificateActions from './actions/CertificateActions';
import * as SecretActions from './actions/SecretActions';
import * as EndpointActions from './actions/EndpointActions';
import * as LogActions from './actions/LogActions';
import * as ServiceActions from './actions/ServiceActions';
import * as SettingsActions from './actions/SettingsActions';
import * as SubscriptionActions from './actions/SubscriptionActions';

export function setLocation(location: string) {
	window.location.hash = location
	let evt = new Event("router_update")
	window.dispatchEvent(evt)
}

export function reload() {
	let evt = new Event("router_update")
	window.dispatchEvent(evt)
}

export function refresh(callback?: () => void) {
	let pathname = window.location.hash.replace(/^#/, '');

	if (pathname === '/users') {
		UserActions.sync().then((): void => {
			if (callback) {
				callback()
			}
		}).catch((): void => {
			if (callback) {
				callback()
			}
		});
	} else if (pathname.startsWith('/user/')) {
		UserActions.reload().then((): void => {
			if (callback) {
				callback()
			}
		}).catch((): void => {
			if (callback) {
				callback()
			}
		});
		SessionActions.reload().then((): void => {
			if (callback) {
				callback()
			}
		}).catch((): void => {
			if (callback) {
				callback()
			}
		});
		DeviceActions.reload().then((): void => {
			if (callback) {
				callback()
			}
		}).catch((): void => {
			if (callback) {
				callback()
			}
		});
		SshcertificateActions.reload().then((): void => {
			if (callback) {
				callback()
			}
		}).catch((): void => {
			if (callback) {
				callback()
			}
		});
		AuditActions.reload().then((): void => {
			if (callback) {
				callback()
			}
		}).catch((): void => {
			if (callback) {
				callback()
			}
		});
	} else if (pathname === '/nodes') {
		ServiceActions.syncNames();
		NodeActions.sync().then((): void => {
			if (callback) {
				callback()
			}
		}).catch((): void => {
			if (callback) {
				callback()
			}
		});
	} else if (pathname === '/policies') {
		ServiceActions.syncNames();
		AuthorityActions.sync();
		SettingsActions.sync();
		PolicyActions.sync().then((): void => {
			if (callback) {
				callback()
			}
		}).catch((): void => {
			if (callback) {
				callback()
			}
		});
	} else if (pathname === '/authorities') {
		AuthorityActions.sync().then((): void => {
			if (callback) {
				callback()
			}
		}).catch((): void => {
			if (callback) {
				callback()
			}
		});
	} else if (pathname === '/certificates') {
		CertificateActions.sync().then((): void => {
			if (callback) {
				callback()
			}
		}).catch((): void => {
			if (callback) {
				callback()
			}
		});
	} else if (pathname === '/secrets') {
		SecretActions.sync().then((): void => {
			if (callback) {
				callback()
			}
		}).catch((): void => {
			if (callback) {
				callback()
			}
		});
	} else if (pathname === '/alerts') {
		AlertActions.sync().then((): void => {
			if (callback) {
				callback()
			}
		}).catch((): void => {
			if (callback) {
				callback()
			}
		});
	} else if (pathname === '/checks') {
		CheckActions.sync().then((): void => {
			if (callback) {
				callback()
			}
		}).catch((): void => {
			if (callback) {
				callback()
			}
		});
	} else if (pathname === '/endpoints') {
		AuthorityActions.sync();
		EndpointActions.sync().then((): void => {
			if (callback) {
				callback()
			}
		}).catch((): void => {
			if (callback) {
				callback()
			}
		});
	} else if (pathname === '/logs') {
		LogActions.sync().then((): void => {
			if (callback) {
				callback()
			}
		}).catch((): void => {
			if (callback) {
				callback()
			}
		});
	} else if (pathname === '/services') {
		AuthorityActions.sync();
		ServiceActions.sync().then((): void => {
			if (callback) {
				callback()
			}
		}).catch((): void => {
			if (callback) {
				callback()
			}
		});
	} else if (pathname === '/settings') {
		SettingsActions.sync().then((): void => {
			if (callback) {
				callback()
			}
		}).catch((): void => {
			if (callback) {
				callback()
			}
		});
	} else if (pathname === '/subscription') {
		SubscriptionActions.sync(true).then((): void => {
			if (callback) {
				callback()
			}
		}).catch((): void => {
			if (callback) {
				callback()
			}
		});
	} else {
		this.setState({
			...this.state,
			disabled: false,
		});
	}
}
