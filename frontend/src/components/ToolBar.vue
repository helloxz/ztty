<template>
    <div class="tools">
        <!-- 右下角按钮 -->
        <Info />
        <!-- 右下角按钮END -->
        <el-tabs
            v-model="editableTabsValue"
            type="card"
            class="tabs"
            editable
            @edit="handleTabsEdit"
        >
            <template #add-icon>
                <el-icon><Plus /></el-icon>
            </template>
            <el-tab-pane
            v-for="item in editableTabs"
            :key="item.name"
            :label="item.title"
            :name="item.name"
            >
            <component :is="item.content" />
            </el-tab-pane>
        </el-tabs>

    </div>
</template>

<script setup lang="ts">
import { ref, markRaw, onMounted, onBeforeUnmount } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import type { TabPaneName } from 'element-plus'
import Terminal from './Terminal.vue'
import Info from './Info.vue'

let tabIndex = 1
const editableTabsValue = ref('1')
const editableTabs = ref([
  {
    title: 'Tab 1',
    name: '1',
    content: markRaw(Terminal),
  }
])

// 处理快捷键的函数
const handleKeyDown = (event: KeyboardEvent) => {
  // Command+N 或 Ctrl+N 添加新标签页
  if ((event.metaKey || event.ctrlKey) && event.key === 'n') {
    event.preventDefault() // 阻止浏览器默认行为
    addNewTab()
  }
  
  // Command+W 或 Ctrl+W 关闭当前标签页
  if ((event.metaKey || event.ctrlKey) && event.key === 'w') {
    event.preventDefault() // 阻止浏览器默认行为
    if (editableTabs.value.length > 1) { // 确保至少保留一个标签页
      removeTab(editableTabsValue.value)
    }
  }
}

const handleTabsEdit = (
  targetName: TabPaneName | undefined,
  action: 'remove' | 'add'
) => {
  if (action === 'add') {
    const newTabName = `${++tabIndex}`
    editableTabs.value.push({
      title: 'Tab ' + tabIndex,
      name: newTabName,
      content: markRaw(Terminal),
    })
    editableTabsValue.value = newTabName
  } else if (action === 'remove') {
    const tabs = editableTabs.value
    let activeName = editableTabsValue.value
    if (activeName === targetName) {
      tabs.forEach((tab, index) => {
        if (tab.name === targetName) {
          const nextTab = tabs[index + 1] || tabs[index - 1]
          if (nextTab) {
            activeName = nextTab.name
          }
        }
      })
    }

    editableTabsValue.value = activeName
    editableTabs.value = tabs.filter((tab) => tab.name !== targetName)
  }
}

// 添加新标签页的函数
const addNewTab = () => {
  const newTabName = `${++tabIndex}`
  editableTabs.value.push({
    title: 'Tab ' + tabIndex,
    name: newTabName,
    content: markRaw(Terminal),
  })
  editableTabsValue.value = newTabName
}

// 关闭标签页的函数
const removeTab = (targetName: TabPaneName | undefined) => {
  if (!targetName) return
  
  const tabs = editableTabs.value
  let activeName = editableTabsValue.value
  
  if (activeName === targetName) {
    tabs.forEach((tab, index) => {
      if (tab.name === targetName) {
        const nextTab = tabs[index + 1] || tabs[index - 1]
        if (nextTab) {
          activeName = nextTab.name
        }
      }
    })
  }

  editableTabsValue.value = activeName
  editableTabs.value = tabs.filter((tab) => tab.name !== targetName)
}

// 组件挂载时添加事件监听
onMounted(() => {
  window.addEventListener('keydown', handleKeyDown)
})

// 组件卸载前移除事件监听
onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleKeyDown)
})
</script>

<style>
.tabs{
    --wails-draggable: drag; /* 启用拖动 */
    /* margin-left:78px; */
}
.tools{
    position: fixed;
    top:0;
    left:0;
    width: 100%;
}
.el-tabs--card>.el-tabs__header{
    border-bottom:none !important;
}


.el-tabs--card>.el-tabs__header .el-tabs__nav{
    border:none;
}

.el-tabs__item.is-active, .el-tabs__item:hover{
    color:#FAFCFF;
    /* border-bottom: 2px solid #409EFF; */
}
.el-tabs--card>.el-tabs__header .el-tabs__item.is-active{
    /* color:#FAFCFF; */
    /* border-bottom: 2px solid #409EFF; */
    border-bottom: none;
    background: rgba(30, 35, 40, 0.9); /* 深灰蓝 */
}

.el-tabs--card>.el-tabs__header .el-tabs__item{
    border-left: none;
}

.el-tabs__new-tab{
    border:none;
    color:#C0C4CC;
    /* margin-right: 16px; */
    padding:12px 12px 12px 0;
}

.el-tabs__new-tab .el-icon{
    font-size: 16px;
    font-weight: bold;
    color:#C0C4CC;
    /**圆角 */
    border-radius: 50%;
    /* background-color: #303133; */
    padding: 3px;
    margin-right: 8px;
    cursor: pointer;
    border:1px solid #606266;
    
}
.el-tabs__nav-wrap{
    margin-left: 78px;
}
.el-tabs__item{
    color:#909399;
    /* width: 92px; */
}


.el-tabs--card>.el-tabs__header{
  background: rgba(40, 40, 45, 0.7) !important;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}
/* :deep(.el-tabs--card>.el-tabs__header){
    border-bottom: none;
} */
</style>