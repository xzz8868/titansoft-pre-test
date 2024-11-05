$(document).ready(function () {
    const GENERATOR_BASE_URL = window._config.GENERATOR_BASE_URL || 'http://localhost:8081';

    // Handle form submission for generating transactions
    $('#generate-transactions-form').submit(function (e) {
        e.preventDefault(); // Prevent default form submission

        // Get input values from the form
        const transactionsNum = $('#transactions_num').val();
        const customersNum = $('#customers_num').val();

        // Verify if the input value exceeds 5000
        if (transactionsNum > 5000) {
            alert('單次最多只能產生5000筆資料');
            return;
        }

        // Construct URL with query parameters
        const urlWithParams = `${GENERATOR_BASE_URL}/generate/transactions?transactions_num=${encodeURIComponent(transactionsNum)}&customers_num=${encodeURIComponent(customersNum)}`;

        // Send AJAX request to generate transactions
        $.ajax({
            url: urlWithParams,
            method: 'POST',
            beforeSend: function() {
                // Show waiting animation
                $('#loading-spinner').show();
            },
            complete: function() {
                // Hide waiting animation
                $('#loading-spinner').hide();
            },
            success: function () {
                alert('資料產生成功'); // Alert success message
                window.location.href = 'index.html'; // Redirect to homepage
            },
            error: function (xhr) {
                try {
                    // Parse and display error message if available
                    var errorResponse = JSON.parse(xhr.responseText);
                    alert('資料產生失敗: ' + errorResponse.error);
                } catch (e) {
                    alert('資料產生失敗'); // General error message
                }
            }
        });
    });
});
