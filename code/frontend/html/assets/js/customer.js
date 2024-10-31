$(document).ready(function() {
    const API_BASE_URL = window._config.API_BASE_URL || 'http://localhost:8080';

    const urlParams = new URLSearchParams(window.location.search);
    const customerId = urlParams.get('id');

    // 獲取客戶詳細資訊
    $.ajax({
        url: `${API_BASE_URL}/customers/${customerId}`,
        method: 'GET',
        success: function(customer) {
            $('#customer-id').val(customer.id);
            $('#name').val(customer.name);
            $('#email').val(customer.email);
            $('#gender').val(customer.gender);
        },
        error: function() {
            alert('無法獲取客戶資料');
        }
    });

    // 提交修改
    $('#customer-form').submit(function(e) {
        e.preventDefault();

        const updatedCustomer = {
            name: $('#name').val(),
            email: $('#email').val(),
            gender: $('#gender').val()
        };

        $.ajax({
            url: `${API_BASE_URL}/customers/${customerId}`,
            method: 'PUT',
            contentType: 'application/json',
            data: JSON.stringify(updatedCustomer),
            success: function() {
                alert('客戶資料已更新');
                window.location.href = 'index.html';
            },
            error: function() {
                alert('更新失敗');
            }
        });
    });
});
