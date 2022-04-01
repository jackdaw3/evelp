<template>
  <div class="OrderTable">
    <el-table
      :data="data"
      :cell-style="tableStyle"
      style="width: 100%"
      :header-row-style="{ color: '#B3B6B7' }"
    >
      <el-table-column
        prop="OrderId"
        :label="orderLabel.orderId"
        min-width="15%"
      ></el-table-column>
      <el-table-column
        prop="SystemName"
        :label="orderLabel.systemName"
        min-width="15%"
      ></el-table-column>
      <el-table-column
        :label="orderLabel.volume"
        :formatter="volumeFormat"
        sortable
        min-width="15%"
      ></el-table-column>
      <el-table-column
        prop="Price"
        :label="orderLabel.price"
        :formatter="stateFormat"
        sortable
        min-width="18%"
      ></el-table-column>
      <el-table-column
        prop="Expiration"
        :label="orderLabel.expiration"
        min-width="22%"
      ></el-table-column>
      <el-table-column
        prop="LastUpdated"
        :label="orderLabel.lastUpdated"
        min-width="15%"
      ></el-table-column>
    </el-table>
  </div>
</template>
<script>
export default {
  props: {
    data: Array,
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
