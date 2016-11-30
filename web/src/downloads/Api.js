import { Map, Record } from 'immutable'

import Request from '../common/request'

export class Downloads {
  static Get() {
    return new Request({syncAPI: true})
      .Get('/api/list-downloads')
      .Query()
      .End()
  }

  static Start() {
    return new Request({syncAPI: true})
      .Get('/api/start')
      .Query()
      .End()
  }

  static Stop() {
    return new Request({syncAPI: true})
      .Get('/api/stop')
      .Query()
      .End()
  }

  static ClearFinished() {
    return new Request({syncAPI: true})
      .Post('/api/clear')
      .End()
  }
}
