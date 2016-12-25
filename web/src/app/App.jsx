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

import _ from 'lodash'

const history = syncHistoryWithStore(browserHistory, store, {
  selectLocationState: (state) => state.get('routing'),
})

window.hostname = ({
  subdomain=null,
  protocol=null,
  path='',
  query=''
} = {}) => {
  if (protocol === false) {
    protocol = ''
  } else {
    protocol = protocol || window.location.protocol
  }

  query = (query) ? `?${query}` : ''
  let hostname = window.location.hostname

  if (_.includes(['127.0.0.1', 'localhost'], hostname)) {
    hostname = 'put.io'
  }

  let splited = hostname.split('.')

  if (subdomain === false) {
    subdomain = ''
  } else {
    subdomain = subdomain || ((splited.length > 2) ? splited[0] : '')
  }

  if (subdomain) {
    subdomain = `${subdomain}.`
  }

  let domain = (splited.length >= 2)
    ? `${splited[splited.length - 2]}.${splited[splited.length - 1]}`
    : 'put.io'

  return `${protocol}//${subdomain}${domain}${path}${query}`
}

ReactDOM.render((
  <Provider store={store}>
    <Router history={history}>
      <Route path="/" component={AppContainerConnected}>
        <IndexRoute component={DownloadsContainerConnected} />
        <Route path="welcome" component={DownloadsContainerConnected} />
        <Route path="settings" component={SettingsContainerConnected} />
      </Route>
    </Router>
  </Provider>
), document.getElementById('app'));

