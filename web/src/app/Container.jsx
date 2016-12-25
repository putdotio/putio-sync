import React from 'react'
import PureRenderMixin from 'react-addons-pure-render-mixin'
import { routerActions } from 'react-router-redux'
import { connect } from 'react-redux'
import $ from 'zepto-modules'
import _ from 'lodash'

import * as Actions from './Actions'
import * as SettingsActions from '../settings/Actions'
import * as DownloadsActions from '../downloads/Actions'

import { HeaderContainerConnected } from '../header/Container'
import { FooterContainerConnected } from '../footer/Container'
import { SettingsContainer } from '../settings/Container'
import DragDrop from '../components/dragdrop'
import Loading from '../components/loading'

export class AppContainer extends React.Component {
  constructor() {
    super()
    this.shouldComponentUpdate = PureRenderMixin.shouldComponentUpdate.bind(this)
  }

  componentWillMount() {
    $('body').on('paste', this.props.HandlePaste)

    let hash = (this.props.location.hash || '').split('=')

    if (hash.length === 2) {
      return this.props.SaveToken(hash[1])
    }

    this.props.GetConfig()
  }

  componentDidMount() {
    if (_.startsWith(this.props.location.pathname, '/welcome')) {
      SettingsContainer
        .Show()
        .then(() => {
          this.props.Start()
          this.props.push('/')
        })
    }
  }

  render() {
    if (!this.props.ready) {
      return <Loading />
    }

    const processing = (this.props.processing) ? <Loading full={true} /> : null

    return (
      <div>
        <DragDrop
          onDrop={this.props.OnFileDrop}
        />

        <HeaderContainerConnected />

        <div id="content">
          {processing}

          <div className="rel">
            {this.props.children}
          </div>
        </div>

        <FooterContainerConnected />

      </div>
    )
  }
}

export const AppContainerConnected = connect(state => ({
  currentUser: state.getIn(['app', 'currentUser']),
  locale: state.getIn(['app', 'locale']),
  processing: state.getIn(['app', 'processing']),
  ready: state.getIn(['app', 'ready']),
}), Object.assign(
  Actions,
  SettingsActions,
  DownloadsActions,
  routerActions
))(AppContainer)
