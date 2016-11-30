import React from 'react'
import _ from 'lodash'
import Jed from 'jed'
import { default as rsr } from 'react-string-replace'

export default class Utils {

  static FILE_TYPE_FOLDER     = 'file_type_folder'
  static FILE_TYPE_VIDEO      = 'file_type_video'
  static FILE_TYPE_TEXT       = 'file_type_text'
  static FILE_TYPE_AUDIO      = 'file_type_audio'
  static FILE_TYPE_PDF        = 'file_type_pdf'
  static FILE_TYPE_IMAGE      = 'file_type_image'
  static FILE_TYPE_COMPRESSED = 'file_type_compressed'
  static FILE_TYPE_OTHER      = 'file_type_other'

  static FileType(contentType) {
    if (contentType === 'application/x-directory') {
      return Utils.FILE_TYPE_FOLDER
    }

    if (_.startsWith(contentType, 'video') || contentType === 'application/ogg') {
      return Utils.FILE_TYPE_VIDEO
    }

    if (_.startsWith(contentType, 'text')) {
      return Utils.FILE_TYPE_TEXT
    }

    if (_.startsWith(contentType, 'audio')) {
      return Utils.FILE_TYPE_AUDIO
    }

    if (_.startsWith(contentType, 'image')) {
      return Utils.FILE_TYPE_IMAGE
    }

    if (contentType === 'application/pdf') {
      return Utils.FILE_TYPE_PDF
    }

    if (_.includes([
      'application/x-rar',
      'application/zip',
    ], contentType)) {
      return Utils.FILE_TYPE_COMPRESSED
    }

    return Utils.FILE_TYPE_OTHER
  }

  static HandleError(err, payload) {
    window.console && console.error && console.error(err)
  }

  static Sprintf(data, ...args) {
    const hasComponent = !!_.find(args, arg => React.isValidElement(arg))

    if (!hasComponent) {
      return Jed.sprintf(data, ...args)
    }

    return rsr(
      data,
      /(\%.*s)/g,
      (match, i) => {
        let elem = args[i-1]

        if (React.isValidElement(elem)) {
          return React.cloneElement(elem, Object.assign({}, elem.props, {
            key: i
          }))
        }

        return elem
      }
    )
  }
}
