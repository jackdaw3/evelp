import Vue from 'vue'
import store from './store'
import App from './App.vue'
import router from './router'
import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';
import i18n from './assets/i18n/i18n'
import axios from 'axios'
import VueAxios from 'vue-axios'
import VueClipboard from 'vue-clipboard2'

Vue.config.productionTip = false
Vue.use(ElementUI);

Vue.use(VueAxios, axios)
Object.defineProperty(Vue.prototype, '$axios', { value: axios });

Vue.use(VueClipboard)

new Vue({
  router,
  store,
  i18n,
  render: h => h(App)
}).$mount('#app')
