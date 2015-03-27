import Router from 'marbles/router';
import WizardComponent from './views/wizard';
import Dispatcher from './dispatcher';

var MainRouter = Router.createClass({
	routes: [
		{ path: '', handler: 'landingPage' },
		{ path: '/install/:install_id', handler: 'landingPage' }
	],

	willInitialize: function () {
		Dispatcher.register(this.handleEvent.bind(this));
	},

	landingPage: function (params, opts, context) {
		var props = {
			dataStore: context.dataStore
		};
		context.render(WizardComponent, props);
	},

	handleEvent: function (event) {
		var installID;
		switch (event.name) {
			case 'INSTALL_ABORT':
				this.history.navigate('/');
			break;

			case 'LAUNCH_INSTALL_SUCCESS':
				installID = event.res.id;
				this.history.navigate('/install/'+ installID);
			break;
		}
	}
});
export default MainRouter;
