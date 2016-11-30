import React from 'react'

export default class FormRowTitle extends React.Component {
  render() {
    return (
      <div className="form-row-title">
        {this.props.title}
      </div>
    )
  }
}

FormRowTitle.propTypes = {
  title: React.PropTypes.string.isRequired,
}
