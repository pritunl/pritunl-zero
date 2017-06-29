/// <reference path="../References.d.ts"/>
import * as React from 'react';
import ReactStripeCheckout from 'react-stripe-checkout';

const css = {
	card: {
		padding: '50px',
		minWidth: '200px',
		maxWidth: '300px',
		margin: '0 auto',
		position: 'absolute',
		top: '50%',
		left: '50%',
		transform: 'translate(-50%, -50%)',
	} as React.CSSProperties,
};

export default class Subscription extends React.Component<{}, {}> {
	render(): JSX.Element {
		return <div>
			<div className="pt-card pt-elevation-2" style={css.card}>
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
