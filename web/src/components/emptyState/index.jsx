import React from 'react'
import _ from 'lodash'
import PureRenderMixin from 'react-addons-pure-render-mixin'

export default class EmptyState extends React.Component {
  constructor() {
    super()
    this.shouldComponentUpdate = PureRenderMixin.shouldComponentUpdate.bind(this)
  }

  render() {
    return (
      <div className="empty-state">
        <div className="empty-state-content">
          <h1>{this.props.title}</h1>
          <p>{this.props.message}</p>
          <div className="cta">
            {this.props.cta}
          </div>
        </div>
      </div>
    )
  }
}

EmptyState.propTypes = {
  title: React.PropTypes.string.isRequired,
  message: React.PropTypes.string,
  cta: React.PropTypes.element,
}
