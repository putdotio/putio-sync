import es6Promise from 'es6-promise'
es6Promise.polyfill()

import RSVP from 'rsvp'
Promise.defer = RSVP.defer

import React from 'react'
import ReactDOM from 'react-dom'
import { Provider } from 'react-redux'
import {
  Router,
  Route,
  RouteHandler,
  Redirect,
  IndexRoute,
  browserHistory
} from 'react-router'

import { syncHistoryWithStore } from 'react-router-redux'
import store from './store'
import { AppContainerConnected } from './Container'
import { DownloadsContainerConnected } from '../downloads/Container'
import { SettingsContainerConnected } from '../settings/Container'

const history = syncHistoryWithStore(browserHistory, store, {
  selectLocationState: (state) => state.get('routing'),
})


ReactDOM.render((
  <Provider store={store}>
    <Router history={history}>
      <Route path="/" component={AppContainerConnected}>
        <IndexRoute component={DownloadsContainerConnected} />
        <Route path="settings" component={SettingsContainerConnected} />
      </Route>
    </Router>
  </Provider>
), document.getElementById('app'));

