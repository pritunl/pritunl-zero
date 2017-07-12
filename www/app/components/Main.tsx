/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as ReactRouter from 'react-router-dom';
import * as SubscriptionTypes from '../types/SubscriptionTypes';
import * as SubscriptionActions from '../actions/SubscriptionActions';
import SubscriptionStore from '../stores/SubscriptionStore';
import Loading from './Loading';
import Subscription from './Subscription';
import Users from './Users';
import UserDetailed from './UserDetailed';
import Settings from './Settings';

document.body.className = 'root pt-dark';

interface State {
	subscription: SubscriptionTypes.SubscriptionRo;
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
				<ReactRouter.Route path="/users" render={() => (
					<Users/>
				)}/>
				<ReactRouter.Route exact path="/user" render={() => (
					<UserDetailed/>
				)}/>
				<ReactRouter.Route path="/user/:userId" render={(props) => (
					<UserDetailed userId={props.match.params.userId}/>
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
