import { RestfulApi, Request } from '../api/decorator/restful'
import Restful from '../api/restful'

export default new
@RestfulApi({
  name: 'services'
})
class Services extends Restful {
  @Request({
    method: 'GET',
    path: ''
  })
  getList () {}

  @Request({
    method: 'POST',
    path: '',
    parameters: ['data']
  })
  create (data) {}

  @Request({
    method: 'PUT',
    path: '/?id=#{id}',
    parameters: [{ id: '-' }, 'data']
  })
  update (id, data) {}

  @Request({
    method: 'DELETE',
    path: '/?id=#{id}',
    parameters: [{ id: '-' }]
  })
  remove (id) {}
}()
