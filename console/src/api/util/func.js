export default {
  getName (descriptor) {
    return descriptor.name
  },
  getParameters (descriptor) {
    return descriptor.toString().match(/.*?\(([^)]*)\)/)[1].replace(/ /g, '').split(',')
  }
}
