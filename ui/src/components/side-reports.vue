<template>
  <div class="card-box">
    <h4 class="header-title">配置统计</h4>
    <p class="text-muted">Heron-Gateway网关路由和微服务统计信息</p>

    <b-row>
      <b-col cols="6">
        <p class="text-muted mb-1">路由表</p>
        <h3 class="mt-0 font-20">{{routesCount}} 条</h3>
      </b-col>
      <b-col cols="6">
        <p class="text-muted mb-1">微服务</p>
        <h3 class="mt-0 font-20">{{servicesCount}} 个</h3>
      </b-col>
    </b-row>
  </div>
</template>

<script>
export default {
  name: 'side-reports',
  props: {
    routesCount: {
      type: Number,
      default: 0
    },
    servicesCount: {
      type: Number,
      default: 0
    }
  }
}
</script>
