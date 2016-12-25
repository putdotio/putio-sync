import _ from 'lodash'
import moment from 'moment'

import * as Actions from './Actions'
import * as AppActions from '../app/Actions'
import { fromJS } from 'immutable'

export default function settings(state = fromJS({
  pollInterval: {
    value: '',
    selected: 1,
    data: [
      { value: '1', label: '1' },
      { value: '2', label: '2' },
      { value: '3', label: '3' },
      { value: '4', label: '4' },
      { value: '5', label: '5' },
      { value: '6', label: '6' },
      { value: '7', label: '7' },
      { value: '8', label: '8' },
      { value: '9', label: '9' },
      { value: '10', label: '10' },
    ],
  },
  source: 0,
  sourceStr: '',
  dest: '',
  simultaneous: 0,
  segments: 0,
  delete_remotefile: false,
}), action) {
  switch (action.type) {

    case AppActions.GET_CONFIG_SUCCESS: {
      const minute = parseInt(action.config['poll-interval'].split('m')[0])

      return state
        .setIn(['pollInterval', 'selected'], minute - 1)
        .setIn(['pollInterval', 'value'], action.config['poll-interval'])
        .set('dest', action.config['download-to'])
        .set('simultaneous', action.config['max-parallel-files'] - 1)
        .set('segments', action.config['segments-per-file'] - 1)
        .set('delete_remotefile', action.config['delete-remotefile'])
    }

    case Actions.SET_SETTING: {
      if (_.includes(['simultaneous', 'segments'], action.key)) {
        action.value = action.value - 1
      }

      return state.set(action.key, action.value)
    }

    case Actions.SET_SETTING_POLL_INTERVAL: {
      return state
        .setIn(['pollInterval', 'value'], action.value)
        .setIn(['pollInterval', 'selected'], action.index)
    }

    case AppActions.GET_SOURCEFOLDER_SUCCESS: {
      return state
        .set('source', action.source)
        .set('sourceStr', action.sourceStr)
    }

    default: {
      return state
    }

  }
}
