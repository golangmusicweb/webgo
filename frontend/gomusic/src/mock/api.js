import axios from 'axios'

// let host = 'http://localhost'
// export const getplaylistrecom = params => { return axios.get(`${host}/api/test/getdatabytime`) }
var instance = axios.create({
  baseURL: 'http://localhost'
})

// 在实例已创建后修改默认值
instance.defaults.headers.common['Authorization'] = 'JWT eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6NDksIlVzZXJuYW1lIjoxODYwMDIzOTQ5NSwiUm9sZSI6ImNpdmlsaWFucyIsImV4cCI6MTUzMTQ2MDc3NywiaXNzIjoiZG9uZ3hpYW95aSIsIm5iZiI6MTUzMTM3MzM3N30.eMPYxpvTlk2yzpXkqnDaPZjZG0vRQDRphxONY4TLi70'
export const getplaylistrecom = params => { return instance.get(`/api/test/getdatabytime`) }
