<template>
  <div class="page">

    <ul>
      <li>支持http代理</li>
      <li>支持https代理</li>
      <li>支持tcp代理（SSH/MySQL/Redis...）</li>
      <li>支持自定义域名（需要备案 -> DNS设置cname -> frp-api后台配置 -> frp-web更新vhost -> 使用）</li>
      <li>...</li>
    </ul>

  </div>
</template>

<script>

</script>

<style scoped>
.page {
  max-width: 1200px;
  margin: 20px auto;
}

:deep(input[type=text]),
:deep(input[type=number]),
:deep(.n-form-item-label__text),
.page {
  font-family: "Source Code Pro", Arial, sans-serif;
}

.cursor,
:deep(.cursor) {
  cursor: pointer;
}

.plr40 {
  padding: 0 40px;
}

.ptb20 {
  padding: 20px 0;
}

.grey {
  color: #c3c3c3;
}

ul {
  list-style-type: square;
  font-size: 20px;
}

</style>