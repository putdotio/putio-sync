import moment from 'moment'
import _ from 'lodash'
import { fromJS } from 'immutable'
import $ from 'zepto-modules'
import Jed from 'jed'

import { SyncApp, User } from '../app/Api'
import { Files } from '../app/Api'
import { int } from '../common'
import * as DownloadsActions from '../downloads/Actions'

export const SYNCAPP_STATUS_STOPPED = 'stopped'
export const SYNCAPP_STATUS_SYNCING = 'syncing'
export const SYNCAPP_STATUS_UPTODATE = 'up-to-date'

export const GET_CONFIG_SUCCESS = 'GET_CONFIG_SUCCESS'
export function GetConfig() {
  return (dispatch, getState) => {
    SyncApp
      .Config()
      .then(response => {
        dispatch({
          type: GET_CONFIG_SUCCESS,
          config: response.body,
        })

        dispatch(Authenticate())
        dispatch(GetSourceFolder())
      })
  }
}

export const AUTHENTICATE_USER = 'AUTHENTICATE_USER'
export const APP_READY = 'APP_READY'
export function Authenticate() {
  return (dispatch, getState) => {
    const config = getState().getIn([
      'app',
      'config',
    ])

    window.token = config.get('oauth2-token')

    if (!token) {
      return dispatch(GrantAccess())
    }

    User
      .Get()
      .then(response => {
        let user = response.body.info

        // get language file
        InitInternalization(user.settings.locale || 'en')
          .then(() => {
            dispatch({
              type: AUTHENTICATE_USER,
              user,
            })

            dispatch({
              type: APP_READY,
              ready: true,
            })
          })
      })
      .catch(err => {
        if (err && err.error_type === 'invalid_grant') {
          return dispatch(GrantAccess())
        }
      })
  }
}

export function GrantAccess() {
  return (dispatch, getState) => {
    const oauth = getState().getIn([
      'app',
      'oauth',
    ])

    window.location.href = ''
    window.location.replace(`https://put.io/v2/oauth2/authenticate?client_id=${oauth.get('id')}&response_type=token&redirect_uri=${oauth.get('callback')}`)
  }
}

export const GET_SOURCEFOLDER_SUCCESS = 'GET_SOURCEFOLDER_SUCCESS'
export function GetSourceFolder() {
  return (dispatch, getState) => {
    const config = getState().getIn([
      'app',
      'config',
    ])

    Files
      .Query(config.get('download-from'), {
        breadcrumbs: true,
      })
      .then(response => {
        let breadcrumbs = response.body.breadcrumbs.map(b => ({
          id: b[0],
          name: b[1],
        }))

        breadcrumbs.push({
          id: response.body.parent.id,
          name: response.body.parent.name,
        })

        let source = _.map(breadcrumbs, b => {
          return b.name
        }).join('/')

        dispatch({
          type: GET_SOURCEFOLDER_SUCCESS,
          source: `/${source}`,
        })
      })
  }
}

export function InitInternalization(locale) {
  return new Promise((resolve, reject) => {
    $.ajax({
      type: 'GET',
      url: `/statics/locale/${locale}.json`,
      dataType: 'json',
      timeout: 1000,
      error: resolve,
      success: data => {
        moment.locale(locale);
        int.init(data);
        resolve()
      },
    })
  })
}

export const SET_PROCESSING = 'SET_PROCESSING'
export function SetProcessing(processing) {
  return (dispatch, getState) => {
    dispatch({
      type: SET_PROCESSING,
      processing,
    })
  }
}

export function Logout() {
  return (dispatch, getState) => {
    const appId = getState().getIn([
      'app',
      'oauth',
      'id',
    ])

    User.Revoke(appId)
      .then(response => {
        User.Logout()
          .then(response => {
            dispatch(GrantAccess())
          })
      })
  }
}

export function OnFileDrop(files) {
  return (dispatch, getState) => {
    let downloadFolder = 0

    const defaultFolder = getState().getIn([
      'app',
      'currentUser',
      'settings',
      'default_download_folder',
    ])

    if (defaultFolder) {
      downloadFolder = defaultFolder
    }

    const saveTo = getState().getIn([
      'transfer',
      'saveto',
    ])

    if (saveTo) {
      downloadFolder = saveTo.get('id')
    }

    //dispatch(TasksActions.AddUploadTask(
    //files,
    //downloadFolder
    //))
  }
}
