/// <reference path="../References.d.ts"/>
import * as React from 'react';
import ReactStripeCheckout from 'react-stripe-checkout';
import * as SubscriptionActions from '../actions/SubscriptionActions';
import * as SubscriptionTypes from '../types/SubscriptionTypes';
import SubscriptionStore from '../stores/SubscriptionStore';
import * as Alert from '../Alert';
import * as MiscUtils from '../utils/MiscUtils';
import ConfirmButton from './ConfirmButton';

interface State {
	subscription: SubscriptionTypes.SubscriptionRo;
	update: boolean;
	message: string;
	license: string;
	disabled: boolean;
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
	card2: {
		padding: '5px',
		minWidth: '310px',
		maxWidth: '380px',
		width: 'calc(100% - 20px)',
		margin: '0',
		position: 'absolute',
		top: '50%',
		left: '50%',
		transform: 'translate(-50%, -50%)',
	} as React.CSSProperties,
	status: {
		width: '180px',
		margin: '20px auto',
		fontSize: '16px',
	} as React.CSSProperties,
	item: {
		margin: '2px 0',
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
		width: '150px',
		margin: '5px',
	} as React.CSSProperties,
	button2: {
		width: '170px',
		margin: '5px',
	} as React.CSSProperties,
	button3: {
		width: '195px',
		margin: '5px',
	} as React.CSSProperties,
	buttons: {
		margin: '0 auto',
	} as React.CSSProperties,
};

