<template>
  <div class="OrderTable">
    <el-table :data="data" stripe :cell-style="tableStyle" style="width: 100%" border>
      <el-table-column prop="OrderId" :label="orderLabel.orderId" min-width="12%"></el-table-column>
      <el-table-column prop="SystemName" :label="orderLabel.systemName" min-width="10%"></el-table-column>
      <el-table-column prop="Issued" :label="orderLabel.issued" min-width="20%"></el-table-column>
      <el-table-column prop="Duration" :label="orderLabel.duration" min-width="8%"></el-table-column>
      <el-table-column
        prop="VolumeRemain"
        :label="orderLabel.volumeRemain"
        :formatter="volumeFormat"
        sortable
        min-width="12%"
      ></el-table-column>
      <el-table-column
        prop="Price"
        :label="orderLabel.price"
        :formatter="stateFormat"
        sortable
        min-width="18%"
      ></el-table-column>
      <el-table-column prop="LastUpdated" :label="orderLabel.lastUpdated" min-width="20%"></el-table-column>
    </el-table>
  </div>
</template>
<script>
export default {
  props: {
    data: Array,
  },
  computed: {
    form: function () {
      return this.$store.state.form;
    },
  },
  data() {
    return {
      orderLabel: this.$t("message.order"),
    };
  },
  methods: {
    stateFormat(row, column, cellValue) {
      cellValue += "";
      if (!cellValue.includes(".")) cellValue += ".";
      return cellValue
        .replace(/(\d)(?=(\d{3})+\.)/g, function ($0, $1) {
          return $1 + ",";
        })
        .replace(/\.$/, "");
    },
    volumeFormat(row) {
      return row.VolumeRemain + "/" + row.VolumeTotal;
    },
    tableStyle() {
      return {
        padding: "0",
        "font-size": "14px",
      };
    },
  },
  watch: {
    "$i18n.locale"() {
      this.orderLabel = this.$t("message.order");
    },
  },
};
</script>
