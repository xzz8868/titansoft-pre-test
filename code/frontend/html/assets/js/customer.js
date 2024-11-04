$(document).ready(function() {
    const SERVER_BASE_URL = window._config.SERVER_BASE_URL || 'http://localhost:8080';

    const urlParams = new URLSearchParams(window.location.search);
    const customerId = urlParams.get('id');

    // Fetch customer details and populate form fields
    $.ajax({
        url: `${SERVER_BASE_URL}/customers/${customerId}`,
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

    // Toggle password fields visibility and requirement based on checkbox
    $('#change-password').change(function() {
        if ($(this).is(':checked')) {
            $('#password-fields').show();
            $('#new-password').attr('required', true);
            $('#confirm-password').attr('required', true);
        } else {
            $('#password-fields').hide();
            $('#new-password').removeAttr('required').val('');
            $('#confirm-password').removeAttr('required').val('');
        }
    });

    // Handle form submission for updating customer details
    $('#customer-form').submit(function(e) {
        e.preventDefault();

        const updatedCustomer = {
            name: $('#name').val(),
            email: $('#email').val(),
            gender: $('#gender').val()
        };

        // Send update request for customer details
        $.ajax({
            url: `${SERVER_BASE_URL}/customers/${customerId}`,
            method: 'PUT',
            contentType: 'application/json',
            data: JSON.stringify(updatedCustomer),
            success: function() {
                // If password change is requested, validate and send password update request
                if ($('#change-password').is(':checked')) {
                    const newPassword = $('#new-password').val();
                    const confirmPassword = $('#confirm-password').val();

                    // Password validation
                    if (newPassword.length < 8) {
                        alert('密碼至少需要8個字元');
                        return;
                    }

                    if (newPassword !== confirmPassword) {
                        alert('新密碼和確認密碼不一致');
                        return;
                    }

                    const passwordData = {
                        password: newPassword,
                        confirm_password: confirmPassword
                    };

                    // Send update request for customer password
                    $.ajax({
                        url: `${SERVER_BASE_URL}/customers/password/${customerId}`,
                        method: 'PUT',
                        contentType: 'application/json',
                        data: JSON.stringify(passwordData),
                        success: function() {
                            alert('客戶資料和密碼已更新');
                            window.location.href = 'index.html';
                        },
                        error: function() {
                            alert('密碼更新失敗');
                        }
                    });
                } else {
                    alert('客戶資料已更新');
                    window.location.href = 'index.html';
                }
            },
            error: function() {
                alert('更新失敗');
            }
        });
    });
});
