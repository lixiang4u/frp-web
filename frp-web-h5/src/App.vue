<script setup>
import {RouterView} from 'vue-router'
</script>

<template>
  <n-message-provider>
    <n-dialog-provider>
      <RouterView/>
    </n-dialog-provider>
  </n-message-provider>
</template>

<style scoped>
</style>
