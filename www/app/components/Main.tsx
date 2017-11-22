/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as ReactRouter from 'react-router-dom';
import * as Theme from '../Theme';
import * as SubscriptionTypes from '../types/SubscriptionTypes';
import SubscriptionStore from '../stores/SubscriptionStore';
import Loading from './Loading';
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
import * as AuditActions from '../actions/AuditActions';
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
		overflowY: 'hidden',
		userSelect: 'none',
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
	loading: {
		margin: '0 5px 0 1px',
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

		if (!this.state.subscription.active) {
			return <Subscription/>;
		}

		return <ReactRouter.HashRouter>
			<div>
				<nav className="pt-navbar layout horizontal" style={css.nav}>
					<div className="pt-navbar-group pt-align-left flex">
						<div className="pt-navbar-heading"
							style={css.heading}
						>Pritunl Zero</div>
						<Loading style={css.loading} size="small"/>
					</div>
					<div className="pt-navbar-group pt-align-right">
						<ReactRouter.Link
							className="pt-button pt-minimal pt-icon-people"
							style={css.link}
							to="/users"
						>
							Users
						</ReactRouter.Link>
						<ReactRouter.Link
							className="pt-button pt-minimal pt-icon-cloud"
							style={css.link}
							to="/services"
						>
							Services
						</ReactRouter.Link>
						<ReactRouter.Link
							className="pt-button pt-minimal pt-icon-layers"
							style={css.link}
							to="/nodes"
						>
							Nodes
						</ReactRouter.Link>
						<ReactRouter.Link
							className="pt-button pt-minimal pt-icon-filter"
							style={css.link}
							to="/policies"
						>
							Policies
						</ReactRouter.Link>
						<ReactRouter.Link
							className="pt-button pt-minimal pt-icon-office"
							style={css.link}
							to="/authorities"
						>
							Authorities
						</ReactRouter.Link>
						<ReactRouter.Link
							className="pt-button pt-minimal pt-icon-endorsed"
							style={css.link}
							to="/certificates"
						>
							Certificates
						</ReactRouter.Link>
						<ReactRouter.Link
							className="pt-button pt-minimal pt-icon-history"
							style={css.link}
							to="/logs"
						>
							Logs
						</ReactRouter.Link>
						<ReactRouter.Link
							className="pt-button pt-minimal pt-icon-cog"
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
								className="pt-button pt-minimal pt-icon-credit-card"
								style={css.link}
								onClick={(): void => {
									SubscriptionActions.sync(true);
								}}
							>Subscription</button>
						</ReactRouter.Link>
						<ReactRouter.Route render={(props) => (
							<button
								className="pt-button pt-minimal pt-icon-refresh"
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
							className="pt-button pt-minimal pt-icon-log-out"
							onClick={() => {
								window.location.href = '/logout';
							}}
						>Logout</button>
						<button
							className="pt-button pt-minimal pt-icon-moon"
							onClick={(): void => {
								Theme.toggle();
								Theme.save();
							}}
						/>
					</div>
				</nav>
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
