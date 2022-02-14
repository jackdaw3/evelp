<template>
  <div class="header">
    <el-dropdown>
      <el-button type="primary" size="small" icon="el-icon-menu">{{ header.menu }}</el-button>
      <el-dropdown-menu slot="dropdown">
        <router-link to="/">
          <el-dropdown-item>
            <i class="el-icon-s-home"></i>
            <span>{{ header.home }}</span>
          </el-dropdown-item>
        </router-link>
      </el-dropdown-menu>
    </el-dropdown>

    <el-divider direction="vertical"></el-divider>
    <span>EVE-LP</span>
    <el-divider direction="vertical"></el-divider>
    <span>{{ header.serverName }}</span>
    <el-divider direction="vertical"></el-divider>
    <span>{{ header.market }}</span>

    <el-dropdown @command="langChange" style="float: right;margin-right: 15px;margin-top:8px;cursor: pointer">
      <span class="el-dropdown-link" >
        {{ language }}<i class="el-icon-arrow-down el-icon--right"></i>
      </span>
      <el-dropdown-menu slot="dropdown">
        <el-dropdown-item command="de">Deutsch</el-dropdown-item>
        <el-dropdown-item command="en">English</el-dropdown-item>
        <el-dropdown-item command="fr">Français</el-dropdown-item>
        <el-dropdown-item command="ja">日本語</el-dropdown-item>
        <el-dropdown-item command="ru">Pусский</el-dropdown-item>
        <el-dropdown-item command="zh">中文</el-dropdown-item>
      </el-dropdown-menu>
    </el-dropdown>
  </div>
</template>

<script>


export default {
  beforeCreate() {
    this.language=this.langLabel(this.$i18n.locale)
  },
  data() {
    return {
      header: this.$t("message.header"),
      language: ""
    };
  },
  methods: {
    langChange(lang) {
      this.$i18n.locale = lang
    },
    langLabel(lang) {
      switch (lang) {
        case "de":
          return "Deutsch";
        case "en":
          return "English";
        case "fr":
          return "Français";
        case "ja":
          return "日本語";
        case "ru":
          return "Pусский";
        case "zh":
          return "中文";
          }
    }
  },
  watch: {
    '$i18n.locale'() {
      this.header = this.$t('message.header')
      this.language=this.langLabel(this.$i18n.locale)
    }
  },
};

</script>

<style>
  A:link {
    color: black;
    text-decoration: none
  }
  .el-dropdown-link {
    cursor: pointer;
    color: #409EFF;
  }
  .el-icon-arrow-down {
    font-size: 12px;
  }
</style>