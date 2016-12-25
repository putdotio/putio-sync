import cookie from 'react-cookie'
import request from 'superagent'

export default class HttpUploader {
  constructor(file) {
    this.file = file

    this.endpoint = hostname({
      protocol: 'https:',
      subdomain: 'upload',
      path: '/v2/files/upload',
    })

    this.req = null
  }

  abort() {
    if (!this.req) {
      return
    }

    this.req.abort()
  }

  upload(onProgress, name, parentId) {
    return new Promise((resolve, reject) => {
      let fd = new FormData()
      fd.append('file', this.file)
      fd.append('filename', name)
      fd.append('parent_id', parentId)

      this.req = request
        .post(this.endpoint)
        .withCredentials()
        .set('Content-Type', undefined)
        .set('Authorization', `token ${cookie.load('oauth_token')}`)
        .send(fd)
        .on('progress', e => {
          if (typeof onProgress === 'function') {
            onProgress(e.percent || 0)
          }
        })
        .end((err, res) => {
          if (err) {
            return reject(err)
          }

          if (typeof onProgress === 'function') {
            onProgress(100)
          }

          resolve(res)
        })
    })
  }
}
