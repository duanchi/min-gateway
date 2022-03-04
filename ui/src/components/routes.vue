<template>
  <div>
  <layout-header />

  <!-- ============================================================== -->
  <!-- Start Page Content here -->
  <!-- ============================================================== -->

  <div class="wrapper">
    <b-container fluid>
      <!-- start page title -->
      <b-row>
        <b-col cols="12">
          <div class="page-title-box">
            <div class="page-title">
              <ol class="breadcrumb m-0">
                <li class="breadcrumb-item"><a href="javascript: void(0);">Heron Gateway</a></li>
                <li class="breadcrumb-item">配置</li>
                <li class="breadcrumb-item active">路由配置</li>
              </ol>
            </div>
          </div>
        </b-col>
      </b-row>
      <!-- end page title -->
      <b-row>
        <b-col cols="9" lg="9" sm="12">
          <div class="card-box">
            <template v-if="!reordering">
              <b-button class="float-right" variant="outline-primary" v-b-modal.modal-add-route><span><i class="mdi mdi-library-plus" /> 新增路由</span></b-button>
              <b-button class="float-right mr-2" variant="outline-secondary" @click="reorder()"><span><i class="mdi mdi-reorder-horizontal" /> 重新排序</span></b-button>
            </template>
            <template v-else>
              <b-button-group class="float-right">
                <b-button variant="outline-primary" @click="reorder()"><span><i class="mdi mdi-reorder-horizontal" /> 确认排序</span></b-button>
                <b-button variant="outline-secondary" @click="resetOrder()"><span><i class="mdi mdi-reorder-horizontal" /> 还原排序</span></b-button>
              </b-button-group>
            </template>
            <b-button class="float-right mr-4" variant="outline-secondary" @click="refresh"><span><i class="mdi mdi-refresh" /> 刷新路由</span></b-button>

            <h4 class="header-title">网关路由列表</h4>
            <p class="sub-header">
              所有Heron Gateway网关配置的路由列表
            </p>

            <div class="table-responsive">
              <b-table ref="table-routes" hover head-variant="light" class="mb-0" :fields="routesFields" :items="routeList" @refreshed="reordering = false">
                <template v-slot:cell(index)="data">
                  {{ data.item.index }}
                </template>
                <template v-slot:cell(rule)="data">
                  {{ data.item.url.match }}
                </template>
                <template v-slot:cell(type)="data">
                  {{ MATCH_TYPE[data.item.url.type] }}
                </template>
                <template v-slot:cell(rewrite)="data">
                  <span class="mr-1" :key="key" v-for="(rewrite, key) in data.item.rewrite">{{rewrite.key}} => "{{rewrite.value}}"</span>
                </template>
                <template v-slot:cell(service)="data">
                  <router-link :to="`/services`">{{getServiceDisplay(data.item.service)}}</router-link>
                </template>
                <template v-slot:cell(method)="data">
                  <b-badge :key="methodName" v-for="methodName in data.item.method" :variant="'light-' +
                  (['GET', 'HEAD', 'OPTIONS', 'CONNECT'].includes(methodName) ? 'info' : '') +
                  (['POST', 'PUT', 'PATCH'].includes(methodName) ? 'success' : '') +
                  (['DELETE'].includes(methodName) ? 'danger' : '') +
                  (['ALL'].includes(methodName) ? 'primary' : '')
                ">{{methodName}}</b-badge>
                </template>
                <template v-slot:cell(authorize)="data">
                  <b-badge v-if="data.item.authorize" variant="light-success">Yes</b-badge>
                  <b-badge v-else variant="light-dark">No</b-badge>
                  <b-badge v-if="data.item.custom_token" variant="light-primary" class="ml-1">CT</b-badge>
                </template>
                <template v-slot:cell(authorize_prefix)="data">
                  <b-badge variant="light-dark">{{ data.item.authorize && (['', '0', '00', '000', '0000'].includes(data.item.authorize_prefix)) ? '默认' : data.item.authorize_prefix }}</b-badge>
                </template>
                <template v-slot:cell(options)="data">
                  <b-button-group size="sm">
                    <b-button variant="outline-secondary" @click="getRouteIntoModal(data.index);$bvModal.show('modal-add-route')">修改</b-button>
                    <b-button v-if="!reordering" variant="outline-danger" @click="remove(data.item.id)"><i class="mdi mdi-delete" /></b-button>
                    <template v-else>
                      <b-button variant="outline-secondary" @click="up(data.index)"><i class="mdi mdi-arrow-up" /></b-button>
                      <b-button variant="outline-secondary" @click="down(data.index)"><i class="mdi mdi-arrow-down" /></b-button>
                    </template>
                  </b-button-group>
                </template>
              </b-table>
            </div>
          </div>
        </b-col>
        <b-col cols="3" lg="3" sm="12">
          <side-reports :routes-count="routeList.length" :services-count="serviceList.length" />
        </b-col>
      </b-row>
    </b-container>
  </div>
  <!-- end wrapper -->
    <layout-footer />
    <b-modal id="modal-add-route" size="lg" title="新增路由">
      <template v-slot:default="{}">
        <b-row>
          <b-col cols="10" offset="1">
            <b-form class="form-horizontal">
              <div class="form-group row">
                <label class="col-sm-2 col-form-label" for="match">路由规则</label>
                <div class="col-sm-10">
                  <input type="text" id="match" class="form-control" placeholder="路由规则" v-model="createRoute.url.match">
                </div>
              </div>
              <div class="form-group row">
                <label class="col-sm-2 col-form-label" for="type">匹配类型</label>
                <b-col cols="4" sm="6">
                  <b-form-select id="type" v-model="createRoute.url.type" :options="MATCH_TYPE_OPTIONS"></b-form-select>
                </b-col>
              </div>
              <div class="form-group row">
                <label class="col-sm-2 col-form-label" for="method">请求类型</label>
                <b-col cols="10">
                  <b-form-group>
                    <b-form-checkbox-group
                      stacked
                      id="method"
                      v-model="createRoute.method"
                      :options="['ALL', 'GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'HEAD', 'OPTIONS', 'CONNECT', 'WEBSOCKET']"
                      name="method"
                      buttons
                      button-variant="outline-primary"
                      size="sm"
                    ></b-form-checkbox-group>
                  </b-form-group>
                </b-col>
              </div>
              <hr>
              <div class="form-group row">
                <label class="col-sm-2 col-form-label" for="service">匹配服务</label>
                <b-col cols="4" sm="6">
                  <b-form-select id="service" v-model="createRoute.service" :options="serviceList"></b-form-select>
                </b-col>
              </div>
              <div class="form-group row">
                <label class="col-sm-2 col-form-label" for="match">路径重写</label>
                <div class="col-10">
                  <b-row :key="index" v-for="(instance, index) in createRoute.rewrite" class="mb-2">
                    <b-col cols="12">
                      <b-input-group>
                        <b-form-input type="text" id="rewrite_key" class="form-control" placeholder="网关路径匹配" v-model="instance.key"></b-form-input>
                        <b-form-input type="text" id="rewrite_value" class="form-control" placeholder="微服务路径重写" v-model="instance.value"></b-form-input>
                        <b-input-group-append>
                          <b-button variant="outline-danger" @click="removeInstancePlaceholder(index)">删除</b-button>
                        </b-input-group-append>
                      </b-input-group>
                    </b-col>
                  </b-row>
                  <b-row>
                    <b-col cols="12">
                      <b-button variant="outline-primary" @click="addInstancePlaceholder()">添加重写规则</b-button>
                    </b-col>
                  </b-row>
                </div>
              </div>
              <hr>
              <div class="form-group row">
                <label class="col-sm-2 col-form-label" for="authorize">需要授权</label>
                <b-col cols="4" sm="6">
                  <b-form-group>
                    <b-form-radio-group
                      id="authorize"
                      v-model="createRoute.authorize"
                      :options="[{value:true, text:'需要'}, {value:false, text:'不需要'}, {value:'authorize', text: '授权路由'}]"
                      buttons
                      button-variant="outline-primary"
                      size="sm"
                      name="authorize"
                    ></b-form-radio-group>
                  </b-form-group>
                </b-col>
              </div>
              <div v-if="createRoute.authorize === 'authorize'" class="form-group row">
                <label class="col-sm-2 col-form-label" for="authorize_key">授权类型字段</label>
                <div class="col-sm-8">
                  <input type="text" id="authorize_key" class="form-control" v-model="createRoute.authorize_type_key" value="HEADER:X-Authorize-Platform">
                  <p class="form-control-plaintext">QUERY:{query key}或HEADER:{header key}</p>
                </div>
              </div>
              <div class="form-group row">
                <label class="col-sm-2 col-form-label" for="authorize_prefix">授权因子</label>
                <div class="col-sm-2">
                  <input type="text" id="authorize_prefix" class="form-control" placeholder="AUTH" v-model="createRoute.authorize_prefix">
                </div>
              </div>
              <div class="form-group row">
                <label class="col-sm-2 col-form-label" for="custom_token">自定义授权</label>
                <b-col cols="4" sm="6">
                  <b-form-group>
                    <b-form-radio-group
                      id="custom_token"
                      v-model="createRoute.custom_token"
                      :options="[{value:false, text:'使用默认JWT'}, {value:true, text:'自定义'}]"
                      buttons
                      button-variant="outline-primary"
                      size="sm"
                      name="custom_token"
                    ></b-form-radio-group>
                  </b-form-group>
                </b-col>
              </div>
            </b-form>
          </b-col>
        </b-row>
      </template>
      <template v-slot:modal-footer="{ ok, cancel, hide }">
        <!-- Emulate built in modal footer ok and cancel button actions -->
        <b-button variant="success" @click="createOrUpdate(ok)">
          {{updateId !== '' ? '确认修改' : '确认添加'}}
        </b-button>
        <b-button variant="link" @click="cancel();getRoutesItems();resetModal()">
          取消
        </b-button>
      </template>
    </b-modal>
  </div>
