<template>
  <n-modal v-model:show="modalConfig.show" :on-mask-click="onModalClose">
    <n-card
        style="width: 600px"
        title="提示"
        preset="dialog"
        :bordered="false"
        size="huge"
        role="dialog"
        aria-modal="true"
    >
      <template #header-extra></template>
      <div v-html="modalConfig.message"></div>
      <template #footer>
        <div class="modal-footer" v-if="modalConfig.isConfirm">
          <n-button class="item" type="warning" @click="onCancelClick">取消</n-button>
          <n-button class="item" type="primary" @click="onConfirmClick">确认</n-button>
        </div>
      </template>
      <template #header>
        <div class="head-title">
          <div class="vcenter">
            <n-icon class="icon" :class="modalConfig.color" size="32" :component="ErrorFilled"/>
            {{ modalConfig.title }}
          </div>
          <div class="vcenter">
            <n-icon class="icon cursor" @click="onModalClose" size="32" :component="CloseRound"/>
          </div>
        </div>
      </template>
    </n-card>
  </n-modal>
</template>

<script>

import {defineComponent, ref} from 'vue'
import ErrorFilled from "@vicons/material/ErrorFilled";
import CloseRound from "@vicons/material/CloseRound";
import {NButton} from "naive-ui";

const modalConfig = ref({
  show: false,
  color: null,
  timeout: null,
  title: null,
  message: null,
  isConfirm: null,
  onCallback: null,
  onConfirmCallback: null,
  onCancelCallback: null,
  defaultConfig: {timeout: 3e3, color: 'default', title: '提示', message: '无消息', isConfirm: false},
})

const countdown = () => {
  if (modalConfig.value.timeout <= 0) {
    modalConfig.value.show = false
    if (modalConfig.value.onCallback) {
      modalConfig.value.onCallback()
    }
  } else {
    modalConfig.value.timeout -= 1e3;
    setTimeout(countdown, 1e3)
  }
}

const showModal = (config) => {
  for (const configKey in modalConfig.value.defaultConfig) {
    modalConfig.value[configKey] = modalConfig.value.defaultConfig[configKey]
  }
  if (config) {
    const keys = ['color', 'timeout', 'title', 'message', 'isConfirm']
    for (const key in keys) {
      if (config.hasOwnProperty(keys[key])) {
        modalConfig.value[keys[key]] = config[keys[key]]
      }
    }
  }
  modalConfig.value.show = true
  if (!modalConfig.value.isConfirm) {
    countdown()
  }
}

const showSuccess = (config, onCallback) => {
  config = Object.assign(config ?? {}, {color: 'success'})
  modalConfig.value.onCallback = onCallback
  showModal(config)
}
const showError = (config) => {
  config = Object.assign(config ?? {}, {color: 'error'})
  showModal(config)
}
const showWarning = (config) => {
  config = Object.assign(config ?? {}, {color: 'warning'})
  showModal(config)
}
const showConfirm = (config, onConfirmCallback, onCancelCallback) => {
  config = Object.assign(config ?? {}, {color: 'warning', isConfirm: true})
  modalConfig.value.onConfirmCallback = onConfirmCallback
  modalConfig.value.onCancelCallback = onCancelCallback
  showModal(config)
}

const onModalClose = () => {
  modalConfig.value.show = false
  modalConfig.value.timeout = 0
}

const onConfirmClick = () => {
  modalConfig.value.show = false
  if (modalConfig.value.onConfirmCallback) {
    modalConfig.value.onConfirmCallback()
  }
}
const onCancelClick = () => {
  modalConfig.value.show = false
  if (modalConfig.value.onCancelCallback) {
    modalConfig.value.onCancelCallback()
  }
}

export default defineComponent({
  components: {
    NButton,
    ErrorFilled,
    CloseRound,
  },
  setup() {
    return {
      ErrorFilled, CloseRound,
      modalConfig,
      showModal,
      showSuccess,
      showError,
      showWarning,
      showConfirm,
      onModalClose,
      onConfirmClick,
      onCancelClick,
    }
  }
})

</script>

<style scoped>

.modal-footer {
  display: flex;
  flex-direction: row;
  justify-content: end;
  align-items: center;

  .item {
    margin: 0 0 0 20px;
  }
}

.head-title {
  display: flex;
  flex-direction: row;
  justify-content: space-between;

  .default {
  }

  .primary {
    color: #18a058;
  }

  .success {
    color: #18a058;
  }

  .warning {
    color: #f0a020;
  }

  .error {
    color: #d03050;
  }

  .vcenter {
    display: flex;
    flex-direction: row;
    align-items: center;

    .icon {
      margin: 0 10px 0 -6px;
    }

    .cursor {
      cursor: pointer;
    }
  }
}

</style>