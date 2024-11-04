$(document).ready(function() {
    const SERVER_BASE_URL = window._config.SERVER_BASE_URL || 'http://localhost:8080';

    const urlParams = new URLSearchParams(window.location.search);
    const customerId = urlParams.get('id');
    $('#customer-id').val(customerId);

    // Default past month for filtering
    let fromDate = new Date();
    fromDate.setMonth(fromDate.getMonth() - 1);
    $('#from-date').val(formatDate(fromDate));

    let toDate = new Date();
    $('#to-date').val(formatDate(toDate));

    // Global variable to store transactions
    let transactionsData = [];

    // Initial transaction load
    loadTransactions();

    // Form submission for filtering
    $('#filter-form').submit(function(e) {
        e.preventDefault();
        filterTransactions();
    });

    function loadTransactions() {
        $.ajax({
            url: `${SERVER_BASE_URL}/customers/${customerId}/transactions`,
            method: 'GET',
            success: function(transactions) {
                transactionsData = transactions;
                displayTransactions(transactionsData);
            },
            error: function() {
                alert('Unable to retrieve transaction records');
            }
        });
    }

    function filterTransactions() {
        let from = $('#from-date').val();
        let to = $('#to-date').val();

        // Filter the stored data
        let filteredData = transactionsData.filter(function(txn) {
            let txnDate = txn.time.substring(0,10);
            return txnDate >= from && txnDate <= to;
        });

        displayTransactions(filteredData);

        // Make a second request to get transactions for the date range
        $.ajax({
            url: `${SERVER_BASE_URL}/customers/${customerId}/transactions/date`,
            method: 'GET',
            contentType: 'application/json',
            data: JSON.stringify({ from: from, to: to }),
            success: function(newTransactions) {
                // Update transactionsData with newTransactions, avoiding duplicates
                newTransactions.forEach(function(newTxn) {
                    let exists = transactionsData.some(function(txn) {
                        return txn.id === newTxn.id;
                    });
                    if (!exists) {
                        transactionsData.push(newTxn);
                    }
                });

                // Filter again with the updated transactionsData
                let updatedFilteredData = transactionsData.filter(function(txn) {
                    let txnDate = txn.time.substring(0,10);
                    return txnDate >= from && txnDate <= to;
                });

                displayTransactions(updatedFilteredData);
            },
            error: function() {
                alert('Unable to retrieve transaction records for the selected date range');
            }
        });
    }

    function displayTransactions(transactions) {
        $('#transaction-table-body').empty();
        let totalAmount = 0;
        transactions.forEach(function(txn) {
            $('#transaction-table-body').append(`
                <tr>
                    <td>${new Date(txn.time).toLocaleString()}</td>
                    <td>${txn.amount.toFixed(2)}</td>
                    <td>${txn.sequence}</td>
                </tr>
            `);
            totalAmount += txn.amount;
        });
        // 更新交易總筆數
        $('#transactions-count').text(`交易總筆數：${transactions.length}`);
        $('#transactions-totalAmount').text(`交易總金額：${totalAmount.toFixed(2)}`)
    }    

    function formatDate(date) {
        let month = '' + (date.getMonth() + 1);
        let day = '' + date.getDate();
        let year = date.getFullYear();

        if (month.length < 2)
            month = '0' + month;
        if (day.length < 2)
            day = '0' + day;

        return [year, month, day].join('-');
    }
});
