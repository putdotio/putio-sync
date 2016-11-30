import { fromJS } from 'immutable'
import * as Actions from './Actions'

export default function downloads(state = fromJS({
  resolved: false,
  fetching: false,
  speed: 0,
  downloads: {
    result: [],
    entities: {},
  },
}), action) {
  switch (action.type) {

    case Actions.GET_DOWNLOADS_START: {
      return state
        .set('fetching', true)
    }

    case Actions.GET_DOWNLOADS_SUCCESS: {
      return state
        .set('downloads', fromJS(action.downloads))
        .set('speed', action.speed)
        .set('resolved', true)
        .set('fetching', false)
    }

    default: {
      return state
    }

  }
}
