import React from 'react'
import PureRenderMixin from 'react-addons-pure-render-mixin'
import ReactDOM from 'react-dom'
import _ from 'lodash'
import $ from 'zepto-modules'

export default class Growl extends React.Component {
  constructor() {
    super()
    this.shouldComponentUpdate = PureRenderMixin.shouldComponentUpdate.bind(this)
  }

  static Show(options) {
    const defered = Promise.defer()

    let $wrapper = $('<div />')
    $('#app').append($wrapper)

    const close = () => {
      ReactDOM.unmountComponentAtNode($wrapper[0])
      setTimeout(function () {
        return $wrapper.remove()
      })
    }

    ReactDOM.render((
      <Growl
        message={options.message}
        scope={options.scope}
        timeout={options.timeout}
        closable={options.closable}
        onClose={() => {
          close()
          defered.resolve()
        }}
      />
    ), $wrapper[0])

    return {
      close: close,
      promise: defered.promise,
    }
  }

  static SCOPE = {
    DEFAULT: 'default',
    ERROR: 'error',
    SUCCESS: 'success',
    INFO: 'info',
    WARNING: 'warning',
  }

  dismiss(e) {
    if (e) {
      e.preventDefault()
    }

    this.props.onClose()
  }

  componentDidMount() {
    if (this.props.timeout) {
      this.timer = setTimeout(() => {
        this.dismiss()
      }, this.props.timeout * 1000)
    }
  }

  componentWillUnmount() {
    if (this.timer) {
      clearTimeout(this.timer)
    }
  }

  render() {
    const className = _.compact([
      'growl',
      this.props.scope || 'default',
      (this.props.message) ? 'visible' : '',
    ]).join(' ')

    const close = (this.props.closable) ? (
      <a
        href="#"
        onClick={this.dismiss.bind(this)}
      >
        <i className="flaticon cancel stroke x-2"></i>
      </a>
    ) : null

    return (
      <div className={className}>
        {this.props.message}
        {close}
      </div>
    )
  }
}

Growl.propTypes = {
  message: React.PropTypes.string.isRequired,
  onClose: React.PropTypes.func.isRequired,
  scope: React.PropTypes.string,
  timeout: React.PropTypes.number,
  closable: React.PropTypes.bool,
}

Growl.defaultProps = {
  scope: 'default',
  timeout: 0,
  closable: true,
}
