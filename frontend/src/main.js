import Vue from 'vue'
import App from './App.vue'
import router from './router'
import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';
import i18n from './assets/i18n/i18n'
import axios from 'axios'
import VueAxios from 'vue-axios'

Vue.config.productionTip = false
Vue.use(ElementUI);

Vue.use(VueAxios, axios)
Object.defineProperty(Vue.prototype, '$axios', { value: axios });

new Vue({
  router,
  i18n,
  render: h => h(App)
}).$mount('#app')
