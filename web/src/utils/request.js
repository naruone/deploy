import axios from 'axios';
import {Loading, Message} from 'element-ui';
import {store} from '../store/index'

const service = axios.create({
    baseURL: process.env.VUE_APP_DEPLOY_BASE_API,
    timeout: process.env.VUE_APP_DEPLOY_TIMEOUT //毫秒
})
let activeAxios = 0 //当前请求个数, 请求一次+1, 当=0时取消遮罩
let loadingInstance = null; //遮罩层
let timer = null;   //loading定时器
const showLoading = () => {
    activeAxios++
    if (timer) {
        clearTimeout(timer)
    }
    timer = setTimeout(() => {
        if (activeAxios > 0) {
            loadingInstance = Loading.service({fullscreen: true, background: 'rgba(0, 0, 0, 0.3)'})
        }
    }, 400);
}

const closeLoading = () => {
    activeAxios--
    if (activeAxios <= 0) {
        clearTimeout(timer)
        loadingInstance && loadingInstance.close()
    }
}
//http request 拦截器
service.interceptors.request.use(config => {
        showLoading()
        const token = store.getters['user/token']
        config.data = JSON.stringify(config.data);
        config.headers = {
            'Content-Type': 'application/json',
            'x-token': token
        }
        return config;
    }, error => {
        closeLoading()
        Message({
            showClose: true,
            message: error,
            type: 'error'
        })
        return Promise.reject(error);
    }
);

//http response 拦截器
service.interceptors.response.use(response => {
        closeLoading()
        if (response.data.code === 0 || response.headers.success === "true") {
            return response.data
        } else {
            Message({
                showClose: true,
                message: response.data.msg || decodeURI(response.headers.msg),
                type: 'error',
            })
            if (response.data.data && response.data.data.reload) {
                store.commit('user/LoginOut')
            }
            return Promise.reject(response.data.msg)
        }
    }, error => {
        closeLoading()
        Message({
            showClose: true,
            message: error,
            type: 'error'
        })
        return Promise.reject(error)
    }
)

export default service