import { RestfulApi, Request } from '../api/decorator/restful'
import Restful from '../api/restful'

export default new
@RestfulApi({
  name: 'routes'
})
class Routes extends Restful {
  @Request({
    method: 'GET',
    path: ''
  })
  getList () {}

  create (route) {
    const rewrite = {}
    for (let n in route.rewrite) {
      rewrite[route.rewrite[n].key] = route.rewrite[n].value
    }
    route.rewrite = rewrite

    return this._create(route)
  }

  @Request({
    method: 'POST',
    path: '',
    parameters: ['data']
  })
  _create (data) {}

  update (id, route) {
    const rewrite = {}
    for (let n in route.rewrite) {
      rewrite[route.rewrite[n].key] = route.rewrite[n].value
    }
    route.rewrite = rewrite

    return this._update(id, route)
  }

  @Request({
    method: 'PUT',
    path: '/?id=#{id}',
    parameters: [{ id: '-' }, 'data']
  })
  _update (id, data) {}

  @Request({
    method: 'DELETE',
    path: '/?id=#{id}',
    parameters: [{ id: '-' }]
  })
  remove (id) {}

  @Request({
    method: 'PUT',
    path: '/?scope=order',
    parameters: ['data']
  })
  reOrder (data) {}
}()
