export default {
    data() {
        return {
            currentPage: 1,
            total: 10,
            pageSize: 10,
            tableData: [],
            searchForm: {
                sCondition: '',
                sValue: ''
            }
        }
    },
    methods: {
        handleSizeChange(size) {
            this.pageSize = size
            this.getTableData().catch(() => {
            })
        },
        handleCurrentChange(cPage) {
            this.currentPage = cPage
            this.getTableData().catch(() => {
            })
        },
        async getTableData(currentPage = this.currentPage, pageSize = this.pageSize) {
            const table = await this.getList({currentPage, pageSize, ...this.searchForm})
            this.tableData = table.data.list
            this.total = table.data.total
            this.currentPage = table.data.currentPage
            this.pageSize = table.data.pageSize
        }
    }
}