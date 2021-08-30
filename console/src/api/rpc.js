import Api from './'
import util from '@/util'

export default class Rpc extends Api {
  _init (target) {
    return target
  }

  _setPrefix (prefix, parameters, values, options) {
    return prefix + '/rpc'
  }

  _setPath (methodName, parameters, values, options) {
    return (options.package || '') + options.basepath.replace('/', '.') + '::' + methodName
  }

  _setMethod (method, parameters, values, options) {
    return 'POST'
  }

  _setHeaders (headers, parameters, values, options) {
    const token = localStorage.getItem('access_token')
    if (token) {
      headers.Authorization = 'Bearer ' + token
    }
    return super._setHeaders(headers, parameters, values, options)
  }

  _setContentType (contentType, parameters, values) {
    if (this.__options.contentType === 'multipart/form-data') {
      return undefined
    }
  }

  _setData (methodName, parameters, values, options) {
    const data = []

    for (const i in values) {
      data.push(values[i])
    }

    return data
  }

  _setResponse (promise, resolve, reject, returns) {
    return promise.then(
      response => {
        if (response?.data) {
          let result = {}
          if (returns && util.typeof(returns) === 'array' && returns.length > 0) {
            if (util.typeof(response.data) === 'array') {
              for (const i in returns) {
                result[returns[i]] = response.data[i] !== undefined ? response.data[i] : undefined
              }
            }
          } else {
            result = util.typeof(response.data) === 'array' && response.data.length > 0 ? response.data[0] : response.data
          }

          resolve(result)
        } else if (response.status === 204) {
          resolve(null)
        } else {
          reject(response)
        }
      }).catch(error => {
      reject(error.response)
    }
    )
  }
}
