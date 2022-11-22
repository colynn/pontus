import request from '@/utils/request'

export function CloudAssetList(query) {
  return request({
    url: '/api/v1/assets/clouds',
    method: 'get',
    params: query
  })
}
export function AssetStatistics(query) {
  return request({
    url: '/api/v1/assets/statistics',
    method: 'get',
    params: query
  })
}
/* ---------------  */
