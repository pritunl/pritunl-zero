/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as ReactRouter from 'react-router-dom';
import * as Theme from '../Theme';
import * as SubscriptionTypes from '../types/SubscriptionTypes';
import SubscriptionStore from '../stores/SubscriptionStore';
import LoadingBar from './LoadingBar';
import Subscription from './Subscription';
import Users from './Users';
import UserDetailed from './UserDetailed';
import Nodes from './Nodes';
import Policies from './Policies';
import Authorities from './Authorities';
import Certificates from './Certificates';
import Endpoints from './Endpoints';
import Alerts from './Alerts';
import Checks from './Checks';
import Logs from './Logs';
import Services from './Services';
import Settings from './Settings';
import * as UserActions from '../actions/UserActions';
import * as SessionActions from '../actions/SessionActions';
import * as DeviceActions from '../actions/DeviceActions';
import * as AlertActions from '../actions/AlertActions';
import * as CheckActions from '../actions/CheckActions';
import * as AuditActions from '../actions/AuditActions';
import * as SshcertificateActions from '../actions/SshcertificateActions';
import * as NodeActions from '../actions/NodeActions';
import * as PolicyActions from '../actions/PolicyActions';
import * as AuthorityActions from '../actions/AuthorityActions';
import * as CertificateActions from '../actions/CertificateActions';
import * as EndpointActions from '../actions/EndpointActions';
import * as LogActions from '../actions/LogActions';
import * as ServiceActions from '../actions/ServiceActions';
import * as SettingsActions from '../actions/SettingsActions';
import * as SubscriptionActions from '../actions/SubscriptionActions';

interface State {
	subscription: SubscriptionTypes.SubscriptionRo;
	disabled: boolean;
}

const css = {
	nav: {
		overflowX: 'auto',
		overflowY: 'auto',
		userSelect: 'none',
		height: 'auto',
	} as React.CSSProperties,
	navTitle: {
		height: 'auto',
	} as React.CSSProperties,
	navGroup: {
		flexWrap: 'wrap',
		height: 'auto',
		padding: '10px 0',
	} as React.CSSProperties,
	link: {
		padding: '0 8px',
		color: 'inherit',
	} as React.CSSProperties,
	sub: {
		color: 'inherit',
	} as React.CSSProperties,
	heading: {
		marginRight: '11px',
		fontSize: '18px',
		fontWeight: 'bold',
		width: '100px',
	} as React.CSSProperties,
};

