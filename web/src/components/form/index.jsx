import React from 'react'
import _ from 'lodash'

export class Form extends React.Component {
  render() {
    const className = _.compact([
      'form',
      this.props.className || '',
    ]).join(' ')

    return (
      <div className={className}>
        {this.props.children}
      </div>
    )
  }
}

Form.propTypes = {
  name: React.PropTypes.string,
  className: React.PropTypes.string,
}

export { default as Row } from './row'
export { default as RowTitle } from './rowTitle'
export { default as RowHelp } from './rowHelp'
export { default as Input } from './input'
export { default as TextArea } from './textarea'
export { default as Select } from './select'
export { default as Radiobutton } from './radiobutton'
export { default as Checkbox } from './checkbox'
