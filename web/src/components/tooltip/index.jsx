import React from 'react'
import PureRenderMixin from 'react-addons-pure-render-mixin'

export default class Tooltip extends React.Component {
  constructor() {
    super()
    this.shouldComponentUpdate = PureRenderMixin.shouldComponentUpdate.bind(this)
  }

  render() {
    return (
      <div className="tooltip">
        {this.props.children}

        <span className="tooltiptext">
          {this.props.text}
          <span className="arrow-up" />
        </span>
      </div>
    )
  }
}

Tooltip.propTypes = {
  text: React.PropTypes.string.isRequired,
}