</template>

<script>
import LayoutHeader from './layout/layout-header'
import routes from '../services/routes'
import services from '../services/services'
import LayoutFooter from './layout/layout-footer'
import SideReports from './side-reports'
export default {
  name: 'routes',
  components: { SideReports, LayoutFooter, LayoutHeader },
  data () {
    return {
      reordering: false,
      updateId: '',
      createRoute: {
        url: {
          match: '',
          type: 'regex'
        },
        description: '',
        method: [],
        rewrite: [],
        authorize: false,
        isAuthRoute: false,
        authorize_prefix: 'AUTH',
        authorize_type_key: 'HEADER:X-Authorize-Platform'
      },
      serviceList: [],
      routeList: [],
      MATCH_TYPE: {
        'regex': '正则匹配',
        'fnmatch': 'FNMATCH匹配',
        'path': '路径匹配'
      },
      MATCH_TYPE_OPTIONS: [
        { value: 'regex', text: '正则匹配' },
        { value: 'fnmatch', text: 'FNMATCH匹配' },
        { value: 'path', text: '路径匹配' }
      ],
      routesFields: [
        {
          key: 'index',
          label: '#'
        },
        {
          key: 'rule',
          label: '路由规则'
        },
        {
          key: 'type',
          label: '匹配类型'
        },
        {
          key: 'rewrite',
          label: '路径重写'
        },
        {
          key: 'service',
          label: '服务名称'
        },
        {
          key: 'method',
          label: '请求类型'
        },
        {
          key: 'authorize',
          label: '需要授权'
        },
        {
          key: 'authorize_prefix',
          label: '授权因子'
        },
        {
          key: 'options',
          label: '',
          class: 'text-right'
        }
      ]
    }
  },
  methods: {
    async getRoutesItems () {
      const routeList = await routes.getList()
      for (let n in routeList) {
        routeList[n].index = parseInt(n) + 1
        if (routeList[n].rewrite !== undefined) {
          const rewrite = []
          for (let k in routeList[n].rewrite) {
            rewrite.push({
              key: k,
              value: routeList[n].rewrite[k]
            })
          }
          routeList[n].rewrite = rewrite
        }
      }
      this.routeList = routeList
      this.$refs['table-routes'].refresh()
      return routeList
    },
    async getServicesItems () {
      const serviceRawList = await services.getList()
      const serviceList = []

      for (let n in serviceRawList) {
        serviceList.push({
          value: serviceRawList[n].name,
          text: serviceRawList[n].display
        })
      }
      return serviceList
    },
    async createOrUpdate (callback) {
      if (this.createRoute.authorize === 'authorize') {
        this.createRoute.isAuthRoute = true
        this.createRoute.authorize = false
      } else {
        this.authorize_type_key = ''
      }

      if (this.updateId !== '') {
        await routes.update(this.updateId, this.createRoute)
      } else {
        await routes.create(this.createRoute)
      }
      callback()
      this.updateId = ''
      this.resetModal()
      this.getRoutesItems()
    },
    resetModal () {
      this.createRoute = {
        url: {
          match: '',
          type: 'regex'
        },
        description: '',
        method: [],
        rewrite: [],
        authorize: false
      }
    },
    async remove (id) {
      const value = await this.$bvModal.msgBoxConfirm('确认删除该路由?', {
        title: '确认删除路由',
        size: 'sm',
        buttonSize: 'sm',
        okVariant: 'danger',
        okTitle: '确认删除',
        cancelTitle: '取消',
        footerClass: 'p-2',
        hideHeaderClose: true,
        centered: true
      })

      if (value) {
        routes.remove(id).then(() => this.getRoutesItems())
      }
    },
    getServiceDisplay (name) {
      for (let n in this.serviceList) {
        if (name === this.serviceList[n].value) {
          return this.serviceList[n].text
        }
      }

      return name
    },
    getRouteIntoModal (index) {
      if (undefined !== this.routeList[index]) {
        this.createRoute = this.routeList[index]
        this.updateId = this.routeList[index].id

        if (this.createRoute.authorize === false && this.createRoute.authorize_type_key !== '') {
          this.createRoute.authorize = 'authorize'
        }
      }
    },
    addInstancePlaceholder () {
      this.createRoute.rewrite.push({ key: '', value: '' })
    },
    removeInstancePlaceholder (index) {
      this.createRoute.rewrite.splice(index, 1)
    },
    refresh () {
      routes.refresh()
    },
    reorder () {
      if (this.reordering) {
        const orders = []
        for (let n in this.routeList) {
          orders.push(this.routeList[n].id)
        }
        routes.reOrder(orders).then(() => this.getRoutesItems())
      }
      this.reordering = !this.reordering
    },
    resetOrder () {
      this.getRoutesItems()
      this.reordering = !this.reordering
    },
    up (index) {
      if (index > 0) {
        const temp = this.routeList[index]
        this.routeList[index] = this.routeList[index - 1]
        this.routeList[index - 1] = temp
        this.$refs['table-routes'].refresh()
      }
    },
    down (index) {
      if (index < this.routeList.length - 1) {
        const temp = this.routeList[index]
        this.routeList[index] = this.routeList[index + 1]
        this.routeList[index + 1] = temp
        this.$refs['table-routes'].refresh()
      }
    }
  },
  async mounted () {
    this.serviceList = await this.getServicesItems()
    this.getRoutesItems()
  }
}
</script>
