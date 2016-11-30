import React from 'react'
import PureRenderMixin from 'react-addons-pure-render-mixin'
import { Modal } from '../index'

export default class ModalTitle extends React.Component {
  constructor() {
    super()
    this.shouldComponentUpdate = PureRenderMixin.shouldComponentUpdate.bind(this)
  }

  render() {
    const cancel = (this.props.cancelable) ? (
      <i
        onClick={() => {
          if (this.props.dismiss) {
            this.props.dismiss(Modal.REASONS.CANCEL_BY_CROSS)
          }
        }}
        className="flaticon solid x-1"
      ></i>
    ) : null

    return (
      <div className="modal-title">
        <h1>
          {this.props.title}
        </h1>

        {cancel}
      </div>
    )
  }
}

ModalTitle.propTypes = {
  title: React.PropTypes.string.isRequired,
  dismiss: React.PropTypes.func,
  cancelable: React.PropTypes.bool,
}
