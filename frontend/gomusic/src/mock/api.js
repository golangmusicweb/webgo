import axios from 'axios'

let host = 'http://localhost:80'

export const getplaylistrecom = params => { return axios.get(`${host}/api/test/getdatabytime`) }
