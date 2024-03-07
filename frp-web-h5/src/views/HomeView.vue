<template>
  <div class="page">

    <n-divider title-placement="left"><h2>1、配置服务器</h2></n-divider>
    <n-form
        class="plr40"
        ref="formServerConfigRef"
        :label-width="100"
        :disabled="formServerConfigDisabled"
        :model="formServerConfigValue"
        :rules="formServerConfigRules">

      <n-grid :cols="24" :x-gap="24">
        <n-form-item-gi :span="6" label="服务器地址" path="host">
          <n-input v-model:value="formServerConfigValue.host" placeholder="请输入服务器地址" :attr-size="45"/>
        </n-form-item-gi>
        <n-form-item-gi :span="4" label="服务器端口" path="bind_port">
          <n-input v-model:value="formServerConfigValue.bind_port" type="number" placeholder="请输入服务器端口"
                   :attr-size="16"/>
        </n-form-item-gi>
        <n-form-item-gi :span="4" label="vhost端口(http)" path="vhost_http_port">
          <n-input v-model:value="formServerConfigValue.vhost_http_port" type="number"
                   placeholder="请输入vhost端口(http)"
                   :attr-size="16"/>
        </n-form-item-gi>
        <n-form-item-gi :span="4" label="vhost端口(https)" path="vhost_https_port">
          <n-input v-model:value="formServerConfigValue.vhost_https_port" type="number"
                   placeholder="请输入vhost端口(https)"
                   :attr-size="16"/>
        </n-form-item-gi>
        <n-form-item-gi :span="4" label="tcpmux端口(tcpmux)" path="tcp_mux_http_connect_port">
          <n-input v-model:value="formServerConfigValue.tcp_mux_http_connect_port" type="number"
                   placeholder="请输入vhost端口(https)"
                   :attr-size="16"/>
        </n-form-item-gi>
        <n-form-item-gi :span="2">
          <n-button :disabled="formServerConfigDisabled" @click="onClickConnectServer" type="primary">连接</n-button>
        </n-form-item-gi>

      </n-grid>

    </n-form>

    <n-divider title-placement="left"><h2>2、虚拟主机配置</h2></n-divider>
    <div class="plr40">
      <n-data-table :columns="vhostColumns" :data="vhosts" :bordered="false"/>
      <n-space class="btn-right">
        <n-button type="primary" @click="onClickShowCreateOrUpdateVhost">添加</n-button>
      </n-space>
    </div>

    <n-divider title-placement="left"><h2>3、启动/重载</h2></n-divider>
    <div class="plr40">
      <n-space class="btn-left">
        <n-button type="primary" :disabled="vhosts.length===0" @click="onClickReloadVhost">启动/重载</n-button>
      </n-space>
    </div>

    <div class="ptb20"></div>
    <div class="ptb20"></div>

    <div class="plr40">
      <div class="grey">
        更新地址：<a class="grey" href="https://github.com/lixiang4u/frp-web/releases/">frp-web</a>
      </div>
    </div>

    <n-modal v-model:show="showModalCreateVhost" preset="dialog" style="width: 880px">
      <template #header>
        <div>添加/修改配置</div>
      </template>
      <div>

        <n-form
            class="plr40"
            ref="formProxyConfigRef"
            :model="formProxyConfigValue"
            :rules="formProxyConfigRules">
          <n-grid :cols="24" :x-gap="24">
            <n-form-item-gi :span="12" label="代理类型" path="type">
              <n-select
                  v-model:value="formProxyConfigValue.type"
                  :options="proxyTypeOptions"
                  placeholder="请选择代理类型"
                  @updateValue="onChangeProxyType"
                  clearable
                  filterable/>
            </n-form-item-gi>
            <n-form-item-gi
                :span="12"
                label="代理名称"
                path="name"
            >
              <n-input v-model:value="formProxyConfigValue.name" placeholder="请输入代理名称" :attr-size="16"/>
            </n-form-item-gi>

            <n-form-item-gi
                :span="12"
                label="本地地址"
                path="local_addr"
            >
              <n-input v-model:value="formProxyConfigValue.local_addr" placeholder="请输入本地地址（只支持ip:port格式）"
                       :attr-size="16"/>
            </n-form-item-gi>

            <n-form-item-gi
                :span="12"
                label="服务器端口(服务器端口会被独占)"
                path="remote_port"
                v-if="['tcp','tcpmux'].includes(formProxyConfigValue.type)"
            >
              <n-input v-model:value="formProxyConfigValue.remote_port"
                       placeholder="请输入服务器端口"
                       :disabled="formProxyConfigValue.custom_domain"
                       :attr-size="16"/>
            </n-form-item-gi>

            <n-form-item-gi
                :span="12"
                label="公网域名"
                path="custom_domain"
                v-if="proxyTypeOptions.map(item=>item.value).includes(formProxyConfigValue.type)"
            >
              <n-input disabled v-model:value="formProxyConfigValue.custom_domain" placeholder="请输入公网域名"
                       :attr-size="45"/>
            </n-form-item-gi>

            <n-form-item-gi
                :span="12"
                label="证书文件(https)"
                path="crt_path"
                v-if="['http','https'].includes(formProxyConfigValue.type)"
            >
              <n-input disabled v-model:value="formProxyConfigValue.crt_path" placeholder="请输入证书文件(https)"
                       :attr-size="16"/>
            </n-form-item-gi>
            <n-form-item-gi
                :span="12"
                label="证书密钥(https)"
                path="key_path"
                v-if="['http','https'].includes(formProxyConfigValue.type)"
            >
              <n-input disabled v-model:value="formProxyConfigValue.key_path" placeholder="请输入证书密钥(https)"
                       :attr-size="16"/>
            </n-form-item-gi>

          </n-grid>
        </n-form>
      </div>
      <template #action>
        <n-button type="primary" @click="onClickVhostSave">提交</n-button>
      </template>
    </n-modal>


    <ModalTipsComponent ref="modalTipsRef"/>
    <ModalWaitingComponent ref="modalWaitingRef"/>

  </div>
