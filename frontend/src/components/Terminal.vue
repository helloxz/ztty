<template>
    <div id="terminal-container">
        <div id="terminal" ref="terminalContainer"></div>
    </div>
</template>

<script lang="ts" setup>
import { onMounted, ref, onBeforeUnmount } from 'vue';
import { Terminal } from 'xterm';
import { FitAddon } from '@xterm/addon-fit';
import { WebLinksAddon } from '@xterm/addon-web-links';
// 正确导入 WebGL 插件
// import { WebglAddon } from 'xterm-addon-webgl';
import 'xterm/css/xterm.css';


// 引用终端容器
const terminalContainer = ref<HTMLElement | null>(null);
const fitAddon = new FitAddon();


// 初始化终端
let terminal: Terminal;
let socket: WebSocket;

// 窗口大小调整处理函数
const handleResize = () => {
  if (fitAddon) {
    setTimeout(() => {
      fitAddon.fit(); // 动态调整终端大小
      const { cols, rows } = terminal; // 获取当前的列数和行数
      socket.send(JSON.stringify({ type: 'resize', cols, rows })); // 通知后端更新尺寸
    }, 100);
  }
};

onMounted(() => {
  // 创建终端实例
  terminal = new Terminal({
    cursorBlink: true,          // 光标闪烁
    fontSize: 14,               // 字体大小
    fontFamily: 'SF Mono, Menlo, Monaco, Courier, monospace', // macOS 风格字体
    theme: {
      background: 'rgba(0, 0, 0, 0.8)', // 背景颜色，透明度0.8
      foreground: 'rgba(245, 245, 245, 0.9)',    // 前景颜色，更亮一些
      cursor: '#f8f8f8',        // 光标颜色
      selectionBackground: 'rgba(255, 255, 255, 0.3)', // 选中背景
      selectionForeground: '#f8f8f8', // 选中前景
    },
    allowTransparency: true,    // 允许背景透明
    scrollback: 10000,          // 增加回滚历史
    macOptionIsMeta: true,      // macOS Option 键作为 Meta 键
    macOptionClickForcesSelection: true, // macOS Option+点击强制选择
    cursorStyle: 'bar',         // macOS 风格光标
    rightClickSelectsWord: true, // 右键选中单词
    disableStdin: false,        // 启用输入
    windowsMode: false,         // 非 Windows 模式
    convertEol: true,           // 转换 EOL
  });

  // 加载 FitAddon 插件
  terminal.loadAddon(fitAddon);

  // 加载 WebLinksAddon 插件
  terminal.loadAddon(new WebLinksAddon());

  // 将终端附加到 DOM 容器
  if (terminalContainer.value) {
    terminal.open(terminalContainer.value);
    fitAddon.fit(); // 自动调整终端大小
    terminal.focus(); // 确保终端立即获得焦点
  }

  // 监听终端容器点击事件，自动对焦
  if (terminalContainer.value) {
    terminalContainer.value.addEventListener('click', () => {
      terminal.focus();
    });
  }

  // 监听窗口大小变化，自动调整终端大小
  // window.addEventListener('resize', handleResize);

    // 调整窗口大小，并通知后端
    handleResize();
  

  // 连接到 WebSocket 后端
  socket = new WebSocket('ws://127.0.0.1:65524/ws');

  // WebSocket 打开时
  socket.onopen = () => {
    //terminal.writeln('WebSocket connected. Type your commands below:');
  };

  // 接收后端消息并写入终端
  socket.onmessage = (event) => {
    terminal.write(event.data); // 将后端输出写入终端
    terminal.scrollToBottom();  // 滚动到底部，确保看到最新内容
  };

  // 处理 WebSocket 错误
  socket.onerror = (error) => {
    console.error('WebSocket error:', error);
    terminal.writeln(`WebSocket error: ${error}`);
  };

  // 处理 WebSocket 关闭
  socket.onclose = () => {
    terminal.writeln('WebSocket connection closed.');
  };
  

  // 监听用户输入并发送到后端
  terminal.onData((input) => {
    socket.send(input);
  });
  
  // 监听终端数据写入事件，确保滚动到底部
  terminal.onLineFeed(() => {
    terminal.scrollToBottom();
  });
  
  // 监听复合输入事件，解决中文输入法问题
  
  // 监听组合输入事件
  
  // 监听窗口大小变化，自动调整终端大小
  window.addEventListener('resize', handleResize);
  // 每隔20s再次获得焦点
  setInterval(() => {
    terminal.focus();
  }, 3000);
});

onBeforeUnmount(() => {
  // 清理事件监听器
  window.removeEventListener('resize', handleResize);
  
  // 关闭 WebSocket 连接
  if (socket && socket.readyState === WebSocket.OPEN) {
    socket.close();
  }
  
  // 销毁终端
  if (terminal) {
    terminal.dispose();
  }
});
</script>

<style scoped>
#terminal-container {
  width: 100%;
  height: calc(100vh - 55px); /* 减去标题栏高度 */
  /* background-color: rgba(0, 0, 0, 0.8); */
  border-radius: 6px;
  overflow: hidden;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  display: flex;
  flex-direction: column;
  color:#C0C4CC;
}

#terminal {
  /* width: 100%; */
  /* height: 100%; */
  padding: 0 14px 14px 14px;
  box-sizing: border-box;
  display: flex; /* 添加 flex 布局 */
  flex-direction: column; /* 垂直排列 */
  flex:1;
  /* 占满整个父元素 */
}

/* .xterm .xterm-screen{
  background: #fff;
} */


:deep(.xterm) {
  flex: 1;
  display: flex;
  flex-direction: column;
  background-color: transparent !important; /* 改为透明而不是白色 */
}

:deep(.xterm-screen) {
  background: transparent !important; /* 显式设置为透明 */
}

:deep(.xterm-viewport) {
  /* 修复底部滚动空白区域问题 */
  overflow-y: auto !important;
  scrollbar-width: thin;
  scrollbar-color: rgba(255, 255, 255, 0.3) transparent;
  flex: 1; /* 确保视口占满空间 */
  background: transparent !important; /* 确保视口也是透明的 */
}

:deep(.xterm-viewport::-webkit-scrollbar) {
  width: 8px;
}

:deep(.xterm-viewport::-webkit-scrollbar-track) {
  background: transparent;
}

:deep(.xterm-viewport::-webkit-scrollbar-thumb) {
  background-color: rgba(255, 255, 255, 0.3);
  border-radius: 4px;
}



</style>