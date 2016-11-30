import React from 'react'
import PureRenderMixin from 'react-addons-pure-render-mixin'

export default class ModalFooterAction extends React.Component {
  constructor() {
    super()
    this.shouldComponentUpdate = PureRenderMixin.shouldComponentUpdate.bind(this)
  }

  render() {
    return (
      <div className="modal-footer-action">
        {this.props.children}
      </div>
    )
  }
}
