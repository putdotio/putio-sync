import React from 'react'
import PureRenderMixin from 'react-addons-pure-render-mixin'
import $ from 'zepto-modules'

import Button from '../button'
import LinkButton from '../button/link'

import {
  Modal,
  ModalContent,
  ModalFooter,
  ModalFooterAction,
} from './index'

export default class Confirmation extends React.Component {
  constructor() {
    super()
    this.shouldComponentUpdate = PureRenderMixin.shouldComponentUpdate.bind(this)
  }

  static Show(options) {
    return new Modal({
      name: options.name || 'prompt-modal',
      small: options.small || false,
      title: options.title || null,
      content: function() {
        return (
          <Confirmation
            message={options.message}
            actionLabel={options.actionLabel}
            onCancel={reason => {
              this.Destroy(reason)
            }}
            onConfirm={response => {
              this.Destroy(null, true)
            }}
          />
        )
      }
    }).Show()
  }

  componentWillUnmount() {
    $(window).off("keydown.confirmation")
  }

  componentDidMount() {
    $(window).on("keydown.confirmation", e => {
      if (e.key !== 'Enter') {
        return
      }

      this.props.onConfirm(this.text)
    })
  }

  render() {
    const message = (typeof this.props.message === 'function') ? this.props.message() : this.props.message

    return (
      <div>
        {message}

        <ModalFooter>
          <ModalFooterAction>
            <LinkButton
              onClick={() => this.props.onCancel(Modal.REASONS.CANCEL_BY_FOOTER) }
              label="Cancel"
            />
          </ModalFooterAction>

          <ModalFooterAction>
            <Button
              label={this.props.actionLabel || "Okay"}
              scope="btn-success"
              onClick={() => this.props.onConfirm(this.text) }
            />
          </ModalFooterAction>
        </ModalFooter>
      </div>
    )
  }
}

Confirmation.propTypes = {
  message: React.PropTypes.oneOfType([
    React.PropTypes.string,
    React.PropTypes.func,
  ]).isRequired,
  actionLabel: React.PropTypes.string,
  onCancel: React.PropTypes.func.isRequired,
  onConfirm: React.PropTypes.func.isRequired,
}

Confirmation.defaultProps = {
  message: () => {},
}
