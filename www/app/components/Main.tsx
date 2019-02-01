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
import Logs from './Logs';
import Services from './Services';
import Settings from './Settings';
import * as UserActions from '../actions/UserActions';
import * as SessionActions from '../actions/SessionActions';
import * as DeviceActions from '../actions/DeviceActions';
import * as AuditActions from '../actions/AuditActions';
import * as SshcertificateActions from '../actions/SshcertificateActions';
import * as NodeActions from '../actions/NodeActions';
import * as PolicyActions from '../actions/PolicyActions';
import * as AuthorityActions from '../actions/AuthorityActions';
import * as CertificateActions from '../actions/CertificateActions';
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
						<div className="bp3-navbar-heading"
							style={css.heading}
						>Pritunl Zero</div>
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
										ServiceActions.sync();
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
										ServiceActions.sync();
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
