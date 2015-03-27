import UserAgent from './css/user-agent';
import Panel from './panel';
import { List, ListItem } from './list';
import { extend } from 'marbles/utils';

var Main = React.createClass({
	getDefaultProps: function () {
		return {
			css: {
				margin: 16,
				display: UserAgent.isSafari() ? '-webkit-flex' : 'flex'
			},
			childrenCSS: {
				flexGrow: 1,
				WebkitFlexGrow: 1
			}
		};
	},

	render: function () {
		return (
			<div style={this.props.css}>
				<div style={extend({}, this.props.childrenCSS, { marginRight: 16 })}>
					<Panel style={{ height: '100%' }}>
						<h2>Clusters</h2>

						<List>
							<ListItem selected={true}>New</ListItem>
						</List>
					</Panel>
				</div>

				<div style={this.props.childrenCSS}>
					{this.props.children}
				</div>
			</div>
		);
	},

	componentDidMount: function () {
		this.props.dataStore.addChangeListener(this.__handleDataChange);
	},

	componentWillUnmount: function () {
		this.props.dataStore.removeChangeListener(this.__handleDataChange);
	},

	getInitialState: function () {
		return this.__getState();
	},

	__getState: function () {
		return this.props.dataStore.state;
	},

	__handleDataChange: function () {
		this.setState(this.__getState());
	}
});
export default Main;
