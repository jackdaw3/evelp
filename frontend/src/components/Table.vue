<template>
  <div class="Table">
    <el-table
      ref="tableList"
      :data="
        tableData
          .filter(
            (data) =>
              !search || data.ItemName.toLowerCase().includes(search.toLowerCase())
          )
          .slice((currentPage - 1) * pageSize, currentPage * pageSize)
      "
      id="table"
      :cell-style="tableStyle"
      :header-row-style="{ color: '#B3B6B7' }"
      @sort-change="sort_change"
      style="width: 100%"
      :row-class-name="handelRowDetail"
    >
      <el-table-column type="expand">
        <template slot-scope="props">
          <el-table
            :data="props.row.Matertials"
            :cell-style="tableStyle"
            style="width: 55%"
            :row-class-name="handelMaterailRowDetail"
            :span-method="objectSpanMethod"
            :header-row-style="{ color: '#B3B6B7' }"
            :header-cell-style="{ padding: '0' }"
          >
            <el-table-column
              :label="tableLabel.material.type"
              min-width="12%"
              align="center"
            >
              <template slot-scope="scope">
                <span v-if="scope.row.IsBluePrint === true">
                  {{ tableLabel.material.bluePrintMaterial }}
                </span>
                <span v-else>{{ tableLabel.material.lpStoreMaterail }}</span>
              </template>
            </el-table-column>
            <el-table-column
              prop="MaterialName"
              :label="tableLabel.material.name"
              min-width="27%"
            >
              <template slot-scope="scope">
                <el-image
                  style="height: 22px; vertical-align: middle"
                  :src="getIcon(scope.row.ItemId)"
                  fit="contain"
                  lazy
                >
                  <div slot="error" class="image-slot">
                    <i class="el-icon-picture-outline"></i>
                  </div>
                </el-image>
                <span>{{ scope.row.MaterialName }}</span>
              </template>
            </el-table-column>
            <el-table-column
              prop="Quantity"
              :label="tableLabel.material.quantity"
              min-width="10%"
            ></el-table-column>
            <el-table-column
              prop="Price"
              :label="tableLabel.material.price"
              :formatter="stateFormat"
              min-width="14%"
            ></el-table-column>
            <el-table-column
              prop="Cost"
              :label="tableLabel.material.cost"
              min-width="14%"
              :formatter="stateFormat"
            ></el-table-column>
            <el-table-column :label="tableLabel.operation" min-width="15%">
              <template v-slot:header>
                <el-button
                  size="mini"
                  plain
                  @click="copyAllMaterials(props.row.Matertials)"
                  >{{ tableLabel.material.copy }}</el-button
                >
              </template>
              <template v-slot="scope">
                <el-button size="mini" plain @click="copyMaterial(scope.row)">{{
                  tableLabel.material.copy
                }}</el-button>
                <el-button
                  v-if="scope.row.Error === true"
                  size="mini"
                  style="margin-left: 0px"
                  plain
                  @click="errorMessage(scope.row.ErrorMessage)"
                  >{{ tableLabel.material.error }}</el-button
                >
              </template>
            </el-table-column>
          </el-table>
        </template>
      </el-table-column>
      <el-table-column prop="ItemName" min-width="20%" :label="tableLabel.name">
        <template slot-scope="scope">
          <el-image
            style="height: 22px; vertical-align: middle"
            :src="getIcon(scope.row.ItemId)"
            fit="contain"
            lazy
          >
            <div slot="error" class="image-slot">
              <i class="el-icon-picture-outline"></i>
            </div>
          </el-image>
          <span>{{ scope.row.ItemName }}</span>
        </template>
      </el-table-column>
      <el-table-column
        prop="Quantity"
        min-width="5%"
        :label="tableLabel.quantity"
        :formatter="stateFormat"
      ></el-table-column>
      <el-table-column
        prop="LpCost"
        min-width="7%"
        :label="tableLabel.lpCost"
        :formatter="stateFormat"
        sortable="custom"
      ></el-table-column>
      <el-table-column
        prop="IskCost"
        min-width="8%"
        :label="tableLabel.iskCost"
        :formatter="stateFormat"
        sortable="custom"
      ></el-table-column>
      <el-table-column
        :label="tableLabel.cost"
        prop="MaterialCost"
        min-width="8%"
        :formatter="stateFormat"
        sortable="custom"
      ></el-table-column>
      <el-table-column
        :label="tableLabel.price"
        prop="Price"
        min-width="8%"
        :formatter="stateFormat"
        sortable="custom"
      ></el-table-column>
      <el-table-column
        :label="this.tableLabel.income"
        prop="Income"
        min-width="8%"
        :formatter="stateFormat"
        sortable="custom"
      ></el-table-column>
      <el-table-column
        :label="tableLabel.profit"
        prop="Profit"
        min-width="8%"
        :formatter="stateFormat"
        sortable="custom"
      ></el-table-column>
      <el-table-column
        :label="tableLabel.volume"
        prop="Volume"
        min-width="6%"
        :formatter="stateFormat"
        sortable="custom"
      ></el-table-column>
      <el-table-column
        :label="tableLabel.saleIndex"
        prop="SaleIndex"
        min-width="6%"
        sortable="custom"
      ></el-table-column>
      <el-table-column
        :label="tableLabel.unitProfit"
        prop="UnitProfit"
        min-width="6%"
        sortable="custom"
      ></el-table-column>
      <el-table-column :label="tableLabel.operation" min-width="10%">
        <template v-slot:header>
          <el-input
            v-model="search"
            size="mini"
            :placeholder="tableLabel.lookUp"
          />
        </template>
        <template v-slot="scope">
          <el-button size="mini" plain @click="orders(scope)">
            {{ tableLabel.orders }}
          </el-button>
          <el-button
            v-if="scope.row.Error === true"
            size="mini"
            plain
            style="margin-left: 0px"
            @click="errorMessage(scope.row.ErrorMessage)"
            >{{ tableLabel.error }}</el-button
          >
        </template>
      </el-table-column>
    </el-table>
    <br />
    <el-pagination
      align="center"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      :current-page="currentPage"
      :page-sizes="[25, 50, 100, 200, 500]"
      :page-size="pageSize"
      layout="total, sizes, prev, pager, next, jumper"
      :total="tableData.length"
    ></el-pagination>
    <br />
  </div>
