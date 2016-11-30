import React from 'react'
import PureRenderMixin from 'react-addons-pure-render-mixin'

export default class FormInput extends React.Component {
  constructor(props) {
    super(props)
    this.shouldComponentUpdate = PureRenderMixin.shouldComponentUpdate.bind(this)
  }

  onChange(e) {
    if (this.props.onChange) {
      this.props.onChange(e.target.value)
    }
  }

  componentDidMount() {
    if (this.props.autofocus) {
      this.refs[this.props.name].focus()
    }
  }

  render() {
    const type = (this.props.password) ? 'password' : 'text'
    const placeholder = this.props.placeholder || ''

    const input = (this.props.disabled) ? (
      <input
        name={this.props.name}
        ref={this.props.name}
        type={this.props.type || type}
        value={this.props.value}
        onChange={this.onChange.bind(this)}
        placeholder={placeholder}
        onKeyDown={this.props.onKeyDown}
        onKeyPress={this.props.onKeyPress}
        onKeyUp={this.props.onKeyUp}
        onFocus={this.props.onFocus}
        onBlur={this.props.onBlur}
        disabled
      />
    ) : (
      <input
        name={this.props.name}
        ref={this.props.name}
        type={type}
        value={this.props.value}
        onChange={this.onChange.bind(this)}
        placeholder={placeholder}
        onKeyDown={this.props.onKeyDown}
        onKeyPress={this.props.onKeyPress}
        onKeyUp={this.props.onKeyUp}
        onFocus={this.props.onFocus}
        onBlur={this.props.onBlur}
      />
    )

    return (
      <div className="form-input">
        {input}
        {this.props.action}
      </div>
    )
  }
}

FormInput.propTypes = {
  name: React.PropTypes.string.isRequired,
  value: React.PropTypes.string,
  type: React.PropTypes.string,
  action: React.PropTypes.element,
  autofocus: React.PropTypes.bool,
  placeholder: React.PropTypes.string,
  password: React.PropTypes.bool,
  onFocus: React.PropTypes.func,
  onBlur: React.PropTypes.func,
  onChange: React.PropTypes.func,
  disabled: React.PropTypes.bool,
}
