import { Schema, arrayOf, normalize } from 'normalizr'
import { Downloads } from './Api'
import * as AppActions from '../app/Actions'

const FilesSchema = new Schema('files', {
  idAttribute: 'file_id',
})

export const GET_DOWNLOADS_START = 'GET_DOWNLOADS_START'
export const GET_DOWNLOADS_SUCCESS = 'GET_DOWNLOADS_SUCCESS'
export const GET_DOWNLOADS_ERROR = 'GET_DOWNLOADS_ERROR'
export function GetDownloads() {
  return (dispatch, getState) => {
    dispatch({
      type: GET_DOWNLOADS_START,
    })

    Downloads
      .Get()
      .then(response => {
        const downloads = normalize(
          response.body.files,
          arrayOf(FilesSchema)
        )

        dispatch({
          type: GET_DOWNLOADS_SUCCESS,
          downloads,
          status: response.body.status,
          speed: response.body.total_speed,
        })
      })
  }
}

export const START_DOWNLOADS_SUCCESS = 'START_DOWNLOADS_SUCCESS'
export function Start() {
  return (dispatch, getState) => {
    dispatch(AppActions.SetProcessing(true))

    Downloads
      .Start()
      .then(response => {
        dispatch({
          type: START_DOWNLOADS_SUCCESS,
        })

        dispatch(AppActions.SetProcessing(false))
      })
  }
}

export const STOP_DOWNLOADS_SUCCESS = 'STOP_DOWNLOADS_SUCCESS'
export function Stop() {
  return (dispatch, getState) => {
    dispatch(AppActions.SetProcessing(true))

    Downloads
      .Stop()
      .then(response => {
        dispatch({
          type: START_DOWNLOADS_SUCCESS,
        })

        dispatch(AppActions.SetProcessing(false))
      })
  }
}

export function ClearFinished() {
  return (dispatch, getState) => {

    Downloads
      .ClearFinished()
      .then(response => {
        dispatch(GetDownloads())
      })
  }
}
