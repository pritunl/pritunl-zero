/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as ServiceTypes from '../types/ServiceTypes';
import * as ServiceActions from '../actions/ServiceActions';
import PageInput from './PageInput';
import PageSave from './PageSave';
import ConfirmButton from './ConfirmButton';
import PageInputButton from './PageInputButton';

interface Props {
	service: ServiceTypes.ServiceRo;
}

interface State {
	disabled: boolean;
	changed: boolean;
	message: string;
	addRole: string;
	service: ServiceTypes.Service;
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
	role: {
		margin: '9px 5px 0 5px',
		height: '20px',
	} as React.CSSProperties,
	group: {
		flex: 1,
		minWidth: '250px',
	} as React.CSSProperties,
	save: {
		paddingBottom: '10px',
	} as React.CSSProperties,
};

export default class Service extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			disabled: false,
			changed: false,
			message: '',
			addRole: '',
			service: null,
		};
	}

	set = (name: string, val: any): void => {
		let service: any;

		if (this.state.changed) {
			service = {
				...this.state.service,
			};
		} else {
			service = {
				...this.props.service,
			};
		}

		service[name] = val;

		this.setState({
			...this.state,
			changed: true,
			service: service,
		});
	}

	onSave = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});
		ServiceActions.commit(this.state.service).then((): void => {
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
						service: null,
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
		ServiceActions.remove(this.props.service.id).then((): void => {
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
		let service: ServiceTypes.Service = this.state.service ||
			this.props.service;

		let roles = [
			...service.roles,
		];

		if (!this.state.addRole) {
			return;
		}

		if (roles.indexOf(this.state.addRole) === -1) {
			roles.push(this.state.addRole);
		}

		roles.sort();

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addRole: '',
			service: {
				...service,
				roles: roles,
			},
		});
	}

	onRemoveRole = (role: string): void => {
		let service: ServiceTypes.Service = this.state.service ||
			this.props.service;

		let roles = [
			...service.roles,
		];

		let i = roles.indexOf(role);
		if (i === -1) {
			return;
		}

		roles.splice(i, 1);

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addRole: '',
			service: {
				...service,
				roles: roles,
			},
		});
	}

	render(): JSX.Element {
		let service: ServiceTypes.Service = this.state.service ||
			this.props.service;

		let roles: JSX.Element[] = [];
		for (let role of service.roles) {
			roles.push(
				<div
					className="pt-tag pt-tag-removable pt-intent-primary"
					style={css.role}
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
							disabled={this.state.disabled}
							onConfirm={this.onDelete}
						/>
					</div>
					<PageInput
						label="Name"
						type="text"
						placeholder="Enter name"
						value={service.name}
						onChange={(val): void => {
							this.set('name', val);
						}}
					/>
					<label className="pt-label">
						Roles
						<div>
							{roles}
						</div>
					</label>
					<PageInputButton
						buttonClass="pt-intent-success"
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
			<PageSave
				style={css.save}
				hidden={!this.state.service}
				message={this.state.message}
				changed={this.state.changed}
				disabled={false}
				light={true}
				onCancel={(): void => {
					this.setState({
						...this.state,
						changed: false,
						service: null,
					});
				}}
				onSave={this.onSave}
			/>
		</div>;
	}
}
