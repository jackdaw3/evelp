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
      <el-tab-pane label="卖单">
        <el-row :gutter="50">
          <el-col :span="12">
            <OrderTable :data="sellOrders"></OrderTable>
          </el-col>
          <el-col :span="12">
            <StatisTable></StatisTable>
          </el-col>
        </el-row>
      </el-tab-pane>
      <el-tab-pane label="买单"></el-tab-pane>
      <el-tab-pane label="历史"></el-tab-pane>
    </el-tabs>
    <br />
  </div>
</template>

<script>
import Header from "@/components/Header.vue";
import OrderTable from "@/components/OrderTable.vue";
import StatisTable from "@/components/StatisTable.vue";

const backend = "http://localhost:9000/";
const the_forge = "10000002";

export default {
  name: "Order",
  components: {
    Header,
    OrderTable,
    StatisTable,
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
    this.getSellOrders();
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
        isBluePrint: false,
        corporationId: 0,
        corporationName: "",
      },
      sellOrders: [],
    };
  },
  methods: {
    getParams() {
      this.order.itemId = this.$route.query.itemId;
      this.order.offerId = this.$route.query.offerId;
      this.order.isBluePrint = this.$route.query.isBluePrint;
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
    getSellOrders() {
      this.axios
        .get(backend + "order", {
          params: {
            regionId: the_forge,
            scope: this.form.scope,
            itemId: this.order.itemId,
            isBluePrint: this.order.isBluePrint,
            isBuyOrder: false,
            lang: this.$i18n.locale,
          },
        })
        .then((response) => {
          this.sellOrders = response.data;
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
