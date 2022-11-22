// 添加日期范围
export function addDateRange(params, dateRange) {
  var search = params
  search.start_date = ''
  search.end_date = ''
  if (dateRange != null && dateRange !== '') {
    search.start_date = dateRange[0]
    search.end_date = dateRange[1]
  }
  return search
}

// 日期-时间格式化
/**
 * Parse the time to string
 * @param {(Object|string|number)} time
 * @param {string} cFormat
 * @returns {string}
 */
export function parseTime(time, pattern) {
  if (arguments.length === 0 || !time) {
    return null
  }
  const format = pattern || '{y}-{m}-{d} {h}:{i}:{s}'
  let date
  if (typeof time === 'object') {
    date = time
  } else {
    if ((typeof time === 'string') && (/^[0-9]+$/.test(time))) {
      time = parseInt(time)
    }
    if ((typeof time === 'number') && (time.toString().length === 10)) {
      time = time * 1000
    }
    date = new Date(time)
  }
  const formatObj = {
    y: date.getFullYear(),
    m: date.getMonth() + 1,
    d: date.getDate(),
    h: date.getHours(),
    i: date.getMinutes(),
    s: date.getSeconds(),
    a: date.getDay()
  }
  const time_str = format.replace(/{(y|m|d|h|i|s|a)+}/g, (result, key) => {
    let value = formatObj[key]
    // Note: getDay() returns 0 on Sunday
    if (key === 'a') { return ['日', '一', '二', '三', '四', '五', '六'][value] }
    if (result.length > 0 && value < 10) {
      value = '0' + value
    }
    return value || 0
  })
  return time_str
}

// 日期格式化
/**
 * Parse the time to string
 * @param {(Object|string|number)} time
 * @param {string} cFormat
 * @returns {string}
 */
export function parseDate(time, pattern) {
  if (arguments.length === 0 || !time) {
    return null
  }
  const format = pattern || '{y}-{m}-{d}'
  let date
  if (typeof time === 'object') {
    date = time
  } else {
    if ((typeof time === 'string') && (/^[0-9]+$/.test(time))) {
      time = parseInt(time)
    }
    if ((typeof time === 'number') && (time.toString().length === 10)) {
      time = time * 1000
    }
    date = new Date(time)
  }
  const formatObj = {
    y: date.getFullYear(),
    m: date.getMonth() + 1,
    d: date.getDate(),
    a: date.getDay()
  }
  const time_str = format.replace(/{(y|m|d|h|i|s|a)+}/g, (result, key) => {
    let value = formatObj[key]
    // Note: getDay() returns 0 on Sunday
    if (key === 'a') { return ['日', '一', '二', '三', '四', '五', '六'][value] }
    if (result.length > 0 && value < 10) {
      value = '0' + value
    }
    return value || 0
  })
  return time_str
}

// 表单重置
export function resetForm(refName) {
  if (this.$refs[refName]) {
    this.$refs[refName].resetFields()
  }
}

// 回显数据字典
export function selectDictLabel(datas, value) {
  var actions = []
  Object.keys(datas).map((key) => {
    if (datas[key].dictValue === ('' + value)) {
      actions.push(datas[key].dictLabel)
      return false
    }
  })
  return actions.join('')
}

// transform to form-data
export function transformRequest(data) {
  let ret = ''
  for (const it in data) {
    ret += encodeURIComponent(it) + '=' + encodeURIComponent(data[it]) + '&'
  }
  return ret
}
