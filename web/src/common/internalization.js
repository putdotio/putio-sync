import Jed from 'jed'

class Internalization {
  init(data) {
    this.i18n = new Jed(data);
  }

  _t(key) {
    let _this = this

    return function() {
      if (!_this.i18n) {
        return key
      }

      return _this.i18n
        .translate(key)
        .fetch()
    }
  }

  _tc(key, context) {
    let _this = this

    return function() {
      if (!_this.i18n) {
        return key
      }

      return _this.i18n
        .translate(key)
        .withContext(context)
        .fetch()
    }
  }

  _tn(singular, plural) {
    let _this = this

    return function (num) {
      if (!_this.i18n) {
        return (num > 1) ? plural : singular
      }

      return _this.i18n
        .translate(singular)
        .ifPlural(num, plural)
        .fetch()
    }
  }

  _tnc(singular, plural, context) {
    let _this = this

    return function (num) {
      if (!_this.i18n) {
        return (num > 1) ? plural : singular
      }

      return _this.i18n
        .translate(singular)
        .withContext(context)
        .ifPlural(num, plural)
        .fetch()
    }
  }
}

export default new Internalization()
