import Vue from "vue";
import Vuex from "vuex";

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    tableData: [],
    form: {
      materialPrice: "sell",
      productPrice: "buy",
      days: "7",
      scope: "0.05",
    },
    corporationName: "",
  },
  mutations: {
    setTableData(state, value) {
      state.tableData = value;
    },
    setForm(state, value) {
      state.form = value;
    },
    setCorporationName(state, value) {
      state.corporationName = value;
    },
  },
  actions: {
    setTableData(context, value) {
      context.commit("setTableData", value);
    },
    setForm(context, value) {
      context.commit("setForm", value);
    },
    setCorporationName(context, value) {
      context.commit("setCorporationName", value);
    },
  },
});