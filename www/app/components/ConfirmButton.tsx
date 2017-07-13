/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as Blueprint from '@blueprintjs/core';
import * as Constants from '../Constants';
import * as MiscUtils from '../utils/MiscUtils';

interface Props {
	style?: React.CSSProperties;
	className?: string;
	hidden?: boolean;
	progressClassName?: string;
	label?: string;
	disabled?: boolean;
	onConfirm?: () => void;
}

interface State {
	dialog: boolean;
	confirm: number;
	confirming: string;
}

const css = {
	actionProgress: {
		position: 'absolute',
		bottom: 0,
		left: 0,
		borderRadius: 0,
		borderBottomLeftRadius: '3px',
		borderBottomRightRadius: '3px',
		width: '100%',
		height: '4px',
	} as React.CSSProperties,
	dialog: {
		width: '180px',
	} as React.CSSProperties,
};

export default class ConfirmButton extends React.Component<Props, State> {
	constructor(props: Props, context: any) {
		super(props, context);
		this.state = {
			dialog: false,
			confirm: 0,
			confirming: null,
		};
	}

	openDialog = (): void => {
		this.setState({
			...this.state,
			dialog: true,
		});
	}

	closeDialog = (): void => {
		this.setState({
			...this.state,
			dialog: false,
		});
	}

	closeDialogConfirm = (): void => {
		this.setState({
			...this.state,
			dialog: false,
		});
		if (this.props.onConfirm) {
			this.props.onConfirm();
		}
	}

	confirm = (evt: React.MouseEvent<{}>): void => {
		let confirmId = MiscUtils.uuid();

		if (evt.shiftKey) {
			if (this.props.onConfirm) {
				this.props.onConfirm();
			}
			return;
		}

		this.setState({
			...this.state,
			confirming: confirmId,
		});

		let i = 10;
		let id = setInterval(() => {
			if (i > 100) {
				clearInterval(id);
				setTimeout(() => {
					if (this.state.confirming === confirmId) {
						this.setState({
							...this.state,
							confirm: 0,
							confirming: null,
						});
						if (this.props.onConfirm) {
							this.props.onConfirm();
						}
					}
				}, 250);
				return;
			} else if (!this.state.confirming) {
				clearInterval(id);
				this.setState({
					...this.state,
					confirm: 0,
					confirming: null,
				});
				return;
			}

			if (i % 10 === 0) {
				this.setState({
					...this.state,
					confirm: i / 10,
				});
			}

			i += 1;
		}, 3);
	}

	clearConfirm = (): void => {
		this.setState({
			...this.state,
			confirm: 0,
			confirming: null,
		});
	}

	render(): JSX.Element {
		let confirmElem: JSX.Element;

		let style = this.props.style || {};
		style.position = 'relative';

		if (Constants.mobile) {
			confirmElem = <Blueprint.Dialog
				title="Confirm"
				style={css.dialog}
				isOpen={this.state.dialog}
				onClose={this.closeDialog}
			>
				<div className="pt-dialog-body">
					Confirm {this.props.label}
				</div>
				<div className="pt-dialog-footer">
					<div className="pt-dialog-footer-actions">
						<button
							className="pt-button"
							type="button"
							onClick={this.closeDialog}
						>Cancel</button>
						<button
							className="pt-button pt-intent-primary"
							type="button"
							onClick={this.closeDialogConfirm}
						>Ok</button>
					</div>
				</div>
			</Blueprint.Dialog>;
		} else {
			if (this.state.confirming) {
				let confirmStyle = {
					width: this.state.confirm * 10 + '%',
					backgroundColor: style.color,
					borderRadius: 0,
					left: 0,
				};

				confirmElem = <div
					className={'pt-progress-bar pt-no-stripes ' + (
						this.props.progressClassName || '')}
					style={css.actionProgress}
				>
					<div className="pt-progress-meter" style={confirmStyle}/>
				</div>;
			}
		}

		let style = this.props.style || {};
		style.position = 'relative';

		return <button
			className={'pt-button ' + (this.props.className || '')}
			style={style}
			type="button"
			hidden={this.props.hidden}
			disabled={this.props.disabled}
			onMouseDown={Constants.mobile ? undefined : this.confirm}
			onMouseUp={Constants.mobile ? undefined : this.clearConfirm}
			onMouseLeave={Constants.mobile ? undefined : this.clearConfirm}
			onClick={Constants.mobile ? this.openDialog : undefined}
		>
			{this.props.label}
			{confirmElem}
		</button>;
	}
}
