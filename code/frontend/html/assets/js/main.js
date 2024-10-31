$(document).ready(function() {
    // 獲取 API 基礎 URL
    const API_BASE_URL = window._config.API_BASE_URL || 'http://localhost:8080';

     // 性別對照表
    const genderMap = {
        'male': '男性',
        'female': '女性',
        'other': '其他'
    };

    // 從後端獲取客戶列表
    $.ajax({
        url: `${API_BASE_URL}/customers`,
        method: 'GET',
        success: function(customers) {
            customers.forEach(function(customer) {
                // 對於每個客戶，獲取其過去一年的交易總額
                $.ajax({
                    url: `${API_BASE_URL}/customers/${customer.id}/transactions?from=${getOneYearAgoDate()}&to=${getTodayDate()}`,
                    method: 'GET',
                    success: function(transactions) {
                        let totalAmount = transactions.reduce((sum, txn) => sum + txn.amount, 0);

                        // 將客戶資料插入表格
                        $('#customer-table-body').append(`
                            <tr>
                                <td>${customer.name}</td>
                                <td>${customer.email}</td>
                                <td>${genderMap[customer.gender]}</td>
                                <td>${new Date(customer.registration_time).toLocaleDateString()}</td>
                                <td>${totalAmount.toFixed(2)}</td>
                                <td>
                                    <a href="customer.html?id=${customer.id}" class="btn btn-sm btn-info">查看/編輯</a>
                                    <a href="transactions.html?id=${customer.id}" class="btn btn-sm btn-secondary">查看交易</a>
                                </td>
                            </tr>
                        `);
                    }
                });
            });
        }
    });

    function getOneYearAgoDate() {
        let date = new Date();
        date.setFullYear(date.getFullYear() - 1);
        return date.toISOString();
    }

    function getTodayDate() {
        let date = new Date();
        return date.toISOString();
    }
});
