<template>
  <div class="order">
    <Header />
    <hr />
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
    this.getParams();
    this.getItemName(this.order.itemId);
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
      },
      title: "1111",
    };
  },
  methods: {
    getParams() {
      this.order.itemId = this.$route.query.itemId;
      this.order.offerId = this.$route.query.offerId;
      this.order.corpName = this.$route.query.corporationId;
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
          console.log(response);
          document.title = response.data.ItemName;
        });
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
