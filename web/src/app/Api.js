import { Map, Record } from 'immutable'

import Request from '../common/request'

export class SyncApp {
  static Config() {
    return new Request({syncAPI: true})
      .Get('/api/config')
      .End()
  }

  static SetConfig(config) {
    return new Request({syncAPI: true})
      .Post('/api/config')
      .Send(config)
      .End()
  }

  static Status() {
    return new Request({syncAPI: true})
      .Get('/api/status')
      .End()
  }

  static Tree(parent) {
    return new Request({syncAPI: true})
      .Get('/api/tree')
      .Query({
        parent,
      })
      .End()
  }
}

export class User {
  static Get(query = {}) {
    return new Request()
      .Get('/account/info')
      .Query(query)
      .End()
  }

  static Revoke(id) {
    return new Request()
      .Post(`/private/api-apps/${id}/delete`)
      .End()
  }

  static Logout() {
    return new Request({syncAPI: true})
      .Post('/api/logout')
      .End()
  }

  static Settings(userId) {
    return new Request()
      .Get('/account/settings')
      .End()
  }

  static SaveSettings(id, settings) {
    return new Request()
      .Post('/account/settings')
      .Send(settings)
      .End()
  }
}

export class File {
  static Get(id) {
    return new Request()
      .Get(`/files/${id}`)
      .End()
  }
}

export class Files {
  static Query(id) {
    return new Request()
      .Get('/files/list')
      .Query({
        parent_id: id,
        breadcrumbs: true,
      })
      .End()
  }
}
