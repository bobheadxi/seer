import Vue from 'vue';
import { __values } from 'tslib';
import 'ant-design-vue/dist/antd.css';

import App from './App.vue';
import store from './store';
import router from './router';

Vue.config.productionTip = false;

new Vue({
  router,
  store,
  render: h => h(App),
  errorCaptured: (err, vm, info) => {
    console.error({ err, vm, info });
    return false;
  },
}).$mount('#app');
