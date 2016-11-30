import React from 'react'
import PureRenderMixin from 'react-addons-pure-render-mixin'

import Spinner from '../../components/spinner'

export default class Loading extends React.Component {
  constructor() {
    super()
    this.shouldComponentUpdate = PureRenderMixin.shouldComponentUpdate.bind(this)
  }

  render() {
    let className = 'content-loading'

    if (this.props.full) {
      className += ' full'
    }

    return (
      <div className={className}>
        <Spinner size="small" />
      </div>
    )
  }
}

Loading.propTypes = {
  full: React.PropTypes.bool
}
