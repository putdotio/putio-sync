import { Map, fromJS } from 'immutable'

import * as Actions from './Actions'
import * as DownloadsActions from '../downloads/Actions'

export default function app(state = fromJS({
  processing: false,
  ready: false,
  oauth: {
    id: 2527,
    callback: 'http://127.0.0.1:3000',
  },
  config: {},
  currentUser: null,
  status: Actions.SYNCAPP_STATUS_STOPPED,
}), action) {
  switch (action.type) {

    case DownloadsActions.GET_DOWNLOADS_SUCCESS: {
      return state
        .set('status', action.status)
    }

    case Actions.GET_CONFIG_SUCCESS: {
      return state
        .set('config', fromJS(action.config))
    }

    case Actions.AUTHENTICATE_USER: {
      return state
        .set('currentUser', fromJS(action.user))
    }

    case Actions.APP_READY: {
      return state
        .set('ready', action.ready)
    }

    case Actions.SET_PROCESSING: {
      return state
        .set('processing', action.processing)
    }

    case Actions.SYNCAPP_STATUS_SUCCESS: {
      return state
        .set('status', action.status)
    }

    default: {
      return state
    }

  }
}
