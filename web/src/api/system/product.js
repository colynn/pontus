// 查询用户列表
import request from '@/utils/request'

export function listProduct(query) {
  return request({
    url: '/api/v1/products',
    method: 'get',
    params: query
  })
}

export function listProductByDeptID(query) {
  return request({
    url: '/api/v1/productList',
    method: 'get',
    params: query
  })
}

export function getProduct(productID) {
  return request({
    url: '/api/v1/products/' + productID,
    method: 'get'
  })
}

export function delProduct(productID) {
  return request({
    url: '/api/v1/products/' + productID,
    method: 'delete'
  })
}

export function updateProduct(productID, data) {
  return request({
    url: '/api/v1/products/' + productID,
    method: 'put',
    data: data
  })
}

export function addProduct(data) {
  return request({
    url: '/api/v1/products',
    method: 'post',
    data: data
  })
}
