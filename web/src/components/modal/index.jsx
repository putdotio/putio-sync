import ReactDOM from 'react-dom'
import _ from 'lodash'
import React from 'react'
import $ from 'zepto-modules'

import Button from '../button'
import ModalContainer from './container'
import ModalTitle from './title'
import ModalContent from './content'
import ModalFooter from './footer'
import ModalFooterAction from './footer/action'

export class Modal {
  constructor(options) {
    this.options = options
    this.content = this.options.content.bind(this) || null
    this.defered = Promise.defer()
  }

  static REASONS = {
    CANCEL_BY_ESC: 'CANCEL_BY_ESC',
    CANCEL_BY_CROSS: 'CANCEL_BY_CROSS',
    CANCEL_BY_BACKDROP: 'CANCEL_BY_BACKDROP',
    CANCEL_BY_CTA: 'CANCEL_BY_CTA',
    CANCEL_BY_FOOTER: 'CANCEL_BY_FOOTER',
  }

  Destroy(reason, data) {
    ReactDOM.unmountComponentAtNode(this.$wrapper[0])
    setTimeout(() => {
      return this.$wrapper.remove()
    })

    if (reason) {
      this.defered.reject(reason)
    } else {
      this.defered.resolve(data)
    }
  }

  Show() {
    this.$wrapper = $('<div />')
    $('#app').append(this.$wrapper)

    const cancelable = (typeof this.options.cancelable === 'undefined') ? true : this.options.cancelable

    const title = (this.options.title) ? (
      <ModalTitle
        title={this.options.title}
        cancelable={cancelable}
        dismiss={reason => {
          this.Destroy(reason)
        }}
      />
    ) : null

    ReactDOM.render((
      <ModalContainer
        onClose={reason => {
          this.Destroy(reason)
        }}
        small={this.options.small || false}
        name={this.options.name}
        cancelable={cancelable}
      >
        {title}

        <ModalContent>
          {this.content()}
        </ModalContent>
      </ModalContainer>
    ), this.$wrapper[0])

    return this.defered.promise
  }
}

export { default as ModalContainer } from './container'
export { default as ModalTitle } from './title'
export { default as ModalContent } from './content'
export { default as ModalFooter } from './footer'
export { default as ModalFooterAction } from './footer/action'
