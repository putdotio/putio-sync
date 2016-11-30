import React from 'react'

export default class FormRow extends React.Component {
  render() {
    return (
      <div className="form-row">
        {this.props.children}
      </div>
    )
  }
}
