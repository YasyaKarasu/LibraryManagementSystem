<template>
    <n-card>
        <n-data-table :single-line="false" :columns="columns" :data="dataRef" :loading="loadingRef"
            :pagination="pagination" />
        <n-space justify="center">
            <n-button @click="handlerAdd">添加借书证</n-button>
            <n-button @click="handlerDownload">导出JSON</n-button>
        </n-space>
    </n-card>
</template>

<script>
import { h, ref } from "vue";
import axios from "axios";
import {
    NCard, NDataTable, NSpace, NButton, NPopconfirm, NRadioGroup, NRadioButton,
    NForm, NFormItem, NInput, useMessage, useDialog
} from "naive-ui";
import FileSaver from "file-saver";

axios.defaults.baseURL = 'http://api.yasyakarasu.tech/lib';

var loadingRef = ref(true);
var dataRef = ref([]);
var countRef = ref(0);

function fetchCardData() {
    loadingRef.value = true
    axios.get('/card/list').then(res => {
        countRef.value = res.data.data.count
        if (res.data.data.count === 0) {
            dataRef.value = []
        } else {
            res.data.data.cards.forEach((item) => {
                var t = item.type
                item.type = t + "（" + (t === "T" ? "教师" : t === "S" ? "学生" : "未知") + "）"
            })
            dataRef.value = res.data.data.cards
        }
        loadingRef.value = false
    }).catch(err => {
        console.log(err)
        loadingRef.value = false
    })
}

const pagination = {
    pageSize: 10,
    prefix: () => h(
        'span',
        "共 " + countRef.value + " 项"
    ),
}

const inputStringValidationStatus = (value, changed) => {
    if ((value === "" || value === null) && changed) {
        return "error";
    }
    return "success";
};

const createStringFeedback = (value, changed, label) => {
    if ((value === "" || value === null) && changed)
        return `${label}不能为空`;
    return undefined;
};

export default {
    setup() {
        const message = useMessage()
        const dialog = useDialog()
        const columns = [
            {
                title: '借书证编号',
                key: 'card_id'
            },
            {
                title: '姓名',
                key: 'name'
            },
            {
                title: '单位',
                key: 'department'
            },
            {
                title: '身份',
                key: 'type'
            },
            {
                title: '删除',
                key: 'remove',
                render: (row) => {
                    return h(
                        NPopconfirm,
                        {
                            onPositiveClick: () => {
                                axios.delete("/card/remove?cid=" + row.card_id).then(res => {
                                    message.success("删除成功！")
                                    countRef.value = 0
                                    loadingRef.value = true
                                    axios.get('/card/list').then(res => {
                                        countRef.value = res.data.data.count
                                        if (res.data.data.count === 0) {
                                            dataRef.value = []
                                        } else {
                                            res.data.data.cards.forEach((item) => {
                                                var t = item.type
                                                item.type = t + "（" + (t === "T" ? "教师" : t === "S" ? "学生" : "未知") + "）"
                                            })
                                            dataRef.value = res.data.data.cards
                                        }
                                        loadingRef.value = false
                                    }).catch(err => {
                                        message.error(err.response.data.msg)
                                    })
                                }).catch(err => {
                                    message.error(err.response.data.msg)
                                })
                            }
                        },
                        {
                            trigger: () => {
                                return h(
                                    NButton,
                                    {
                                        tertiary: true,
                                        size: 'small',
                                        type: 'error',
                                    },
                                    '删除'
                                )
                            },
                            default: () => {
                                return h(
                                    'span',
                                    '确认删除？'
                                )
                            }
                        }
                    )
                }
            }
        ]
        const handlerAdd = () => {
            const formValue = ref({
                name: null,
                department: null,
                type: null
            })
            const rules = {
                name: { required: true },
                department: { required: true },
                type: { required: true }
            }
            var nameChanged = false;
            var departmentChanged = false;
            var typeChanged = false;

            dialog.create({
                title: '添加借书证',
                content: () => h(
                    NCard,
                    [
                        h(
                            NForm,
                            {
                                model: formValue,
                                rules: rules,
                            },
                            [
                                h(
                                    NFormItem,
                                    {
                                        label: "姓名",
                                        path: "name",
                                        validationStatus: inputStringValidationStatus(formValue.value.name, nameChanged),
                                        feedback: createStringFeedback(formValue.value.name, nameChanged, '姓名')
                                    },
                                    h(
                                        NInput,
                                        {
                                            maxlength: 63,
                                            showCount: true,
                                            clearable: true,
                                            onUpdateValue: (value) => {
                                                formValue.value.name = value
                                                nameChanged = true
                                            }
                                        }
                                    )
                                ),
                                h(
                                    NFormItem,
                                    {
                                        label: "部门",
                                        path: "department",
                                        validationStatus: inputStringValidationStatus(formValue.value.department, departmentChanged),
                                        feedback: createStringFeedback(formValue.value.department, departmentChanged, '标题')
                                    },
                                    h(
                                        NInput,
                                        {
                                            maxlength: 63,
                                            showCount: true,
                                            clearable: true,
                                            onUpdateValue: (value) => {
                                                formValue.value.department = value
                                                departmentChanged = true
                                            }
                                        }
                                    )
                                ),
                                h(
                                    NFormItem,
                                    {
                                        label: "身份",
                                        path: "type",
                                        validationStatus: inputStringValidationStatus(formValue.value.type, typeChanged),
                                        feedback: createStringFeedback(formValue.value.type, typeChanged, '身份')
                                    },
                                    h(
                                        NRadioGroup,
                                        {
                                            value: formValue.value.type,
                                            onUpdateValue: (value) => {
                                                formValue.value.type = value
                                                typeChanged = true
                                            }
                                        },
                                        [
                                            h(
                                                NRadioButton,
                                                {
                                                    label: "教师",
                                                    value: "T",
                                                }
                                            ),
                                            h(
                                                NRadioButton,
                                                {
                                                    label: "学生",
                                                    value: "S",
                                                }
                                            )
                                        ]
                                    )
                                ),
                            ]
                        )
                    ]
                ),
                positiveText: '确定',
                onPositiveClick: () => {
                    if (inputStringValidationStatus(formValue.value.name, true) === 'error') {
                        message.error('请输入姓名')
                        return false
                    }
                    if (inputStringValidationStatus(formValue.value.department, true) === 'error') {
                        message.error('请输入部门')
                        return false
                    }
                    if (inputStringValidationStatus(formValue.value.type, true) === 'error') {
                        message.error('请选择身份')
                        return false
                    }
                    axios.post('/card/create', formValue.value).then(res => {
                        if (res.data.code === 0) {
                            message.success('添加成功')
                            fetchCardData()
                        } else {
                            message.error(res.data.msg)
                        }
                    }).catch(err => {
                        message.error(err.response.data.msg)
                    })
                },
                negativeText: '取消',
            })
        };
        const handlerDownload = () => {
            const card = JSON.parse(JSON.stringify(dataRef.value))
            for (var i = 0; i < card.length; i++) {
                card[i].type = card[i].type[0]
            }
            const data = JSON.stringify(card, null, 2);
            const blob = new Blob([data], { type: 'text/plain' });
            FileSaver.saveAs(blob, 'card.json');
        };
        return {
            columns,
            handlerAdd,
            handlerDownload
        }
    },
    created() {
        fetchCardData()
    },
    data() {
        return {
            loadingRef,
            dataRef,
            pagination
        }
    },
    components: {
        NCard,
        NDataTable,
        NSpace,
        NButton,
        NPopconfirm
    }
}
</script>