import _ from 'lodash'
import { Schema, arrayOf, normalize } from 'normalizr'

import { SyncApp } from '../app/Api'

const TreeSchema = new Schema('tree', {
  idAttribute: 'id',
})

TreeSchema.define({
  children: arrayOf(TreeSchema),
})

export const LOCAL_FILETREE_GET_CHILDREN_SUCCESS = 'LOCAL_FILETREE_GET_CHILDREN_SUCCESS'
export const LOCAL_FILETREE_GET_CHILDREN_ERROR = 'LOCAL_FILETREE_GET_CHILDREN_ERROR'
export function GetChildren(parent) {
  return (dispatch, getState) => {
    SyncApp
      .Tree(parent)
      .then(response => {
        let children = _.map(response.body.folders, c => ({
          id: c.path,
          name: c.name,
          parentId: c.parent,
        }))

        let parentOfParent = _.slice(
          response.body.parent.split('/'),
          0,
          response.body.parent.split('/').length - 1
        ).join('/')

        let parentName = _.slice(
          response.body.parent.split('/'),
          response.body.parent.split('/').length - 1,
        )

        const tree = normalize({
          id: response.body.parent,
          name: parentName,
          unfolded: true,
          parentId: parentOfParent,
          children,
        }, TreeSchema)

        dispatch({
          type: LOCAL_FILETREE_GET_CHILDREN_SUCCESS,
          tree,
          root: (!parent) ? response.body.parent : getState().getIn([
            'localfiletree',
            'tree',
            'result',
          ])
        })
      })
      .catch(err => {
        dispatch({
          type: LOCAL_FILETREE_GET_CHILDREN_ERROR,
          err,
          parent,
        })
      })
  }
}

export const LOCAL_FILETREE_HIGHLIGHT_FOLDER = 'LOCAL_FILETREE_HIGHLIGHT_FOLDER'
export function HighlightFolder(id) {
  return (dispatch, getState) => {
    let newHighlight = (getState().getIn([
      'localfiletree',
      'tree',
      'entities',
      'tree',
      id.toString(),
      'highlighted',
    ])) ? false : true

    dispatch({
      type: LOCAL_FILETREE_HIGHLIGHT_FOLDER,
      id,
      newHighlight,
    })
  }
}

export const LOCAL_FILETREE_SET_FOLDER_FOLD = 'LOCAL_FILETREE_SET_FOLDER_FOLD'
export function ToggleFold(id) {
  return (dispatch, getState) => {
    let newFold = (getState().getIn([
      'localfiletree',
      'tree',
      'entities',
      'tree',
      id.toString(),
      'unfolded',
    ])) ? false : true

    dispatch({
      type: LOCAL_FILETREE_SET_FOLDER_FOLD,
      id,
      newFold,
    })

    if (newFold) {
      dispatch(GetChildren(id))
    }
  }
}
