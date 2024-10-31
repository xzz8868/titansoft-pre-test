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

    // 監聽修改密碼的勾選框
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

    // 提交修改
    $('#customer-form').submit(function(e) {
        e.preventDefault();

        const updatedCustomer = {
            name: $('#name').val(),
            email: $('#email').val(),
            gender: $('#gender').val()
        };

        // 更新基本資料
        $.ajax({
            url: `${API_BASE_URL}/customers/${customerId}`,
            method: 'PUT',
            contentType: 'application/json',
            data: JSON.stringify(updatedCustomer),
            success: function() {
                // 如果勾選了修改密碼，發送密碼更新請求
                if ($('#change-password').is(':checked')) {
                    const newPassword = $('#new-password').val();
                    const confirmPassword = $('#confirm-password').val();

                    // 密碼驗證
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

                    $.ajax({
                        url: `${API_BASE_URL}/customers/password/${customerId}`,
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
