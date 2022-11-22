// 查询用户列表
import request from '@/utils/request'

export function listRole(query) {
  return request({
    url: '/api/v1/roles',
    method: 'get',
    params: query
  })
}

export function listAllRole() {
  return request({
    url: '/api/v1/allroles',
    method: 'get'
  })
}

export function getRole(roleID) {
  return request({
    url: '/api/v1/roles/' + roleID,
    method: 'get'
  })
}

export function delRole(roleID) {
  return request({
    url: '/api/v1/roles/' + roleID,
    method: 'delete'
  })
}

export function updateRole(roleID, data) {
  return request({
    url: '/api/v1/roles/' + roleID,
    method: 'put',
    data: data
  })
}

// 角色状态修改
export function changeRoleStatus(roleID, status) {
  const data = {
    roleID,
    status
  }
  return request({
    url: '/api/v1/roles/' + roleID,
    method: 'put',
    data: data
  })
}

export function addRole(data) {
  return request({
    url: '/api/v1/roles',
    method: 'post',
    data: data
  })
}

export function getRoutes() {
  return request({
    url: '/api/v1/role/menus',
    method: 'get'
  })
}

// 根据角色ID查询部门树结构
export function roleDeptTreeselect(roleId) {
  return request({
    url: '/api/v1/roleDeptTreeselect/' + roleId,
    method: 'get'
  })
}