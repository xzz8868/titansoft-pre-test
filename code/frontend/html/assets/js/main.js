$(document).ready(function() {
    const SERVER_BASE_URL = window._config.SERVER_BASE_URL || 'http://localhost:8080';

    // Gender display mapping
    const genderMap = {
        'male': '男性',
        'female': '女性',
        'other': '其他'
    };

    // Fetch customer list from backend
    $.ajax({
        url: `${SERVER_BASE_URL}/customers`,
        method: 'GET',
        success: function(customers) {
            customers.forEach(function(customer) {
                let totalAmount = customer.total_transaction_amount || 0;

                // Insert customer data into table
                $('#customer-table-body').append(
                    `<tr>
                        <td>${customer.name}</td>
                        <td>${customer.email}</td>
                        <td>${genderMap[customer.gender]}</td>
                        <td>${totalAmount.toFixed(2)}</td>
                        <td>
                            <a href="customer.html?id=${customer.id}" class="btn btn-sm btn-info">查看/編輯</a>
                            <a href="transactions.html?id=${customer.id}" class="btn btn-sm btn-secondary">查看交易</a>
                        </td>
                    </tr>`
                );
            });
            // Update total customer count
            $('#customer-count').text('客戶總數：' + customers.length);
        },
        error: function(error) {
            console.error('Failed to fetch customer list:', error);
        }
    });

    // Handle reset button click
    $('#reset_button').click(function() {
        if (confirm('確定要清除所有資料嗎？')) {
            $.ajax({
                url: `${SERVER_BASE_URL}/customers/reset`,
                method: 'DELETE',
                success: function() {
                    alert('資料清除成功!');
                    window.location.href = 'index.html';
                },
                error: function() {
                    alert('資料清除失敗!');
                }
            });
        }
    });
});
