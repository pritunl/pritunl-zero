/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as NodeTypes from '../types/NodeTypes';
import * as ServiceTypes from '../types/ServiceTypes';
import * as CertificateTypes from '../types/CertificateTypes';
import * as NodeActions from '../actions/NodeActions';
import * as MiscUtils from '../utils/MiscUtils';
import ServicesStore from '../stores/ServicesStore';
import PageInput from './PageInput';
import PageSwitch from './PageSwitch';
import PageInputSwitch from './PageInputSwitch';
import PageSelect from './PageSelect';
import PageSelectButton from './PageSelectButton';
import PageInfo from './PageInfo';
import PageSave from './PageSave';
import ConfirmButton from './ConfirmButton';
import Help from './Help';

interface Props {
	node: NodeTypes.NodeRo;
	services: ServiceTypes.ServicesRo;
	certificates: CertificateTypes.CertificatesRo;
}

interface State {
	disabled: boolean;
	changed: boolean;
	message: string;
	node: NodeTypes.Node;
	addService: string;
	forwardedChecked: boolean;
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

export default class Node extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			disabled: false,
			changed: false,
			message: '',
			node: null,
			addService: null,
			forwardedChecked: false,
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

	toggleType(typ: string): void {
		let node: NodeTypes.Node = this.state.node || this.props.node;

		let vals = (node.type || '').split('_');

		let i = vals.indexOf(typ);
		if (i === -1) {
			vals.push(typ);
		} else {
			vals.splice(i, 1);
		}

		vals.sort();

		let val = vals.join('_');
		if (val === '') {
			val = 'management';
		}

		this.set('type', val);
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
			return;
		}

		let serviceId = this.state.addService || this.props.services[0].id;

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
			node: node,
		});
	}

	render(): JSX.Element {
		let node: NodeTypes.Node = this.state.node || this.props.node;
		let active = node.requests_min !== 0 || node.memory !== 0 ||
				node.load1 !== 0 || node.load5 !== 0 || node.load15 !== 0;

		let services: JSX.Element[] = [];
		for (let serviceId of node.services) {
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

		let certificates: JSX.Element[] = [
			<option key="null" value="">Self Signed</option>,
		];
		if (this.props.certificates.length) {
			for (let certificate of this.props.certificates) {
				certificates.push(
					<option key={certificate.id} value={certificate.id}>
						{certificate.name}
					</option>,
				);
			}
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
							disabled={active || this.state.disabled}
							onConfirm={this.onDelete}
						/>
					</div>
					<PageInput
						label="Name"
						help="Name of node"
						type="text"
						placeholder="Enter name"
						value={node.name}
						onChange={(val): void => {
							this.set('name', val);
						}}
					/>
					<PageSwitch
						label="Management"
						help="Provides access to the admin console."
						checked={node.type.indexOf('management') !== -1}
						onToggle={(): void => {
							this.toggleType('management');
						}}
					/>
					<PageSwitch
						label="User"
						help="Provides access to the user console for SSH certificates."
						checked={node.type.indexOf('user') !== -1}
						onToggle={(): void => {
							this.toggleType('user');
						}}
					/>
					<PageSwitch
						label="Proxy"
						help="Provides access to the services added to the node."
						checked={node.type.indexOf('proxy') !== -1}
						onToggle={(): void => {
							this.toggleType('proxy');
						}}
					/>
					<PageInput
						hidden={node.type.indexOf('_') === -1 ||
							node.type.indexOf('management') === -1}
						label="Management Domain"
						help="Domain that will be used to access the management interface."
						type="text"
						placeholder="Enter management domain"
						value={node.management_domain}
						onChange={(val): void => {
							this.set('management_domain', val);
						}}
					/>
					<PageInput
						hidden={node.type.indexOf('_') === -1 ||
							node.type.indexOf('user') === -1}
						label="User Domain"
						help="Domain that will be used to access the user interface"
						type="text"
						placeholder="Enter user domain"
						value={node.user_domain}
						onChange={(val): void => {
							this.set('user_domain', val);
						}}
					/>
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
					<label
						className="pt-label"
						style={css.label}
						hidden={node.type.indexOf('proxy') === -1}
					>
						Services
						<Help
							title="Services"
							content="Services that can be accessed from this node. The nodes certificate must be valid for all the service domains. The node also needs to be able to access all the interal servers of the services."
						/>
						<div>
							{services}
						</div>
					</label>
					<PageSelectButton
						hidden={node.type.indexOf('proxy') === -1}
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
				</div>
				<div style={css.group}>
					<PageInfo
						fields={[
							{
								label: 'ID',
								value: node.id || 'None',
							},
							{
								valueClass: active ? '' : 'pt-text-intent-danger',
								label: 'Timestamp',
								value: MiscUtils.formatDate(node.timestamp) || 'Inactive',
							},
							{
								label: 'Requests',
								value: node.requests_min + '/min',
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
					<PageSelect
						label="Certificate"
						help="The certificate to use for this nodes web server. The certificate must be valid for all the domains that this node provides access to. This includes the management domain and any service domains."
						value={node.certificate}
						hidden={node.protocol === 'http'}
						onChange={(val): void => {
							this.set('certificate', val);
						}}
					>
						{certificates}
					</PageSelect>
					<PageInputSwitch
						label="Forwarded for header"
						help="Enable when using a load balancer. This header value will be used to get the users IP address. It is important to only enable this when a load balancer is used. If it is enabled without a load balancer users can spoof their IP address by providing a value for the header that will not be overwritten by a load balancer. Additionally the nodes firewall should be configured to only accept requests from the load balancer to prevent requests being sent directly to the node bypassing the load balancer."
						type="text"
						placeholder="Forwarded for header"
						value={node.forwarded_for_header}
						checked={this.state.forwardedChecked}
						defaultValue="X-Forwarded-For"
						onChange={(state: boolean, val: string): void => {
							let nde: NodeTypes.Node;

							if (this.state.changed) {
								nde = {
									...this.state.node,
								};
							} else {
								nde = {
									...this.props.node,
								};
							}

							nde.forwarded_for_header = val;

							this.setState({
								...this.state,
								changed: true,
								forwardedChecked: state,
								node: nde,
							});
						}}
					/>
				</div>
			</div>
			<PageSave
				style={css.save}
				hidden={!this.state.node}
				message={this.state.message}
				changed={this.state.changed}
				disabled={this.state.disabled}
				light={true}
				onCancel={(): void => {
					this.setState({
						...this.state,
						changed: false,
						forwardedChecked: false,
						node: null,
					});
				}}
				onSave={this.onSave}
			/>
		</div>;
	}
}
