<template>
    <n-card>
        <n-input-group>
            <n-input-number v-model:value="inputCID" placeholder="借书卡编号" :show-button="false"
                :keyboard="{ ArrowUp: false, ArrowDown: false }" clearable style="width: 20%" />
            <n-button strong secondary type="primary" @click="query">查询</n-button>
        </n-input-group>
        <n-descriptions label-placement="left" style="margin-top: 20px">
            <n-descriptions-item label="姓名">{{ cardName }}</n-descriptions-item>
            <n-descriptions-item label="部门">{{ cardDepartment }}</n-descriptions-item>
            <n-descriptions-item label="身份">{{ cardType }}</n-descriptions-item>
        </n-descriptions>
        <n-data-table :single-line="false" :columns="columns" :data="dataRef" :loading="loadingRef" :pagination="pagination"
            style="margin-top: 20px" />
        <n-space justify="center">
            <n-button v-bind:disabled="currentCID==null" @click="handlerBorrow">借阅图书</n-button>
            <n-button @click="handlerDownload">导出JSON</n-button>
        </n-space>
    </n-card>
</template>

<script>
import { h, ref } from "vue";
import {
    NCard, NInputGroup, NInputNumber, NInput, NButton, NDescriptions, NDescriptionsItem,
    useDialog, useMessage, NDataTable, NPopconfirm, NSpace, NForm, NFormItem, NSelect, NSpin
} from "naive-ui";
import axios from "axios";
import FileSaver from "file-saver";

axios.defaults.baseURL = "http://api.yasyakarasu.tech/lib";

const inputCID = ref(null);
const currentCID = ref(null);
const cardName = ref(null);
const cardDepartment = ref(null);
const cardType = ref(null);

var loadingRef = ref(false);
var dataRef = ref([]);
var countRef = ref(0);

function formatTime(timestamp) {
    var date = new Date(timestamp);
    var Y = date.getFullYear() + '-';
    var M = (date.getMonth() + 1 < 10 ? '0' + (date.getMonth() + 1) : date.getMonth() + 1) + '-';
    var D = date.getDate() + ' ';
    var h = (date.getHours() < 10 ? '0' + date.getHours() : date.getHours()) + ':';
    var m = (date.getMinutes() < 10 ? '0' + date.getMinutes() : date.getMinutes()) + ':';
    var s = (date.getSeconds() < 10 ? '0' + date.getSeconds() : date.getSeconds());
    return Y + M + D + h + m + s;
}

