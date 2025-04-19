/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as EndpointTypes from '../types/EndpointTypes';
import * as AuthorityTypes from "../types/AuthorityTypes";
import * as EndpointActions from '../actions/EndpointActions';
import * as PageInfos from './PageInfo';
import * as MiscUtils from '../utils/MiscUtils';
import PageInput from './PageInput';
import PageCreate from './PageCreate';
import PageInfo from './PageInfo';
import ConfirmButton from './ConfirmButton';
import PageInputButton from './PageInputButton';
import EndpointCharts from './EndpointCharts';
import EndpointKmsg from './EndpointKmsg';
import Help from './Help';
import PageSwitch from "./PageSwitch";

interface Props {
	authorities: AuthorityTypes.AuthoritiesRo;
	onClose: () => void;
}

interface State {
	closed: boolean;
	disabled: boolean;
	changed: boolean;
	message: string;
	addRole: string;
	showCharts: boolean;
	endpoint: EndpointTypes.Endpoint;
}

const css = {
	row: {
		display: 'table-row',
		width: '100%',
		padding: 0,
		boxShadow: 'none',
		position: 'relative',
	} as React.CSSProperties,
	card: {
		position: 'relative',
		padding: '10px 10px 0 10px',
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
		position: 'absolute',
		top: '5px',
		right: '5px',
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
};

export default class EndpointNew extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			closed: false,
			disabled: false,
			changed: false,
			message: '',
			addRole: '',
			showCharts: true,
			endpoint: {},
		};
	}

	set(name: string, val: any): void {
		let endpoint: any = {
			...this.state.endpoint,
		};

		endpoint[name] = val;

		this.setState({
			...this.state,
			changed: true,
			endpoint: endpoint,
		});
	}

	onCreate = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});

		let endpoint: any = {
			...this.state.endpoint,
		};

		EndpointActions.create(endpoint).then((): void => {
			this.setState({
				...this.state,
				message: 'Endpoint created successfully',
				changed: false,
			});

			setTimeout((): void => {
				this.setState({
					...this.state,
					disabled: false,
					changed: true,
				});
			}, 2000);
		}).catch((): void => {
			this.setState({
				...this.state,
				message: '',
				disabled: false,
			});
		});
	}

	onAddRole = (): void => {
		let endpoint: EndpointTypes.Endpoint;

		endpoint = {
			...this.state.endpoint,
		};

		let roles = [
			...(endpoint.roles || []),
		];

		if (!this.state.addRole) {
			return;
		}

		if (roles.indexOf(this.state.addRole) === -1) {
			roles.push(this.state.addRole);
		}

		roles.sort();

		endpoint.roles = roles;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addRole: '',
			endpoint: endpoint,
		});
	}

	onRemoveRole(role: string): void {
		let endpoint: EndpointTypes.Endpoint;

		endpoint = {
			...this.state.endpoint,
		};

		let roles = [
			...(endpoint.roles || []),
		];

		let i = roles.indexOf(role);
		if (i === -1) {
			return;
		}

		roles.splice(i, 1);

		endpoint.roles = roles;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addRole: '',
			endpoint: endpoint,
		});
	}

	render(): JSX.Element {
		let endpoint: EndpointTypes.Endpoint = this.state.endpoint;

		let roles: JSX.Element[] = [];
		for (let role of (endpoint.roles || [])) {
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

		return <div
			className="bp5-card bp5-row"
			style={css.row}
		>
			<td
				className="bp5-cell"
				colSpan={3}
				style={css.card}
			>
				<div className="layout horizontal wrap">
					<div style={css.group}>
						<div style={css.buttons}>
						</div>
						<PageInput
							label="Name"
							help="Name of endpoint"
							type="text"
							placeholder="Enter name"
							value={endpoint.name}
							onChange={(val): void => {
								this.set('name', val);
							}}
						/>
						<label className="bp5-label">
							Roles
							<Help
								title="Roles"
								content="The user roles that will be allowed access to this endpoint. At least one role must match for the user to access the endpoint."
							/>
							<div>
								{roles}
							</div>
						</label>
						<PageInputButton
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
					</div>
					<div style={css.group}>
					</div>
				</div>
				<PageCreate
					style={css.save}
					hidden={!this.state.endpoint}
					message={this.state.message}
					changed={this.state.changed}
					disabled={this.state.disabled}
					closed={this.state.closed}
					light={true}
					onCancel={this.props.onClose}
					onCreate={this.onCreate}
				/>
			</td>
		</div>;
	}
}
