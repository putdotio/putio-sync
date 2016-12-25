import cookie from 'react-cookie'
import _ from 'lodash'
import tus from 'tus-js-client'

export default class TusUploader {
  constructor(file) {
    this.file = file

    this.endpoint = hostname({
      protocol: 'https:',
      subdomain: 'upload',
      path: '/files/',
    })

    this.tusUpload = null
  }

  abort() {
    if (!this.tusUpload) {
      return
    }

    this.tusUpload.abort()
  }

  upload(onProgress, name, parentId) {
    return new Promise((resolve, reject) => {
      this.tusUpload = new tus.Upload(this.file, {
        endpoint: this.endpoint,
        withCredentials: true,
        chunkSize: 4194304,
        headers: {
          Authorization: `token ${cookie.load('oauth_token')}`,
        },
        metadata: _.extend({
          name: name,
          type: this.file.type,
          parent_id: '',
          callback_url: '',
        }, {
          parent_id: parentId,
        }),

        onError: () => {
          reject(err)
        },

        onProgress: (bytesUploaded, bytesTotal) => {
          if (typeof onProgress === 'function') {
            let percent = (bytesUploaded / bytesTotal * 100).toFixed(2);
            onProgress(percent)
          }
        },

        onSuccess: () => {
          resolve({
            name: this.file.name,
            url: this.endpoint,
          });
        },
      })

      this.tusUpload.start()
    })
  }
}
