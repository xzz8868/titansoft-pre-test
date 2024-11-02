$(document).ready(function() {
    const SERVER_BASE_URL = window._config.SERVER_BASE_URL || 'http://localhost:8080';

    $('#new-customer-form').submit(function(e) {
        e.preventDefault();

        const password = $('#password').val();
        const confirmPassword = $('#confirm-password').val();

        // 密碼驗證
        if (password.length < 8) {
            alert('密碼至少需要8個字元');
            return;
        }

        if (password !== confirmPassword) {
            alert('密碼和確認密碼不一致');
            return;
        }

        const newCustomer = {
            name: $('#name').val(),
            email: $('#email').val(),
            password: password,
            gender: $('#gender').val()
        };

        $.ajax({
            url: `${SERVER_BASE_URL}/customers`,
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
