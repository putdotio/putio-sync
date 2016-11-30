import _ from 'lodash'
import { List, Map } from 'immutable'
import React from 'react'
import PureRenderMixin from 'react-addons-pure-render-mixin'
import { connect } from 'react-redux'
import store from '../app/store'
import Button from '../components/button'
import LinkButton from '../components/button/link'
import * as Actions from './Actions'
import Loading from '../components/loading'

import {
  Modal,
  ModalContent,
  ModalFooter,
  ModalFooterAction,
} from '../components/modal'
import FileTreeItem from '../filetree/item'
import { translations } from '../common'

export class LocalFileTree extends React.Component {
  constructor() {
    super()
    this.shouldComponentUpdate = PureRenderMixin.shouldComponentUpdate.bind(this)
  }

  static Show() {
    return new Modal({
      title: translations.filetree_title(),
      name: 'filetree-modal',
      content: function() {
        return (
          <LocalFileTreeConnected
            store={store}
            onCancel={(reason) => {
              this.Destroy(reason)
            }}
            onSelect={folder => {
              this.Destroy(null, folder)
            }}
          />
        )
      }
    }).Show()
  }

  onClick(folder, e) {
    e.preventDefault()
    this.props.SelectFolder(folder)
    this.props.onSelected(folder)
  }

  componentWillUpdate(props) {
    //this.props.GetChildren(0)
  }

  componentWillMount() {
    this.props.GetChildren()
  }

  renderChilds(root) {
    return (
      <ul>
        {root.get('children').map((id, i) => {
          let child = this.props.tree.getIn([
            'entities',
            'tree',
            id.toString(),
          ])

          const className = (child.get('selected'))
            ? 'selected'
            : null

          const children = (child.get('children') && child.get('unfolded'))
            ? this.renderChilds(child)
            : null

          return (
            <li
              className={className}
              key={child.get('id')}
            >
              <FileTreeItem
                item={child}
                onHighlight={() => {
                  this.props.HighlightFolder(child.get('id'))
                }}
                onToggleFold={() => {
                  this.props.ToggleFold(child.get('id'))
                }}
              />
              {children}
            </li>
          )
        })}
      </ul>
    )
  }

  render() {
    if (this.props.tree.get('result') === null) {
      return <Loading full={true} />
    }

    let root = this.props.tree.getIn([
      'entities',
      'tree',
      this.props.tree.get('result').toString(),
    ])

    return (
      <div>
        <div id="filetree">
          <ul>
            <li
              className="selected"
              key={root.get('id')}
            >
              <FileTreeItem
                item={root}
                onHighlight={() => {
                  this.props.HighlightFolder(root.get('id'))
                }}
                onToggleFold={() => {
                  this.props.ToggleFold(root.get('id'))
                }}
              />

              {this.renderChilds(root)}
            </li>
          </ul>
        </div>

        <ModalFooter>
          <ModalFooterAction>
            <LinkButton
              onClick={() => this.props.onCancel(Modal.REASONS.CANCEL_BY_FOOTER) }
              label={translations.filetree_action_cancel_label()}
            />
          </ModalFooterAction>

          <ModalFooterAction>
            <Button
              label="Select"
              scope="btn-success"
              disabled={this.props.tree.get('highlighted') === null}
              onClick={() => {
                const selected = this.props.tree.getIn([
                  'entities',
                  'tree',
                  this.props.tree.get('highlighted').toString(),
                ])

                this.props.onSelect(selected)
              }}
            />
          </ModalFooterAction>
        </ModalFooter>
      </div>
    )
  }
}

export const LocalFileTreeConnected = connect((state) => ({
  tree: state.getIn(['localfiletree', 'tree']),
}), Actions)(LocalFileTree)
