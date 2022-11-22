import request from '@/utils/request'

export function PCList(query) {
  return request({
    url: '/api/v1/assets/pc/list',
    method: 'get',
    params: query
  })
}

export function PhysicalList(query) {
  return request({
    url: '/api/v1/assets/physical/list',
    method: 'get',
    params: query
  })
}
export function pcDetail(instanceID) {
  return request({
    url: '/api/v1/assets/tangibles/' + instanceID,
    method: 'get'
  })
}
// pc 的变更日志
export function pcAuditLog(instanceID, query) {
  return request({
    url: '/api/v1/assets/tangibles/' + instanceID + '/audits',
    method: 'get',
    params: query
  })
}

//  pcUpdate 编辑信息
export function pcUpdate(instanceID, data) {
  return request({
    url: '/api/v1/assets/tangibles/' + instanceID,
    method: 'put',
    data: data
  })
}

//  pcDelete 编辑信息
export function pcDelete(instanceID) {
  return request({
    url: '/api/v1/assets/tangibles/' + instanceID,
    method: 'delete'
  })
}
