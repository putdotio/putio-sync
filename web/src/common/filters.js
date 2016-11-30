import moment from 'moment'

class Filters {
  static Number(input, precision = 0) {
    try {
      let num = parseFloat(input)
      return num.toFixed(precision)
    } catch (e) {
      return NaN
    }
  }

  static ToFileSize(input, withUnit, precision) {
    var i = -1;

    if (typeof precision === 'undefined') {
      precision = 1
    }

    var byteUnits = [
      ' kB',
      ' MB',
      ' GB',
      ' TB',
      'PB',
      'EB',
      'ZB',
      'YB',
    ]

    do {
      input = input / 1024
      i++
    } while (input > 1024)

    var s = (i > 1)
      ? Math.max(input, 0.1).toFixed(precision)
      : Math.max(input, 0.1).toFixed(0)

    if (typeof withUnit === 'undefined' || withUnit === true) {
      s = s + byteUnits[i]
    }

    return s
  }

  static GetUTCOffset() {
    let offset = new Date().getTimezoneOffset()
    return (offset / -60)
  }

  static FormatDate(input, format = 'LL') {
    return moment
      .utc(input)
      .add(Filters.GetUTCOffset(), 'hours')
      .format(format)
  }

  static DaysDiff(date1, date2) {
    date1 = moment(date1).format('YYYY-MM-DD')
    date2 = moment(date2).format('YYYY-MM-DD')
    return moment(date1).diff(moment(date2), 'days')
  }

  static ToTimeAgo(input, utc = false) {
    //@TODO: remove undefined check
    if (input === undefined) {
      return moment('2014-06-20 00:14:23', 'YYYY-MM-DD hh:mm:ss').fromNow()
    }

    if (utc) {
      return moment.utc(input, 'YYYY-MM-DD HH:mm:ss').fromNow()
    }

    return moment(input, 'YYYY-MM-DD HH:mm:ss').fromNow()
  }
}

export default Filters
