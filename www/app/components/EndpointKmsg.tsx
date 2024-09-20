/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as Theme from '../Theme';
import * as EndpointTypes from '../types/EndpointTypes';
import * as EndpointActions from '../actions/EndpointActions';
import * as CheckActions from '../actions/CheckActions';
import {Ace} from "ace-builds";
import AceEditor from "react-ace";

import "ace-builds/src-noconflict/mode-text";
import "ace-builds/src-noconflict/theme-dracula";
import "ace-builds/src-noconflict/theme-eclipse";

interface Props {
	endpoint?: string;
	check?: string;
	disabled: boolean;
}

interface State {
	data: string;
	loading: boolean;
	cancelable: boolean;
}

const css = {
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
	button: {
		margin: '8px 0 0 8px',
	} as React.CSSProperties,
	buttons: {
		marginTop: '8px',
	} as React.CSSProperties,
	group: {
		flex: 1,
		minWidth: '250px',
		margin: '0 10px',
	} as React.CSSProperties,
	editorGroup: {
		marginBottom: '15px',
	} as React.CSSProperties,
};

export default class EndpointKmsg extends React.Component<Props, State> {
	loaded: boolean;
	editor: Ace.Editor;

	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			data: '',
			loading: false,
			cancelable: false,
		};
	}

	componentDidMount(): void {
		Theme.addChangeListener(this.onChange);
	}

	componentWillUnmount(): void {
		Theme.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		this.setState({
			...this.state,
		});
	}

	setLoading(): void {
		this.setState({
			...this.state,
			loading: true,
			cancelable: true,
		});
	}

	setLoaded(): void {
		this.setState({
			...this.state,
			loading: false,
			cancelable: false,
		});
	}

	update(): void {
		let loading = true;
		this.setLoading();

		let logResp: Promise<any>

		if (this.props.endpoint) {
			logResp = EndpointActions.log(
				this.props.endpoint,
				'kmsg',
			)
		} else {
			logResp = CheckActions.log(
				this.props.check,
				'check',
			)
		}

		logResp.then((data: EndpointTypes.LogData): void => {
			if (loading) {
				loading = false;
				this.setLoaded();
			}

			this.setState({
				...this.state,
				data: data.join(''),
			});

			setTimeout((): void => {
				if (this.editor) {
					this.editor.scrollToLine(Number.POSITIVE_INFINITY,
						false, true, null);
				}
			}, 100);
			setTimeout((): void => {
				if (this.editor) {
					this.editor.scrollToLine(Number.POSITIVE_INFINITY,
						false, true, null);
				}
			}, 200);
			setTimeout((): void => {
				if (this.editor) {
					this.editor.scrollToLine(Number.POSITIVE_INFINITY,
						false, true, null);
				}
			}, 300);
			setTimeout((): void => {
				if (this.editor) {
					this.editor.scrollToLine(Number.POSITIVE_INFINITY,
						false, true, null);
				}
			}, 400);
		}).catch((): void => {
			if (loading) {
				loading = false;
				this.setLoaded();
			}
		});
	}

	render(): JSX.Element {
		if (this.props.disabled) {
			return <div/>;
		}

		if (!this.loaded) {
			this.loaded = true;
			setTimeout((): void => {
				this.update();
			});
		}

		let title = ""
		if (this.props.endpoint) {
			title = "Dmesg"
		} else {
			title = "Error Log"
		}

		let refreshDisabled = false;
		let refreshLabel = '';
		let refreshClass = 'bp5-button';
		if (Object.entries(this.state.cancelable).length) {
			refreshLabel = 'Cancel';
			refreshClass += ' bp5-intent-warning bp5-icon-delete'
		} else {
			if (Object.entries(this.state.loading).length) {
				refreshDisabled = true;
			}
			refreshLabel = 'Refresh';
			refreshClass += ' bp5-intent-success bp5-icon-refresh'
		}

		return <div>
			<div className="layout horizontal wrap bp5-border" style={css.header}>
				<h3 style={css.heading}>{title}</h3>
				<div className="flex"/>
				<div style={css.buttons}>
					<button
						className={refreshClass}
						style={css.button}
						disabled={refreshDisabled}
						type="button"
						onClick={(): void => {
							if (Object.entries(this.state.cancelable).length) {
								if (this.props.endpoint) {
									EndpointActions.dataCancel();
								} else {
									CheckActions.dataCancel();
								}
							} else {
								this.update();
							}
						}}
					>
						{refreshLabel}
					</button>
				</div>
			</div>
			<div className="layout horizontal wrap" style={css.editorGroup}>
				<AceEditor
					name={this.props.endpoint + "-dmesg"}
					theme={Theme.editorTheme()}
					height="400px"
					width="100%"
					mode="text"
					fontSize="12px"
					wrapEnabled={true}
					showPrintMargin={false}
					showGutter={true}
					readOnly={true}
					value={this.state.data}
					editorProps={{
						$blockScrolling: true,
					}}
					setOptions={{
						showFoldWidgets: false,
					}}
					onLoad={(editor: Ace.Editor): void => {
						this.editor = editor;
						this.editor.scrollToLine(Number.POSITIVE_INFINITY,
							false, false, null);
					}}
					onChange={(): void => {
						if (this.editor) {
							this.editor.scrollToLine(Number.POSITIVE_INFINITY,
								false, true, null);
						}
					}}
				/>
			</div>
		</div>;
	}
}
