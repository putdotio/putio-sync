import React from 'react'
import PureRenderMixin from 'react-addons-pure-render-mixin'
import { connect } from 'react-redux'
import { routerActions } from 'react-router-redux'
import { Line } from 'rc-progress'

import EmptyState from '../components/emptyState'
import Button from '../components/button'
import { Filters, translations } from '../common'
import { SettingsContainer } from '../settings/Container'

import * as Actions from './Actions'
import * as AppActions from '../app/Actions'

export class DownloadsContainer extends React.Component {
  constructor() {
    super()
    this.shouldComponentUpdate = PureRenderMixin.shouldComponentUpdate.bind(this)
  }

  pollDownloads() {
    this.props.GetDownloads()

    setTimeout(() => {
      this.pollDownloads()
    }, 3000)
  }

  componentWillMount() {
    this.pollDownloads()
  }

  render() {
    let ids = this.props.downloads.get('result')

    if (!ids.count()) {
      return (
        <div>
          <EmptyState
            title={translations.downloads_empty_state_title()}
            message={translations.downloads_empty_state_message()}
            cta={(
              <Button
                onClick={() => {
                  SettingsContainer.Show()
                }}
                label={translations.downloads_empty_state_cta_label()}
                scope="btn-success"
              />
            )}
          />
        </div>
      )
    }

    return (
      <div id="downloads-container">
        <div className="downloads">
          {this.props.downloads.get('result').map((id, i) => {
            let file = this.props.downloads.getIn([
              'entities',
              'files',
              id.toString(),
            ])

            const chunksAll = parseInt(file.getIn(['bitfield', 'bit_count_all']))
            const chunksDone = parseInt(file.getIn(['bitfield', 'bit_count_set']))
            const perc = (!chunksAll) ? 0 : chunksDone * 100 / chunksAll

            const progress = (perc === 100) ? (
              <Line
                percent={perc}
                strokeWidth="1"
                strokeColor="#1fae7d"
              />
            ) : (
              <Line
                percent={perc}
                strokeWidth="1"
                trailWidth="1"
                strokeColor="#FDCE45"
              />
            )

            const name = (
              <div className="file-name">
                {file.get('file_name')}
              </div>
            )

            const timeAgo = (file.get('download_status') === 'completed') ? (
              <span className="info-item">
                {Filters.ToTimeAgo(file.get('download_finished_at'))}
              </span>
            ) : null

            const go2File = (file.get('download_status') === 'completed') ? (
              <span className="info-item">
                <a
                  href="#"
                  onClick={e => {
                    e.preventDefault()
                    e.stopPropagation()
                    console.log(`go to file for ${file.get('file_name')}`);
                  }}>
                  Show in finder
                </a>
              </span>
            ) : null

            const downloadSpeed = (file.get('download_status') !== 'completed') ? (
              <span className="info-item">
                {Filters.Number(file.get('download_speed') / 1024, 1)} MB/s
              </span>
            ) : null

            const size = (
              <span className="info-item">
                {Filters.ToFileSize(file.get('file_length'))}
              </span>
            )

            return (
              <div
                key={file.get('file_id')}
                className="download"
              >
                {name}
                {downloadSpeed}
                {timeAgo}
                {size}
                {progress}
              </div>
            )
          })}
        </div>
      </div>
    )
  }
}

export const DownloadsContainerConnected = connect(state => ({
  currentUser: state.getIn(['app', 'currentUser']),
  status: state.getIn(['app', 'status']),
  downloads: state.getIn(['downloads', 'downloads']),
}), Object.assign(
  Actions,
  routerActions,
))(DownloadsContainer)
