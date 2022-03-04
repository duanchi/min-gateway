module.exports = {
  lintOnSave: false,
  publicPath: '/',
  devServer: {
    proxy: {
      '/api': {
        pathRewrite: {
          '/api': '/_api'
        },
        target: 'http://127.0.0.1:9080',
        ws: true,
        changeOrigin: true
      }
    }
  }
}
