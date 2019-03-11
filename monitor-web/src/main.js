// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'

import VCharts from 'v-charts'


Vue.use(VCharts)

new Vue({
  el: '#app',
  data: function () {
    return {
      chartData: {
        columns: ['time', 'tcp 現在連線數', 'external api 處理中數量', 'external api 等待處理數'],
        rows: [
        ]
      },
      status: {
        now_connect_number: 0,
        total_connect_number: 0,
        finish_external_api_request: 0,
        external_api_request_ing: 0,
        wait_external_api_request: 0,
        all_external_api_request: 0,
      }
    }
  },
  created: function(){
    // this.updateData()
    this.timer()
  },
  methods: {
    timer: function() {
      setInterval(() => { 
        this.updateData();
      }, 1000)
    },
    updateData: function() {
      let self = this;
      axios({
        methods: 'get',
        url: 'http://localhost:8081/monitor'
      })
      .then((resp) => {
        console.log('resp=>',resp)
        // self.teachers = resp.data;
        if (self.chartData.rows.length > 15) {
          self.chartData.rows.shift()
        }
        var now = new Date();
        self.chartData.rows.push({
          'time': now.toLocaleTimeString(),
          'tcp 現在連線數': resp.data.now_connect_number,
          'external api 處理中數量': resp.data.external_api_request_ing,
          'external api 等待處理數': resp.data.wait_external_api_request,
        })
        self.status.now_connect_number = resp.data.now_connect_number
        self.status.total_connect_number = resp.data.total_connect_number
        self.status.finish_external_api_request = resp.data.finish_external_api_request
        self.status.external_api_request_ing = resp.data.external_api_request_ing
        self.status.wait_external_api_request = resp.data.wait_external_api_request
        self.status.all_external_api_request = resp.data.all_external_api_request
      });
    }
  }
})