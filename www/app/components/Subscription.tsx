/// <reference path="../References.d.ts"/>
import * as React from 'react';
import ReactStripeCheckout from 'react-stripe-checkout';
import * as SubscriptionActions from '../actions/SubscriptionActions';
import * as Alert from '../Alert';

interface State {
	message: string;
	license: string;
}

const css = {
	card: {
		padding: '10px',
		minWidth: '310px',
		maxWidth: '350px',
		width: 'calc(100% - 20px)',
		margin: '0',
		position: 'absolute',
		top: '50%',
		left: '50%',
		transform: 'translate(-50%, -50%)',
	} as React.CSSProperties,
	message: {
		margin: '0 0 10px 0',
	} as React.CSSProperties,
	license: {
		width: '100%',
		height: '130px',
		margin: '0 0 10px 0',
		resize: 'none',
		fontFamily: '"Lucida Console", Monaco, monospace',
	} as React.CSSProperties,
	button: {
		marginRight: '10px',
	} as React.CSSProperties,
	buttons: {
		margin: '0 auto',
	} as React.CSSProperties,
};

export default class Subscription extends React.Component<{}, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			message: '',
			license: '',
		};
	}

	render(): JSX.Element {
		return <div>
			<div className="pt-card pt-elevation-2" style={css.card}>
				<div
					className="pt-callout pt-intent-success"
					style={css.message}
					hidden={!this.state.message}
				>
					{this.state.message}
				</div>
				<textarea
					className="pt-input"
					style={css.license}
					value={this.state.license}
					onChange={(evt): void => {
						this.setState({
							...this.state,
							license: evt.target.value,
						})
					}}
				/>
				<div className="layout horizontal center-justified">
					<button
						className="pt-button pt-icon-endorsed"
						style={css.button}
					>Activate License Key</button>
					<ReactStripeCheckout
						label="Pritunl Zero"
						image="//s3.amazonaws.com/pritunl-static/logo_stripe.png"
						allowRememberMe={false}
						zipCode={true}
						amount={5000}
						name="Pritunl Zero"
						description="Subscribe to Zero ($50/month)"
						panelLabel="Subscribe"
						token={(token): void => {
							SubscriptionActions.checkout(
								'zero',
								token.id,
								token.email,
							).then((message: string): void => {
								this.setState({
									...this.state,
									message: message,
								});
							});
						}}
						onScriptError={(err): void => {
							Alert.error('Failed to load Stripe Checkout');
						}}
						stripeKey="pk_test_4YSuzxPmd08oSV2s4kLi7zU2"
					>
						<button
							className="pt-button pt-icon-credit-card"
						>Subscribe</button>
					</ReactStripeCheckout>
				</div>
			</div>
		</div>;
	}
}
