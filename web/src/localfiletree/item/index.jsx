import React from 'react'
import _ from 'lodash'
import { Map } from 'immutable'

export default class FileTreeItem extends React.Component {
  render() {
    const className = _.compact([
      'folder-name',
      this.props.item.get('highlighted') ? 'highlighted' : null,
    ]).join(' ')

    const foldIcon = (this.props.item.get('unfolded')) ? (
      <div>
        - <i className="flaticon solid open-folder-3"></i>
      </div>
    ) : (
      <div>
        + <i className="flaticon solid folder-1"></i>
      </div>
    )

    /*<span us-spinner="{lines:8, radius:2, width:2, length:4, speed:2}"></span>*/

    return (
      <div className="item-content">
        <div className="indicators">
          <a
            href="#"
            onClick={e => {
              e.preventDefault()
              this.props.onToggleFold()
            }}
          >
            {foldIcon}
          </a>
        </div>

        <div className={className}>
          <a
            href="#"
            onClick={e => {
              e.preventDefault()
              this.props.onHighlight()
            }}
          >
            {this.props.item.get('name')}
          </a>
        </div>
      </div>
    )
  }
}

FileTreeItem.propTypes = {
  item: React.PropTypes.instanceOf(Map),
  onHighlight: React.PropTypes.func.isRequired,
  onToggleFold: React.PropTypes.func.isRequired,
}
