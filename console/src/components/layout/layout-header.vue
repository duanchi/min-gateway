<template>
  <header id="topnav">

    <!-- Topbar Start -->
    <div class="navbar-custom">
      <!-- LOGO -->
      <div class="logo-box">
        <a href="/" class="logo text-center">
            <span class="logo-lg">
              <span class="logo-lg-text-light">Heron Gateway</span>
            </span>
          <span class="logo-sm">
              <span class="logo-sm-text-dark">Heron Gateway</span>
            </span>
        </a>
      </div>

      <ul class="list-unstyled topnav-menu float-right m-0">
        <li class="dropdown notification-list">
          <a href="#" @click.prevent="resetAuthorize" class="nav-link right-bar-toggle waves-effect waves-light">
            更换授权
          </a>
        </li>
        <li class="dropdown notification-list">
          <a href="#" @click.prevent="logout" class="nav-link right-bar-toggle waves-effect waves-light">
            <i class="fe-log-out"></i> 注销
          </a>
        </li>
      </ul>
      <ul class="list-unstyled topnav-menu topnav-menu-left m-0">
        <li class="active"><a class="nav-link dropdown-toggle waves-effect waves-light" href="#" @click.prevent role="button" aria-haspopup="false" aria-expanded="false">配置</a></li>
      </ul>
      <div class="clearfix"></div>
    </div>
    <!-- end Topbar -->

    <div class="topbar-menu">
      <div class="container-fluid">
        <div id="navigation">
          <!-- Navigation Menu-->
          <ul class="navigation-menu">
            <router-link to="/routes" v-slot="{ href, route, navigate, isActive, isExactActive }">
              <li :class="['has-submenu', {'active': isActive}]"><a :class="{'active': isActive}" :href="href" @click="navigate"><i class="remixicon-dashboard-line"></i>路由配置</a></li>
            </router-link>
            <router-link to="/services" v-slot="{ href, route, navigate, isActive, isExactActive }">
              <li :class="['has-submenu', {'active': isActive}]"><a :class="{'active': isActive}" :href="href" @click="navigate"><i class="remixicon-stack-line"></i>服务配置</a></li>
            </router-link>
          </ul>
          <!-- End navigation menu -->

          <div class="clearfix"></div>
        </div>
        <!-- end #navigation -->
      </div>
      <!-- end container -->
    </div>
    <!-- end navbar-custom -->

  </header>
  <!-- End Navigation Bar-->
</template>

<script>
import { Vue, Component } from 'vue-property-decorator'
import AuthorizeService from '../../services/authorize'

export default
@Component
class LayoutHeader extends Vue {
  async checkAuthorize () {
    await AuthorizeService.check()
  }

  async logout () {
    await AuthorizeService.removeToken()
    this.$router.push('/authorize')
  }

  async resetAuthorize () {
    await AuthorizeService.removeToken()
    this.$router.push('/authorize')
  }
}
</script>
