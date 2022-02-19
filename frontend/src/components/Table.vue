<template>
  <div class="Table">
    <el-table
      :data="tableData.filter(data => !search ||
       (data.Name.toLowerCase().includes(search.toLowerCase())))
        .slice((currentPage-1)*pageSize,currentPage*pageSize)"
      id="table"
      :cell-style="tableStyle"
      stripe
      style="width: 100%"
      :row-class-name="handelRowDetail"
    >
      <el-table-column type="expand">
        <template slot-scope="props">
          <el-table
            :data="props.row.Matertials"
            stripe
            :cell-style="tableStyle"
            style="width: 55%"
            border
            :row-class-name="handelMaterailRowDetail"
            :span-method="objectSpanMethod"
            :header-cell-style="{padding: '0'}"
          >
            <el-table-column :label="tableLabel.material.type" min-width="12%" align="center">
              <template slot-scope="scope">
                <span
                  v-if="scope.row.IsBluePrint === true"
                >{{tableLabel.material.bluePrintMaterial}}</span>
                <span v-else>{{tableLabel.material.lpStoreMaterail}}</span>
              </template>
            </el-table-column>
            <el-table-column prop="Name" :label="tableLabel.material.name" min-width="27%">
              <template slot-scope="scope">
                <el-image
                  style="height: 22px;vertical-align: middle"
                  :src="getIcon(scope.row.ItemId)"
                  fit="contain"
                  lazy
                >
                  <div slot="error" class="image-slot">
                    <i class="el-icon-picture-outline"></i>
                  </div>
                </el-image>
                <span>{{ scope.row.Name }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="Quantity" :label="tableLabel.material.quantity" min-width="10%"></el-table-column>
            <el-table-column
              prop="Price"
              :label="tableLabel.material.price"
              :formatter="stateFormat"
              min-width="15%"
            ></el-table-column>
            <el-table-column
              prop="Cost"
              :label="tableLabel.material.cost"
              min-width="15%"
              :formatter="stateFormat"
            ></el-table-column>
            <el-table-column :label="tableLabel.operation" min-width="9%">
              <template v-slot="scope">
                <el-button
                  v-if="scope.row.Error === true"
                  size="mini"
                  type="warning"
                  @click="errorMessage(scope.row.ErrorMessage)"
                >{{tableLabel.material.error}}</el-button>
              </template>
            </el-table-column>
          </el-table>
        </template>
      </el-table-column>
      <el-table-column prop="Name" min-width="20%" :label="tableLabel.name">
        <template slot-scope="scope">
          <el-image
            style="height: 22px;vertical-align: middle"
            :src="getIcon(scope.row.ItemId)"
            fit="contain"
            lazy
          >
            <div slot="error" class="image-slot">
              <i class="el-icon-picture-outline"></i>
            </div>
          </el-image>
          <span>{{ scope.row.Name }}</span>
        </template>
      </el-table-column>
      <el-table-column
        prop="Quantity"
        min-width="6%"
        :label="tableLabel.quantity"
        :formatter="stateFormat"
      ></el-table-column>
      <el-table-column
        prop="LpCost"
        min-width="7%"
        :label="tableLabel.lpCost"
        :formatter="stateFormat"
      ></el-table-column>
      <el-table-column
        prop="IskCost"
        min-width="8%"
        :label="tableLabel.iskCost"
        :formatter="stateFormat"
      ></el-table-column>
      <el-table-column
        :label="tableLabel.sumCost"
        prop="MaterialCost"
        min-width="8%"
        :formatter="stateFormat"
      ></el-table-column>
      <el-table-column
        :label="tableLabel.sumGain"
        prop="Income"
        min-width="8%"
        :formatter="stateFormat"
      ></el-table-column>
      <el-table-column
        :label="tableLabel.sumProfit"
        prop="Profit"
        min-width="8%"
        :formatter="stateFormat"
      ></el-table-column>
      <el-table-column
        :label="tableLabel.volume"
        prop="Volume"
        min-width="7%"
        :formatter="stateFormat"
      ></el-table-column>
      <el-table-column :label="tableLabel.saleIndex" prop="SalaIndex" min-width="7%"></el-table-column>
      <el-table-column :label="tableLabel.unitProfit" prop="UnitProfit" min-width="7%"></el-table-column>
      <el-table-column :label="tableLabel.operation" min-width="9%">
        <template v-slot:header>
          <el-input v-model="search" size="mini" :placeholder="tableLabel.lookUp" />
        </template>

        <template v-slot="scope">
          <el-button size="mini" type="primary" @click="orders(scope)">{{ tableLabel.orders }}</el-button>
          <el-button
            v-if="scope.row.Error === true"
            size="mini"
            type="warning"
            @click="errorMessage(scope.row.ErrorMessage)"
          >{{tableLabel.error}}</el-button>
        </template>
      </el-table-column>
    </el-table>
    <br />
    <el-pagination
      align="center"
      background
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      :current-page="currentPage"
      :page-sizes="[25,50,100,200,500]"
      :page-size="pageSize"
      layout="total, sizes, prev, pager, next, jumper"
      :total="tableData.length"
    ></el-pagination>
    <br />
  </div>
</template>

<script>
const iconServer = "https://imageserver.eveonline.com/";
export default {
  props: {
    tableData: Array,
    form: Object,
    corporationName: String,
  },
  data() {
    return {
      currentPage: 1,
      total: 20,
      pageSize: 25,
      search: "",
      tableLabel: this.$t("message.table"),
    };
  },
  methods: {
    tableStyle() {
      return {
        padding: "0",
        "font-size": "14px",
      };
    },
    handleSizeChange: function (val) {
      this.currentPage = 1;
      this.pageSize = val;
    },
    handleCurrentChange: function (val) {
      this.currentPage = val;
    },
    getIcon(typeId) {
      return iconServer + "Type/" + typeId + "_32.png";
    },
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
    handelRowDetail: function (row) {
      var value = "";
      var expand = true;
      if (row.row.Matertials === null) {
        value += "row-expand-cover";
        expand = false;
      }
      if (row.row.Error == true) {
        if (expand == true) {
          value += "warning-row";
        } else {
          value += " warning-row";
        }
      }
      return value;
    },
    handelMaterailRowDetail: function (row) {
      if (row.row.Error == true) {
        return "warning-row";
      }
    },
    objectSpanMethod({ row, column, rowIndex, columnIndex }) {
      column = rowIndex;
      rowIndex = column;
      if (columnIndex === 0) {
        if (row.length != 0) {
          return {
            rowspan: row.length,
            colspan: 1,
          };
        } else {
          return {
            rowspan: 0,
            colspan: 0,
          };
        }
      }
    },
    errorMessage(message) {
      this.$message({
        message: message,
        type: "warning",
        showClose: true,
        duration: 8000,
      });
    },
  },
  watch: {
    "$i18n.locale"() {
      this.tableLabel = this.$t("message.table");
    },
  },
};
</script>
<style>
.row-expand-cover .el-table__expand-column .cell {
  display: none;
}
.el-table .warning-row {
  background: oldlace;
}
</style>