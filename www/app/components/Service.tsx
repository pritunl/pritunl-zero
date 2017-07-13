/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as ServiceTypes from '../types/ServiceTypes';
import * as ServiceActions from '../actions/ServiceActions';
import PageInput from './PageInput';
import PageSave from './PageSave';
import ConfirmButton from './ConfirmButton';

interface Props {
	service: ServiceTypes.ServiceRo;
}

interface State {
	disabled: boolean;
	changed: boolean;
	message: string;
	service: ServiceTypes.Service;
}

const css = {
	card: {
		position: 'relative',
		padding: '10px',
		marginBottom: '5px',
	} as React.CSSProperties,
	remove: {
		position: 'absolute',
		top: '5px',
		right: '5px',
	} as React.CSSProperties,
};

export default class Service extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			disabled: false,
			changed: false,
			message: '',
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

	render(): JSX.Element {
		let service: ServiceTypes.Service = this.state.service ||
			this.props.service;

		return <div
			className="pt-card"
			style={css.card}
		>
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
			<PageSave
				hidden={!this.state.service}
				message={this.state.message}
				changed={this.state.changed}
				disabled={false}
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
