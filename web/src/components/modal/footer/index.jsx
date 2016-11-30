import React from 'react'
import PureRenderMixin from 'react-addons-pure-render-mixin'

import LinkButton from '../../button/link'
import ModalFooterAction from './action'

export default class ModalFooter extends React.Component {
  constructor() {
    super()
    this.shouldComponentUpdate = PureRenderMixin.shouldComponentUpdate.bind(this)
  }

  render() {
    return (
      <div className="modal-footer">
        <div className="modal-footer-actions">
          {this.props.children}
        </div>
      </div>
    )
  }
}
