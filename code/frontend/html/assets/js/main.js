$(document).ready(function() {
    // 獲取 API 基礎 URL
    const SERVER_BASE_URL = window._config.SERVER_BASE_URL || 'http://localhost:8080';

    // 性別對照表
    const genderMap = {
        'male': '男性',
        'female': '女性',
        'other': '其他'
    };

    // 從後端獲取客戶列表
    $.ajax({
        url: `${SERVER_BASE_URL}/customers`,
        method: 'GET',
        success: function(customers) {
            customers.forEach(function(customer) {

                let totalAmount = customer.total_transaction_amount || 0;

                // 將客戶資料插入表格
                $('#customer-table-body').append(`
                    <tr>
                        <td>${customer.name}</td>
                        <td>${customer.email}</td>
                        <td>${genderMap[customer.gender]}</td>
                        <td>${totalAmount.toFixed(2)}</td>
                        <td>
                            <a href="customer.html?id=${customer.id}" class="btn btn-sm btn-info">查看/編輯</a>
                            <a href="transactions.html?id=${customer.id}" class="btn btn-sm btn-secondary">查看交易</a>
                        </td>
                    </tr>
                `);
            });
        },
        error: function(xhr, status, error) {
            console.error('获取客户列表失败:', error);
        }
    });
});
