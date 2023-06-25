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
      tax: 0,
      corporations: [],
    },
    corporationData: [],
    corporationId: "",
    corporationName: "",
  },
  mutations: {
    setTableData(state, value) {
      state.tableData = value;
    },
    setForm(state, value) {
      state.form = value;
    },
    setCorporationData(state, value){
      state.corporationData = value;
    },
    setCorporationId(state, value) {
      state.corporationId = value;
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
    setCorporationData(context, value) {
      context.commit("setCorporationData", value);
    },
    setCorporationId(context, value) {
      context.commit("setCorporationId", value);
    },
    setCorporationName(context, value) {
      context.commit("setCorporationName", value);
    },
  },
});
