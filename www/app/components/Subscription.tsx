/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as SubscriptionActions from '../actions/SubscriptionActions';
import * as SubscriptionTypes from '../types/SubscriptionTypes';
import SubscriptionStore from '../stores/SubscriptionStore';
import * as Alert from '../Alert';
import * as MiscUtils from '../utils/MiscUtils';
import ConfirmButton from './ConfirmButton';
import * as Theme from '../Theme';

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
		margin: '30px auto',
	} as React.CSSProperties,
	card2: {
		padding: '5px',
		minWidth: '310px',
		maxWidth: '380px',
		width: 'calc(100% - 20px)',
		margin: '30px auto',
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
		fontSize: Theme.monospaceSize,
		fontFamily: Theme.monospaceFont,
		fontWeight: Theme.monospaceWeight,
	} as React.CSSProperties,
	button: {
		width: '150px',
		margin: '5px',
	} as React.CSSProperties,
	button2: {
		width: '160px',
		margin: '5px',
	} as React.CSSProperties,
	button3: {
		width: '195px',
		margin: '5px',
	} as React.CSSProperties,
	buttons: {
		margin: '0 auto',
	} as React.CSSProperties,
	buttonsEnd: {
		marginBottom: '10px',
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
			<div className="bp5-card bp5-elevation-2" style={css.card}>
				<div
					className="bp5-callout bp5-intent-success"
					style={css.message}
					hidden={!this.state.message}
				>
					{this.state.message}
				</div>
				<textarea
					className="bp5-input"
					style={css.license}
					disabled={this.state.disabled}
					placeholder="New License Key"
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
						className="bp5-button bp5-intent-danger bp5-icon-cross"
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
						className="bp5-button bp5-intent-primary bp5-icon-endorsed"
						style={css.button}
						disabled={this.state.disabled}
						onClick={(): void => {
							this.setState({
								...this.state,
								disabled: true,
							});
							SubscriptionActions.activate(
								this.state.license,
							).then((): void => {
								this.setState({
									...this.state,
									disabled: false,
									update: false,
									license: '',
								});
							}).catch((): void => {
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
			<div className="bp5-card bp5-elevation-2" style={css.card}>
				<div
					className="bp5-callout bp5-intent-success"
					style={css.message}
					hidden={!this.state.message}
				>
					{this.state.message}
				</div>
				<textarea
					className="bp5-input"
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
						className="bp5-button bp5-intent-primary bp5-icon-endorsed"
						style={css.button}
						disabled={this.state.disabled}
						onClick={(): void => {
							this.setState({
								...this.state,
								disabled: true,
							});
							SubscriptionActions.activate(
								this.state.license,
							).then((): void => {
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
						}}
					>Activate License</button>
					<a
						className="bp5-button bp5-intent-success bp5-icon-credit-card"
						style={css.button2}
						target="_blank"
						href="https://app.pritunl.com/checkout/zero"
					>Create Account</a>
				</div>
			</div>
		</div>;
	}

	reactivate(): JSX.Element {
		let sub = this.state.subscription;
		let canceling = sub.cancel_at_period_end || sub.status === 'canceled';
		let status = sub.cancel_at_period_end ? 'canceled' : sub.status;
		let periodEnd = MiscUtils.formatDateShort(sub.period_end);
		let trialEnd = MiscUtils.formatDateShort(sub.trial_end);

		let balance: string;
		let balanceLabel: string;
		if (sub.balance < 0) {
			balance = MiscUtils.formatAmount(sub.balance * -1);
			balanceLabel = 'Credit';
		} else {
			balance = MiscUtils.formatAmount(sub.balance);
			balanceLabel = 'Balance';
		}

		return <div>
			<div className="bp5-card bp5-elevation-2" style={css.card2}>
				<div
					className="bp5-callout bp5-intent-success"
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
						<div className="flex">{balanceLabel}:</div>
						<div>
							{balance}
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
					<a
						className="bp5-button bp5-intent-success bp5-icon-cog"
						style={css.button2}
						target="_blank"
						href={"https://app.pritunl.com/subscription/" +
							this.state.subscription.url_key}
					>Manage Account</a>
				</div>
				<div
					className="layout horizontal center-justified"
					style={css.buttonsEnd}
				>
					<ConfirmButton
						className="bp5-intent-danger bp5-icon-delete"
						progressClassName="bp5-intent-danger"
						style={css.button2}
						disabled={this.state.disabled}
						safe={true}
						label="Remove License"
						onConfirm={(): void => {
							this.setState({
								...this.state,
								disabled: true,
							});
							SubscriptionActions.activate('').then((): void => {
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
						}}
					/>
					<button
						className="bp5-button bp5-intent-primary bp5-icon-endorsed"
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
