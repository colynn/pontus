import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/login',
    method: 'post',
    data: data
  })
}

export function refreshtoken(token) {
  return request({
    url: '/refresh_token',
    method: 'get',
    params: { 'token': token }
  })
}

export function getInfo(token) {
  return request({
    url: '/api/v1/userinfo',
    method: 'get',
    params: { token }
  })
}

export function logout() {
  return request({
    url: '/api/v1/logout',
    method: 'post'
  })
}

// 查询用户个人信息
export function getUserProfile() {
  return request({
    url: '/api/v1/user/profile',
    method: 'get'
  })
}

export function updateUserProfile(data) {
  return request({
    url: '/api/v1/user/profile',
    method: 'put',
    data
  })
}

export function updateUserPwd(oldPassword, newPassword) {
  const data = {
    oldPassword,
    newPassword
  }
  return request({
    url: '/api/v1/logout',
    method: 'put',
    data: data
  })
}

export function updateUser(userID, data) {
  return request({
    url: '/api/v1/users/' + userID,
    method: 'put',
    data
  })
}

// admin user  //

// 查询用户列表
export function listUser(query) {
  return request({
    url: '/api/v1/users',
    method: 'get',
    params: query
  })
}

// 查询用户详细
export function getUser(userID) {
  return request({
    url: '/api/v1/users/' + userID,
    method: 'get'
  })
}

// 新增用户
export function addUser(data) {
  return request({
    url: '/api/v1/users',
    method: 'post',
    data: data
  })
}

// 删除用户
export function delUser(userId) {
  return request({
    url: '/api/v1/users/' + userId,
    method: 'delete'
  })
}
