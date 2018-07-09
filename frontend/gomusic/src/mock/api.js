import axios from 'axios'

let host = 'http://localhost:8080'

export const getplaylistrecom = params => { return axios.get(`${host}/api/test/getdatabytime`) }
