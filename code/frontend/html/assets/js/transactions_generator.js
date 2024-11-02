$(document).ready(function() {
    const GENERATOR_BASE_URL = window._config.GENERATOR_BASE_URL || 'http://localhost:8081';

    $('#generate-transactions-form').submit(function(e) {
        e.preventDefault();

        // 從表單中獲取輸入值
        const transactionsNum = $('#transactions_num').val();
        const customersNum = $('#customers_num').val();

        // 構建帶參數的 URL
        const urlWithParams = `${GENERATOR_BASE_URL}/generate/transactions?transactions_num=${encodeURIComponent(transactionsNum)}&customers_num=${encodeURIComponent(customersNum)}`;

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
