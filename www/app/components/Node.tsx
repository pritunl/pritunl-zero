/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as NodeTypes from '../types/NodeTypes';
import * as NodeActions from '../actions/NodeActions';
import * as MiscUtils from '../utils/MiscUtils';
import PageInput from './PageInput';
import PageSelect from './PageSelect';
import PageInfo from './PageInfo';
import PageSave from './PageSave';
import ConfirmButton from './ConfirmButton';

interface Props {
	node: NodeTypes.NodeRo;
}

interface State {
	disabled: boolean;
	changed: boolean;
	message: string;
	node: NodeTypes.Node;
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
};

export default class Node extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			disabled: false,
			changed: false,
			message: '',
			node: null,
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

	render(): JSX.Element {
		let node: NodeTypes.Node = this.state.node ||
			this.props.node;

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
				</div>
				<div style={css.group}>
					<PageInfo
						fields={[
							{
								label: 'ID',
								value: node.id || 'None',
							},
							{
								label: 'Type',
								value: node.type,
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
