<template>
    <div class="ItemDialog">
        <el-dialog :visible.sync="itemDialogVisible" v-if="itemDialogVisible" class="custom-dialog">
            <template slot="title">
                <div style="display: flex; align-items: center; max-height: 50%;">
                    <el-image style="margin-right: 8px; height: 52px;" :src="getIcon()" v-if="itemDialogVisible"></el-image>
                    <div>
                        <div style="color: #CDD0D1; font-size: 18px;">{{ itemDialogData.ItemName }}</div>
                        <br class="short-br">
                        <div style="color: #CDD0D1; font-size: 18px;">{{ itemDialogData.Volume }}</div>
                    </div>
                </div>
            </template>
            <el-divider></el-divider>
            <span v-html="removeHtmlTags(itemDialogData.Description)"></span>
        </el-dialog>
    </div>
</template>
<script>
import { ICON_SERVER } from '@/constants';
export default {
    computed: {
        itemDialogVisible: {
            get() {
                return this.$store.state.itemDialogVisible;
            },
            set(value) {
                this.$store.commit('setItemDialogVisible', value);
            }
        },
        itemDialogData: {
            get() {
                return this.$store.state.itemDialogData;
            },
            set(value) {
                this.$store.commit('setItemDialogData', value);
            }
        }
    },
    data() {
        return {
            
        };
    },
    methods: {
        getIcon() {
            return ICON_SERVER + "Type/" + this.itemDialogData.ItemId + "_64.png";
        },
        removeHtmlTags(str) {
            if (str != undefined) {
                str = str.replace(/(\n+)/g, '<br><br>');
                str = str.replace(/(<br>\s*){3,}/g, '<br><br>')
                return str.replace(/<[^>]+>/g, function (match) {
                    if (match === '<br>' || match === '<br/>') {
                        return match;
                    } else {
                        return '';
                    }
                });
            } else {
                return "";
            }
        },
    },
};
</script>
<style>
.custom-dialog .el-dialog__header {
  padding: 10px 20px; 
  max-height: 0px; 
}
/* CSS样式 */
.short-br {
  display: block;
  content: "";
  margin-top: 6px;
}

</style>