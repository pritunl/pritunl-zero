/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as Blueprint from '@blueprintjs/core';
import * as Constants from '../Constants';
import * as MiscUtils from '../utils/MiscUtils';

interface Props {
	style?: React.CSSProperties;
	grouped?: boolean;
	className?: string;
	hidden?: boolean;
	progressClassName?: string;
	label?: string;
	confirmMsg?: string;
	disabled?: boolean;
	onConfirm?: () => void;
}

interface State {
	dialog: boolean;
	confirm: number;
	confirming: string;
}

const css = {
	box: {
		display: 'inline-flex',
	} as React.CSSProperties,
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
	squareActionProgress: {
		position: 'absolute',
		bottom: 0,
		left: 0,
		borderRadius: 0,
		borderBottomLeftRadius: '1px',
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
		let style = {
			...this.props.style,
		};
		style.position = 'relative';

		let className = this.props.className || '';
		if (!this.props.label) {
			className += ' bp5-button-empty';
		}

		if (Constants.mobile) {
			let confirmMsg = this.props.confirmMsg ? this.props.confirmMsg :
				'Confirm ' + (this.props.label || '');

			return <div style={css.box}>
				<button
					className={'bp5-button ' + className}
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
				</button>
				<Blueprint.Dialog
					title="Confirm"
					style={css.dialog}
					isOpen={this.state.dialog}
					usePortal={true}
					portalContainer={document.body}
					onClose={this.closeDialog}
				>
					<div className="bp5-dialog-body">
						{confirmMsg}
					</div>
					<div className="bp5-dialog-footer">
						<div className="bp5-dialog-footer-actions">
							<button
								className="bp5-button"
								type="button"
								onClick={this.closeDialog}
							>Cancel</button>
							<button
								className="bp5-button bp5-intent-primary"
								type="button"
								onClick={this.closeDialogConfirm}
							>Ok</button>
						</div>
					</div>
				</Blueprint.Dialog>
			</div>
		} else {
			let confirmElem: JSX.Element;

			if (this.state.confirming) {
				let confirmStyle = {
					width: this.state.confirm * 10 + '%',
					backgroundColor: style.color,
					borderRadius: 0,
					left: 0,
				};

				let progressStyle: React.CSSProperties;
				if (this.props.grouped) {
					progressStyle = css.squareActionProgress;
				} else {
					progressStyle = css.actionProgress;
				}

				confirmElem = <div
					className={'bp5-progress-bar bp5-no-stripes ' + (
						this.props.progressClassName || '')}
					style={progressStyle}
				>
					<div className="bp5-progress-meter" style={confirmStyle}/>
				</div>;
			}

			return <button
				className={'bp5-button ' + className}
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
}
