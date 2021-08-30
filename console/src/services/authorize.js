import { RestfulApi, Request } from '../api/decorator/restful'
import Restful from '../api/restful'

export default new
@RestfulApi({
  name: 'authorize'
})
class Authorize extends Restful {
  setToken (token, remember = false) {
    window.localStorage.setItem('gateway-native-api-remember-access-token', remember ? 'true' : 'false')

    if (remember) {
      window.localStorage.setItem('gateway-native-api-access-token', token)
    } else {
      window.sessionStorage.setItem('gateway-native-api-access-token', token)
    }
  }

  getToken () {
    const isRemember = window.localStorage.getItem('gateway-native-api-remember-access-token')
    return isRemember === 'true' ? window.localStorage.getItem('gateway-native-api-access-token') : window.sessionStorage.getItem('gateway-native-api-access-token')
  }

  removeToken () {
    window.localStorage.removeItem('gateway-native-api-access-token')
    window.sessionStorage.removeItem('gateway-native-api-access-token')
  }

  @Request({
    method: 'GET',
    path: ''
  })
  check () {}
}()
