import Colors from './css/colors';
import { extend } from 'marbles/utils';

var List = React.createClass({
	getDefaultProps: function () {
		return {
			baseCSS: {
				listStyle: 'none',
				margin: 0,
				padding: 0
			}
		};
	},

	render: function () {
		return (
			<ul style={extend({}, this.props.baseCSS, this.props.style || {})}>
				{this.props.children}
			</ul>
		);
	}
});

var ListItem = React.createClass({
	getDefaultProps: function () {
		return {
			baseCSS: {
				padding: '0.5em 1em',
			},

			selectedCSS: {
				backgroundColor: Colors.greenColor,
				color: Colors.whiteColor
			}
		};
	},

	getCSS: function () {
		return extend({},
			this.props.baseCSS,
			this.props.selected ? this.props.selectedCSS : {},
			this.props.style || {});
	},

	render: function () {
		return (
			<li style={this.getCSS()}>
				{this.props.children}
			</li>
		);
	}
});

export { List, ListItem };
