/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as CheckTypes from '../types/CheckTypes';
import * as AuthorityTypes from "../types/AuthorityTypes";
import * as CheckActions from '../actions/CheckActions';
import * as PageInfos from './PageInfo';
import PageInput from './PageInput';
import PageSave from './PageSave';
import PageInfo from './PageInfo';
import ConfirmButton from './ConfirmButton';
import PageInputButton from './PageInputButton';
import CheckCharts from './CheckCharts';
import Help from './Help';
import PageSwitch from "./PageSwitch";
import PageSelect from "./PageSelect";
import CheckHeader from "./CheckHeader";
import EndpointKmsg from "./EndpointKmsg";

interface Props {
	check: CheckTypes.CheckRo;
	authorities: AuthorityTypes.AuthoritiesRo;
	selected: boolean;
	onSelect: (shift: boolean) => void;
	onClose: () => void;
}

interface State {
	disabled: boolean;
	changed: boolean;
	message: string;
	addRole: string;
	addTarget: string;
	check: CheckTypes.Check;
}

const css = {
	card: {
		position: 'relative',
		padding: '48px 10px 0 10px',
		width: '100%',
	} as React.CSSProperties,
	remove: {
		position: 'absolute',
		top: '5px',
		right: '5px',
	} as React.CSSProperties,
	item: {
		margin: '9px 5px 0 5px',
		height: '20px',
	} as React.CSSProperties,
	itemsLabel: {
		display: 'block',
	} as React.CSSProperties,
	itemsAdd: {
		margin: '8px 0 15px 0',
	} as React.CSSProperties,
	group: {
		flex: 1,
		minWidth: '250px',
		margin: '0 10px',
	} as React.CSSProperties,
	controlButton: {
		marginRight: '10px',
	} as React.CSSProperties,
	save: {
		paddingBottom: '10px',
	} as React.CSSProperties,
	button: {
		height: '30px',
	} as React.CSSProperties,
	buttons: {
		cursor: 'pointer',
		position: 'absolute',
		top: 0,
		left: 0,
		right: 0,
		padding: '4px',
		height: '39px',
	} as React.CSSProperties,
	label: {
		width: '100%',
		maxWidth: '280px',
	} as React.CSSProperties,
	status: {
		margin: '6px 0 0 1px',
	} as React.CSSProperties,
	icon: {
		marginRight: '3px',
	} as React.CSSProperties,
	inputGroup: {
		width: '100%',
	} as React.CSSProperties,
	protocol: {
		flex: '0 1 auto',
	} as React.CSSProperties,
	port: {
		flex: '1',
	} as React.CSSProperties,
	select: {
		margin: '7px 0px 0px 6px',
		paddingTop: '3px',
	} as React.CSSProperties,
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
	checksButtons: {
		marginTop: '8px',
	} as React.CSSProperties,
	checksAdd: {
		margin: '8px 0 0 8px',
	} as React.CSSProperties,
};

