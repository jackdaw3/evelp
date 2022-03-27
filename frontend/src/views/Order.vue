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
      <el-tab-pane :label="orderLabel.buyOrder">
        <el-row :gutter="35">
          <el-col :span="12">
            <OrderTable :data="buyOrders"></OrderTable>
          </el-col>
          <el-col :span="12">
            <StatisTable :data="buyStatis"></StatisTable>
          </el-col>
        </el-row>
      </el-tab-pane>
      <el-tab-pane :label="orderLabel.sellOrder">
        <el-row :gutter="35">
          <el-col :span="12">
            <OrderTable :data="sellOrders"></OrderTable>
          </el-col>
          <el-col :span="12">
            <StatisTable :data="sellStatis"></StatisTable>
          </el-col>
        </el-row>
      </el-tab-pane>
      <el-tab-pane :label="orderLabel.history" :lazy="stockLazy">
        <Stock :style="{ height: stockHeight }" :history="history"></Stock>
      </el-tab-pane>
    </el-tabs>
    <br />
  </div>
</template>

<script>
import Header from "@/components/Header.vue";
import OrderTable from "@/components/OrderTable.vue";
import StatisTable from "@/components/StatisTable.vue";
import Stock from "@/components/Stock.vue";

const backend = "https://eve-lp.com/api/";
const the_forge = "10000002";

export default {
  name: "Order",
  components: {
    Header,
    OrderTable,
    StatisTable,
    Stock,
  },
  created() {
    if (localStorage.lang == null) {
      localStorage.lang = "en";
    }
    if (localStorage.lang) {
      this.$i18n.locale = localStorage.lang;
    }
    this.initial();
  },
  computed: {
    stockHeight() {
      return window.screen.height * 0.68 + "px";
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
        materialPrice: "sell",
        scope: "0.05",
        tax: 0,
      },
      history: {
        average: [],
        average5d: [],
        average20d: [],
        minAndmax: [],
        minAndmax5d: [],
        volume: [],
        borderWidth: 0,
        title: "",
        label: this.$t("message.stock"),
      },
      stockLazy: true,
      sellOrders: [],
      sellStatis: [],
      buyOrders: [],
      buyStatis: [],
      orderLabel: this.$t("message.order"),
    };
  },
  methods: {
    getParams() {
      this.order.itemId = this.$route.query.itemId;
      this.order.offerId = this.$route.query.offerId;
      this.order.isBluePrint = this.$route.query.isBluePrint;
      this.order.corporationId = this.$route.query.corporationId;
      this.order.materialPrice = this.$route.query.materialPrice;
      this.order.scope = this.$route.query.scope;
      this.order.tax = this.$route.query.tax;
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
          this.history.title = response.data.ItemName;
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
    initial() {
      this.getParams();
      this.getItemName(this.order.itemId);
      this.getCorporationName(this.order.corporationId);

      this.getOrders(true);
      this.getStatis(true);
      this.getOrders(false);
      this.getStatis(false);
      this.getHistory();
    },
    getOrders(isBuyOrder) {
      this.axios
        .get(backend + "order", {
          params: {
            regionId: the_forge,
            scope: this.order.scope,
            tax: this.order.tax,
            itemId: this.order.itemId,
            isBluePrint: this.order.isBluePrint,
            isBuyOrder: isBuyOrder,
            lang: this.$i18n.locale,
          },
        })
        .then((response) => {
          if (!isBuyOrder) {
            this.sellOrders = response.data;
          } else {
            this.buyOrders = response.data;
          }
        });
    },
    getStatis(isBuyOrder) {
      this.axios
        .get(backend + "statis", {
          params: {
            offerId: this.order.offerId,
            regionId: the_forge,
            scope: this.order.scope,
            materialPrice: this.order.materialPrice,
            tax: this.order.tax,
            isBuyOrder: isBuyOrder,
            lang: this.$i18n.locale,
          },
        })
        .then((response) => {
          if (!isBuyOrder) {
            this.sellStatis = response.data;
          } else {
            this.buyStatis = response.data;
          }
        });
    },
    getHistory() {
      this.axios
        .get(backend + "history", {
          params: {
            regionId: the_forge,
            itemId: this.order.itemId,
            isBluePrint: this.order.isBluePrint,
          },
        })
        .then((response) => {
          for (let i = 0; i < response.data.length; i++) {
            this.history.average.push([
              Date.parse(new Date(response.data[i].Date)),
              response.data[i].Average,
            ]);
            this.history.volume.push([
              Date.parse(new Date(response.data[i].Date)),
              response.data[i].Volume,
            ]);
            this.history.minAndmax.push([
              Date.parse(new Date(response.data[i].Date)),
              response.data[i].Lowest,
              response.data[i].Highest,
            ]);
            this.history.average5d.push([
              Date.parse(new Date(response.data[i].Date)),
              response.data[i].Average5d,
            ]);
            this.history.average20d.push([
              Date.parse(new Date(response.data[i].Date)),
              response.data[i].Average20d,
            ]);
            this.history.minAndmax5d.push([
              Date.parse(new Date(response.data[i].Date)),
              response.data[i].Lowest5d,
              response.data[i].Highest5d,
            ]);
          }
          this.history.borderWidth =
            document.documentElement.clientWidth * 0.0388;
        });
    },
  },
  watch: {
    "$i18n.locale"() {
      this.orderLabel = this.$t("message.order");
      this.history.label = this.$t("message.stock");
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
