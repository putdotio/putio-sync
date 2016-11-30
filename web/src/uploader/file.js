import _ from 'lodash'
import HttpUploader from './http'
import TusUploader from './tus'

export default class File {
  constructor(file, to, options) {
    this.file = file
    this.name = options.name || this.file.name
    this.parentId = to
    this.small = (this.file.size < 4000000)
    this.uploader  = (this.small)
      ? new HttpUploader(file)
      : new TusUploader(file)
  }

  abort() {
    if (!this.uploader) {
      return
    }

    this.uploader.abort()
  }

  upload(onProgress) {
    return this.uploader.upload(
      (progress) => {
        onProgress(this, progress)
      },

      this.name,
      this.parentId
    )
  }
}
