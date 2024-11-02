$(document).ready(function() {
    const GENERATOR_BASE_URL = window._config.GENERATOR_BASE_URL || 'http://localhost:8081';

    $('#generate-customer-form').submit(function(e) {
        e.preventDefault();

        // 從表單中獲取輸入值
        const num = $('#num').val();

        // 構建帶參數的 URL
        const urlWithParams = `${GENERATOR_BASE_URL}/generate/customer?num=${encodeURIComponent(num)}`;

        $.ajax({
            url: urlWithParams,
            method: 'POST',
            success: function() {
                alert('資料產生成功');
                window.location.href = 'index.html';
            },
            error: function() {
                alert('資料產生失敗');
            }
        });
    });
});
