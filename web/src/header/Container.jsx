import React from 'react'
import PureRenderMixin from 'react-addons-pure-render-mixin'
import { connect } from 'react-redux'
import $ from 'zepto-modules'

import * as Actions from './Actions'
import * as AppActions from '../app/Actions'
import * as DownloadsActions from '../downloads/Actions'

import Button from '../components/button'
import Tooltip from '../components/tooltip'
import Confirmation from '../components/modal/confirmation'
import { translations } from '../common'
import { SettingsContainer } from '../settings/Container'

export class HeaderContainer extends React.Component {
  constructor() {
    super()
    this.shouldComponentUpdate = PureRenderMixin.shouldComponentUpdate.bind(this)
  }

  render() {
    let status = null
    let action = null

    if (this.props.status === AppActions.SYNCAPP_STATUS_STOPPED) {
      action = (
        <div className="item">
          <Tooltip text={translations.app_header_hint_play()}>
            <a
              href="#"
              onClick={e => {
                e.preventDefault()
                e.stopPropagation()

                this.props.Start()
              }}
            >
              <i className="flaticon solid play-2" />
            </a>
          </Tooltip>
        </div>
      )

      status = (
        <div>
          Stopped
        </div>
      )
    }

    if (this.props.status === AppActions.SYNCAPP_STATUS_SYNCING) {
      action = (
        <div className="item">
          <Tooltip text={translations.app_header_hint_pause()}>
            <a
              href="#"
              onClick={e => {
                e.preventDefault()
                e.stopPropagation()

                this.props.Stop()
              }}
            >
              <i className="flaticon solid pause-2" />
            </a>
          </Tooltip>
        </div>
      )

      status = (
        <div>
          Syncing
        </div>
      )
    }

    if (this.props.status === AppActions.SYNCAPP_STATUS_UPTODATE) {
      action = (this.props.config.get('is-paused')) ? (
        <div className="item">
          <Tooltip text={translations.app_header_hint_play()}>
            <a
              href="#"
              onClick={e => {
                e.preventDefault()
                e.stopPropagation()

                this.props.Start()
              }}
            >
              <i className="flaticon solid play-2" />
            </a>
          </Tooltip>
        </div>
      ) : (
        <div className="item">
          <Tooltip text={translations.app_header_hint_pause()}>
            <a
              href="#"
              onClick={e => {
                e.preventDefault()
                e.stopPropagation()

                this.props.Stop()
              }}
            >
              <i className="flaticon solid pause-2" />
            </a>
          </Tooltip>
        </div>
      )

      status = (
        <div>
          <div className="status status-syncing">
            <span className="message">
              <i className="flaticon solid checkmark-2" /> All finished
            </span>
          </div>
        </div>
      )
    }

    return (
      <header>
        <div className="left">
          {status}
        </div>

        <div className="right">
          <div className="item">
            <Tooltip text={translations.app_header_hint_clear()}>
              <a
                href="#"
                onClick={e => {
                  e.preventDefault()
                  e.stopPropagation()

                  Confirmation.Show({
                    small: true,
                    message: 'Are you sure to clear finished downloads?',
                  })
                    .then(r => {
                      this.props.ClearFinished()
                    })
                }}
              >
                <i className="flaticon solid magic-wand-1" />
              </a>
            </Tooltip>
          </div>

          {action}

          <div className="item">
            <Tooltip text={translations.app_header_hint_settings()}>
              <a
                href="#"
                onClick={e => {
                  e.preventDefault()
                  e.stopPropagation()

                  SettingsContainer.Show()
                }}
              >
                <i className="flaticon stroke settings-2" />
              </a>
            </Tooltip>
          </div>

          <div className="item">
            <Tooltip text={translations.app_header_hint_logout()}>
              <a
                href="#"
                onClick={e => {
                  e.preventDefault()
                  e.stopPropagation()

                  Confirmation.Show({
                    small: true,
                    message: 'Are you sure to logout?',
                  })
                    .then(r => {
                      this.props.Logout()
                    })
                }}
              >
                <i className="flaticon stroke logout-1" />
              </a>
            </Tooltip>
          </div>
        </div>
      </header>
    )
  }
}

export const HeaderContainerConnected = connect((state) => ({
  currentUser: state.getIn(['app', 'currentUser']),
  config: state.getIn(['app', 'config']),
  status: state.getIn(['app', 'status']),
}), Object.assign(
  Actions,
  AppActions,
  DownloadsActions,
))(HeaderContainer)
