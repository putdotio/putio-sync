import { fromJS } from 'immutable'
import * as Actions from './Actions'

export default function header(state = fromJS({
}), action) {
  switch (action.type) {

    default: {
      return state
    }

  }
}
