import React from 'react'
import PureRenderMixin from 'react-addons-pure-render-mixin'
import $ from 'zepto-modules'

import { Modal } from '../index'

export default class ModalContainer extends React.Component {
  constructor() {
    super()
    this.shouldComponentUpdate = PureRenderMixin.shouldComponentUpdate.bind(this)
  }

  dismiss(event) {
    this.props.onClose(event)
  }

  componentDidMount() {
    $(document).keydown(e => {
      if (e.which === 27 && this.props.cancelable) {
        this.dismiss(Modal.REASONS.CANCEL_BY_ESC)
      }
    })
  }

  componentWillUnmount() {
    $(document).unbind('keydown')
  }

  render() {
    const className = _.compact([
      'modal-container',
      (this.props.small) ? 'small' : null,
      this.props.name,
    ]).join(' ')

    return (
      <div className="modal">
        <div
          className="modal-backdrop"
          onClick={() => {
            if (!this.props.cancelable) {
              return
            }

            this.dismiss(Modal.REASONS.CANCEL_BY_BACKDROP)
          }}
        >
        </div>

        <div className={className}>
          {this.props.children}
        </div>
      </div>
    )
  }
}

ModalContainer.propTypes = {
  name: React.PropTypes.string,
  onClose: React.PropTypes.func.isRequired,
  small: React.PropTypes.bool,
  cancelable: React.PropTypes.bool,
}
