import util from './util'
import axios from 'axios'

const CONTENT_TYPE = {
  APPLICATION_JSON: 'application/json',
  FORM_DATA: 'application/x-www-form-urlencoded',
  TEXT_PLAIN: 'text/plain',
  MULTIPART_FORM_DATA: 'multipart/form-data'
}

const REQUEST_METHOD = {
  GET: 'GET',
  POST: 'POST',
  PUT: 'PUT',
  DELETE: 'DELETE',
  HEAD: 'HEAD',
  PATCH: 'PATCH',
  OPTIONS: 'OPTIONS'
}

const _DEFAULTS = {
  prefix: '',
  method: 'GET',
  contentType: CONTENT_TYPE.APPLICATION_JSON,
  fallback: undefined,
  returnType: 'json',
  headers: {
    // 'X-Requested-With': 'XmlHttpRequest'
  }
}

const container = {}

function proxy (object) {
  /**
   * 判断是否有注解
   */
  // const hasDecorator = (undefined !== object.__map)

  // const methods = (hasDecorator ? object.__map : getApiMethod(object))
  const methods = object.__map || {}

  return new Proxy(object, {
    get (target, property, receiver) {
      if (undefined !== methods && undefined !== methods[property]) {
        return function () {
          const args = arguments
          // const parameters = hasDecorator ? (object.__map[property].parameters || []) : util.func.getParameters(methods[property])
          const parameters = object.__map[property].parameters || []
          const options = Object.assign({}, object.__options)
          Object.assign(options, object.__map[property])
          const executeOptions = Object.assign(options, {
            path: parseParameterMapping(
              object._setPath(property, parameters, args, options),
              parameters,
              args
            ),
            prefix: object._setPrefix(options.prefix || _DEFAULTS.prefix, parameters, args, options),
            method: object._setMethod(options.method || _DEFAULTS.method, parameters, args, options),
            contentType: object._setContentType(options.contentType || _DEFAULTS.contentType, parameters, args, options),
            fallback: options.fallback || _DEFAULTS.fallback,
            returnType: object._setReturnType(options.returnType || _DEFAULTS.returnType, parameters, args, options),
            headers: object._setHeaders(options.headers || _DEFAULTS.headers || {}, parameters, args, options)
          })
          const rawResult = execute(
            object._setData(property, parameters, args, options),
            executeOptions,
            target
          )

          let result = new Promise(function (resolve, reject) {
            target._setResponse(rawResult, resolve, reject, executeOptions.returns || [])
          })

          if (util.typeof(executeOptions.fallback) === 'function') {
            result = result.catch(e => {
              return new Promise((resolve, reject) => {
                executeOptions.fallback(e, resolve, reject)
              })
            })
          }

          const rawDescriptor = Reflect.apply(Reflect.get(target, property, receiver), target, args)

          switch (util.typeof(rawDescriptor)) {
            case 'promise':
              result = result.then(data => {
                return rawDescriptor.then((func) => func(data))
              })
              break

            case 'function':
              result = result.then(data => {
                const d = rawDescriptor(data)
                return d
              })
              break
          }

          return result
        }
      } else {
        let descriptor = Reflect.get(target, property, receiver)
        if (!descriptor) {
          descriptor = Reflect.get(Object.getPrototypeOf(target).constructor, property, receiver)
        }
        return descriptor
      }
    }
  })
}

function parseParameterMapping (strings, parameters, values) {
  strings = strings.replace(/[$#:]\{(.+?)\}/g, (match, p1) => {
    /* const position = parameters.indexOf(p1)

    console.log(match, p1)

    if (position !== -1) {
      return values[position]
    }
    return '' */
    const res = getValueFromParameters(p1, parameters, values)
    return res
  })

  return strings
}

function getValueFromParameters (key, parameters, value) {
  for (const n in parameters) {
    if (util.typeof(parameters[n]) === 'object' && parameters[n][key]) {
      return value[n]
    } else if (parameters[n] === key) {
      return value[n]
    }
  }

  return ''
}

function execute (data, options, target) {
  let url = ((options.prefix.endsWith('/') ? options.prefix.substring(0, options.prefix.length - 1) : options.prefix) + '/' + (options.path.startsWith('/') ? options.path.substring(1) : options.path)).replace(/\/\//g, '/')
  if (Object.keys(data).length > 0 && !['POST', 'PUT', 'PATCH'].includes(options.method)) {
    const joinString = (url.indexOf('?') !== -1 ? '&' : '?')
    url += joinString + buildQuery(data)
  }

  return axios({
    url,
    method: options.method,
    responseType: options.returnType,
    data,
    headers: options.headers
  })
}

function buildQuery (data) {
  if (!data) return ''
  return cleanArray(Object.keys(data).map(key => {
    if (data[key] === undefined) return ''
    return encodeURIComponent(key) + '=' +
      encodeURIComponent(data[key])
  })).join('&')
}

function cleanArray (actual) {
  const newArray = []
  for (let i = 0; i < actual.length; i++) {
    if (actual[i]) {
      newArray.push(actual[i])
    }
  }
  return newArray
}

/* function getApiMethod (object) {
  const properties = Object.getOwnPropertyDescriptors(Object.getPrototypeOf(object))
  const methods = {}
  for (let n in properties) {
    if (!n.startsWith('_') && n !== 'constructor' && util.typeof(properties[n].value) === 'function') {
      methods[n] = properties[n].value
    }
  }
  return methods
} */

/* function getExtendsMethod (methodName, target, last = false) {
  if (target.constructor.name === 'Api') {
    last = true
  }
  if (undefined === target[methodName]) {
    if (last) {
      return () => {
      }
    } else {
      return getExtendsMethod(methodName, Object.getPrototypeOf(methodName), last)
    }
  } else {
    return target[methodName]
  }
} */

class Api {
  constructor (options) {
    options = options || {}
    this.__name = this.__name || Object.getPrototypeOf(this).constructor.name
    const protoOptions = Object.getPrototypeOf(this).__options
    // console.log(this, protoOptions)
    this.__options = {
      path: '',
      basepath: options.path || protoOptions.path || ('/' + this.__name),
      prefix: options.prefix || protoOptions.prefix,
      method: options.method || protoOptions.method,
      contentType: options.contentType || protoOptions.contentType,
      returnType: options.returnType || protoOptions.returnType
    }

    this.__options = Object.assign(protoOptions, this.__options)

    this._init()

    const name = this.__name

    container[name] = container[name] || proxy(this)

    return container[name]
  }

  _init () {
  }

  _setPrefix (prefix, parameters, values, options) {
    return prefix
  }

  _setPath (methodName, parameters, values, options) {
    return ''
  }

  _setData (methodName, parameters, values, options) {
    return ''
  }

  _setHeaders (headers, parameters, values, options) {
    return headers || {}
  }

  _setMethod (method, parameters, values, options) {
    return method
  }

  _setContentType (contentType, parameters, values, options) {
    return contentType
  }

  _setReturnType (returnType, parameters, values, options) {
    return returnType
  }

  _setResponse (response, resolve, reject, returns) {
    return resolve(response)
  }

  _setFallback () {

  }
}

Api.prototype.__options = {
  path: '',
  basepath: '',
  prefix: _DEFAULTS.prefix,
  method: _DEFAULTS.method,
  contentType: _DEFAULTS.contentType,
  returnType: _DEFAULTS.returnType
}

function setDefaults (defaultFunc) {
  defaultFunc(_DEFAULTS)
}

export { _DEFAULTS, CONTENT_TYPE, REQUEST_METHOD, setDefaults }

export default Api
