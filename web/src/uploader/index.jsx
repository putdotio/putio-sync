import React from 'react'
import _ from 'lodash'
import { List } from 'immutable'
import File from './file'

export default class Uploader {
  constructor() {
    this.files = new List()
  }

  add(file, to, options = {}) {
    this.files = this.files.push(new File(file, to, options))
    return this
  }

  abort() {
    this.files.map(f => {
      f.abort()
    })
  }

  start(onProgress) {
    return Promise.all(
      this.files.map(f => f.upload(onProgress))
    )
  }

  remove(file) {
    file.abort()
    this.files = this.files.filter(f => (r.name !== file.name))
  }

  count() {
    return this.files.size
  }
}
