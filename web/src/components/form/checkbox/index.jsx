import React from 'react'
import $ from 'zepto-modules'

export default class Checkbox extends React.Component {
  onChange(event, checked) {
    event.stopPropagation()

    if (this.props.onChange) {
      this.props.onChange(checked, event)
    }

    return true
  }

  render() {
    const id = `checkbox-${this.props.id}`

    const disabled = (this.props.disabled) ? (
      <span className="disabled"></span>
    ) : null

    const effectiveArea = (this.props.effectiveArea) ? (
      <span
        className="effective-area"
        onClick={e => {
          this.onChange(e, !this.refs.input.checked)
        }}
      ></span>
    ) : null

    return (
      <div className="checkbox">
        {disabled}
        {effectiveArea}

        <input
          type="checkbox"
          ref="input"
          id={id}
          onChange={e => {
            this.onChange(e, this.refs.input.checked)
          }}
          checked={this.props.checked}
        />

        <label htmlFor={id}>
          {this.props.label || ''}
        </label>
      </div>
    )
  }
}

Checkbox.propTypes = {
  id: React.PropTypes.string.isRequired,
  onChange: React.PropTypes.func,
  checked: React.PropTypes.bool,
  label: React.PropTypes.string,
  disabled: React.PropTypes.bool,
  effectiveArea: React.PropTypes.bool,
}

Checkbox.defaultProps = {
  effectiveArea: false,
}
