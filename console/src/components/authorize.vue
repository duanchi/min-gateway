<template>
  <div class="account-pages mt-5 mb-5">
    <div class="container">
      <div class="row justify-content-center">
        <div class="col-md-8 col-lg-6 col-xl-5">
          <div class="card">

            <div class="card-body p-4">

              <div class="text-center w-75 m-auto">
                <router-link to="/authorize">
                  <span>HERON GATEWAY</span>
                </router-link>
                <p class="text-muted mb-4 mt-3">请将AccessToken置入输入框进行认证</p>
              </div>

              <form action="#">

                <div class="form-group mb-3">
                  <textarea class="form-control" id="access-token" rows="3" v-model="token" placeholder="ACCESS TOKEN可以在网关部署时的环境变量中获取"></textarea>
                </div>
                <div class="form-group form-check mb-3 text-center">
                  <input type="checkbox" class="form-check-input" v-model="isRemember" id="remember">
                  <label class="form-check-label" for="remember">记住授权</label>
                </div>

                <div class="form-group mb-0 text-center">
                  <button class="btn btn-primary btn-block" @click.prevent="authorize"> 认证 </button>
                </div>
              </form>

            </div> <!-- end card-body -->
          </div>
          <!-- end card -->

        </div> <!-- end col -->
      </div>
      <!-- end row -->
    </div>
    <!-- end container -->
  </div>
</template>

<script>
import { Component, Vue } from 'vue-property-decorator'
import AuthorizeService from '../services/authorize'

export default
@Component
class Login extends Vue {
  token = ''

  isRemember = false

  authorize () {
    AuthorizeService.setToken(this.token, this.isRemember)
    this.$router.push('/')
  }
}
</script>
