module.exports = {
    devServer: {
        proxy: {
            "/api": {
                target: 'http://127.0.0.1:8085',
                ws: false,
                pathRewrite: {
                    '^/api': '/'
                }
            }
        }
    }
}