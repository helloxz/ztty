<template>
    <!-- 内容区域 -->
    <div class="info-content">
        <el-dialog
            v-model="dialogVisible"
            title="About"
            width="450"
            class="info-dialog"
        >
            <div>
                <p><b>Ztty v1.0.0</b></p>
                <p>A lightweight, tabbed terminal built with Golang & Wails.</p>
                <p><b>Github:</b> <el-link :underline="false" @click="BrowserOpenURL('https://github.com/helloxz/ztty')" type="primary">https://github.com/helloxz/ztty</el-link></p>
                <p><b>License:</b> GNU AGPLv3</p>
            </div>
            
        </el-dialog>
    </div>
    <!-- 内容区域END -->
    <div class="info-btn">
        <el-link @click="dialogVisible = true" :underline="false" type="info"><el-icon><InfoFilled /></el-icon></el-link>
    </div>
</template>

<script lang="ts" setup>
import {ref} from 'vue';
import { InfoFilled } from '@element-plus/icons-vue';
import { BrowserOpenURL } from '../../wailsjs/runtime';
const dialogVisible = ref(false);
</script>

<style>
.info-btn{
    position: fixed;
    bottom: 4px;
    right: 26px;
    z-index: 102;
}
.info-dialog{
    background-color: rgba(60, 60, 60, 0.95);
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
    backdrop-filter: blur(2px); /* 弹窗背景模糊（部分浏览器支持） */
    border-radius: 16px;
    /* 或带轻微透明度 */
   
    
}

.el-dialog__body,.el-dialog__title{
    color: rgba(255, 255, 255, 0.95);
}
</style>