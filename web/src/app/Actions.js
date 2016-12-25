import moment from 'moment'
import _ from 'lodash'
import { fromJS } from 'immutable'
import $ from 'zepto-modules'
import Jed from 'jed'

import { SyncApp, User } from '../app/Api'
import { Files, Transfers } from '../app/Api'
import { int, translations } from '../common'
import Uploader from '../uploader'
import Growl from '../components/growl'
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

export function HandlePaste(e) {
  return (dispatch, getState) => {
    let text

    if (window.clipboardData && window.clipboardData.getData) {
      text = window.clipboardData.getData('Text');
    } else if (e.clipboardData && e.clipboardData.getData) {
      text = e.clipboardData.getData('text/plain');
    }

    if (text.match(/magnet:\?xt=urn/i) === null && text.match(/^https?.*$/ig) === null) {
      return
    }

    dispatch(SetProcessing(true))
    dispatch(StartTransfers(text))
  }
}

export const ANALYSIS_SUCCESS = 'ANALYSIS_SUCCESS'
export function StartTransfers(text) {
  return (dispatch, getState) => {
    const links = _.compact(text.split('\n'))

    if (!links.length) {
      Growl.Show({
        message: translations.new_transfer_no_link_error(),
        scope: Growl.SCOPE.ERROR,
        timeout: 2,
      })

      return dispatch(SetProcessing(false))
    }

    Transfers
      .Analysis(links)
      .then(response => {
        let eFiles = _.filter(response.body.ret, f => f.error)
        let hasError = (eFiles.length === response.body.ret.length)
        let someError = (eFiles.length && eFiles.length < response.body.ret.length)

        if (hasError) {
          Growl.Show({
            message: translations.new_transfer_invalid_link_error(),
            scope: Growl.SCOPE.ERROR,
            timeout: 2,
          })

          return dispatch(SetProcessing(false))
        }

        if (someError) {
          Growl.Show({
            message: 'Some of the hashes couldn\'t add',
            scope: Growl.SCOPE.ERROR,
            timeout: 2,
          })
        }

        const files = _.chain(response.body.ret)
          .filter(f => !f.error)
          .map(f => ({
            id: f.url,
            name: f.name,
            size: f.file_size,
            email_when_complete: false,
            type: 'magnet',
            _source: f,
          }))
          .value()

        dispatch(StartFetching(files))
      })
      .catch(err => {
        dispatch(SetProcessing(false))
      })
  }
}

export function StartFetching(files) {
  return (dispatch, getState) => {
    const magnets = fromJS(files)
      .map(m => ({
        url: m.get('id'),
        email_when_complete: m.get('email_when_complete'),
        extract: m.get('extract'),
        save_parent_id: 0,
      })).toJS()

    Transfers
      .StartFetching(magnets)
      .then(response => {
        Growl.Show({
          message: translations.new_transfer_success_message(),
          scope: Growl.SCOPE.SUCCESS,
          timeout: 2,
        })

        dispatch(SetProcessing(false))
      })
      .catch(err => {
        dispatch(SetProcessing(false))
      })
  }
}

export function OnFileDrop(files) {
  return (dispatch, getState) => {
    dispatch(SetProcessing(true))

    let uploader = new Uploader()

    _.each(files, f => {
      uploader.add(f, 0)
    })

    uploader.start()
      .then(() => {
        dispatch(SetProcessing(false))

        Growl.Show({
          message: translations.app_drag_drop_success(),
          scope: Growl.SCOPE.SUCCESS,
          timeout: 3,
        })
      })
      .catch(err => {
        dispatch(SetProcessing(false))

        Growl.Show({
          message: translations.app_drag_drop_error(),
          scope: Growl.SCOPE.ERROR,
          timeout: 3,
        })
      })
  }
}
