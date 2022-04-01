<template>
  <div class="StatisTable">
    <el-table
      :data="data"
      :cell-style="tableStyle"
      style="width: 100%"
      :header-row-style="{ color: '#B3B6B7' }"
    >
      <el-table-column type="expand">
        <template slot-scope="props">
          <el-table
            :data="props.row.Orderwrappers"
            :cell-style="tableStyle"
            :header-cell-style="{
              padding: '0',
            }"
            :header-row-style="{ color: '#B3B6B7' }"
          >
            <el-table-column
              prop="OrderDTO.SystemName"
              :label="orderLabel.systemName"
              min-width="15%"
            ></el-table-column>
            <el-table-column
              :label="orderLabel.volume"
              :formatter="volumeFormat"
              min-width="15%"
            ></el-table-column>
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
      <el-table-column
        prop="UnitProfitRange"
        :label="statisLabel.lpRange"
        min-width="25%"
      ></el-table-column>
      <el-table-column
        prop="Quantity"
        :label="statisLabel.number"
        min-width="11%"
      ></el-table-column>
      <el-table-column
        prop="Income"
        :label="statisLabel.income"
        :formatter="stateFormat"
        sortable
        min-width="16%"
      ></el-table-column>
      <el-table-column
        prop="Cost"
        :label="statisLabel.cost"
        :formatter="stateFormat"
        sortable
        min-width="16%"
      ></el-table-column>
      <el-table-column
        prop="Profit"
        :label="statisLabel.profit"
        :formatter="stateFormat"
        sortable
        min-width="16%"
      ></el-table-column>
      <el-table-column
        prop="AveUnitProfit"
        :label="statisLabel.aveLpPrice"
        sortable
        min-width="16%"
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
      statisLabel: this.$t("message.statis"),
      orderLabel: this.$t("message.order"),
    };
  },
  methods: {
    stateFormat(row, column, cellValue) {
      cellValue = Math.floor(cellValue);
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
