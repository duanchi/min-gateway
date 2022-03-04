
function RestfulApi (options) {
  return (target) => {
    const returnOptions = Object.assign({}, target.prototype.__options || {})
    options = options || {}

    Object.assign(returnOptions || {}, options || {})

    returnOptions.name = returnOptions.name || target.name

    if (!Object.prototype.hasOwnProperty.call(target, '__options') || undefined === target.__options) {
      target.prototype.__options = {}
      target.prototype.__name = returnOptions.name

      Object.assign(target.prototype.__options, returnOptions)
    }

    /* if (!target.hasOwnProperty('__map') || undefined === target.__map) {
      target.prototype.__map = {}
    } */
  }
}

function Request (options) {
  return (target, name, descriptor) => {
    if (undefined === target.__map) {
      target.__map = {}
    }

    Object.assign(target.__map, {
      [name]: options
    })
  }
}

export {
  RestfulApi, Request
}
