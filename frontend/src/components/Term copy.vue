<template>
    <div id="terminal" ref="terminalContainer"></div>
</template>

<script lang="ts" setup>
import { onMounted, ref } from 'vue';
import { Terminal } from 'xterm';
import { FitAddon } from '@xterm/addon-fit';
import { WebLinksAddon } from '@xterm/addon-web-links';

// 引用终端容器
const terminalContainer = ref<HTMLElement | null>(null);
const fitAddon = new FitAddon();

// 初始化终端
let terminal: Terminal;
let socket: WebSocket;

onMounted(() => {
  // 创建终端实例
  terminal = new Terminal({
    cursorBlink: true, // 光标闪烁
    fontSize: 14,      // 字体大小
    theme: {
      background: 'rgba(27, 38, 54, 1)', // 背景颜色
      foreground: '#fff', // 前景颜色
    },
  });

  // 加载 FitAddon 插件
  const fitAddon = new FitAddon();
  terminal.loadAddon(fitAddon);

  // 加载 WebLinksAddon 插件（可选）
  terminal.loadAddon(new WebLinksAddon());

  // 将终端附加到 DOM 容器
  if (terminalContainer.value) {
    terminal.open(terminalContainer.value);
    fitAddon.fit(); // 自动调整终端大小
  }

  // 连接到 WebSocket 后端
  socket = new WebSocket('ws://127.0.0.1:9080/ws');

  // WebSocket 打开时
  socket.onopen = () => {
    //terminal.writeln('WebSocket connected. Type your commands below:');
  };

  // 接收后端消息并写入终端
  socket.onmessage = (event) => {
    terminal.write(event.data); // 将后端输出写入终端
    // terminal.scrollToBottom(); // 滚动到底部
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
});


</script>

<style scoped>
#terminal {
  width: auto;
  height: 100vh;
  padding:16px;
}
</style>