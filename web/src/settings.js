module.exports = {

  title: '运维统一平台',

  devServer: {
    proxy: 'https://localhost:8000'
  },
  /**
   * @type {boolean} true | false
   * @description Whether fix the header
   */
  fixedHeader: false,

  /**
   * @type {boolean} true | false
   * @description Whether show the logo in sidebar
   */
  sidebarLogo: true
}
