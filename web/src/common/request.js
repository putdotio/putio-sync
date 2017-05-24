import _ from 'lodash'
import request from 'superagent'
import { Iterable } from 'immutable'
import cookie from 'react-cookie'

export default class Request {
  constructor({
    syncAPI = false
  } = {}) {
    this._url = 'https://put.io/v2'

    if (syncAPI) {
      this._url = window.location.origin;
    }

    this._query = {}
    this._oauthToken = window.token || cookie.load('oauth_token')
    this._options = {
      method: 'get',
      credentials: 'include',
      headers: {
        Accept: 'application/json',
        //Pragma: 'no-cache',
        //'If-Modified-Since': 'Mon, 26 Jul 1997 05:00:00 GMT',
        'Cache-Control': 'no-cache',
        Authorization: `token ${this._oauthToken}`,
      },
    }
  }

  Get(url) {
    this._url = `${this._url}${url}`
    this._options.method = 'get'

    return this
  }

  Post(url) {
    this._url = `${this._url}${url}`
    this._options.method = 'post'

    return this
  }

  Query(query) {
    this._query = query
    return this
  }

  Send(body) {
    if (Iterable.isIterable(body)) {
      body = body.toJS()
    }

    this._options.body = JSON.stringify(body)
    this.Set('content-type', 'application/json')

    return this
  }

  Set(header, value) {
    this._options.headers[header] = value
    return this
  }

  End(debug = false) {
    return new Promise((resolve, reject) => {
      let r = null

      if (this._options.method === 'get') {
        r = request.get(this._url)
      } else if (this._options.method === 'post') {
        r = request.post(this._url)
      }

      r = r.query(this._query)

      _.each(this._options.headers, (k, v) => {
        r = r.set(v, k)
      })

      if (this._options.body) {
        r = r.send(this._options.body)
      }

      r.end((err, res) => {
        if (debug) {
          console.log(err, res);
        }

        if (err) {
          return reject(res ? res.body : err)
        }

        res.body = JSON.parse(res.text)
        resolve(res)
      })
    })
  }
}
