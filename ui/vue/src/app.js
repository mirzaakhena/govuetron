import Vue from 'vue'
import App from './App'

import router from './router'
import BootstrapVue from 'bootstrap-vue'
import axios from 'axios'
import moment from 'moment'
import numeral from 'numeral'
import Icon from 'vue-awesome/components/Icon'
import VueSweetAlert from 'vue-sweetalert2'
import money from 'v-money'
import VueNativeSock from 'vue-native-websocket'

import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import 'vue-awesome/icons/remove'
import 'vue-awesome/icons/pencil'

Vue.use(BootstrapVue)
Vue.use(VueSweetAlert)
Vue.use(money, {
  decimal: ',',
  thousands: '.',
  prefix: 'Rp ',
  suffix: ' ',
  precision: 0,
  masked: true
})
Vue.use(VueNativeSock, 'ws://localhost:8081/ws', { format: 'json' })

Vue.component('icon', Icon)

Vue.config.productionTip = false

Vue.prototype.$http = axios.create({
  baseURL: 'http://localhost:8081'
})

Vue.filter('digitgroup', function (value) {
  if (value === '') {
    return ''
  }
  return numeral(value).format('0,0').replace(',', '.')
})

Vue.filter('dateformat', value => {
  if (value) {
    return moment(String(value)).format('MM/DD/YYYY hh:mm')
  }
})

window.onload = function () {
  var topic = document.createElement("div");
  topic.id = 'vue-app';
  document.body.insertBefore(topic, document.body.childNodes[0])  

  new Vue({
    el: '#vue-app',
    router,
    components: { App },
    template: '<App/>'
  })
}

