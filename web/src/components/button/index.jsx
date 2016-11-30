import React from 'react'
import PureRenderMixin from 'react-addons-pure-render-mixin'
import _ from 'lodash'
import $ from 'zepto-modules'

export default class Button extends React.Component {
  constructor() {
    super()
    this.shouldComponentUpdate = PureRenderMixin.shouldComponentUpdate.bind(this)
  }

  onClick(e) {
    if (!this.props.href) {
      e.preventDefault()
    }

    if (this.props.onClick) {
      this.props.onClick()
    }
  }

  render() {
    let className = _.compact([
      'btn',
      this.props.scope || 'btn-default',
      (this.props.disabled) ? 'btn-disabled' : null,
      (this.props.fixed) ? 'btn-fixed' : null,
      (this.props.className) ? this.props.className : null,
      (this.props.mini) ? 'btn-mini' : null,
      (this.props.link) ? 'btn-link' : null,
      (this.props.align) ? 'align-left' : null,
      (!this.props.label) ? 'btn-icon-only' : null,
    ]).join(' ')

    const icon = (this.props.icon) ? (
      <i className={this.props.icon}></i>
    ) : null

    const iconRight = (this.props.iconRight) ? (
      <i className={this.props.iconRight}></i>
    ) : null

    const label = (this.props.label) ? (
      <span className="btn-label">{this.props.label}</span>
    ) : null

    if (this.props.href) {
      return (
        <a
          className={className}
          onClick={this.onClick.bind(this)}
          href={this.props.href}
          target={this.props.target || '_self'}
        >
          {icon}
          {label}
          {iconRight}
        </a>
      )
    }

    return (
      <a
        className={className}
        onClick={this.onClick.bind(this)}
        target={this.props.target || '_self'}
      >
        {icon}
        {label}
        {iconRight}
      </a>
    )
  }
}

Button.propTypes = {
  onClick: React.PropTypes.func,
  href: React.PropTypes.string,
  target: React.PropTypes.string,
  className: React.PropTypes.string,
  label: React.PropTypes.string,
  fixed: React.PropTypes.bool,
  icon: React.PropTypes.string,
  scope: React.PropTypes.string,
  disabled: React.PropTypes.bool,
  mini: React.PropTypes.bool,
  link: React.PropTypes.bool,
}