function fetchBorrowData() {
    loadingRef.value = true
    axios.get('/borrow/list?cid=' + currentCID.value).then(res => {
        countRef.value = res.data.data.count
        if (res.data.data.count === 0) {
            dataRef.value = []
        } else {
            res.data.data.items.forEach((item) => {
                item.price = item.price.toFixed(2)
                item.borrow_time = formatTime(item.borrow_time * 1000)
                if (item.return_time === 0) {
                    item.return_time = '暂未归还'
                } else {
                    item.return_time = formatTime(item.return_time * 1000)
                }
            })
            dataRef.value = res.data.data.items
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

export default {
    setup() {
        const message = useMessage();
        const dialog = useDialog();

        const columns = [
            {
                title: '图书编号',
                key: 'book_id',
                width: 100
            },
            {
                title: '分类',
                key: 'category',
                resizable: true
            },
            {
                title: '标题',
                key: 'title',
                resizable: true
            },
            {
                title: '出版社',
                key: 'press',
                resizable: true
            },
            {
                title: '出版年份',
                key: 'publish_year'
            },
            {
                title: '作者',
                key: 'author',
                resizable: true
            },
            {
                title: '价格',
                key: 'price'
            },
            {
                title: '借书时间',
                key: 'borrow_time'
            },
            {
                title: '还书时间',
                key: 'return_time'
            },
            {
                title: '操作',
                key: 'actions',
                render(row) {
                    return h(
                        NPopconfirm,
                        {
                            onPositiveClick: () => {
                                const returnForm = {
                                    book_id: row.book_id,
                                    card_id: currentCID.value,
                                    borrow_time: Math.round(new Date(row.borrow_time.replace(/-/g, '/')).getTime() / 1000),
                                    return_time: Math.round(new Date().getTime() / 1000)
                                };
                                axios.put("/borrow/return", returnForm).then(res => {
                                    if (res.data.code === 0) {
                                        message.success("还书成功");
                                        fetchBorrowData()
                                    } else {
                                        message.error("还书失败");
                                    }
                                }).catch(err => {
                                    message.error(err.response.data.msg);
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
                                        disabled: row.return_time !== "暂未归还"
                                    },
                                    '还书'
                                )
                            },
                            default: () => {
                                return h(
                                    'span',
                                    '确认还书？'
                                )
                            }
                        }
                    )
                }
            }
        ]

        const query = () => {
            if (inputCID.value == null) {
                message.error("请输入借书卡编号");
                return;
            }
            currentCID.value = inputCID.value;
            axios.get("/card/get?cid=" + currentCID.value).then(res => {
                cardName.value = res.data.data.name;
                cardDepartment.value = res.data.data.department;
                cardType.value = res.data.data.type;
            }).catch(err => {
                message.error(err.response.data.msg)
            })
            fetchBorrowData()
        }

        const handlerDownload = () => {
            const borrow = JSON.parse(JSON.stringify(dataRef.value))
            for (var i = 0; i < borrow.length; i++) {
                borrow[i].card_id = currentCID.value;
                borrow[i].price = parseFloat(borrow[i].price);
                borrow[i].borrow_time = Math.round(new Date(borrow[i].borrow_time.replace(/-/g, '/')).getTime() / 1000);
                if (borrow[i].return_time === "暂未归还") {
                    borrow[i].return_time = 0;
                } else {
                    borrow[i].return_time = Math.round(new Date(borrow[i].return_time.replace(/-/g, '/')).getTime() / 1000);
                }
            }
            const data = JSON.stringify(borrow, null, 2);
            const blob = new Blob([data], { type: 'text/plain' });
            FileSaver.saveAs(blob, 'borrow.json');
        }

        const handlerBorrow = () => {
            const formValue = ref({
                book_id: null
            })
            const spinShow = ref(true)
            const disableSelect = ref(false)
            let options = []
            axios.post('/book/list').then(res => {
                if (res.data.data.count === 0) {
                    options = []
                } else {
                    res.data.data.results.forEach((item) => {
                        item.price = item.price.toFixed(2)
                        options.push({
                            label: item.title + "，" + item.author + "，" + item.press,
                            value: item.book_id
                        })
                    })
                }
                spinShow.value = false
            }).catch(err => {
                message.error(err.response.data.msg)
                spinShow.value = false
                disableSelect.value = true
            })
            dialog.create({
                title: '借阅图书',
                content: () => h(
                    NCard,
                    [
                        h(
                            NForm,
                            {
                                model: formValue
                            },
                            [
                                h(
                                    NFormItem,
                                    {
                                        label: "图书编号"
                                    },
                                    h(
                                        NInputNumber,
                                        {
                                            showButton: false,
                                            keyboard: { ArrowUp: false, ArrowDown: false },
                                            clearable: true,
                                            value: formValue.value.book_id,
                                            onUpdateValue: (value) => {
                                                formValue.value.book_id = value
                                            }
                                        }
                                    )
                                ),
                                h(
                                    NFormItem,
                                    {
                                        label: "选择图书"
                                    },
                                    h(
                                        NSpin,
                                        {
                                            show: spinShow.value,
                                            size: "small",
                                            style: "width: 340.62px"
                                        },
                                        h(
                                            NSelect,
                                            {
                                                disabled: disableSelect.value,
                                                filterable: true,
                                                options: options,
                                                value: formValue.value.book_id,
                                                onUpdateValue: (value) => {
                                                    formValue.value.book_id = value
                                                }
                                            }
                                        )
                                    )
                                )
                            ]
                        ),
                    ]
                ),
                positiveText: "确定",
                negativeText: "取消",
                onPositiveClick: () => {
                    const borrowForm = {
                        book_id: formValue.value.book_id,
                        card_id: currentCID.value,
                        borrow_time: Math.round(new Date().getTime() / 1000),
                        return_time: 0
                    }
                    axios.put("/borrow/borrow", borrowForm).then(res => {
                        if (res.data.code === 0) {
                            message.success("借阅成功");
                            fetchBorrowData()
                        } else {
                            message.error(res.data.msg);
                        }
                    }).catch(err => {
                        message.error(err.response.data.msg);
                    })
                }
            })
        }

        return {
            cardName,
            cardDepartment,
            cardType,
            columns,
            query,
            handlerDownload,
            handlerBorrow
        }
    },
    components: {
        NCard,
        NInputGroup,
        NInputNumber,
        NInput,
        NButton,
        NDescriptions,
        NDescriptionsItem,
        NDataTable,
        NSpace,
        NSpin,
        NSelect
    },
    data() {
        return {
            inputCID,
            currentCID,
            dataRef,
            loadingRef,
            pagination
        }
    }
}
</script>