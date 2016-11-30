import React from 'react'
import PureRenderMixin from 'react-addons-pure-render-mixin'
import { connect } from 'react-redux'

import { routerActions } from 'react-router-redux'
import { Filters } from '../common'

export class FooterContainer extends React.Component {
  constructor() {
    super()
    this.shouldComponentUpdate = PureRenderMixin.shouldComponentUpdate.bind(this)
  }

  render() {
    return (
      <footer>
        <div className="total-speed">
          {Filters.Number(this.props.speed / 1024, 1)} MB/s
        </div>
      </footer>
    )
  }
}

export const FooterContainerConnected = connect(state => ({
  currentUser: state.getIn(['app', 'currentUser']),
  speed: state.getIn(['downloads', 'speed']),
}), Object.assign(
  routerActions,
))(FooterContainer)
