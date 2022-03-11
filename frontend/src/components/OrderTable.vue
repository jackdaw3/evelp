<template>
  <div class="OrderTable">
    <el-table
        :data="data"
        stripe
        v-loading="loading"
        :cell-style="tableStyle"
        style="width: 100%">
      <el-table-column
          prop="orderId"
          :label="ordersTableLabel.orderId">
      </el-table-column>
      <el-table-column
          v-if="$i18n.locale === 'cn'"
          prop="systemNameZh"
          :label="ordersTableLabel.systemName">
      </el-table-column>
      <el-table-column
          v-else
          prop="systemName"
          :label="ordersTableLabel.systemName">
      </el-table-column>
      <el-table-column
          prop="issued"
          :label="ordersTableLabel.issued"
          sortable>
      </el-table-column>
      <el-table-column
          prop="volumeRemain"
          :label="ordersTableLabel.volumeRemain"
          sortable>
      </el-table-column>
      <el-table-column
          prop="duration"
          :label="ordersTableLabel.duration"
          sortable>
      </el-table-column>
      <el-table-column
          prop="price"
          :label="ordersTableLabel.price"
          :formatter="stateFormat"
          sortable>
      </el-table-column>

    </el-table>
  </div>
</template>
<script>
export default {
  props: {
    data: Array,  //所传递数据的类型
    loading: Boolean
  },
  methods: {
    stateFormat(row, column, cellValue) {
      cellValue += '';
      if (!cellValue.includes('.')) cellValue += '.';
      return cellValue.replace(/(\d)(?=(\d{3})+\.)/g, function ($0, $1) {
        return $1 + ',';
      }).replace(/\.$/, '');
    },
    tableStyle(){
      if(localStorage.tableStyle=="default"){
        return {'font-size': localStorage.tableFontSize+'px'}
      }
      if(localStorage.tableStyle=="flat"){
        return {
          padding: '0',
          'font-size': localStorage.tableFontSize+'px'
        }
      }
      return ''
    }
  },
  data() {
    return {
      ordersTableLabel: this.$t('message.ordersTable')
    }
  },
  watch: {
    '$i18n.locale'() {
      this.ordersTableLabel = this.$t('message.ordersTable')
    }
  }
}
</script>
