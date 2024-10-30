$(document).ready(function() {
    const API_BASE_URL = window._config.API_BASE_URL || 'http://localhost:8080';

    $('#new-customer-form').submit(function(e) {
        e.preventDefault();

        const newCustomer = {
            name: $('#name').val(),
            email: $('#email').val(),
            password: $('#password').val(),
            gender: $('#gender').val()
        };

        $.ajax({
            url: `${API_BASE_URL}/customers`,
            method: 'POST',
            contentType: 'application/json',
            data: JSON.stringify(newCustomer),
            success: function() {
                alert('客戶新增成功');
                window.location.href = 'index.html';
            },
            error: function() {
                alert('新增失敗');
            }
        });
    });
});