</template>

<script>
import FileSaver from "file-saver";
import XLSX from "xlsx";

const iconServer = "https://imageserver.eveonline.com/";
export default {
  data() {
    return {
      currentPage: 1,
      total: 20,
      pageSize: 25,
      search: "",
      tableLabel: this.$t("message.table"),
    };
  },
  computed: {
    tableData: {
      get() {
        return this.$store.state.tableData;
      },
      set() {},
    },
    form: function () {
      return this.$store.state.form;
    },
    corporationId: function () {
      return this.$store.state.corporationId;
    },
    corporationName: function () {
      return this.$store.state.corporationName;
    },
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
        style: "backgound-color:red",
        showClose: true,
        duration: 5000,
      });
    },
    copyMaterial(row) {
      var value = row.Quantity + " " + row.MaterialName + "\n";
      this.$copyText(value).then(
        () => {
          this.$message({
            message: this.tableLabel.material.copySuccess,
            type: "success",
            duration: 2000,
          });
        },
        function (e) {
          this.$message({
            message: this.tableLabel.material.copyFail,
            type: "error",
            duration: 2000,
          });
          console.log(e);
        }
      );
    },
    copyAllMaterials(list) {
      let value = "";
      for (let i = 0; i < list.length; ++i) {
        value += list[i].Quantity + " " + list[i].MaterialName + "\n";
      }

      this.$copyText(value).then(
        () => {
          this.$message({
            message: this.tableLabel.material.copySuccess,
            type: "success",
            duration: 2000,
          });
        },
        function (e) {
          this.$message({
            message: this.tableLabel.material.copyFail,
            type: "error",
            duration: 2000,
          });
          console.log(e);
        }
      );
    },
    exportExcel() {
      if (this.tableData.length === 0) {
        this.$message({
          message: this.$t("message.noData"),
          type: "warning",
          duration: 1500,
        });
        return;
      }
      const pages = this.pageSize;
      this.pageSize = 500;
      this.currentPage = 1;
      let date = new Date();
      let name =
        this.corporationName +
        "-" +
        this.dateFormat("YYYY_mm_dd-HH_MM_SS", date) +
        "-" +
        this.form.tax +
        "%";
      this.$nextTick(function () {
        let wb = XLSX.utils.table_to_book(document.getElementById("table"));
        let wbout = XLSX.write(wb, {
          bookType: "xlsx",
          bookSST: true,
          type: "array",
        });
        try {
          FileSaver.saveAs(
            new Blob([wbout], { type: "application/octet-stream" }),
            name + ".xlsx"
          );
        } catch (e) {
          if (typeof console !== "undefined") console.log(e, wbout);
        }
        this.pageSize = pages;
        return wbout;
      });
    },
    dateFormat(fmt, date) {
      let ret;
      const opt = {
        "Y+": date.getFullYear().toString(),
        "m+": (date.getMonth() + 1).toString(),
        "d+": date.getDate().toString(),
        "H+": date.getHours().toString(),
        "M+": date.getMinutes().toString(),
        "S+": date.getSeconds().toString(),
      };
      for (let k in opt) {
        ret = new RegExp("(" + k + ")").exec(fmt);
        if (ret) {
          fmt = fmt.replace(
            ret[1],
            ret[1].length == 1 ? opt[k] : opt[k].padStart(ret[1].length, "0")
          );
        }
      }
      return fmt;
    },
    sort_change(column) {
      this.current_page = 1;
      if (column.prop === "SaleIndex") {
        if (column.order === "descending") {
          this.tableData = this.tableData.sort(this.saleIndexDescSort);
        } else if (column.order === "ascending") {
          this.tableData = this.tableData.sort(this.saleIndexAscSort);
        }
      } else if (column.prop === "UnitProfit") {
        if (column.order === "descending") {
          this.tableData = this.tableData.sort(this.unitProfitDescSort);
        } else if (column.order === "ascending") {
          this.tableData = this.tableData.sort(this.unitProfitAscSort);
        }
      } else if (column.prop === "Profit") {
        if (column.order === "descending") {
          this.tableData = this.tableData.sort(this.profitDescSort);
        } else if (column.order === "ascending") {
          this.tableData = this.tableData.sort(this.profitAscSort);
        }
      } else if (column.prop === "Volume") {
        if (column.order === "descending") {
          this.tableData = this.tableData.sort(this.volumeDescSort);
        } else if (column.order === "ascending") {
          this.tableData = this.tableData.sort(this.volumeAscSort);
        }
      } else if (column.prop === "LpCost") {
        if (column.order === "descending") {
          this.tableData = this.tableData.sort(this.lpCostDescSort);
        } else if (column.order === "ascending") {
          this.tableData = this.tableData.sort(this.lpCostAscSort);
        }
      } else if (column.prop === "IskCost") {
        if (column.order === "descending") {
          this.tableData = this.tableData.sort(this.iskCostDescSort);
        } else if (column.order === "ascending") {
          this.tableData = this.tableData.sort(this.iskCostAscSort);
        }
      } else if (column.prop === "MaterialCost") {
        if (column.order === "descending") {
          this.tableData = this.tableData.sort(this.materialCostDescSort);
        } else if (column.order === "ascending") {
          this.tableData = this.tableData.sort(this.materialCostAscSort);
        }
      } else if (column.prop === "Income") {
        if (column.order === "descending") {
          this.tableData = this.tableData.sort(this.incomeDescSort);
        } else if (column.order === "ascending") {
          this.tableData = this.tableData.sort(this.incomeAscSort);
        }
      } else if (column.prop === "Price") {
        if (column.order === "descending") {
          this.tableData = this.tableData.sort(this.priceDescSort);
        } else if (column.order === "ascending") {
          this.tableData = this.tableData.sort(this.priceAscSort);
        }
      }
      this.showed_data = this.tableData.slice(0, this.pageSize);
    },
    saleIndexDescSort(a, b) {
      if (a.SaleIndex > b.SaleIndex) {
        return -1;
      } else if (a.SaleIndex < b.SaleIndex) {
        return 1;
      } else {
        return 0;
      }
    },
    saleIndexAscSort(a, b) {
      if (a.SaleIndex < b.SaleIndex) {
        return -1;
      } else if (a.SaleIndex > b.SaleIndex) {
        return 1;
      } else {
        return 0;
      }
    },
    profitDescSort(a, b) {
      if (a.Profit > b.Profit) {
        return -1;
      } else if (a.Profit < b.Profit) {
        return 1;
      } else {
        return 0;
      }
    },
    profitAscSort(a, b) {
      if (a.Profit < b.Profit) {
        return -1;
      } else if (a.Profit > b.Profit) {
        return 1;
      } else {
        return 0;
      }
    },
    volumeDescSort(a, b) {
      if (a.Volume > b.Volume) {
        return -1;
      } else if (a.Volume < b.Volume) {
        return 1;
      } else {
        return 0;
      }
    },
    volumeAscSort(a, b) {
      if (a.Volume < b.Volume) {
        return -1;
      } else if (a.Volume > b.Volume) {
        return 1;
      } else {
        return 0;
      }
    },
    unitProfitDescSort(a, b) {
      if (a.UnitProfit > b.UnitProfit) {
        return -1;
      } else if (a.UnitProfit < b.UnitProfit) {
        return 1;
      } else {
        return 0;
      }
    },
    unitProfitAscSort(a, b) {
      if (a.UnitProfit < b.UnitProfit) {
        return -1;
      } else if (a.UnitProfit > b.UnitProfit) {
        return 1;
      } else {
        return 0;
      }
    },
    lpCostDescSort(a, b) {
      if (a.LpCost > b.LpCost) {
        return -1;
      } else if (a.LpCost < b.LpCost) {
        return 1;
      } else {
        return 0;
      }
    },
    lpCostAscSort(a, b) {
      if (a.LpCost < b.LpCost) {
        return -1;
      } else if (a.LpCost > b.LpCost) {
        return 1;
      } else {
        return 0;
      }
    },
    iskCostDescSort(a, b) {
      if (a.IskCost > b.IskCost) {
        return -1;
      } else if (a.IskCost < b.IskCost) {
        return 1;
      } else {
        return 0;
      }
    },
    iskCostAscSort(a, b) {
      if (a.IskCost < b.IskCost) {
        return -1;
      } else if (a.IskCost > b.IskCost) {
        return 1;
      } else {
        return 0;
      }
    },
    materialCostDescSort(a, b) {
      if (a.MaterialCost > b.MaterialCost) {
        return -1;
      } else if (a.MaterialCost < b.MaterialCost) {
        return 1;
      } else {
        return 0;
      }
    },
    materialCostAscSort(a, b) {
      if (a.MaterialCost < b.MaterialCost) {
        return -1;
      } else if (a.MaterialCost > b.MaterialCost) {
        return 1;
      } else {
        return 0;
      }
    },
    incomeDescSort(a, b) {
      if (a.Income > b.Income) {
        return -1;
      } else if (a.Income < b.Income) {
        return 1;
      } else {
        return 0;
      }
    },
    incomeAscSort(a, b) {
      if (a.Income < b.Income) {
        return -1;
      } else if (a.Income > b.Income) {
        return 1;
      } else {
        return 0;
      }
    },
    priceDescSort(a, b) {
      if (a.Price > b.Price) {
        return -1;
      } else if (a.Price < b.Price) {
        return 1;
      } else {
        return 0;
      }
    },
    priceAscSort(a, b) {
      if (a.Price < b.Price) {
        return -1;
      } else if (a.Price > b.Price) {
        return 1;
      } else {
        return 0;
      }
    },
    orders(scope) {
      let routeUrl = this.$router.resolve({
        name: "Order",
        query: {
          itemId: scope.row.ItemId,
          offerId: scope.row.OfferId,
          isBluePrint: scope.row.IsBluePrint,
          corporationId: this.corporationId,
          materialPrice: this.form.materialPrice,
          scope: this.form.scope,
          tax: this.form.tax,
        },
      });
      window.open(routeUrl.href, "_blank");
    },
  },
  watch: {
    "$i18n.locale"() {
      this.tableLabel = this.$t("message.table");
    },
    tableData() {
      this.currentPage = 1;
      if (this.tableData.length == 0) {
        this.$refs.tableList.clearSort();
      }
    },
  },
};
</script>
<style>
.row-expand-cover .el-table__expand-column .cell {
  display: none;
}
.el-table .warning-row {
  background: #251b07;
}
</style>
