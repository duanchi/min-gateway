import { vue } from '@/main'
import _typeof from './typeof'

const toast = {
  _options (options, ttl = 0) {
    return Object.assign({
      closeButton: true,
      debug: false,
      newestOnTop: true,
      progressBar: true,
      positionClass: 'toast-top-center',
      preventDuplicates: true,
      onclick: null,
      showDuration: 300,
      hideDuration: 100,
      timeOut: ttl,
      extendedTimeOut: 0,
      showEasing: 'swing',
      hideEasing: 'linear',
      showMethod: 'fadeIn',
      hideMethod: 'fadeOut',
      tapToDismiss: false
    }, options)
  },
  show (title, content, type = 'info', ttl = 2000, options = {}) {
    if (_typeof(title) === 'object') {
      type = title.type || 'info'
      options = title
      title = options.title || ''
      content = options.content || ''
    }
    options = this._options(options, ttl)

    switch (type) {
      case 'info':
        this.info(content, title, options)
        break
      case 'warning':
        this.warning(content, title, options)
        break
      case 'error':
        this.error(content, title, options)
        break
      case 'success':
        this.success(content, title, options)
        break
    }
  },
  info (title, content, ttl = 2000, options = {}) {
    options = this._options(options, ttl)
    vue.$toastr.info(content, title, options)
  },
  success (title, content, ttl = 2000, options = {}) {
    options = this._options(options, ttl)
    vue.$toastr.success(content, title, options)
  },
  error (title, content, ttl = 0, options = {}) {
    options = this._options(options, ttl)
    vue.$toastr.error(content, title, options)
  },
  warning (title, content, ttl = 0, options = {}) {
    options = this._options(options, ttl)
    vue.$toastr.warning(content, title, options)
  }
}

export default {
  show (title, content, options = {}) {
    if (_typeof(title) !== 'object') {
      options = Object.assign({
        title: title,
        text: content
      }, options)
    }
    return vue.$swal(options)
  },
  warning (title, content, ttl = undefined, options = {}) {
    return vue.$swal(Object.assign({
      type: 'warning',
      title: title,
      text: content,
      timer: ttl,
      showConfirmButton: false,
      showCancelButton: true,
      cancelButtonText: '关闭'
    }, options))
  },
  error (title, content, ttl = undefined, options = {}) {
    return vue.$swal(Object.assign({
      type: 'error',
      title: title,
      text: content,
      timer: ttl,
      showConfirmButton: false,
      showCancelButton: true,
      cancelButtonText: '关闭'
    }, options))
  },
  confirm (title, content, ttl = undefined, options = {}) {
    return vue.$swal(Object.assign({
      type: 'question',
      title: title,
      text: content,
      timer: ttl,
      showConfirmButton: true,
      showCancelButton: true,
      cancelButtonText: '取消',
      confirmButtonText: '确定'
    }, options))
  },
  success (title, content, ttl = undefined, options = {}) {
    return vue.$swal(Object.assign({
      type: 'success',
      title: title,
      text: content,
      timer: ttl,
      showConfirmButton: false,
      showCancelButton: true,
      cancelButtonText: '关闭'
    }, options))
  },
  info (title, content, ttl = undefined, options = {}) {
    return vue.$swal(Object.assign({
      type: 'info',
      title: title,
      text: content,
      timer: ttl,
      showConfirmButton: false,
      showCancelButton: true,
      cancelButtonText: '关闭'
    }, options))
  },
  toast
}
