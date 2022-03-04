import i from 'i'

const stringUtils = i()

stringUtils.replaceAll = function (search, replace, string) {
  if (string != null) { string = string.replace(new RegExp(search, 'g'), replace) }
  return string
}

export default stringUtils
