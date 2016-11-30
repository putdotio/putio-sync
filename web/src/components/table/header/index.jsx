import React from 'react'

export default class TableHeader extends React.Component {
  render() {
    const className = _.compact([
      'table-header',
      this.props.className,
    ]).join(' ')

    return (
      <div className={className}>
        {this.props.children}
      </div>
    )
  }
}

TableHeader.propTypes = {
  className: React.PropTypes.string,
}
