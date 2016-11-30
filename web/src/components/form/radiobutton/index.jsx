import React from 'react'

export default class Radiobutton extends React.Component {
  onChange() {
    if (this.props.onChange) {
      this.props.onChange(this.refs.input.checked)
    }
  }

  render() {
    const name = this.props.group || 'radio-group'

    return (
      <div className="radiobutton">
        <input
          type="radio"
          ref="input"
          id={this.props.id}
          name={name}
          checked={this.props.selected}
          onChange={this.onChange.bind(this)}
        />

        <label htmlFor={this.props.id}>
          <span>
            {this.props.label || ''}
          </span>
        </label>
      </div>
    )
  }
}

Radiobutton.propTypes = {
  id: React.PropTypes.string.isRequired,
  group: React.PropTypes.string,
  onChange: React.PropTypes.func,
  selected: React.PropTypes.bool,
  label: React.PropTypes.string,
}
