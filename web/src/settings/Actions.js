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

export const SET_SETTING = 'SET_SETTING'
export const SET_SETTING_POLL_INTERVAL = 'SET_SETTING_POLL_INTERVAL'
export function SetSettings(key, value) {
  return (dispatch, getState) => {
    let config = getState().getIn([
      'app',
      'config',
    ])

    if (key === 'pollInterval') {
      const duration = value + 'm0s'

      dispatch({
        type: SET_SETTING_POLL_INTERVAL,
        index: value - 1,
        value: duration,
      })

      value = duration
    } else {
      dispatch({
        type: SET_SETTING,
        key,
        value,
      })
    }
  }
}

export function SaveSettings(silent = false) {
  return (dispatch, getState) => {
    let config = getState().getIn([
      'app',
      'config',
    ])

    const keyMap = {
      token:             'oauth2-token',
      source:            'download-from',
      dest:              'download-to',
      simultaneous:      'max-parallel-files',
      segments:          'segments-per-file',
      delete_remotefile: 'delete-remotefile',
      pollInterval:      'poll-interval',
    }

    config = config
      .set(keyMap['source'], getState().getIn(['settings', 'source']))
      .set(keyMap['dest'], getState().getIn(['settings', 'dest']))
      .set(keyMap['simultaneous'], getState().getIn(['settings', 'simultaneous']) + 1)
      .set(keyMap['segments'], getState().getIn(['settings', 'segments']) + 1)
      .set(keyMap['delete_remotefile'], getState().getIn(['settings', 'delete_remotefile']))
      .set(keyMap['pollInterval'], getState().getIn(['settings', 'pollInterval', 'value']))

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
