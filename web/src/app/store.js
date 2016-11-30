import {
  createStore,
  applyMiddleware,
  compose
} from 'redux'
import { browserHistory } from 'react-router'
import { combineReducers } from 'redux-immutable'
import { routerReducer, routerMiddleware } from 'react-router-redux'
import thunk from 'redux-thunk'
import createLogger from 'redux-logger'

// Import Reducers
import app from './Reducer'
import header from '../header/Reducer'
import downloads from '../downloads/Reducer'
import settings from '../settings/Reducer'
import filetree from '../filetree/Reducer'
import localfiletree from '../localfiletree/Reducer'

const logger = createLogger({
  stateTransformer: (state) => state.toJS(),
})

const dev = (window.__env__ === 'development' && window.devToolsExtension) ? window.devToolsExtension() : f => f

const store = createStore(
  combineReducers({
    app,
    header,
    downloads,
    settings,
    filetree,
    localfiletree,
    routing: routerReducer,
  }),
  compose(
    applyMiddleware(
      thunk,
      routerMiddleware(browserHistory),
    ),
    dev
  )
);

export default store
