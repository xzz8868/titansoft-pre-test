$(document).ready(function() {
    const GENERATOR_BASE_URL = window._config.GENERATOR_BASE_URL || 'http://localhost:8081';

    // Handle form submission for customer data generation
    $('#generate-customer-form').submit(function(e) {
        e.preventDefault(); // Prevent default form submission

        // Get the input value for the number of customers to generate
        const num = $('#num').val();

        // Verify if the input value exceeds 2000
        if (num > 2000) {
            alert('單次最多只能產生2000筆資料');
            return;
        }

        // Construct URL with query parameters
        const urlWithParams = `${GENERATOR_BASE_URL}/generate/customer?num=${encodeURIComponent(num)}`;

        // Send AJAX request to generate customer data
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
            success: function() {
                alert('資料產生成功'); // Alert on success
                window.location.href = 'index.html'; // Redirect to index page
            },
            error: function() {
                alert('資料產生失敗'); // Alert on error
            }
        });
    });
});
