import React from 'react'
import Spinjs from 'spin.js'

export default class Spinner extends React.Component {
  componentDidMount() {
    let config = {
      lines: 11,
      length: 5,
      width: 3,
      radius: 6,
      scale: 1,
      corners: 1,
      color: '#000',
      opacity: 0.25,
      rotate: 0,
      direction: 1,
      speed: 1.5,
      trail: 60,
      fps: 20,
      zIndex: 2e9,
      className: 'spinner',
      top: '50%',
      left: '50%',
      shadow: false,
      hwaccel: false,
      position: 'absolute',
    }

    if (this.props.size === 'small') {
      config.lines = 9
      config.length = 5
      config.width = 2
      config.radius = 2
    }

    if (this.props.size === 'tiny') {
      config.lines = 6
      config.length = 3
      config.width = 2
      config.radius = 1
    }

    this.spinner = new Spinjs(config);

    if (!this.props.stopped) {
      this.spinner.spin(this.refs.container)
    }
  }

  componentWillReceiveProps(newProps) {
    if (newProps.stopped === true && !this.props.stopped) {
      this.spinner.stop()
    } else if (!newProps.stopped && this.props.stopped === true) {
      this.spinner.spin(this.refs.container)
    }
  }

  componentWillUnmount() {
    this.spinner.stop()
  }

  render() {
    const size = this.props.size || '';

    return (
      <span
        ref="container"
        className={`spinner-container ${size}`}
      >
      </span>
    )
  }
}

Spinner.propTypes = {
  size: React.PropTypes.string,
}
