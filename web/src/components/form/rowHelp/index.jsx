import React from 'react'

export default class FormRowHelp extends React.Component {
  render() {
    return (
      <div className="form-row-help">
        {this.props.help}
      </div>
    )
  }
}

FormRowHelp.propTypes = {
  help: React.PropTypes.string.isRequired,
}
