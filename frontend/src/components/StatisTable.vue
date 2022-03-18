<template>
  <div class="StatisTable">
    <el-table :data="data" stripe :cell-style="tableStyle" style="width: 100%" v-loading="loading">
      <el-table-column type="expand">
        <template slot-scope="props">
          <el-table
            :data="props.row.Orderwrappers"
            :cell-style="tableStyle"
            :header-cell-style="{
                background:'#F2F3F4',
                padding: '0'
              }"
          >
            <el-table-column
              prop="OrderDTO.SystemName"
              :label="orderLabel.systemName"
              min-width="15%"
            ></el-table-column>
            <el-table-column :label="orderLabel.volume" :formatter="volumeFormat" min-width="15%"></el-table-column>
            <el-table-column
              prop="OrderDTO.Price"
              :label="orderLabel.price"
              :formatter="stateFormat"
              sortable
              min-width="20%"
            ></el-table-column>
            <el-table-column
              prop="Income"
              :label="statisLabel.income"
              :formatter="stateFormat"
              sortable
              min-width="20%"
            ></el-table-column>
            <el-table-column
              prop="Cost"
              :label="statisLabel.cost"
              :formatter="stateFormat"
              sortable
              min-width="20%"
            ></el-table-column>
            <el-table-column
              prop="Profit"
              :label="statisLabel.profit"
              :formatter="stateFormat"
              sortable
              min-width="20%"
            ></el-table-column>
            <el-table-column
              prop="UnitProfit"
              :label="statisLabel.unitProfit"
              sortable
              min-width="15%"
            ></el-table-column>
          </el-table>
        </template>
      </el-table-column>
      <el-table-column prop="UnitProfitRange" :label="statisLabel.lpRange"></el-table-column>
      <el-table-column prop="Quantity" :label="statisLabel.number"></el-table-column>
      <el-table-column prop="Income" :label="statisLabel.income" :formatter="stateFormat" sortable></el-table-column>
      <el-table-column prop="Cost" :label="statisLabel.cost" :formatter="stateFormat" sortable></el-table-column>
      <el-table-column prop="Profit" :label="statisLabel.profit" :formatter="stateFormat" sortable></el-table-column>
      <el-table-column prop="AveUnitProfit" :label="statisLabel.aveLpPrice" sortable></el-table-column>
    </el-table>
  </div>
</template>
<script>
export default {
  props: {
    data: Array,
    loading: Boolean,
  },
  data() {
    return {
      statisLabel: this.$t("message.statis"),
      orderLabel: this.$t("message.order"),
    };
  },
  methods: {
    stateFormat(row, column, cellValue) {
      cellValue = Math.round(cellValue);
      cellValue += "";
      if (!cellValue.includes(".")) cellValue += ".";
      return cellValue
        .replace(/(\d)(?=(\d{3})+\.)/g, function ($0, $1) {
          return $1 + ",";
        })
        .replace(/\.$/, "");
    },
    tableStyle() {
      return {
        padding: "0",
        "font-size": "14px",
      };
    },
    volumeFormat(row) {
      return row.OrderDTO.VolumeRemain + "/" + row.OrderDTO.VolumeTotal;
    },
  },
  watch: {
    "$i18n.locale"() {
      this.statisLabel = this.$t("message.statis");
      this.orderLabel = this.$t("message.order");
    },
  },
};
</script>
