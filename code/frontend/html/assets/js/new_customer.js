$(document).ready(function() {
    const SERVER_BASE_URL = window._config.SERVER_BASE_URL || 'http://localhost:8080';

    // Form submission event for creating a new customer
    $('#new-customer-form').submit(function(e) {
        e.preventDefault();

        const password = $('#password').val();
        const confirmPassword = $('#confirm-password').val();

        // Password validation
        if (password.length < 8) {
            alert('密碼至少需要8個字元');
            return;
        }

        if (password !== confirmPassword) {
            alert('密碼和確認密碼不一致');
            return;
        }

        // Construct new customer data object
        const newCustomer = {
            name: $('#name').val(),
            email: $('#email').val(),
            password: password,
            gender: $('#gender').val()
        };

        // Send new customer data to the server via AJAX
        $.ajax({
            url: `${SERVER_BASE_URL}/customers`,
            method: 'POST',
            contentType: 'application/json',
            data: JSON.stringify(newCustomer),
            success: function() {
                alert('客戶新增成功');
                window.location.href = 'index.html'; // Redirect on success
            },
            error: function() {
                alert('新增失敗');
            }
        });
    });
});
