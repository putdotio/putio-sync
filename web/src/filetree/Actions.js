import _ from 'lodash'
import { Schema, arrayOf, normalize } from 'normalizr'

import { Files } from '../app/Api'

const TreeSchema = new Schema('tree', {
  idAttribute: 'id',
});

TreeSchema.define({
  children: arrayOf(TreeSchema),
});

export const GET_CHILDREN_START = 'GET_CHILDREN_START'
export const GET_CHILDREN_SUCCESS = 'GET_CHILDREN_SUCCESS'
export const GET_CHILDREN_ERROR = 'GET_CHILDREN_ERROR'
export function GetChildren(parent) {
  return (dispatch, getState) => {
    Files
    .Query(parent)
    .then(response => {
      let children = response.body.files

      children = _.filter(children, c =>
        (c.content_type === 'application/x-directory')
      )

      children = _.map(children, c => ({
        id: c.id,
        name: c.name,
        parentId: c.parent_id,
      }))

      response.body.parent.children = children

      const tree = normalize({
        id: response.body.parent.id,
        name: response.body.parent.name,
        parentId: response.body.parent.parent_id,
        unfolded: true,
        children: children,
      }, TreeSchema)

      dispatch({
        type: GET_CHILDREN_SUCCESS,
        tree,
      })
    })
    .catch(err => {
      dispatch({
        type: GET_CHILDREN_ERROR,
        err,
        parent,
      })
    })
  }
}

export const HIGHLIGHT_FOLDER = 'HIGHLIGHT_FOLDER'
export function HighlightFolder(id) {
  return (dispatch, getState) => {
    let newHighlight = (getState().getIn([
      'filetree',
      'tree',
      'entities',
      'tree',
      id.toString(),
      'highlighted',
    ])) ? false : true

    dispatch({
      type: HIGHLIGHT_FOLDER,
      id,
      newHighlight,
    })
  }
}

export const SET_FOLDER_FOLD = 'SET_FOLDER_FOLD'
export function ToggleFold(id) {
  return (dispatch, getState) => {
    let newFold = (getState().getIn([
      'filetree',
      'tree',
      'entities',
      'tree',
      id.toString(),
      'unfolded',
    ])) ? false : true

    dispatch({
      type: SET_FOLDER_FOLD,
      id,
      newFold,
    })

    if (newFold) {
      dispatch(GetChildren(id))
    }
  }
}
