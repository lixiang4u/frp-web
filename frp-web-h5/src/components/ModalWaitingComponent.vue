<template>
  <n-modal v-model:show="modalConfig.show" :mask-closable="false" style="box-shadow: 0 0 0">
    <n-spin size="large" style="background-color: transparent;"/>
  </n-modal>
</template>

<script>

import {defineComponent, ref} from 'vue'
import {NSpin} from "naive-ui";

const modalConfig = ref({
  show: false,
})

const showModal = () => {
  modalConfig.value.show = true
}
const closeModal = () => {
  modalConfig.value.show = false
}

export default defineComponent({
  components: {
    NSpin,
  },
  setup() {
    return {
      modalConfig,
      showModal,
      closeModal,
    }
  }
})

</script>

<style scoped>

</style>