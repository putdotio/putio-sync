import React from 'react'

export default class TableRow extends React.Component {
  render() {
    const className = _.compact([
      'table-row',
      this.props.className,
    ]).join(' ')

    return (
      <div
        className={className}
        id={this.props.id}
      >
        {this.props.children}
      </div>
    )
  }
}

TableRow.propTypes = {
  className: React.PropTypes.string,
  id: React.PropTypes.string,
}
