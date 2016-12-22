import Growl from '../components/growl'
import { translations } from '../common'
import { SyncApp } from '../app/Api'

import * as AppActions from '../app/Actions'

export const RESET = 'RESET'
export function Reset() {
  return (dispatch, getState) => {
    dispatch({
      type: RESET,
      currentUser: getState().getIn(['app', 'currentUser']),
      defaultDownloadFolder: getState().getIn(['app', 'defaultDownloadFolder']),
    })
  }
}

export const SET_SETTING_POLL_INTERVAL = 'SET_SETTING_POLL_INTERVAL'
export function SetSettings(key, value, silent = false) {
  return (dispatch, getState) => {
    let config = getState().getIn([
      'app',
      'config',
    ])

    if (key === 'pollInterval') {
      const duration = value + 'm0s'

      dispatch({
        type: SET_SETTING_POLL_INTERVAL,
        index: value,
        value: duration,
      })

      value = duration
    }

    const keyMap = {
      token:             'oauth2-token',
      source:            'download-from',
      dest:              'download-to',
      simultaneous:      'max-parallel-files',
      segments:          'segments-per-file',
      delete_remotefile: 'delete-remotefile',
      pollInterval:      'poll-interval',
    }

    config = config.set(keyMap[key], value)

    SyncApp
      .SetConfig(config.toJS())
      .then(response => {
        dispatch(AppActions.GetConfig())

        if (!silent) {
          Growl.Show({
            message: translations.settings_save_success_message(),
            scope: Growl.SCOPE.SUCCESS,
            timeout: 3,
          })
        }
      })
  }
}
