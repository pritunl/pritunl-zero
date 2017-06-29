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
	license: {
		width: '100%',
		height: '130px',
		margin: '10px 0',
		fontFamily: '"Lucida Console", Monaco, monospace',
	} as React.CSSProperties,
};

export default class Subscription extends React.Component<{}, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			message: 'test',
			license: '',
		};
	}

	render(): JSX.Element {
		return <div>
			<div className="pt-card pt-elevation-2" style={css.card}>
				<div
					className="pt-callout pt-intent-success"
					hidden={!this.state.message}
				>
					{this.state.message}
				</div>
				<textarea
					className="pt-input"
					style={css.license}
					onChange={(evt): void => {
						this.setState({
							...this.state,
							license: evt.target.value,
						})
					}}
				>{this.state.license}</textarea>
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
						console.log(token);
					}}
					onScriptError={(err): void => {
						console.log(err);
					}}
					stripeKey="pk_test_4YSuzxPmd08oSV2s4kLi7zU2"
				>
					<button
						className="pt-button pt-minimal pt-icon-checkout"
					>Subscribe</button>
				</ReactStripeCheckout>
			</div>
		</div>;
	}
}
