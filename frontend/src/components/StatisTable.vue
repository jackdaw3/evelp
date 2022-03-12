<template>
  <div class="StatisTable">
    <el-table :data="data" stripe  :cell-style="tableStyle" style="width: 100%">
      <el-table-column prop="UnitProfitRange" :label="statisLabel.lpRange"></el-table-column>
      <el-table-column prop="Quantity" :label="statisLabel.number"></el-table-column>
      <el-table-column prop="Cost" :label="statisLabel.sumCost" sortable></el-table-column>
      <el-table-column prop="Income" :label="statisLabel.sumIncome" sortable></el-table-column>
      <el-table-column prop="Profit" :label="statisLabel.sumProfit" sortable></el-table-column>
      <el-table-column prop="AveUnitProfit" :label="statisLabel.aveLpPrice" sortable></el-table-column>
    </el-table>
  </div>
</template>
<script>
export default {
  data() {
    return {
      statisLabel: this.$t("message.statis"),
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
    tableStyle() {
      if (localStorage.tableStyle == "default") {
        return { "font-size": localStorage.tableFontSize + "px" };
      }
      if (localStorage.tableStyle == "flat") {
        return {
          padding: "0",
          "font-size": localStorage.tableFontSize + "px",
        };
      }
      return "";
    },
  },
  watch: {
    "$i18n.locale"() {
      this.statisLabel = this.$t("message.statis");
    },
  },
};
</script>
