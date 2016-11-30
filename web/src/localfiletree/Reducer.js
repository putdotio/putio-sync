import { Map, List, fromJS } from 'immutable'
import * as Actions from './Actions'

export default function localfiletree(state = fromJS({
  tree: {
    result: null,
    entities: Map(),
    highlighted: null,
  },
}), action) {
  switch (action.type) {

    case Actions.LOCAL_FILETREE_GET_CHILDREN_SUCCESS: {
      return state
        .setIn(['tree', 'result'], action.root)
        .mergeDeepIn([
          'tree',
          'entities',
        ], fromJS(action.tree.entities))
    }

    case Actions.LOCAL_FILETREE_SET_FOLDER_FOLD: {
      return state
        .setIn([
          'tree',
          'entities',
          'tree',
          action.id.toString(),
          'unfolded',
        ], action.newFold)
    }

    case Actions.LOCAL_FILETREE_HIGHLIGHT_FOLDER: {
      if (state.getIn([
        'tree',
        'highlighted',
      ]) !== null) {
        state = state
          .setIn([
            'tree',
            'entities',
            'tree',
            state.getIn([
              'tree',
              'highlighted',
            ]).toString(),
            'highlighted',
          ], false)
      }

      return state
        .setIn([
          'tree',
          'entities',
          'tree',
          action.id.toString(),
          'highlighted',
        ], action.newHighlight)
        .setIn(['tree', 'highlighted'], (action.newHighlight) ? action.id : null)
    }

    default: {
      return state
    }

  }
}
