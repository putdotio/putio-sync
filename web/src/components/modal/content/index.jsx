import React from 'react'
import PureRenderMixin from 'react-addons-pure-render-mixin'

export default class ModalContent extends React.Component {
  constructor() {
    super()
    this.shouldComponentUpdate = PureRenderMixin.shouldComponentUpdate.bind(this)
  }

  render() {
    return (
      <div className="modal-content">
        {this.props.children}
      </div>
    )
  }
}