</template>

<script>
import {defineComponent, h, onBeforeMount, ref} from "vue";
import {NButton, NSpace, useDialog, useMessage} from "naive-ui";
import api from "@/api/api.js";
import ModalTipsComponent from "@/components/ModalTipsComponent.vue";
import ModalWaitingComponent from "@/components/ModalWaitingComponent.vue";


let message = null// message弹出框
let dialog = null// dialog弹出框
const modalTipsRef = ref(null)
const modalWaitingRef = ref(null)
const formServerConfigDisabled = ref(false)
const formServerConfigRef = ref(null)

const showModalCreateVhost = ref(false)

const vhosts = ref([])

const createVhostColumns = () => {
  return [
    {
      title: "#",
      key: "index",
      render(row, index) {
        return h(
            NSpace,
            {},
            {default: () => index + 1,}
        )
      }
    },
    {
      title: "类型",
      key: "type",
      resizable: true,
    },
    {
      title: "名称",
      key: "name",
      resizable: true,
    },
    {
      title: "公网域名",
      key: "custom_domain",
      resizable: true,
      render(row) {
        if (['tcp', 'tcpmux'].includes(row.type)) {
          const tmpUrl = `${row.custom_domain}:${row.remote_port}`
          return [
            h('div', {class: 'lh150'}, h('a', {}, tmpUrl)),
          ]
        } else {
          const tmpUrl = `${row.type}://${row.custom_domain}`
          return [
            h('div', {class: 'lh150'}, h('a', {href: tmpUrl, target: '_blank'}, tmpUrl)),
            // h('div', {class: 'lh150'}, h('a', {href: `${tmpUrl}/files/`, target: '_blank'}, `${tmpUrl}/files/`)),
          ]
        }
      },
    },
    {
      title: "本地地址",
      key: "local_addr",
      resizable: true,
      render(row) {
        if (['tcp', 'tcpmux'].includes(row.type)) {
          const tmpUrl = row.local_addr
          return [
            h('div', {class: 'lh150'}, h('a', {}, tmpUrl)),
          ]

        } else {
          const tmpUrl = `${row.type}://${row.local_addr}`
          return [
            h('div', {class: 'lh150'}, h('a', {href: tmpUrl, target: '_blank'}, tmpUrl)),
            // h('div', {class: 'lh150'}, h('a', {href: `${tmpUrl}/files/`, target: '_blank'}, `${tmpUrl}/files/`)),
          ]
        }
      },
    },
    {
      title: "操作",
      key: "actions",
      maxWidth: 160,
      width: 160,
      render(row, index) {
        return [
          h(
              NButton,
              {
                type: "warning",
                onClick: () => {
                  showModalCreateVhost.value = true
                  formProxyConfigValue.value = Object.assign({index: index}, row)
                }
              },
              {default: () => '修改'}
          ),
          h("span", {}, " "),
          h(
              NButton,
              {
                type: "error",
                onClick: () => {
                  dialog.error({
                    title: '警告',
                    content: `确定删除 ${row.name}？`,
                    positiveText: '确定',
                    negativeText: '不确定',
                    onPositiveClick: () => {
                      removeVhostFromList(row, index)
                    },
                    onNegativeClick: () => {
                    }
                  })
                }
              },
              {default: () => '删除'}
          )
        ]
      }
    },
  ]
}

