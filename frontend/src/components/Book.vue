<script setup>
import { h } from "vue";
import {
    NCard, NDataTable, NSpace, NButton, useDialog, NForm, NFormItem, NInput, NDatePicker,
    NInputNumber, NPopconfirm, NUpload, NUploadDragger, NIcon, NText, NH4, NSpin, useMessage,
    NRadioGroup, NRadioButton, NSlider, NSelect
} from "naive-ui";
import { ArchiveOutline as ArchiveIcon } from "@vicons/ionicons5";
import FileSaver from "file-saver";
const parseCurrency = (input) => {
    const nums = input.replace(/(,|¥|\s)/g, "").trim();
    if (/^\d+(\.(\d+)?)?$/.test(nums))
        return Number(nums);
    return nums === "" ? null : Number.NaN;
};
const formatCurrency = (value) => {
    if (value === null)
        return "";
    return `${value.toLocaleString("en-US")} ¥`;
}

const inputStringValidationStatus = (value, changed) => {
    if ((value === "" || value === null) && changed) {
        return "error";
    }
    return "success";
};

const inputNumberValidationStatus = (value, changed) => {
    if ((value === 0 || value === null) && changed) {
        return "error";
    }
    return "success";
};

const createStringFeedback = (value, changed, label) => {
    if ((value === "" || value === null) && changed)
        return `${label}不能为空`;
    return undefined;
};

const createNumberFeedback = (value, changed, label) => {
    if ((value === 0 || value === null) && changed)
        return `${label}不能为空`;
    return undefined;
};

