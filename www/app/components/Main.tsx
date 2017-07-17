/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as ReactRouter from 'react-router-dom';
import * as SubscriptionTypes from '../types/SubscriptionTypes';
import SubscriptionStore from '../stores/SubscriptionStore';
import Loading from './Loading';
import Subscription from './Subscription';
import Users from './Users';
import UserDetailed from './UserDetailed';
import Nodes from './Nodes';
import Services from './Services';
import Settings from './Settings';
import * as UserActions from '../actions/UserActions';
import * as NodeActions from '../actions/NodeActions';
import * as ServiceActions from '../actions/ServiceActions';
import * as SettingsActions from '../actions/SettingsActions';
import * as SubscriptionActions from '../actions/SubscriptionActions';

document.body.className = 'root pt-dark';

interface State {
	subscription: SubscriptionTypes.SubscriptionRo;
	disabled: boolean;
}

const css = {
	nav: {
		overflowX: 'auto',
		overflowY: 'hidden',
	} as React.CSSProperties,
	link: {
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
						<Loading size="small"/>
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
							className="pt-button pt-minimal pt-icon-satellite"
							style={css.link}
							to="/nodes"
						>
							Nodes
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
							style={css.link}
						>
							<button
								className="pt-button pt-minimal pt-icon-credit-card"
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
									} else if (pathname === '/user') {
										// UserActions.load();
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
								let className = 'root';

								if (document.body.className.indexOf('pt-dark') === -1) {
									className += ' pt-dark';
								}

								document.body.className = className;
							}}
						/>
					</div>
				</nav>
				<ReactRouter.Redirect from="/" to="/users"/>
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
