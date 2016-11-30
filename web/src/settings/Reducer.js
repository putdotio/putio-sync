import _ from 'lodash'

import * as Actions from './Actions'
import * as AppActions from '../app/Actions'
import { fromJS } from 'immutable'

export default function settings(state = fromJS({
  source: '',
  dest: '',
  simultaneous: 0,
  segments: 0,
}), action) {
  switch (action.type) {

    case AppActions.GET_CONFIG_SUCCESS: {
      return state
        .set('dest', action.config['download-to'])
        .set('simultaneous', action.config['max-parallel-files'] - 1)
        .set('segments', action.config['segments-per-file'] - 1)
    }

    case AppActions.GET_SOURCEFOLDER_SUCCESS: {
      return state
        .set('source', fromJS(action.source))
    }

    default: {
      return state
    }

  }
}
