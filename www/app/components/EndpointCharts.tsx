/// <reference path="../References.d.ts"/>
import * as React from 'react';
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
	chartBoxRef: React.RefObject<HTMLDivElement>;

	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			sync: 0,
			period: 1440,
			interval: 30,
			loading: {},
			cancelable: {},
		};

		this.loading = {};
		this.chartBoxRef = React.createRef();
	}

	getDefaultInterval(period: number): number {
		switch (period) {
			case 60:
				return 1;
			case 180:
				return 5;
			case 360:
				return 5;
			case 720:
				return 30;
			case 1440:
				return 30;
			case 4320:
				return 60;
			case 10080:
				return 120;
			case 20160:
				return 360;
			case 43200:
				return 720;
			case 86400:
				return 1440;
			case 129600:
				return 1440;
			case 172800:
				return 4320;
			default:
				return 360;
		}
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

		let intervalMin = 0;
		let intervalMax = 0;
		if (this.state.period > 43200) {
			intervalMin = 120;
		} else if (this.state.period > 20160) {
			intervalMin = 30;
		} else if (this.state.period > 4320) {
			intervalMin = 5;
		}

		if (this.state.period <= 60) {
			intervalMax = 30;
		} else if (this.state.period <= 180) {
			intervalMax = 60;
		} else if (this.state.period <= 360) {
			intervalMax = 120;
		} else if (this.state.period <= 720) {
			intervalMax = 360;
		} else if (this.state.period <= 1440) {
			intervalMax = 720;
		} else if (this.state.period <= 4320) {
			intervalMax = 1440;
		} else if (this.state.period <= 10080) {
			intervalMax = 4320;
		} else {
			intervalMax = 10080;
		}

		let refreshDisabled = false;
		let refreshLabel = '';
		let refreshClass = 'bp5-button';
		if (Object.entries(this.state.cancelable).length) {
			refreshLabel = 'Cancel';
			refreshClass += ' bp5-intent-warning bp5-icon-delete'
		} else {
			if (Object.entries(this.state.loading).length) {
				refreshDisabled = true;
			}
			refreshLabel = 'Refresh';
			refreshClass += ' bp5-intent-success bp5-icon-refresh'
		}

		return <div ref={this.chartBoxRef}>
			<div className="layout horizontal wrap bp5-border" style={css.header}>
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
								EndpointActions.dataCancel();
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
							let period = parseInt(val, 10);
							this.setState({
								...this.state,
								period: period,
								interval: this.getDefaultInterval(period),
							});
						}}
					>
						<option value="60">1 hour</option>
						<option value="180">3 hours</option>
						<option value="360">6 hours</option>
						<option value="720">12 hours</option>
						<option value="1440">24 hours</option>
						<option value="4320">3 days</option>
						<option value="10080">7 days</option>
						<option value="20160">14 days</option>
						<option value="43200">30 days</option>
						<option value="86400">60 days</option>
						<option value="129600" hidden={true}>90 days</option>
						<option value="172800" hidden={true}>120 days</option>
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
						<option
							value="1"
							hidden={1 < intervalMin || 1 > intervalMax}
						>1 minute</option>
						<option
							value="5"
							hidden={5 < intervalMin || 5 > intervalMax}
						>5 minutes</option>
						<option
							value="30"
							hidden={30 < intervalMin || 30 > intervalMax}
						>30 minutes</option>
						<option
							value="60"
							hidden={60 < intervalMin || 60 > intervalMax}
						>1 hour</option>
						<option
							value="120"
							hidden={120 < intervalMin || 120 > intervalMax}
						>2 hours</option>
						<option
							value="360"
							hidden={360 < intervalMin || 360 > intervalMax}
						>6 hours</option>
						<option
							value="720"
							hidden={720 < intervalMin || 720 > intervalMax}
						>12 hours</option>
						<option
							value="1440"
							hidden={1440 < intervalMin || 1440 > intervalMax}
						>24 hours</option>
						<option
							value="4320"
							hidden={4320 < intervalMin || 4320 > intervalMax}
						>3 days</option>
						<option
							value="10080"
							hidden={10080 < intervalMin || 10080 > intervalMax}
						>7 days</option>
					</PageSelect>
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
						left={true}
						onLoading={(): void => {
							this.setLoading('system');
						}}
						onLoaded={(): void => {
							this.setLoaded('system');
						}}
						getBoxRect={(): DOMRect => {
							return this.chartBoxRef.current.getBoundingClientRect();
						}}
					/>
				</div>
				<div style={css.chartGroup}>
					<EndpointChart
						endpoint={this.props.endpoint}
						resource={'load'}
						sync={this.state.sync}
						period={this.state.period}
						interval={this.state.interval}
						left={false}
						onLoading={(): void => {
							this.setLoading('load');
						}}
						onLoaded={(): void => {
							this.setLoaded('load');
						}}
						getBoxRect={(): DOMRect => {
							return this.chartBoxRef.current.getBoundingClientRect();
						}}
					/>
				</div>
			</div>
			<div className="layout horizontal wrap">
				<div style={css.chartGroup}>
					<EndpointChart
						endpoint={this.props.endpoint}
						resource={'disk'}
						sync={this.state.sync}
						period={this.state.period}
						interval={this.state.interval}
						left={true}
						onLoading={(): void => {
							this.setLoading('disk');
						}}
						onLoaded={(): void => {
							this.setLoaded('disk');
						}}
						getBoxRect={(): DOMRect => {
							return this.chartBoxRef.current.getBoundingClientRect();
						}}
					/>
				</div>
				<div style={css.chartGroup}>
					<EndpointChart
						endpoint={this.props.endpoint}
						resource={'network'}
						sync={this.state.sync}
						period={this.state.period}
						interval={this.state.interval}
						left={false}
						onLoading={(): void => {
							this.setLoading('network');
						}}
						onLoaded={(): void => {
							this.setLoaded('network');
						}}
						getBoxRect={(): DOMRect => {
							return this.chartBoxRef.current.getBoundingClientRect();
						}}
					/>
				</div>
			</div>
			<div className="layout horizontal wrap">
				<div style={css.chartGroup}>
					<EndpointChart
						endpoint={this.props.endpoint}
						resource={'diskio0'}
						sync={this.state.sync}
						period={this.state.period}
						interval={this.state.interval}
						left={true}
						onLoading={(): void => {
							this.setLoading('diskio0');
						}}
						onLoaded={(): void => {
							this.setLoaded('diskio0');
						}}
						getBoxRect={(): DOMRect => {
							return this.chartBoxRef.current.getBoundingClientRect();
						}}
					/>
				</div>
				<div style={css.chartGroup}>
					<EndpointChart
						endpoint={this.props.endpoint}
						resource={'diskio1'}
						sync={this.state.sync}
						period={this.state.period}
						interval={this.state.interval}
						left={false}
						onLoading={(): void => {
							this.setLoading('diskio1');
						}}
						onLoaded={(): void => {
							this.setLoaded('diskio1');
						}}
						getBoxRect={(): DOMRect => {
							return this.chartBoxRef.current.getBoundingClientRect();
						}}
					/>
				</div>
			</div>
		</div>;
	}
}
