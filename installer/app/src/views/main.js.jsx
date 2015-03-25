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
							<ListItem>flynn-1427144430</ListItem>
							<ListItem>flynn-1421342430</ListItem>
							<ListItem>flynn-1420322333</ListItem>
						</List>
					</Panel>
				</div>

				<div style={this.props.childrenCSS}>
					{this.props.children}
				</div>
			</div>
		);
	}
});
export default Main;
