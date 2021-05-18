/// <reference path="../References.d.ts"/>
import * as React from 'react';
import PageNumInput from './PageNumInput';
import PageSelect from './PageSelect';
import EndpointChart from './EndpointChart';
import * as EndpointActions from '../actions/EndpointActions';

interface Props {
	endpoint: string;
	disabled: boolean;
}

interface State {
	sync: number;
	period: number;
	interval: number;
	loading: {[key: string]: boolean};
	cancelable: {[key: string]: boolean};
}

const css = {
	header: {
		fontSize: '20px',
		marginTop: '-10px',
		paddingBottom: '2px',
		marginBottom: '10px',
		borderBottomStyle: 'solid',
	} as React.CSSProperties,
	heading: {
		margin: '19px 0 0 0',
	} as React.CSSProperties,
	button: {
		margin: '8px 0 0 8px',
	} as React.CSSProperties,
	buttons: {
		marginTop: '8px',
	} as React.CSSProperties,
	group: {
		flex: 1,
		minWidth: '250px',
		margin: '0 10px',
	} as React.CSSProperties,
	chartGroup: {
		flex: 1,
		minWidth: '250px',
		margin: '0 10px',
		marginBottom: '15px',
	} as React.CSSProperties,
};

export default class EndpointCharts extends React.Component<Props, State> {
	loading: {[key: string]: boolean};

	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			sync: 0,
			period: 1440,
			interval: 1,
			loading: {},
			cancelable: {},
		};

		this.loading = {};
	}

	setLoading(resource: string): void {
		this.loading[resource] = true;

		let loading = {
			...this.state.loading,
		};
		loading[resource] = true;

		setTimeout((): void => {
			if (this.loading[resource]) {
				let cancelable = {
					...this.state.cancelable,
				};
				cancelable[resource] = true;

				this.setState({
					...this.state,
					cancelable: cancelable,
				});
			}
		}, 3000);

		this.setState({
			...this.state,
			loading: loading,
		});
	}

	setLoaded(resource: string): void {
		delete this.loading[resource];

		let loading = {
			...this.state.loading,
		};
		delete loading[resource];

		let cancelable = {
			...this.state.cancelable,
		};
		delete cancelable[resource];

		this.setState({
			...this.state,
			loading: loading,
			cancelable: cancelable,
		});
	}

	render(): JSX.Element {
		if (this.props.disabled) {
			return <div/>;
		}

		let refreshDisabled = false;
		let refreshLabel = '';
		let refreshClass = 'bp3-button';
		if (Object.entries(this.state.cancelable).length) {
			refreshLabel = 'Cancel';
			refreshClass += ' bp3-intent-warning bp3-icon-delete'
		} else {
			if (Object.entries(this.state.loading).length) {
				refreshDisabled = true;
			}
			refreshLabel = 'Refresh';
			refreshClass += ' bp3-intent-success bp3-icon-refresh'
		}

		return <div>
			<div className="layout horizontal wrap bp3-border" style={css.header}>
				<h3 style={css.heading}>Charts</h3>
				<div className="flex"/>
				<div style={css.buttons}>
					<button
						className={refreshClass}
						style={css.button}
						disabled={refreshDisabled}
						type="button"
						onClick={(): void => {
							if (Object.entries(this.state.cancelable).length) {
								EndpointActions.chartCancel();
							} else {
								this.setState({
									...this.state,
									sync: this.state.sync + 1,
								});
							}
						}}
					>
						{refreshLabel}
					</button>
				</div>
			</div>
			<div className="layout horizontal wrap">
				<div style={css.group}>
					<PageSelect
						label="Time Range"
						help="Select chart time range."
						value={this.state.period.toString()}
						onChange={(val: string): void => {
							this.setState({
								...this.state,
								period: parseInt(val, 10),
							});
						}}
					>
						<option value="1440">24 hours</option>
						<option value="4320">3 days</option>
						<option value="10080">7 days</option>
						<option value="20160">14 days</option>
						<option value="43200">30 days</option>
						<option value="86400">60 days</option>
						<option value="129600">90 days</option>
						<option value="172800">120 days</option>
					</PageSelect>
				</div>
				<div style={css.group}>
					<PageSelect
						label="Interval"
						help="Select chart interval."
						value={this.state.interval.toString()}
						onChange={(val: string): void => {
							this.setState({
								...this.state,
								interval: parseInt(val, 10),
							});
						}}
					>
						<option value="1">1 minute</option>
						<option value="5">5 minutes</option>
						<option value="30">30 minutes</option>
						<option value="60">1 hour</option>
						<option value="120">2 hours</option>
						<option value="360">6 hours</option>
						<option value="720">12 hours</option>
						<option value="1440">24 hours</option>
						<option value="4320">3 days</option>
						<option value="10080">7 days</option>
					</PageSelect>
					<PageNumInput
						label="Hours Select"
						help="Select time range to view."
						min={1}
						minorStepSize={1}
						stepSize={6}
						majorStepSize={12}
						selectAllOnFocus={true}
						onChange={(val: number): void => {
							this.setState({
								...this.state,
								period: val * 60,
							});
						}}
						value={this.state.period * 60}
					/>
				</div>
			</div>
			<div className="layout horizontal wrap">
				<div style={css.chartGroup}>
					<EndpointChart
						endpoint={this.props.endpoint}
						resource={'system'}
						sync={this.state.sync}
						period={this.state.period}
						interval={this.state.interval}
						onLoading={(): void => {
							this.setLoading('chart1');
						}}
						onLoaded={(): void => {
							this.setLoaded('chart1');
						}}
					/>
					<EndpointChart
						endpoint={this.props.endpoint}
						resource={'system'}
						sync={this.state.sync}
						period={this.state.period}
						interval={this.state.interval}
						onLoading={(): void => {
							this.setLoading('chart2');
						}}
						onLoaded={(): void => {
							this.setLoaded('chart2');
						}}
					/>
				</div>
				<div style={css.chartGroup}>
					<EndpointChart
						endpoint={this.props.endpoint}
						resource={'system'}
						sync={this.state.sync}
						period={this.state.period}
						interval={this.state.interval}
						onLoading={(): void => {
							this.setLoading('chart3');
						}}
						onLoaded={(): void => {
							this.setLoaded('chart3');
						}}
					/>
					<EndpointChart
						endpoint={this.props.endpoint}
						resource={'system'}
						sync={this.state.sync}
						period={this.state.period}
						interval={this.state.interval}
						onLoading={(): void => {
							this.setLoading('chart4');
						}}
						onLoaded={(): void => {
							this.setLoaded('chart4');
						}}
					/>
				</div>
			</div>
		</div>;
	}
}
