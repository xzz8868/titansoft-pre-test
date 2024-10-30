$(document).ready(function() {
    const API_BASE_URL = window._config.API_BASE_URL || 'http://localhost:8080';
    
    const urlParams = new URLSearchParams(window.location.search);
    const customerId = urlParams.get('id');
    $('#customer-id').val(customerId);

    // 預設為過去一個月
    let fromDate = new Date();
    fromDate.setMonth(fromDate.getMonth() - 1);
    $('#from-date').val(formatDate(fromDate));

    let toDate = new Date();
    $('#to-date').val(formatDate(toDate));

    // 初始載入交易紀錄
    loadTransactions();

    // 篩選表單提交事件
    $('#filter-form').submit(function(e) {
        e.preventDefault();
        loadTransactions();
    });

    function loadTransactions() {
        let from = $('#from-date').val();
        let to = $('#to-date').val();

        $.ajax({
            url: `${API_BASE_URL}/customers/${customerId}/transactions?from=${from}&to=${to}`,
            method: 'GET',
            success: function(transactions) {
                $('#transaction-table-body').empty();
                transactions.forEach(function(txn) {
                    $('#transaction-table-body').append(`
                        <tr>
                            <td>${new Date(txn.time).toLocaleString()}</td>
                            <td>${txn.amount.toFixed(2)}</td>
                            <td>${txn.sequence}</td>
                        </tr>
                    `);
                });
            },
            error: function() {
                alert('無法獲取交易紀錄');
            }
        });
    }

    function formatDate(date) {
        let month = '' + (date.getMonth() + 1);
        let day = '' + date.getDate();
        let year = date.getFullYear();

        if (month.length < 2)
            month = '0' + month;
        if (day.length < 2)
            day = '0' + day;

        return [year, month, day].join('-');
    }
});
