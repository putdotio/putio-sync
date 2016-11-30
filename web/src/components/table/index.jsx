import React from 'react'
import _ from 'lodash'

export class Table extends React.Component {
  render() {
    const className = _.compact([
      'table',
      this.props.className,
    ]).join(' ')

    const ratios = this.props.ratio.split(':')

    const children = React.Children.map(this.props.children, row => {
      if (!row) {
        return row
      }

      const columns = React.Children.map(row.props.children, (column, index) =>
        React.cloneElement(column, Object.assign({}, column.props, {
          className: `col-${ratios[index]} ${column.props.className || ''}`,
        }))
      )

      return React.cloneElement(row, row.props, columns)
    })

    return (
      <div className={className}>
        {children}
      </div>
    )
  }
}

Table.propTypes = {
  name: React.PropTypes.string,
  className: React.PropTypes.string,
  ratio: React.PropTypes.string,
}

export { default as TableRow } from './row'
export { default as TableColumn } from './column'
export { default as TableHeader } from './header'
