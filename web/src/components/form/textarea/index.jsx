import React from 'react'
import $ from 'zepto-modules'

export default class FormTextArea extends React.Component {
  onChange(e) {
    if (this.props.onChange) {
      this.props.onChange(e.target.value)
    }
  }

  componentDidMount() {
    if (this.props.autofocus) {
      this.refs[this.props.name].focus()
    }

    if (this.props.autoselect) {
      this.refs[this.props.name].focus()
      this.refs[this.props.name].select()
    }
  }

  componentDidUpdate() {
    if (this.props.autofocus) {
      this.refs[this.props.name].focus()
    }

    if (this.props.autoselect) {
      this.refs[this.props.name].focus()
      this.refs[this.props.name].select()
    }
  }

  render() {
    const textarea = (this.props.readOnly) ? (
      <textarea
        name={this.props.name}
        ref={this.props.name}
        value={this.props.value}
        onChange={this.onChange.bind(this)}
        onFocus={this.props.onFocus}
        onBlur={this.props.onBlur}
        readOnly
      >
      </textarea>
    ) : (
      <textarea
        name={this.props.name}
        ref={this.props.name}
        value={this.props.value}
        onChange={this.onChange.bind(this)}
        onFocus={this.props.onFocus}
        onBlur={this.props.onBlur}
      >
      </textarea>
    )

    return (
      <div className="form-textarea">
        {textarea}
      </div>
    )
  }
}

FormTextArea.propTypes = {
  name: React.PropTypes.string.isRequired,
  value: React.PropTypes.string,
  readonly: React.PropTypes.bool.isRequired,
  autofocus: React.PropTypes.bool.isRequired,
  autoselect: React.PropTypes.bool.isRequired,
  onFocus: React.PropTypes.func,
  onBlur: React.PropTypes.func,
  onChange: React.PropTypes.func,
}

FormTextArea.defaultProps = {
  readonly: false,
  autofocus: false,
  autoselect: false,
}
