import React from 'react'
import PureRenderMixin from 'react-addons-pure-render-mixin'
import { routerActions } from 'react-router-redux'
import _ from 'lodash'
import { connect } from 'react-redux'
import { fromJS } from 'immutable'
import $ from 'zepto-modules'

import store from '../app/store'
import * as Actions from './Actions'
import * as AppActions from '../app/Actions'
import { translations } from '../common'

import {
  Form,
  Row,
  RowTitle,
  RowHelp,
  Input,
  Select,
  Checkbox,
} from '../components/form'
import Button from '../components/button'
import { Modal } from '../components/modal'
import { FileTree } from '../filetree/Container'
import { LocalFileTree } from '../localfiletree/Container'

export class SettingsContainer extends React.Component {
  constructor(props) {
    super(props)
    this.shouldComponentUpdate = PureRenderMixin.shouldComponentUpdate.bind(this)
  }

  static Show() {
    return new Modal({
      name: 'settings-modal',
      title: 'Settings',
      content: function() {
        return (
          <SettingsContainerConnected
            store={store}
            onCancel={reason => {
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

  render() {
    return (
      <div className="settings-container">
        <Form className="settings-form">
          <Row>
            <RowTitle title={translations.settings_poll_interval_message()} />

            <Select
              name="poll-interval"
              options={this.props.pollInterval.get('data')}
              selected={this.props.pollInterval.get('selected')}
              required={true}
              onSelect={index => {
                this.props.SetSettings('pollInterval', index + 1)
              }}
            />
          </Row>

          <Row>
            <RowTitle title={translations.settings_source_folder_label()} />
            <Input
              name="source-folder"
              value={this.props.sourceStr}
              disabled={true}
              action={(
                <Button
                  label={translations.settings_source_folder_action_label()}
                  className="input-action"
                  scope="btn-success"
                  onClick={() => {
                    FileTree
                      .Show()
                      .then(folder => {
                        this.props.GetSourceFolder(folder.get('id'))
                      })
                  }}
                />
              )}
            />
          </Row>

          <Row>
            <RowTitle title={translations.settings_dest_folder_label()} />
            <Input
              name="dest-folder"
              value={this.props.dest}
              disabled={true}
              action={(
                <Button
                  label={translations.settings_dest_folder_action_label()}
                  className="input-action"
                  scope="btn-success"
                  onClick={() => {
                    LocalFileTree
                      .Show()
                      .then(folder => {
                        this.props.SetSettings('dest', folder.get('id'))
                      })
                  }}
                />
              )}
            />
          </Row>

          <Row>
            <RowTitle title={translations.settings_simultaneous_download_label()} />
            <Select
              name="simultaneous_download"
              required={true}
              options={fromJS([
                {label: '1', value: '1'},
                {label: '2', value: '2'},
                {label: '3', value: '3'},
                {label: '4', value: '4'},
                {label: '5', value: '5'},
              ])}
              selected={this.props.simultaneous}
              onSelect={index => {
                this.props.SetSettings('simultaneous', index + 1)
              }}
            />
          </Row>

          <Row>
            <RowTitle title={translations.settings_segments_perfile_label()} />
            <Select
              name="segments-per-file"
              required={true}
              options={fromJS([
                {label: '1', value: '1'},
                {label: '2', value: '2'},
                {label: '3', value: '3'},
                {label: '4', value: '4'},
                {label: '5', value: '5'},
              ])}
              selected={this.props.segments}
              onSelect={index => {
                this.props.SetSettings('segments', index + 1)
              }}
            />
          </Row>

          <Row>
            <Checkbox
              id="settings-delete-remote-file"
              label={translations.settings_delete_source_message()}
              checked={this.props.delete_remotefile}
              onChange={r => {
                this.props.SetSettings('delete_remotefile', r)
              }}
            />

            <RowHelp
              help={translations.settings_delete_source_hint()}
            />
          </Row>

          <Row>
            <Button
              onClick={this.props.SaveSettings}
              label={translations.settings_save_label()}
              scope="btn-success"
            />
          </Row>
        </Form>
      </div>
    )
  }
}

export const SettingsContainerConnected = connect(state => ({
  currentUser: state.getIn(['app', 'currentUser']),
  pollInterval: state.getIn(['settings', 'pollInterval']),
  source: state.getIn(['settings', 'source']),
  sourceStr: state.getIn(['settings', 'sourceStr']),
  dest: state.getIn(['settings', 'dest']),
  simultaneous: state.getIn(['settings', 'simultaneous']),
  segments: state.getIn(['settings', 'segments']),
  delete_remotefile: state.getIn(['settings', 'delete_remotefile']),
}), Object.assign(
  Actions,
  AppActions,
  routerActions
))(SettingsContainer);
