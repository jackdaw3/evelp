<template>
  <div class="order">
    <Header />
    <hr />
    <div style="margin-top: -11px">
      <el-alert
        :title="order.itemName"
        type="info"
        :closable="false"
        :description="order.corporationName"
      ></el-alert>
    </div>

    <el-tabs type="card">
      <el-tab-pane label="卖单"></el-tab-pane>
      <el-tab-pane label="买单"></el-tab-pane>
      <el-tab-pane label="历史"></el-tab-pane>
    </el-tabs>
    <br />
    <el-row :gutter="62">
      <el-col :span="12">
        <el-table :data="tableData" style="width: 100%" stripe>
          <el-table-column prop="date" label="日期" width="180"></el-table-column>
          <el-table-column prop="name" label="姓名" width="180"></el-table-column>
          <el-table-column prop="address" label="地址"></el-table-column>
        </el-table>
      </el-col>

      <el-col :span="12">
        <el-table :data="tableData" style="width: 100%;margin-right: 0px" stripe>
          <el-table-column prop="date" label="日期" width="180"></el-table-column>
          <el-table-column prop="name" label="姓名" width="180"></el-table-column>
          <el-table-column prop="address" label="地址"></el-table-column>
        </el-table>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import Header from "@/components/Header.vue";

const backend = "http://localhost:9000/";

export default {
  name: "Order",
  components: {
    Header,
  },
  created() {
    if (localStorage.lang == null) {
      localStorage.lang = "en";
    }
    if (localStorage.lang) {
      this.$i18n.locale = localStorage.lang;
    }
    this.getParams();
    this.getItemName(this.order.itemId);
    this.getCorporationName(this.order.corporationId);
  },
  computed: {
    form: function () {
      return this.$store.state.form;
    },
  },
  data() {
    return {
      order: {
        itemId: 0,
        itemName: "",
        offerId: 0,
        corporationId: 0,
        corporationName: "",
      },
      tableData: [
        {
          date: "2016-05-02",
          name: "王小虎",
          address: "上海市普陀区金沙江路 1518 弄",
        },
        {
          date: "2016-05-04",
          name: "王小虎",
          address: "上海市普陀区金沙江路 1517 弄",
        },
        {
          date: "2016-05-01",
          name: "王小虎",
          address: "上海市普陀区金沙江路 1519 弄",
        },
        {
          date: "2016-05-03",
          name: "王小虎",
          address: "上海市普陀区金沙江路 1516 弄",
        },
      ],
    };
  },
  methods: {
    getParams() {
      this.order.itemId = this.$route.query.itemId;
      this.order.offerId = this.$route.query.offerId;
      this.order.corporationId = this.$route.query.corporationId;
    },
    getItemName(itemId) {
      this.axios
        .get(backend + "item", {
          params: {
            itemId: itemId,
            lang: this.$i18n.locale,
          },
        })
        .then((response) => {
          document.title = response.data.ItemName;
          this.order.itemName = response.data.ItemName;
        });
    },
    getCorporationName(corporationId) {
      this.axios
        .get(backend + "corporation", {
          params: {
            corporationId: corporationId,
            lang: this.$i18n.locale,
          },
        })
        .then((response) => {
          this.order.corporationName = response.data.CorporationName;
        });
    },
  },
  watch: {
    "$i18n.locale"() {
      this.getItemName(this.order.itemId);
      this.getCorporationName(this.order.corporationId);
    },
  },
};
</script>

<style>
hr {
  display: block;
  height: 1px;
  border: 0;
  border-top: 1px solid #ccc;
  margin: 0.618em 0;
  padding: 0;
}
</style>
