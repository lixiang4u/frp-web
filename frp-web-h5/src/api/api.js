import request from "../utils/request";

export default {
    getConfig: () => {
        return request.get(
            '/api/config'
        )
    },
    newVhost: (data) => {
        return request.post(
            '/api/vhost', data
        )
    },
    getVhosts: (params) => {
        return request.get(
            '/api/vhosts', params
        )
    },

}
