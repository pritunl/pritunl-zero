/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as PolicyTypes from '../types/PolicyTypes';
import * as ServiceTypes from '../types/ServiceTypes';
import * as PolicyActions from '../actions/PolicyActions';
import * as MiscUtils from '../utils/MiscUtils';
import ServicesStore from '../stores/ServicesStore';
import PageInput from './PageInput';
import PageSelect from './PageSelect';
import PageSelectButton from './PageSelectButton';
import PageInputButton from './PageInputButton';
import PageInfo from './PageInfo';
import PageSave from './PageSave';
import ConfirmButton from './ConfirmButton';

interface Props {
	policy: PolicyTypes.PolicyRo;
	services: ServiceTypes.ServicesRo;
}

interface State {
	disabled: boolean;
	changed: boolean;
	message: string;
	policy: PolicyTypes.Policy;
	addService: string;
	addRole: string;
}

const css = {
	card: {
		position: 'relative',
		padding: '10px 10px 0 10px',
		marginBottom: '5px',
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
	} as React.CSSProperties,
	save: {
		paddingBottom: '10px',
	} as React.CSSProperties,
	label: {
		width: '100%',
		maxWidth: '280px',
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
};

export default class Policy extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			disabled: false,
			changed: false,
			message: '',
			policy: null,
			addService: null,
			addRole: null,
		};
	}

	set(name: string, val: any): void {
		let policy: any;

		if (this.state.changed) {
			policy = {
				...this.state.policy,
			};
		} else {
			policy = {
				...this.props.policy,
			};
		}

		policy[name] = val;

		this.setState({
			...this.state,
			changed: true,
			policy: policy,
		});
	}

	onSave = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});
		PolicyActions.commit(this.state.policy).then((): void => {
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
						message: '',
						changed: false,
						policy: null,
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
		PolicyActions.remove(this.props.policy.id).then((): void => {
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

	onAddService = (): void => {
		let policy: PolicyTypes.Policy;

		if (!this.state.addService && !this.props.services.length) {
			return;
		}

		let serviceId = this.state.addService || this.props.services[0].id;

		if (this.state.changed) {
			policy = {
				...this.state.policy,
			};
		} else {
			policy = {
				...this.props.policy,
			};
		}

		let services = [
			...policy.services,
		];

		if (services.indexOf(serviceId) === -1) {
			services.push(serviceId);
		}

		services.sort();

		policy.services = services;

		this.setState({
			...this.state,
			changed: true,
			policy: policy,
		});
	}

	onRemoveService = (service: string): void => {
		let policy: PolicyTypes.Policy;

		if (this.state.changed) {
			policy = {
				...this.state.policy,
			};
		} else {
			policy = {
				...this.props.policy,
			};
		}

		let services = [
			...policy.services,
		];

		let i = services.indexOf(service);
		if (i === -1) {
			return;
		}

		services.splice(i, 1);

		policy.services = services;

		this.setState({
			...this.state,
			changed: true,
			policy: policy,
		});
	}

	onAddRole = (): void => {
		let policy: PolicyTypes.Policy;

		if (this.state.changed) {
			policy = {
				...this.state.policy,
			};
		} else {
			policy = {
				...this.props.policy,
			};
		}

		let roles = [
			...policy.roles,
		];

		if (!this.state.addRole) {
			return;
		}

		if (roles.indexOf(this.state.addRole) === -1) {
			roles.push(this.state.addRole);
		}

		roles.sort();

		policy.roles = roles;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addRole: '',
			policy: policy,
		});
	}

	onRemoveRole(role: string): void {
		let policy: PolicyTypes.Policy;

		if (this.state.changed) {
			policy = {
				...this.state.policy,
			};
		} else {
			policy = {
				...this.props.policy,
			};
		}

		let roles = [
			...policy.roles,
		];

		let i = roles.indexOf(role);
		if (i === -1) {
			return;
		}

		roles.splice(i, 1);

		policy.roles = roles;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addRole: '',
			policy: policy,
		});
	}

	render(): JSX.Element {
		let policy: PolicyTypes.Policy = this.state.policy ||
			this.props.policy;

		let services: JSX.Element[] = [];
		for (let serviceId of policy.services) {
			let service = ServicesStore.service(serviceId);
			if (!service) {
				continue;
			}

			services.push(
				<div
					className="pt-tag pt-tag-removable pt-intent-primary"
					style={css.item}
					key={service.id}
				>
					{service.name}
					<button
						className="pt-tag-remove"
						onMouseUp={(): void => {
							this.onRemoveService(service.id);
						}}
					/>
				</div>,
			);
		}

		let servicesSelect: JSX.Element[] = [];
		if (this.props.services.length) {
			for (let service of this.props.services) {
				servicesSelect.push(
					<option key={service.id} value={service.id}>{service.name}</option>,
				);
			}
		} else {
			servicesSelect.push(<option key="null" value="">None</option>);
		}

		let roles: JSX.Element[] = [];
		for (let role of policy.roles) {
			roles.push(
				<div
					className="pt-tag pt-tag-removable pt-intent-primary"
					style={css.item}
					key={role}
				>
					{role}
					<button
						className="pt-tag-remove"
						onMouseUp={(): void => {
							this.onRemoveRole(role);
						}}
					/>
				</div>,
			);
		}

		return <div
			className="pt-card"
			style={css.card}
		>
			<div className="layout horizontal wrap">
				<div style={css.group}>
					<div style={css.remove}>
						<ConfirmButton
							className="pt-minimal pt-intent-danger pt-icon-cross"
							progressClassName="pt-intent-danger"
							confirmMsg="Confirm policy remove"
							disabled={this.state.disabled}
							onConfirm={this.onDelete}
						/>
					</div>
					<PageInput
						label="Name"
						type="text"
						placeholder="Enter name"
						value={policy.name}
						onChange={(val): void => {
							this.set('name', val);
						}}
					/>
					<label
						className="pt-label"
						style={css.label}
					>
						Services
						<div>
							{services}
						</div>
					</label>
					<PageSelectButton
						label="Add Service"
						value={this.state.addService}
						disabled={!this.props.services.length}
						buttonClass="pt-intent-success"
						onChange={(val: string): void => {
							this.setState({
								...this.state,
								addService: val,
							});
						}}
						onSubmit={this.onAddService}
					>
						{servicesSelect}
					</PageSelectButton>
					<label className="pt-label">
						Roles
						<div>
							{roles}
						</div>
					</label>
					<PageInputButton
						buttonClass="pt-intent-success pt-icon-add"
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
					<PageInfo
						fields={[
							{
								label: 'ID',
								value: policy.id || 'None',
							},
						]}
					/>
				</div>
			</div>
			<PageSave
				style={css.save}
				hidden={!this.state.policy}
				message={this.state.message}
				changed={this.state.changed}
				disabled={this.state.disabled}
				light={true}
				onCancel={(): void => {
					this.setState({
						...this.state,
						changed: false,
						policy: null,
					});
				}}
				onSave={this.onSave}
			/>
		</div>;
	}
}