const vhostColumns = createVhostColumns()

const proxyTypeOptions = ref([
  {
    label: "代理本地http",
    value: 'http',
    default_local_addr: '127.0.0.1:8000',
    default_remote_port: '',
  },
  {
    label: '代理本地https',
    value: 'https',
    default_local_addr: '127.0.0.1:8000',
    default_remote_port: '',
  },
  {
    label: '代理本地tcp(适用ssh)',
    value: 'tcp',
    default_local_addr: '127.0.0.1:22',
  },
  {
    label: '代理本地tcp(mux)(只适用httpconnect)',
    value: 'tcpmux',
    default_local_addr: '127.0.0.1:6000',
  }
])

const onChangeProxyType = (v) => {
  const find = proxyTypeOptions.value.find(item => {
    return item.value === v
  })
  formProxyConfigValue.value.local_addr = find.default_local_addr
  formProxyConfigValue.value.remote_port = formServerConfigValue.value.tcp_mux_http_connect_port

}

const onClickConnectServer = () => {
}

const onClickVhostSave = () => {
  formProxyConfigRef.value?.validate(errors => {
    if (errors) {
      console.log('[errors]', errors)
      return
    }
    api.newVhost({
      id: formProxyConfigValue.value.id,
      type: formProxyConfigValue.value.type,
      name: formProxyConfigValue.value.name,
      local_addr: formProxyConfigValue.value.local_addr,
      remote_port: +formProxyConfigValue.value.remote_port,
    }).then(resp => {
      console.log('[newVhost-resp]', resp)

      getVhostList()

    }).catch(err => {
      console.log('[err]', err)
      modalTipsRef.value.showError({'message': err.msg ?? ('系统错误：' + err)})
    }).finally(() => {
      showModalCreateVhost.value = false
    })
    // if (errors) {
    //   message.error("配置错误")
    // } else {
    //   addVhostToList(formProxyConfigValue.value)
    //   clearFormProxyConfigValue()
    //   showModalCreateVhost.value = false
    // }
  })

}

const addVhostToList = (data) => {
  console.log('[dataX]', data)
  if (data.index || data.index === 0) {
    vhosts.value[data.index] = data
    delete data.index
  } else {
    vhosts.value.push(data)
  }
}

const removeVhostFromList = (data, index) => {
  api.removeVhost(data.id).then(resp => {
    console.log('[removeVhost-resp]', resp)

    getVhostList()

  }).catch(err => {
    console.log('[err]', err)
    modalTipsRef.value.showError({'message': err.msg ?? ('系统错误：' + err)})
  })
  // vhosts.value.splice(index, 1)
}

const formProxyConfigRules = {
  type: {
    required: true,
    message: "请选择类型",
    trigger: ['blur', 'change'],
  },
  name: {
    required: true,
    message: "请输入名称",
    trigger: ['blur', 'change'],
  },
  local_addr: {
    required: true,
    message: "请输入本地地址",
    trigger: ['blur', 'change'],
  },
  remote_port: {
    required: true,
    message: "请输入服务器端口",
    trigger: ['blur', 'change'],
  },
}

const formServerConfigValue = ref({
  host: "",
  bind_port: "",
  vhost_http_port: "",
  vhost_https_port: "",
  tcp_mux_http_connect_port: "",
})

const formServerConfigRules = {
  host: {
    required: true,
    message: "请输入服务器地址",
  },
  bind_port: {
    required: true,
    message: "请输入服务器端口",
    type: 'number',
  },
  vhost_http_port: {
    required: true,
    message: "请输入vhost端口(http)",
    type: 'number',
  },
  vhost_https_port: {
    required: true,
    message: "请输入vhost端口(https)",
    type: 'number',
  },
}

const defaultProxyConfig = {
  id: null,
  type: null,
  name: null,
  custom_domain: null,
  custom_domains: [],
  local_port: null,
  local_addr: null,
  remote_port: null,
  crt_path: null,
  key_path: null,
}

const formProxyConfigRef = ref(null)
const formProxyConfigValue = ref(Object.assign({}, defaultProxyConfig))

const clearFormProxyConfigValue = () => {
  formProxyConfigValue.value = Object.assign({}, defaultProxyConfig)
}