export default class Main extends React.Component<{}, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			subscription: SubscriptionStore.subscription,
			disabled: false,
		};
	}

	componentDidMount(): void {
		SubscriptionStore.addChangeListener(this.onChange);
		SubscriptionActions.sync(false);
	}

	componentWillUnmount(): void {
		SubscriptionStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		this.setState({
			...this.state,
			subscription: SubscriptionStore.subscription,
		});
	}

	render(): JSX.Element {
		if (!this.state.subscription) {
			return <div/>;
		}

		return <ReactRouter.HashRouter>
			<div>
				<nav className="bp3-navbar layout horizontal" style={css.nav}>
					<div
						className="bp3-navbar-group bp3-align-left flex"
						style={css.navTitle}
					>
						<img className="logo-light bp3-navbar-heading" style={css.heading} src="data:image/png;base64,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" alt="Pritunl Zero"/>
						<img className="logo-dark bp3-navbar-heading" style={css.heading} src="data:image/png;base64,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" alt="Pritunl Zero"/>
					</div>
					<div className="bp3-navbar-group bp3-align-right" style={css.navGroup}>
						<ReactRouter.Link
							className="bp3-button bp3-minimal bp3-icon-people"
							style={css.link}
							to="/users"
						>
							Users
						</ReactRouter.Link>
						<ReactRouter.Link
							className="bp3-button bp3-minimal bp3-icon-cloud"
							style={css.link}
							to="/services"
						>
							Services
						</ReactRouter.Link>
						<ReactRouter.Link
							className="bp3-button bp3-minimal bp3-icon-layers"
							style={css.link}
							to="/nodes"
						>
							Nodes
						</ReactRouter.Link>
						<ReactRouter.Link
							className="bp3-button bp3-minimal bp3-icon-filter"
							style={css.link}
							to="/policies"
						>
							Policies
						</ReactRouter.Link>
						<ReactRouter.Link
							className="bp3-button bp3-minimal bp3-icon-office"
							style={css.link}
							to="/authorities"
						>
							Authorities
						</ReactRouter.Link>
						<ReactRouter.Link
							className="bp3-button bp3-minimal bp3-icon-endorsed"
							style={css.link}
							to="/certificates"
						>
							Certificates
						</ReactRouter.Link>
						<ReactRouter.Link
							className="bp3-button bp3-minimal bp3-icon-shield"
							style={css.link}
							to="/endpoints"
						>
							Endpoints
						</ReactRouter.Link>
						<ReactRouter.Link
							className="bp3-button bp3-minimal bp3-icon-notifications"
							style={css.link}
							to="/alerts"
						>
							Alerts
						</ReactRouter.Link>
						<ReactRouter.Link
							className="bp3-button bp3-minimal bp3-icon-lifesaver"
							style={css.link}
							to="/checks"
						>
							Health Checks
						</ReactRouter.Link>
						<ReactRouter.Link
							className="bp3-button bp3-minimal bp3-icon-history"
							style={css.link}
							to="/logs"
						>
							Logs
						</ReactRouter.Link>
						<ReactRouter.Link
							className="bp3-button bp3-minimal bp3-icon-cog"
							style={css.link}
							to="/settings"
						>
							Settings
						</ReactRouter.Link>
						<ReactRouter.Link
							to="/subscription"
							style={css.sub}
						>
							<button
								className="bp3-button bp3-minimal bp3-icon-credit-card"
								style={css.link}
								onClick={(): void => {
									SubscriptionActions.sync(true);
								}}
							>Subscription</button>
						</ReactRouter.Link>
						<ReactRouter.Route render={(props) => (
							<button
								className="bp3-button bp3-minimal bp3-icon-refresh"
								disabled={this.state.disabled}
								onClick={() => {
									let pathname = props.location.pathname;

									this.setState({
										...this.state,
										disabled: true,
									});

									if (pathname === '/users') {
										UserActions.sync().then((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										}).catch((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										});
									} else if (pathname.startsWith('/user/')) {
										UserActions.reload().then((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										}).catch((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										});
										SessionActions.reload().then((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										}).catch((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										});
										DeviceActions.reload().then((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										}).catch((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										});
										SshcertificateActions.reload().then((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										}).catch((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										});
										AuditActions.reload().then((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										}).catch((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										});
									} else if (pathname === '/nodes') {
										ServiceActions.syncNames();
										NodeActions.sync().then((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										}).catch((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										});
									} else if (pathname === '/policies') {
										ServiceActions.syncNames();
										AuthorityActions.sync();
										SettingsActions.sync();
										PolicyActions.sync().then((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										}).catch((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										});
									} else if (pathname === '/authorities') {
										AuthorityActions.sync().then((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										}).catch((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										});
									} else if (pathname === '/certificates') {
										CertificateActions.sync().then((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										}).catch((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										});
									} else if (pathname === '/alerts') {
										AlertActions.sync().then((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										}).catch((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										});
									} else if (pathname === '/checks') {
										CheckActions.sync().then((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										}).catch((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										});
									} else if (pathname === '/endpoints') {
										AuthorityActions.sync();
										EndpointActions.sync().then((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										}).catch((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										});
									} else if (pathname === '/logs') {
										LogActions.sync().then((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										}).catch((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										});
									} else if (pathname === '/services') {
										AuthorityActions.sync();
										ServiceActions.sync().then((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										}).catch((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										});
									} else if (pathname === '/settings') {
										SettingsActions.sync().then((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										}).catch((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										});
									} else if (pathname === '/subscription') {
										SubscriptionActions.sync(true).then((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										}).catch((): void => {
											this.setState({
												...this.state,
												disabled: false,
											});
										});
									} else {
										this.setState({
											...this.state,
											disabled: false,
										});
									}
								}}
							>Refresh</button>
						)}/>
						<button
							className="bp3-button bp3-minimal bp3-icon-log-out"
							onClick={() => {
								window.location.href = '/logout';
							}}
						>Logout</button>
						<button
							className="bp3-button bp3-minimal bp3-icon-moon"
							onClick={(): void => {
								Theme.toggle();
								Theme.save();
							}}
						/>
					</div>
				</nav>
				<LoadingBar intent="primary"/>
				<ReactRouter.Route path="/" exact={true} render={() => (
					<Users/>
				)}/>
				<ReactRouter.Route path="/users" render={() => (
					<Users/>
				)}/>
				<ReactRouter.Route exact path="/user" render={() => (
					<UserDetailed/>
				)}/>
				<ReactRouter.Route path="/user/:userId" render={(props) => (
					<UserDetailed userId={props.match.params.userId}/>
				)}/>
				<ReactRouter.Route path="/nodes" render={() => (
					<Nodes/>
				)}/>
				<ReactRouter.Route path="/policies" render={() => (
					<Policies/>
				)}/>
				<ReactRouter.Route path="/authorities" render={() => (
					<Authorities/>
				)}/>
				<ReactRouter.Route path="/certificates" render={() => (
					<Certificates/>
				)}/>
				<ReactRouter.Route path="/alerts" render={() => (
					<Alerts/>
				)}/>
				<ReactRouter.Route path="/checks" render={() => (
					<Checks/>
				)}/>
				<ReactRouter.Route path="/endpoints" render={() => (
					<Endpoints/>
				)}/>
				<ReactRouter.Route path="/logs" render={() => (
					<Logs/>
				)}/>
				<ReactRouter.Route path="/services" render={() => (
					<Services/>
				)}/>
				<ReactRouter.Route path="/settings" render={() => (
					<Settings/>
				)}/>
				<ReactRouter.Route path="/subscription" render={() => (
					<Subscription/>
				)}/>
			</div>
		</ReactRouter.HashRouter>;
	}
}