export default class Subscription extends React.Component<{}, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			subscription: SubscriptionStore.subscription,
			update: false,
			message: '',
			license: '',
			disabled: false,
		};
	}

	componentDidMount(): void {
		SubscriptionStore.addChangeListener(this.onChange);
		if (!this.state.subscription.active) {
			SubscriptionActions.sync(true);
		}
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

	update(): JSX.Element {
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
					disabled={this.state.disabled}
					placeholder="New License Key"
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
						className="pt-button pt-intent-danger pt-icon-cross"
						style={css.button}
						disabled={this.state.disabled}
						onClick={(): void => {
							this.setState({
								...this.state,
								update: false,
								license: '',
							});
						}}
					>Cancel</button>
					<button
						className="pt-button pt-intent-primary pt-icon-endorsed"
						style={css.button}
						disabled={this.state.disabled}
						onClick={(): void => {
							this.setState({
								...this.state,
								disabled: true,
							});
							SubscriptionActions.activate(this.state.license).then(
								(): void => {
									this.setState({
										...this.state,
										disabled: false,
										update: false,
										license: '',
									});
								}
							).catch((): void => {
								this.setState({
									...this.state,
									disabled: false,
								});
							});
						}}
					>Update License</button>
				</div>
			</div>
		</div>;
	}

	activate(): JSX.Element {
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
					disabled={this.state.disabled}
					placeholder="License Key"
					value={this.state.license}
					onChange={(evt): void => {
						this.setState({
							...this.state,
							license: evt.target.value,
						});
					}}
				/>
				<div className="layout horizontal center-justified">
					<button
						className="pt-button pt-intent-primary pt-icon-endorsed"
						style={css.button}
						disabled={this.state.disabled}
						onClick={(): void => {
							this.setState({
								...this.state,
								disabled: true,
							});
							SubscriptionActions.activate(this.state.license).then(
								(): void => {
									this.setState({
										...this.state,
										disabled: false,
									});
								}
							).catch((): void => {
								this.setState({
									...this.state,
									disabled: false,
								});
							});
						}}
					>Activate License</button>
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
							this.setState({
								...this.state,
								disabled: true,
							});
							SubscriptionActions.checkout(
								'zero',
								token.id,
								token.email,
							).then((message: string): void => {
								this.setState({
									...this.state,
									disabled: false,
									message: message,
								});
							}).catch((): void => {
								this.setState({
									...this.state,
									disabled: false,
								});
							});
						}}
						onScriptError={(err): void => {
							Alert.error('Failed to load Stripe Checkout');
						}}
						stripeKey="pk_test_4YSuzxPmd08oSV2s4kLi7zU2"
					>
						<button
							className="pt-button pt-intent-success pt-icon-credit-card"
							style={css.button}
							disabled={this.state.disabled}
						>Subscribe</button>
					</ReactStripeCheckout>
				</div>
			</div>
		</div>;
	}

	reactivate(): JSX.Element {
		let sub = this.state.subscription;
		let canceling = sub.cancel_at_period_end || sub.status === 'canceled';
		let status = sub.cancel_at_period_end ? 'Canceled' : sub.status;
		let periodEnd = MiscUtils.formatDateShort(sub.period_end);
		let trialEnd = MiscUtils.formatDateShort(sub.trial_end);

		return <div>
			<div className="pt-card pt-elevation-2" style={css.card2}>
				<div
					className="pt-callout pt-intent-success"
					style={css.message}
					hidden={!this.state.message}
				>
					{this.state.message}
				</div>
				<div className="layout vertical" style={css.status}>
					<div className="layout horizontal">
						<div className="flex">Status:</div>
						<div>
							{MiscUtils.capitalize(status)}
						</div>
					</div>
					<div className="layout horizontal" style={css.item}>
						<div className="flex">Plan:</div>
						<div>
							{MiscUtils.capitalize(sub.plan)}
						</div>
					</div>
					<div className="layout horizontal" style={css.item}>
						<div className="flex">Amount:</div>
						<div>
							{MiscUtils.formatAmount(sub.amount)}
						</div>
					</div>
					<div className="layout horizontal" style={css.item}>
						<div className="flex">Quantity:</div>
						<div>
							{sub.quantity}
						</div>
					</div>
					<div
						className="layout horizontal"
						style={css.item}
						hidden={!sub.balance}
					>
						<div className="flex">Balance:</div>
						<div>
							{MiscUtils.formatAmount(sub.balance)}
						</div>
					</div>
					<div
						className="layout horizontal"
						style={css.item}
						hidden={periodEnd === ''}
					>
						<div className="flex">
							{canceling ? 'Ends' : 'Renew'}:
						</div>
						<div>
							{periodEnd}
						</div>
					</div>
					<div
						className="layout horizontal"
						style={css.item}
						hidden={trialEnd === ''}
					>
						<div className="flex">Trial Ends:</div>
						<div>
							{trialEnd}
						</div>
					</div>
				</div>
				<div className="layout horizontal center-justified">
					<ConfirmButton
						className="pt-intent-danger pt-icon-delete"
						progressClassName="pt-intent-danger"
						style={css.button2}
						disabled={this.state.disabled}
						hidden={canceling}
						label="Cancel Subscription"
						onConfirm={(): void => {
							this.setState({
								...this.state,
								disabled: true,
							});
							SubscriptionActions.cancel(
								this.state.subscription.url_key,
							).then((): void => {
									this.setState({
										...this.state,
										disabled: false,
									});
								}
							).catch((): void => {
								this.setState({
									...this.state,
									disabled: false,
								});
							});
						}}
					/>
					<ReactStripeCheckout
						label="Pritunl Zero"
						image="//s3.amazonaws.com/pritunl-static/logo_stripe.png"
						allowRememberMe={false}
						zipCode={true}
						amount={canceling && sub.status !== 'active' ? 5000 : 0}
						name="Pritunl Zero"
						description={canceling ?
							'Reactivate Subscription ($50/month)' :
							'Update Payment Information'
						}
						panelLabel={canceling ? 'Reactivate' : 'Update'}
						token={(token): void => {
							this.setState({
								...this.state,
								disabled: true,
							});
							SubscriptionActions.payment(
								this.state.subscription.url_key,
								'zero',
								token.id,
								token.email,
							).then(
								(): void => {
									this.setState({
										...this.state,
										disabled: false,
									});
								}
							).catch((): void => {
								this.setState({
									...this.state,
									disabled: false,
								});
							});;
						}}
						onScriptError={(err): void => {
							Alert.error('Failed to load Stripe Checkout');
						}}
						stripeKey="pk_test_4YSuzxPmd08oSV2s4kLi7zU2"
					>
						<button
							className="pt-button pt-intent-success pt-icon-credit-card"
							style={canceling ? css.button3 : css.button2}
							disabled={this.state.disabled}
						>
							{canceling ? 'Reactivate Subscription' : 'Update Payment'}
						</button>
					</ReactStripeCheckout>
				</div>
				<div className="layout horizontal center-justified">
					<ConfirmButton
						className="pt-intent-danger pt-icon-delete"
						progressClassName="pt-intent-danger"
						style={css.button2}
						disabled={this.state.disabled}
						label="Remove License"
						onConfirm={(): void => {
							this.setState({
								...this.state,
								disabled: true,
							});
							SubscriptionActions.activate('').then(
								(): void => {
									this.setState({
										...this.state,
										disabled: false,
									});
								}
							).catch((): void => {
								this.setState({
									...this.state,
									disabled: false,
								});
							});;
						}}
					/>
					<button
						className="pt-button pt-intent-primary pt-icon-endorsed"
						style={css.button2}
						disabled={this.state.disabled}
						onClick={(): void => {
							this.setState({
								...this.state,
								update: true,
							});
						}}
					>Update License</button>
				</div>
			</div>
		</div>;
	}

	render(): JSX.Element {
		if (this.state.update) {
			return this.update();
		} else if (this.state.subscription.status) {
			return this.reactivate();
		} else {
			return this.activate();
		}
	}
}
