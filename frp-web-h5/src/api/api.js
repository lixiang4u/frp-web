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
    getVhosts: () => {
        return request.get(
            '/api/vhosts'
        )
    },
    removeVhost: (vhostId) => {
        return request.delete(
            `/api/vhost/${vhostId}`
        )
    },
    reloadVhost: () => {
        return request.post(
            `/api/frp/reload`
        )
    },
    getUsePort: () => {
        return request.post(
            `/api/use-port-check`
        )
    },

}
