import React from 'react'
import PureRenderMixin from 'react-addons-pure-render-mixin'
import _ from 'lodash'
import $ from 'zepto-modules'
import { translations } from '../../common'

export default class DragDrop extends React.Component {
  constructor() {
    super()
    this.shouldComponentUpdate = PureRenderMixin.shouldComponentUpdate.bind(this)
  }

  isFile(event) {
    let dt = event.dataTransfer
      ? event.dataTransfer
      : event.originalEvent.dataTransfer

    if(dt.types === null) {
      return
    }

    let isFile = false

    if (dt.types.indexOf) {
      isFile = dt.types.indexOf('Files') !== -1
    } else {
      isFile = dt.types.contains('application/x-moz-file')
    }

    return isFile
  }

  componentDidMount() {
    const element = $(this.refs.dragArea)

    if (!element) {
      return
    }

    $('body').bind('dragover', event => {
      if (this.isFile(event)) {
        element.addClass('visible')
      }
    })

    element.bind('drag dragstart dragend dragover dragenter dragleave drop', event => {
      event.preventDefault()
      event.stopPropagation()
    })

    element.bind('dragleave dragend drop', event => {
      element.removeClass('visible')
    })

    element.bind('drop', event => {
      var transfer = event.dataTransfer
        ? event.dataTransfer
        : event.originalEvent.dataTransfer

      transfer.dropEffect = 'copy'

      this.props.onDrop(transfer.files)
    })
  }

  render() {
    return (
      <div
        ref="dragArea"
        className="drop-all-over-indicator"
      >
        <div className="cloud">
          <i className="flaticon solid cloud-1"></i>
          <h1>{translations.app_drag_drop_message()}</h1>
        </div>
      </div>
    )
  }
}

DragDrop.propTypes = {
  onDrop: React.PropTypes.func.isRequired,
}
