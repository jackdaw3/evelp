<template>
  <div class="order" v-title data-title="title">
    <Header />
    <hr />
  </div>
</template>

<script>
import Header from "@/components/Header.vue";
export default {
  name: "Order",
  components: {
    Header,
  },
  computed: {
    form: function () {
      return this.$store.state.form;
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
