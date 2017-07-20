/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as NodeTypes from '../types/NodeTypes';
import * as ServiceTypes from '../types/ServiceTypes';
import * as NodeActions from '../actions/NodeActions';
import * as MiscUtils from '../utils/MiscUtils';
import ServicesStore from '../stores/ServicesStore';
import PageInput from './PageInput';
import PageSelect from './PageSelect';
import PageSelectButton from './PageSelectButton';
import PageInfo from './PageInfo';
import PageSave from './PageSave';
import ConfirmButton from './ConfirmButton';

interface Props {
	node: NodeTypes.NodeRo;
	services: ServiceTypes.ServicesRo;
}

interface State {
	disabled: boolean;
	changed: boolean;
	message: string;
	node: NodeTypes.Node;
	addService: string;
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

export default class Node extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			disabled: false,
			changed: false,
			message: '',
			node: null,
			addService: null,
		};
	}

	set(name: string, val: any): void {
		let node: any;

		if (this.state.changed) {
			node = {
				...this.state.node,
			};
		} else {
			node = {
				...this.props.node,
			};
		}

		node[name] = val;

		this.setState({
			...this.state,
			changed: true,
			node: node,
		});
	}

	onSave = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});
		NodeActions.commit(this.state.node).then((): void => {
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
						node: null,
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
		NodeActions.remove(this.props.node.id).then((): void => {
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
		let node: NodeTypes.Node;

		if (!this.state.addService && !this.props.services.length) {
			return
		}

		let serviceId = this.state.addService || this.props.services[0].id;

		console.log('***************************************************');
		console.log(serviceId);
		console.log('***************************************************');

		if (this.state.changed) {
			node = {
				...this.state.node,
			};
		} else {
			node = {
				...this.props.node,
			};
		}

		let services = [
			...node.services,
		];

		if (services.indexOf(serviceId) === -1) {
			services.push(serviceId);
		}

		services.sort();

		node.services = services;

		this.setState({
			...this.state,
			changed: true,
			addService: null,
			node: node,
		});
	}

	onRemoveService = (service: string): void => {
		let node: NodeTypes.Node;

		if (this.state.changed) {
			node = {
				...this.state.node,
			};
		} else {
			node = {
				...this.props.node,
			};
		}

		let services = [
			...node.services,
		];

		let i = services.indexOf(service);
		if (i === -1) {
			return;
		}

		services.splice(i, 1);

		node.services = services;

		this.setState({
			...this.state,
			changed: true,
			addService: null,
			node: node,
		});
	}

	render(): JSX.Element {
		let node: NodeTypes.Node = this.state.node ||
			this.props.node;

		let services: JSX.Element[] = [];
		for (let service of this.props.services) {
			services.push(
				<option key={service.id} value={service.id}>{service.name}</option>
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
							confirmMsg="Confirm node remove"
							disabled={this.state.disabled}
							onConfirm={this.onDelete}
						/>
					</div>
					<PageInput
						label="Name"
						type="text"
						placeholder="Enter name"
						value={node.name}
						onChange={(val): void => {
							this.set('name', val);
						}}
					/>
					<PageSelect
						label="Type"
						value={node.type}
						onChange={(val): void => {
							this.set('type', val);
						}}
					>
						<option value="management">Management</option>
						<option value="proxy">Proxy</option>
						<option value="management_proxy">Management + Proxy</option>
					</PageSelect>
					<label className="pt-label" style={css.label}>
						Protocol and Port
						<div className="pt-control-group" style={css.inputGroup}>
							<div className="pt-select" style={css.protocol}>
								<select
									value={node.protocol || 'https'}
									onChange={(evt): void => {
										this.set('protocol', evt.target.value);
									}}
								>
									<option value="http">HTTP</option>
									<option value="https">HTTPS</option>
								</select>
							</div>
							<input
								className="pt-input"
								style={css.port}
								type="text"
								autoCapitalize="off"
								spellCheck={false}
								placeholder="Port"
								value={node.port || 443}
								onChange={(evt): void => {
									this.set('port', parseInt(evt.target.value, 10));
								}}
							/>
						</div>
					</label>
					<PageSelectButton
						label="Add Service"
						value={this.state.addService}
						buttonClass="pt-intent-success"
						onChange={(val: string): void => {
							this.setState({
								...this.state,
								serviceId: val,
							});
						}}
						onSubmit={this.onAddService}
					>
						{services}
					</PageSelectButton>
				</div>
				<div style={css.group}>
					<PageInfo
						fields={[
							{
								label: 'ID',
								value: node.id || 'None',
							},
							{
								label: 'Timestamp',
								value: MiscUtils.formatDate(node.timestamp) || 'Inactive',
							},
						]}
						bars={[
							{
								progressClass: 'pt-no-stripes pt-intent-primary',
								label: 'Memory',
								value: node.memory,
							},
							{
								progressClass: 'pt-no-stripes pt-intent-success',
								label: 'Load1',
								value: node.load1,
							},
							{
								progressClass: 'pt-no-stripes pt-intent-warning',
								label: 'Load5',
								value: node.load5,
							},
							{
								progressClass: 'pt-no-stripes pt-intent-danger',
								label: 'Load15',
								value: node.load15,
							},
						]}
					/>
				</div>
			</div>
			<PageSave
				style={css.save}
				hidden={!this.state.node}
				message={this.state.message}
				changed={this.state.changed}
				disabled={false}
				light={true}
				onCancel={(): void => {
					this.setState({
						...this.state,
						changed: false,
						node: null,
					});
				}}
				onSave={this.onSave}
			/>
		</div>;
	}
}