export default class CheckDetailed extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			disabled: false,
			changed: false,
			message: '',
			addRole: '',
			addTarget: '',
			check: null,
		};
	}

	set(name: string, val: any): void {
		let check: any;

		if (this.state.changed) {
			check = {
				...this.state.check,
			};
		} else {
			check = {
				...this.props.check,
			};
		}

		check[name] = val;

		this.setState({
			...this.state,
			changed: true,
			check: check,
		});
	}

	onSave = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});
		CheckActions.commit(this.state.check).then((): void => {
			this.setState({
				...this.state,
				message: 'Your changes have been saved',
				changed: false,
				disabled: false,
			});

			setTimeout((): void => {
				if (!this.state.changed) {
					this.setState({
						...this.state,
						check: null,
						changed: false,
					});
				}
			}, 1000);

			setTimeout((): void => {
				if (!this.state.changed) {
					this.setState({
						...this.state,
						message: '',
					});
				}
			}, 3000);
		}).catch((): void => {
			this.setState({
				...this.state,
				message: '',
				disabled: false,
			});
		});
	}

	onDelete = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});
		CheckActions.remove(this.props.check.id).then((): void => {
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

	onAddRole = (): void => {
		let check: CheckTypes.Check;

		if (this.state.changed) {
			check = {
				...this.state.check,
			};
		} else {
			check = {
				...this.props.check,
			};
		}

		let roles = [
			...check.roles,
		];

		if (!this.state.addRole) {
			return;
		}

		if (roles.indexOf(this.state.addRole) === -1) {
			roles.push(this.state.addRole);
		}

		roles.sort();

		check.roles = roles;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addRole: '',
			check: check,
		});
	}

	onRemoveRole(role: string): void {
		let check: CheckTypes.Check;

		if (this.state.changed) {
			check = {
				...this.state.check,
			};
		} else {
			check = {
				...this.props.check,
			};
		}

		let roles = [
			...check.roles,
		];

		let i = roles.indexOf(role);
		if (i === -1) {
			return;
		}

		roles.splice(i, 1);

		check.roles = roles;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addRole: '',
			check: check,
		});
	}

	onAddTarget = (): void => {
		let check: CheckTypes.Check;

		if (this.state.changed) {
			check = {
				...this.state.check,
			};
		} else {
			check = {
				...this.props.check,
			};
		}

		let targets = [
			...(check.targets || []),
		];

		if (!this.state.addTarget) {
			return;
		}

		if (targets.indexOf(this.state.addTarget) === -1) {
			targets.push(this.state.addTarget);
		}

		targets.sort();

		check.targets = targets;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addTarget: '',
			check: check,
		});
	}

	onRemoveTarget(target: string): void {
		let check: CheckTypes.Check;

		if (this.state.changed) {
			check = {
				...this.state.check,
			};
		} else {
			check = {
				...this.props.check,
			};
		}

		let targets = [
			...(check.targets || []),
		];

		let i = targets.indexOf(target);
		if (i === -1) {
			return;
		}

		targets.splice(i, 1);

		check.targets = targets;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addTarget: '',
			check: check,
		});
	}

	onAddHeader = (): void => {
		let check: CheckTypes.Check;

		if (this.state.changed) {
			check = {
				...this.state.check,
			};
		} else {
			check = {
				...this.props.check,
			};
		}

		let headers = [
			...check.headers,
			{},
		];

		check.headers = headers;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			check: check,
		});
	}

	onChangeHeader(i: number, state: CheckTypes.Header): void {
		let check: CheckTypes.Check;

		if (this.state.changed) {
			check = {
				...this.state.check,
			};
		} else {
			check = {
				...this.props.check,
			};
		}

		let headers = [
			...check.headers,
		];

		headers[i] = state;

		check.headers = headers;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			check: check,
		});
	}

	onRemoveHeader(i: number): void {
		let check: CheckTypes.Check;

		if (this.state.changed) {
			check = {
				...this.state.check,
			};
		} else {
			check = {
				...this.props.check,
			};
		}

		let headers = [
			...check.headers,
		];

		headers.splice(i, 1);

		check.headers = headers;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			check: check,
		});
	}

	render(): JSX.Element {
		let check: CheckTypes.Check = this.state.check ||
			this.props.check;

		let fields: PageInfos.Field[] = [
			{
				label: 'ID',
				value: this.props.check.id || 'None',
			},
		];

		let roles: JSX.Element[] = [];
		for (let role of check.roles) {
			roles.push(
				<div
					className="bp5-tag bp5-tag-removable bp5-intent-primary"
					style={css.item}
					key={role}
				>
					{role}
					<button
						className="bp5-tag-remove"
						onMouseUp={(): void => {
							this.onRemoveRole(role);
						}}
					/>
				</div>,
			);
		}

		let targets: JSX.Element[] = [];
		for (let target of (check.targets || [])) {
			targets.push(
				<div
					className="bp5-tag bp5-tag-removable bp5-intent-primary"
					style={css.item}
					key={target}
				>
					{target}
					<button
						className="bp5-tag-remove"
						onMouseUp={(): void => {
							this.onRemoveTarget(target);
						}}
					/>
				</div>,
			);
		}

		let targetLabel = '';
		let targetTitle = '';
		let targetHelp = '';

		targetLabel = 'Targets';
		targetTitle = 'Targets';
		targetHelp = 'Targets for health check. For most configurations each ' +
			'target should be placed in a separate check.';

		let headers: JSX.Element[] = [];
		if (check.type === "http") {
			for (let i = 0; i < check.headers.length; i++) {
				let index = i;

				headers.push(
					<CheckHeader
						key={"check-header-" + index}
						header={check.headers[index]}
						onChange={(state: CheckTypes.Header): void => {
							this.onChangeHeader(index, state);
						}}
						onRemove={(): void => {
							this.onRemoveHeader(index);
						}}
					/>,
				);
			}
		}

		return <td
			className="bp5-cell"
			colSpan={2}
			style={css.card}
		>
			<div className="layout horizontal wrap">
				<div style={css.group}>
					<div
						className="layout horizontal tab-close"
						style={css.buttons}
						onClick={(evt): void => {
							let target = evt.target as HTMLElement;

							if (target.className.indexOf('tab-close') !== -1) {
								this.props.onClose();
							}
						}}
					>
						<div>
							<label
								className="bp5-control bp5-checkbox"
								style={css.select}
							>
								<input
									type="checkbox"
									checked={this.props.selected}
									onChange={(evt): void => {
									}}
									onClick={(evt): void => {
										this.props.onSelect(evt.shiftKey);
									}}
								/>
								<span className="bp5-control-indicator"/>
							</label>
						</div>
						<div className="flex tab-close"/>
						<ConfirmButton
							safe={true}
							className="bp5-minimal bp5-intent-danger bp5-icon-trash"
							progressClassName="bp5-intent-danger"
							dialogClassName="bp5-intent-danger bp5-icon-delete"
							dialogLabel="Delete Check"
							confirmMsg="Permanently delete this check"
							confirmInput={true}
							items={[check.name]}
							disabled={this.state.disabled}
							onConfirm={this.onDelete}
						/>
					</div>
					<PageInput
						disabled={this.state.disabled}
						label="Name"
						help="Name of check"
						type="text"
						placeholder="Enter name"
						value={check.name}
						onChange={(val): void => {
							this.set('name', val);
						}}
					/>
					<label className="bp5-label">
						Roles
						<Help
							title="Roles"
							content="The roles used to match to endpoints. Endpoints that have a matching role will perform checks."
						/>
						<div>
							{roles}
						</div>
					</label>
					<PageInputButton
						disabled={this.state.disabled}
						buttonClass="bp5-intent-success bp5-icon-add"
						label="Add"
						type="text"
						placeholder="Add role"
						value={this.state.addRole}
						onChange={(val): void => {
							this.setState({
								...this.state,
								addRole: val,
							});
						}}
						onSubmit={this.onAddRole}
					/>
					<PageSelect
						disabled={this.state.disabled}
						label="Check Type"
						help="Type of check"
						value={check.type}
						onChange={(val): void => {
							this.set('type', val);
						}}
					>
						<option
							value="http"
						>HTTP Request</option>
					</PageSelect>
					<label className="bp5-label">
						{targetLabel}
						<Help
							title={targetTitle}
							content={targetHelp}
						/>
						<div>
							{targets}
						</div>
					</label>
					<PageInputButton
						disabled={this.state.disabled}
						buttonClass="bp5-intent-success bp5-icon-add"
						label="Add"
						type="text"
						placeholder="Add target"
						value={this.state.addTarget}
						onChange={(val): void => {
							this.setState({
								...this.state,
								addTarget: val,
							});
						}}
						onSubmit={this.onAddTarget}
					/>
				</div>
				<div style={css.group}>
					<PageInfo
						fields={fields}
					/>
					<PageInput
						disabled={this.state.disabled}
						label="Check Frequency"
						help="Minimum duration in seconds between repeat checks."
						type="text"
						placeholder="Enter frequency"
						value={check.frequency}
						onChange={(val): void => {
							this.set('frequency', parseInt(val, 10));
						}}
					/>
					<PageInput
						disabled={this.state.disabled}
						label="Check Timeout"
						help="Time in seconds before check times out."
						type="text"
						placeholder="Enter timeout"
						value={check.timeout}
						onChange={(val): void => {
							this.set('timeout', parseInt(val, 10));
						}}
					/>
					<PageInput
						disabled={this.state.disabled}
						label="HTTP Status Code"
						help="Expected status code to receive."
						type="text"
						placeholder="Enter status code"
						hidden={check.type !== "http"}
						value={check.status_code}
						onChange={(val): void => {
							this.set('status_code', parseInt(val, 10));
						}}
					/>
					<label style={css.itemsLabel} hidden={check.type !== "http"}>
						HTTP Headers
						<Help
							title="HTTP Headers"
							content="Headers to include when sending HTTP health check request."
						/>
					</label>
					{headers}
					<button
						className="bp5-button bp5-intent-success bp5-icon-add"
						style={css.itemsAdd}
						hidden={check.type !== "http"}
						type="button"
						onClick={this.onAddHeader}
					>
						Add Header
					</button>
				</div>
			</div>
			<PageSave
				style={css.save}
				hidden={!this.state.check && !this.state.message}
				message={this.state.message}
				changed={this.state.changed}
				disabled={this.state.disabled}
				light={true}
				onCancel={(): void => {
					this.setState({
						...this.state,
						changed: false,
						check: null,
					});
				}}
				onSave={this.onSave}
			/>
			<CheckCharts
				check={check.id}
				disabled={false}
			/>
			<EndpointKmsg
				check={check.id}
				disabled={false}
			/>
		</td>;
	}
}