const message = useMessage();
const dialog = useDialog();
const handlerFilter = () => {
    const sortByOptions = [
        {
            label: "图书编号",
            value: "book_id"
        },
        {
            label: "分类",
            value: "category"
        },
        {
            label: "标题",
            value: "title"
        },
        {
            label: "出版社",
            value: "press"
        },
        {
            label: "出版年份",
            value: "publish_year"
        },
        {
            label: "作者",
            value: "author"
        },
        {
            label: "价格",
            value: "price"
        },
        {
            label: "库存",
            value: "stock"
        }
    ]
    const sortOrderOptions = [
        {
            label: "升序",
            value: "ASC"
        },
        {
            label: "降序",
            value: "DESC"
        }
    ]
    dialog.create({
        title: '设置筛选器',
        content: () => h(
            NCard,
            [
                h(
                    NForm,
                    {
                        model: filterRef
                    },
                    [
                        h(
                            NFormItem,
                            {
                                label: "分类"
                            },
                            h(
                                NInput,
                                {
                                    value: filterRef.value.category,
                                    maxlength: 63,
                                    showCount: true,
                                    clearable: true,
                                    onUpdateValue: (value) => {
                                        if (value === "") {
                                            filterRef.value.category = null;
                                        } else {
                                            filterRef.value.category = value;
                                        }
                                    }
                                }
                            )
                        ),
                        h(
                            NFormItem,
                            {
                                label: "标题"
                            },
                            h(
                                NInput,
                                {
                                    value: filterRef.value.title,
                                    maxlength: 63,
                                    showCount: true,
                                    clearable: true,
                                    onUpdateValue: (value) => {
                                        if (value === "") {
                                            filterRef.value.title = null;
                                        } else {
                                            filterRef.value.title = value;
                                        }
                                    }
                                }
                            )
                        ),
                        h(
                            NFormItem,
                            {
                                label: "出版社"
                            },
                            h(
                                NInput,
                                {
                                    value: filterRef.value.press,
                                    maxlength: 63,
                                    showCount: true,
                                    clearable: true,
                                    onUpdateValue: (value) => {
                                        if (value === "") {
                                            filterRef.value.press = null;
                                        } else {
                                            filterRef.value.press = value;
                                        }
                                    }
                                }
                            )
                        ),
                        h(
                            NFormItem,
                            {
                                label: "出版年份"
                            },
                            h(
                                NSpace,
                                {
                                    vertical: true
                                },
                                [
                                    h(
                                        NSlider,
                                        {
                                            value: publish_year_value,
                                            range: true,
                                            min: 1900,
                                            max: 2100,
                                            onUpdateValue: ([miny, maxy]) => {
                                                if (maxy < miny) {
                                                    return;
                                                }
                                                publish_year_value.value = [miny, maxy];
                                                filterRef.value.min_publish_year = miny;
                                                filterRef.value.max_publish_year = maxy;
                                            }
                                        }
                                    ),
                                    h(
                                        NSpace,
                                        [
                                            h(
                                                NInputNumber,
                                                {
                                                    value: publish_year_value.value[0],
                                                    size: "small",
                                                    onUpdateValue: (value) => {
                                                        if (value > publish_year_value.value[1]) {
                                                            return;
                                                        }
                                                        publish_year_value.value = [value, publish_year_value.value[1]];
                                                        filterRef.value.min_publish_year = value;
                                                    }
                                                }
                                            ),
                                            h(
                                                NInputNumber,
                                                {
                                                    value: publish_year_value.value[1],
                                                    size: "small",
                                                    onUpdateValue: (value) => {
                                                        if (value < publish_year_value.value[0]) {
                                                            return;
                                                        }
                                                        publish_year_value.value = [publish_year_value.value[0], value];
                                                        filterRef.value.max_publish_year = value;
                                                    }
                                                }
                                            )
                                        ]
                                    )
                                ]
                            )
                        ),
                        h(
                            NFormItem,
                            {
                                label: "作者"
                            },
                            h(
                                NInput,
                                {
                                    value: filterRef.value.author,
                                    maxlength: 63,
                                    showCount: true,
                                    clearable: true,
                                    onUpdateValue: (value) => {
                                        if (value === "") {
                                            filterRef.value.author = null;
                                        } else {
                                            filterRef.value.author = value;
                                        }
                                    }
                                }
                            )
                        ),
                        h(
                            NFormItem,
                            {
                                label: "价格"
                            },
                            h(
                                NSpace,
                                {
                                    vertical: true
                                },
                                [
                                    h(
                                        NSlider,
                                        {
                                            value: price_value,
                                            range: true,
                                            min: 0,
                                            max: 1000,
                                            step: 0.1,
                                            onUpdateValue: ([minp, maxp]) => {
                                                if (maxp < minp) {
                                                    return;
                                                }
                                                price_value.value = [minp, maxp];
                                                filterRef.value.min_price = minp;
                                                filterRef.value.max_price = maxp;
                                            }
                                        }
                                    ),
                                    h(
                                        NSpace,
                                        [
                                            h(
                                                NInputNumber,
                                                {
                                                    value: price_value.value[0],
                                                    size: "small",
                                                    onUpdateValue: (value) => {
                                                        if (value > price_value.value[1]) {
                                                            return;
                                                        }
                                                        price_value.value = [value, price_value.value[1]];
                                                        filterRef.value.min_price = value;
                                                    }
                                                }
                                            ),
                                            h(
                                                NInputNumber,
                                                {
                                                    value: price_value.value[1],
                                                    size: "small",
                                                    onUpdateValue: (value) => {
                                                        if (value < price_value.value[0]) {
                                                            return;
                                                        }
                                                        price_value.value = [price_value.value[0], value];
                                                        filterRef.value.max_price = value;
                                                    }
                                                }
                                            )
                                        ]
                                    )
                                ]
                            )
                        ),
                        h(
                            NFormItem,
                            {
                                label: "排序列"
                            },
                            h(
                                NSelect,
                                {
                                    value: filterRef.value.sort_by,
                                    options: sortByOptions,
                                    onUpdateValue: (value) => {
                                        filterRef.value.sort_by = value;
                                    }
                                }
                            )
                        ),
                        h(
                            NFormItem,
                            {
                                label: "升降序"
                            },
                            h(
                                NSelect,
                                {
                                    value: filterRef.value.sort_order,
                                    options: sortOrderOptions,
                                    onUpdateValue: (value) => {
                                        filterRef.value.sort_order = value;
                                    }
                                }
                            )
                        )
                    ]
                ),
                // h('pre', JSON.stringify(filterRef.value, null, 2))
            ]
        ),
        positiveText: "确定",
        negativeText: "取消",
        onPositiveClick: () => {
            fetchBookData()
        }
    })
};
const handlerAdd = () => {
    const formValue = ref({
        category: null,
        title: null,
        press: null,
        publish_year: null,
        author: null,
        price: null,
        stock: null
    })
    const rules = {
        category: { required: true },
        title: { required: true },
        press: { required: true },
        publish_year: { required: true },
        author: { required: true },
        price: { required: true },
        stock: { required: true }
    }
    var categoryChanged = false;
    var titleChanged = false;
    var pressChanged = false;
    var publishYearChanged = false;
    var authorChanged = false;
    var priceChanged = false;
    var stockChanged = false;

    dialog.create({
        title: '添加图书',
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
                                label: "分类",
                                path: "category",
                                validationStatus: inputStringValidationStatus(formValue.value.category, categoryChanged),
                                feedback: createStringFeedback(formValue.value.category, categoryChanged, '分类')
                            },
                            h(
                                NInput,
                                {
                                    maxlength: 63,
                                    showCount: true,
                                    clearable: true,
                                    onUpdateValue: (value) => {
                                        formValue.value.category = value
                                        categoryChanged = true
                                    }
                                }
                            )
                        ),
                        h(
                            NFormItem,
                            {
                                label: "标题",
                                path: "title",
                                validationStatus: inputStringValidationStatus(formValue.value.title, titleChanged),
                                feedback: createStringFeedback(formValue.value.title, titleChanged, '标题')
                            },
                            h(
                                NInput,
                                {
                                    maxlength: 63,
                                    showCount: true,
                                    clearable: true,
                                    onUpdateValue: (value) => {
                                        formValue.value.title = value
                                        titleChanged = true
                                    }
                                }
                            )
                        ),
                        h(
                            NFormItem,
                            {
                                label: "出版社",
                                path: "press",
                                validationStatus: inputStringValidationStatus(formValue.value.press, pressChanged),
                                feedback: createStringFeedback(formValue.value.press, pressChanged, '出版社')
                            },
                            h(
                                NInput,
                                {
                                    maxlength: 63,
                                    showCount: true,
                                    clearable: true,
                                    onUpdateValue: (value) => {
                                        formValue.value.press = value
                                        pressChanged = true
                                    }
                                }
                            )
                        ),
                        h(
                            NFormItem,
                            {
                                label: "出版年份",
                                path: "publish_year",
                                validationStatus: inputNumberValidationStatus(formValue.value.publish_year, publishYearChanged),
                                feedback: createNumberFeedback(formValue.value.publish_year, publishYearChanged, '出版年份')
                            },
                            h(
                                NDatePicker,
                                {
                                    type: "year",
                                    clearable: true,
                                    style: "width: 340.62px",
                                    onClear: () => {
                                        formValue.value.publish_year = 0
                                        publishYearChanged = true
                                    },
                                    onUpdateValue: (value) => {
                                        if (value === null) {
                                            formValue.value.publish_year = 0
                                            publishYearChanged = true
                                            return
                                        }
                                        var year = new Date(value).getFullYear()
                                        formValue.value.publish_year = year
                                        publishYearChanged = true
                                    }
                                }
                            )
                        ),
                        h(
                            NFormItem,
                            {
                                label: "作者",
                                path: "author",
                                validationStatus: inputStringValidationStatus(formValue.value.author, authorChanged),
                                feedback: createStringFeedback(formValue.value.author, authorChanged, '作者')
                            },
                            h(
                                NInput,
                                {
                                    maxlength: 63,
                                    showCount: true,
                                    clearable: true,
                                    onUpdateValue: (value) => {
                                        formValue.value.author = value
                                        authorChanged = true
                                    }
                                }
                            )
                        ),
                        h(
                            NFormItem,
                            {
                                label: "价格",
                                path: "price",
                                validationStatus: inputNumberValidationStatus(formValue.value.price, priceChanged),
                                feedback: createNumberFeedback(formValue.value.price, priceChanged, '价格')
                            },
                            h(
                                NInputNumber,
                                {
                                    parse: parseCurrency,
                                    format: formatCurrency,
                                    clearable: true,
                                    precision: 2,
                                    min: 0,
                                    style: "width: 340.62px",
                                    onUpdateValue: (value) => {
                                        formValue.value.price = value
                                        priceChanged = true
                                    }
                                }
                            )
                        ),
                        h(
                            NFormItem,
                            {
                                label: "库存",
                                path: "stock",
                                validationStatus: inputNumberValidationStatus(formValue.value.stock, stockChanged),
                                feedback: createNumberFeedback(formValue.value.stock, stockChanged, '库存')
                            },
                            h(
                                NInputNumber,
                                {
                                    clearable: true,
                                    precision: 0,
                                    min: 0,
                                    style: "width: 340.62px",
                                    onUpdateValue: (value) => {
                                        formValue.value.stock = value
                                        stockChanged = true
                                    }
                                }
                            )
                        )
                    ]
                )
            ]
        ),
        positiveText: '确定',
        onPositiveClick: () => {
            if (inputStringValidationStatus(formValue.value.category, true) === 'error') {
                message.error('请输入分类')
                return false
            }
            if (inputStringValidationStatus(formValue.value.title, true) === 'error') {
                message.error('请输入标题')
                return false
            }
            if (inputStringValidationStatus(formValue.value.press, true) === 'error') {
                message.error('请输入出版社')
                return false
            }
            if (inputNumberValidationStatus(formValue.value.publish_year, true) === 'error') {
                message.error('请输入出版年份')
                return false
            }
            if (inputStringValidationStatus(formValue.value.author, true) === 'error') {
                message.error('请输入作者')
                return false
            }
            if (inputNumberValidationStatus(formValue.value.price, true) === 'error') {
                message.error('请输入价格')
                return false
            }
            if (inputNumberValidationStatus(formValue.value.stock, true) === 'error') {
                message.error('请输入库存')
                return false
            }
            axios.post('/book/create', formValue.value).then(res => {
                if (res.data.code === 0) {
                    message.success('添加成功')
                    fetchBookData()
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
const handlerBatchAdd = () => {
    const uploadJSON = ref(null)
    dialog.create({
        title: '批量添加',
        content: () => h(
            NCard,
            [
                h(
                    NH4,
                    "导入图书信息(JSON)"
                ),
                h(
                    NUpload,
                    {
                        directoryDnd: true,
                        defaultUpload: false,
                        onChange: (files) => {
                            if (files.fileList.length === 0) {
                                return
                            }
                            const file = files.fileList[0]
                            if (file.status === 'error') {
                                message.error('上传失败')
                                return
                            }
                            if (file.status === 'pending') {
                                const reader = new FileReader()
                                reader.onload = (e) => {
                                    const data = e.target.result
                                    try {
                                        uploadJSON.value = JSON.parse(data)
                                    } catch (e) {
                                        message.error('文件格式错误')
                                    }
                                }
                                reader.readAsText(file.file)
                            }
                        }
                    },
                    h(
                        NUploadDragger,
                        [
                            h(
                                'div',
                                {
                                    style: "margin-bottom: 12px"
                                },
                                h(
                                    NIcon,
                                    {
                                        size: 48,
                                        depth: 3
                                    },
                                    h(
                                        ArchiveIcon
                                    )
                                )
                            ),
                            h(
                                NText,
                                {
                                    style: "font-size: 16px"
                                },
                                "点击或者拖动文件到该区域来上传"
                            )
                        ]
                    )
                )
            ]
        ),
        positiveText: '确定',
        negativeText: '取消',
        onPositiveClick: () => {
            let uploadData = []
            if (uploadJSON.value.length == undefined) {
                uploadData.push(uploadJSON.value)
                uploadData[0].book_id = null
            } else {
                uploadData = uploadJSON.value
                uploadData.forEach((item) => {
                    item.book_id = null
                })
            }
            axios.post("/book/create/batch", uploadData).then(res => {
                if (res.data.code === 0) {
                    message.success('添加成功')
                    fetchBookData()
                } else {
                    message.error(res.data.msg)
                }
            }).catch(err => {
                message.error(err.response.data.msg)
            })
        }
    })
};
const handlerEdit = (row) => {
    const formValue = ref({
        book_id: row.book_id,
        category: row.category,
        title: row.title,
        press: row.press,
        publish_year: row.publish_year,
        author: row.author,
        price: parseFloat(row.price),
    })
    const rules = {
        category: { required: true },
        title: { required: true },
        press: { required: true },
        publish_year: { required: true },
        author: { required: true },
        price: { required: true }
    }

    console.log(formValue.value)
    dialog.create({
        title: '添加图书',
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
                                label: "分类",
                                path: "category",
                                validationStatus: inputStringValidationStatus(formValue.value.category, true),
                                feedback: createStringFeedback(formValue.value.category, true, '分类')
                            },
                            h(
                                NInput,
                                {
                                    value: formValue.value.category,
                                    maxlength: 63,
                                    showCount: true,
                                    clearable: true,
                                    onUpdateValue: (value) => {
                                        formValue.value.category = value
                                    }
                                }
                            )
                        ),
                        h(
                            NFormItem,
                            {
                                label: "标题",
                                path: "title",
                                validationStatus: inputStringValidationStatus(formValue.value.title, true),
                                feedback: createStringFeedback(formValue.value.title, true, '标题')
                            },
                            h(
                                NInput,
                                {
                                    value: formValue.value.title,
                                    maxlength: 63,
                                    showCount: true,
                                    clearable: true,
                                    onUpdateValue: (value) => {
                                        formValue.value.title = value
                                    }
                                }
                            )
                        ),
                        h(
                            NFormItem,
                            {
                                label: "出版社",
                                path: "press",
                                validationStatus: inputStringValidationStatus(formValue.value.press, true),
                                feedback: createStringFeedback(formValue.value.press, true, '出版社')
                            },
                            h(
                                NInput,
                                {
                                    value: formValue.value.press,
                                    maxlength: 63,
                                    showCount: true,
                                    clearable: true,
                                    onUpdateValue: (value) => {
                                        formValue.value.press = value
                                    }
                                }
                            )
                        ),
                        h(
                            NFormItem,
                            {
                                label: "出版年份",
                                path: "publish_year",
                                validationStatus: inputNumberValidationStatus(formValue.value.publish_year, true),
                                feedback: createNumberFeedback(formValue.value.publish_year, true, '出版年份')
                            },
                            h(
                                NDatePicker,
                                {
                                    value: formValue.value.publish_year === 0 ? null : new Date(formValue.value.publish_year, 0, 1),
                                    type: "year",
                                    clearable: true,
                                    style: "width: 340.62px",
                                    onClear: () => {
                                        formValue.value.publish_year = 0
                                    },
                                    onUpdateValue: (value) => {
                                        if (value === null) {
                                            formValue.value.publish_year = 0
                                            publishYearChanged = true
                                            return
                                        }
                                        var year = new Date(value).getFullYear()
                                        formValue.value.publish_year = year
                                    }
                                }
                            )
                        ),
                        h(
                            NFormItem,
                            {
                                label: "作者",
                                path: "author",
                                validationStatus: inputStringValidationStatus(formValue.value.author, true),
                                feedback: createStringFeedback(formValue.value.author, true, '作者')
                            },
                            h(
                                NInput,
                                {
                                    value: formValue.value.author,
                                    maxlength: 63,
                                    showCount: true,
                                    clearable: true,
                                    onUpdateValue: (value) => {
                                        formValue.value.author = value
                                    }
                                }
                            )
                        ),
                        h(
                            NFormItem,
                            {
                                label: "价格",
                                path: "price",
                                validationStatus: inputNumberValidationStatus(formValue.value.price, true),
                                feedback: createNumberFeedback(formValue.value.price, true, '价格')
                            },
                            h(
                                NInputNumber,
                                {
                                    value: formValue.value.price,
                                    parse: parseCurrency,
                                    format: formatCurrency,
                                    clearable: true,
                                    precision: 2,
                                    min: 0,
                                    style: "width: 340.62px",
                                    onUpdateValue: (value) => {
                                        formValue.value.price = value
                                    }
                                }
                            )
                        )
                    ]
                )
            ]
        ),
        positiveText: '确定',
        onPositiveClick: () => {
            if (inputStringValidationStatus(formValue.value.category, true) === 'error') {
                message.error('请输入分类')
                return false
            }
            if (inputStringValidationStatus(formValue.value.title, true) === 'error') {
                message.error('请输入标题')
                return false
            }
            if (inputStringValidationStatus(formValue.value.press, true) === 'error') {
                message.error('请输入出版社')
                return false
            }
            if (inputNumberValidationStatus(formValue.value.publish_year, true) === 'error') {
                message.error('请输入出版年份')
                return false
            }
            if (inputStringValidationStatus(formValue.value.author, true) === 'error') {
                message.error('请输入作者')
                return false
            }
            if (inputNumberValidationStatus(formValue.value.price, true) === 'error') {
                message.error('请输入价格')
                return false
            }
            axios.put('/book/update', formValue.value).then(res => {
                if (res.data.code === 0) {
                    message.success('修改成功')
                    fetchBookData()
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
const handlerIncStock = (row) => {
    const formValue = ref({
        option: null,
        delta: null
    })
    const rules = {
        option: { required: true },
        delta: { required: true }
    }
    var optionChanged = false;
    var deltaChanged = false;

    dialog.create({
        title: '修改库存',
        content: () => h(
            NCard,
            [
                h(
                    NForm,
                    {
                        model: formValue,
                        rules: rules
                    },
                    [
                        h(
                            NFormItem,
                            {
                                label: "操作",
                                path: "option",
                                validationStatus: inputStringValidationStatus(formValue.value.option, optionChanged),
                                feedback: createStringFeedback(formValue.value.option, optionChanged, '操作')
                            },
                            h(
                                NRadioGroup,
                                {
                                    value: formValue.value.option,
                                    onUpdateValue: (value) => {
                                        formValue.value.option = value
                                        optionChanged = true
                                    }
                                },
                                [
                                    h(
                                        NRadioButton,
                                        {
                                            label: "增加",
                                            value: "inc",
                                        }
                                    ),
                                    h(
                                        NRadioButton,
                                        {
                                            label: "减少",
                                            value: "dec",
                                        }
                                    )
                                ]
                            )
                        ),
                        h(
                            NFormItem,
                            {
                                label: "数量",
                                path: "delta",
                                validationStatus: inputNumberValidationStatus(formValue.value.delta, deltaChanged),
                                feedback: createNumberFeedback(formValue.value.delta, deltaChanged, '数量')
                            },
                            h(
                                NInputNumber,
                                {
                                    clearable: true,
                                    precision: 0,
                                    min: 0,
                                    style: "width: 340.62px",
                                    onUpdateValue: (value) => {
                                        formValue.value.delta = value
                                        deltaChanged = true
                                    }
                                }
                            )
                        )
                    ]
                ),
            ]
        ),
        positiveText: '确认',
        negativeText: '取消',
        onPositiveClick: () => {
            if (inputStringValidationStatus(formValue.value.option, optionChanged) === 'error') {
                message.error('请选择操作')
                return false
            }
            if (inputNumberValidationStatus(formValue.value.delta, deltaChanged) === 'error') {
                message.error('请输入数量')
                return false
            }
            axios.put('/book/stock/update', {
                book_id: row.book_id,
                option: formValue.value.option,
                delta: formValue.value.delta
            }).then(res => {
                if (res.data.code === 0) {
                    message.success('修改成功')
                    fetchBookData()
                } else {
                    message.error(res.data.msg)
                }
            }).catch(err => {
                message.error(err.response.data.msg)
            })
        }
    })
}
const handlerDownload = () => {
    spinShow.value = true;
    console.log(spinShow.value)
    const book = JSON.parse(JSON.stringify(dataRef.value));
    for (var i = 0; i < book.length; i++) {
        book[i].price = parseFloat(book[i].price);
    }
    const data = JSON.stringify(book, null, 2);
    const blob = new Blob([data], { type: 'text/plain' });
    FileSaver.saveAs(blob, 'book.json');
    spinShow.value = false;
    console.log(spinShow.value)
}
const spinShow = ref(false);

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
        title: '库存',
        key: 'stock'
    },
    {
        title: '编辑图书',
        key: 'edit',
        render: (row) => {
            return h(
                NButton,
                {
                    tertiary: true,
                    size: 'small',
                    onClick: () => {
                        handlerEdit(row)
                    }
                },
                '编辑'
            )
        }
    },
    {
        title: '修改库存',
        key: 'incStock',
        render: (row) => {
            return h(
                NButton,
                {
                    tertiary: true,
                    size: 'small',
                    onClick: () => {
                        handlerIncStock(row)
                    }
                },
                '修改'
            )
        }
    },
    {
        title: '删除',
        key: 'remove',
        render: (row) => {
            return h(
                NPopconfirm,
                {
                    onPositiveClick: () => {
                        axios.delete("/book/remove?bid=" + row.book_id).then(res => {
                            message.success("删除成功！")
                            countRef.value = 0
                            loadingRef.value = true
                            fetchBookData()
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
</script>

<template>
    <n-card>
        <n-data-table :single-line="false" :columns="columns" :data="dataRef" :loading="loadingRef"
            :pagination="pagination" />
        <n-space justify="center">
            <n-button @click="handlerFilter">设置筛选器</n-button>
            <n-button @click="handlerAdd">添加图书</n-button>
            <n-button @click="handlerBatchAdd">批量添加</n-button>
            <n-spin :show="spinShow">
                <n-button @click="handlerDownload">导出JSON</n-button>
            </n-spin>
        </n-space>
    </n-card>
</template>

<script>
import axios from "axios";
import { h, ref } from "vue";
axios.defaults.baseURL = 'http://api.yasyakarasu.tech/lib';

var filterRef = ref({
    category: null,
    title: null,
    press: null,
    min_publish_year: null,
    max_publish_year: null,
    author: null,
    min_price: null,
    max_price: null,
    sort_by: null,
    sort_order: null
})
var publish_year_value = ref([0, 3000])
var price_value = ref([0, 1000])
var loadingRef = ref(true)
var dataRef = ref([])
var countRef = ref(0)

function fetchBookData() {
    loadingRef.value = true
    axios.post('/book/list', filterRef.value).then(res => {
        countRef.value = res.data.data.count
        if (res.data.data.count === 0) {
            dataRef.value = []
        } else {
            res.data.data.results.forEach((item) => {
                item.price = item.price.toFixed(2)
            })
            dataRef.value = res.data.data.results
        }
        loadingRef.value = false
    }).catch(err => {
        console.log(err.data.msg)
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
    created() {
        fetchBookData()
    },
    data() {
        return {
            loadingRef,
            dataRef,
            pagination
        }
    }
}
</script>