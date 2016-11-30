import React from 'react'

export default class TableColumn extends React.Component {
  render() {
    const className = _.compact([
      'table-column',
      this.props.className,
    ]).join(' ')

    return (
      <div className={className}>
        {this.props.children}
      </div>
    )
  }
}

TableColumn.propTypes = {
  className: React.PropTypes.string,
}