const onBeforeMountHandler = () => {
  api.getConfig().then(resp => {
    resp.data.config.bind_port = '' + resp.data.config.bind_port
    resp.data.config.vhost_http_port = '' + resp.data.config.vhost_http_port
    resp.data.config.vhost_https_port = '' + resp.data.config.vhost_https_port
    resp.data.config.tcp_mux_http_connect_port = '' + resp.data.config.tcp_mux_http_connect_port
    formServerConfigValue.value = resp.data.config
    formServerConfigDisabled.value = true
    console.log('[getConfig-resp]', resp)

    // getVhostList()
    api.getVhosts().then(resp => {
      console.log('[getVhostList-resp]', resp)
      vhosts.value = resp.data.vhosts

      api.reloadVhost().then(resp => {
        console.log('[reloadVhost-resp]', resp)
      }).catch(err => {
        console.log('[err]', err)
        modalTipsRef.value.showError({'message': err.msg ?? ('系统错误：' + err)})
      })

    }).catch(err => {
      console.log('[err]', err)
      modalTipsRef.value.showError({'message': err.msg ?? ('系统错误：' + err)})
    })

  }).catch(err => {
    console.log('[err]', err)
    modalTipsRef.value.showError({'message': err.msg ?? ('系统错误：' + err)})
  })
}

const getVhostList = () => {
  api.getVhosts().then(resp => {
    console.log('[getVhostList-resp]', resp)
    vhosts.value = resp.data.vhosts
  }).catch(err => {
    console.log('[err]', err)
    modalTipsRef.value.showError({'message': err.msg ?? ('系统错误：' + err)})
  })
}

const newVhost = () => {
  //newVhost
  api.newVhost({}).then(resp => {
    console.log('[newVhost-resp]', resp)
  }).catch(err => {
    modalTipsRef.value.showError({'message': err.msg ?? ('系统错误：' + err)})
  })
}

const loadVhostOrCreate = () => {
  api.newVhost({
    type: '',
    machine_id: 'machine_id',
  }).then(resp => {
    console.log('[newVhost-resp]', resp)
  }).catch(err => {
    modalTipsRef.value.showError({'message': err.msg ?? ('系统错误：' + err)})
  })
}


const onClickShowCreateOrUpdateVhost = () => {
  clearFormProxyConfigValue()
  showModalCreateVhost.value = true
}

const onClickReloadVhost = () => {
  modalWaitingRef.value.showModal()
  api.reloadVhost().then(resp => {
    console.log('[reloadVhost-resp]', resp)
  }).catch(err => {
    modalTipsRef.value.showError({'message': err.msg ?? ('系统错误：' + err)})
  }).finally(() => {
    modalWaitingRef.value.closeModal()
  })
}

export default defineComponent({
  components: {
    ModalTipsComponent,
    ModalWaitingComponent,
  },
  setup() {
    message = useMessage()
    dialog = useDialog()

    onBeforeMount(onBeforeMountHandler)

    return {
      formProxyConfigRef,
      formServerConfigValue,
      formServerConfigRules,
      formProxyConfigRules,
      onClickConnectServer,
      formServerConfigRef,
      proxyTypeOptions,
      vhosts,
      vhostColumns,
      formServerConfigDisabled,
      showModalCreateVhost,
      onClickShowCreateOrUpdateVhost,
      onClickReloadVhost,
      onClickVhostSave,
      formProxyConfigValue,
      modalTipsRef,
      modalWaitingRef,
      onChangeProxyType,
    };
  }
});

</script>

<style scoped>
.page {
  max-width: 1200px;
  margin: 20px auto;
}

.plr40 {
  padding: 0 40px;
}

.ptb20 {
  padding: 20px 0;
}

.light-green {
  height: 108px;
  background-color: rgba(0, 128, 0, 0.12);
}

.green {
  height: 108px;
  background-color: rgba(0, 128, 0, 0.24);
}

.item-host-title {
  font-weight: bold;
}

.item-host {
  line-height: 58px;
  font-size: 120%;
  padding: 0 10px;
  border-bottom: 1px solid rgb(239, 239, 245);
}

.item-host:hover {
  background-color: #f8f8f8;
}

:deep(thead .n-data-table-th__title) {
  font-weight: bold;
}

.btn-right {
  margin: 20px 0;
  display: flex;
  justify-content: flex-end !important;
}

.btn-left {
  margin: 20px 0;
  display: flex;
}

:deep(.lh150) {
  line-height: 200%;
}

.grey {
  color: #c3c3c3;
}

</style>