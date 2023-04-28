<script setup>
import { NDialogProvider, NMessageProvider, NConfigProvider, zhCN, dateZhCN } from "naive-ui"
import Global from "./components/Global.vue"
const locale = zhCN
const dateLocale = dateZhCN
</script>

<template>
    <n-config-provider :locale="locale" :date-locale="dateLocale">
        <n-dialog-provider>
            <n-message-provider>
                <Global />
            </n-message-provider>
        </n-dialog-provider>
    </n-config-provider>
</template>