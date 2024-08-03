let pieChart = null;
var ctx = document.getElementById("pieChart").getContext("2d");

function createPieChart(chartData) {
    const ctx = document.getElementById('pieChart').getContext('2d');

    if (pieChart) {
        pieChart.destroy();
    }

    pieChart = new Chart(ctx, {
        type: 'pie',
        data: {
            labels: ['Netpay', 'Payee', 'Health', 'Housing', 'Pension', 'Emp. Pension'],
            datasets: [{
                data: [
                    chartData.payslip2,
                    chartData.payee2,
                    chartData.health2,
                    chartData.housing2,
                    chartData.pension2,
                    chartData.emp_pension2,
                ],
                backgroundColor: [
                    'rgba(255, 99, 132, 0.7)',
                    'rgba(54, 162, 235, 0.7)',
                    'rgba(255, 206, 86, 0.7)',
                    'rgba(75, 192, 192, 0.7)',
                    'rgba(153, 102, 255, 0.7)',
                    'rgba(255, 159, 64, 0.7)',
                    'rgba(201, 203, 207, 0.7)'
                ],
                borderWidth: 1
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                title: {
                    display: true,
                    text: 'Payroll Summary'
                }
            }
        }
    });
}

$(document).ready(function () {
    $("#calculate-button").click(function () {
        var formData = {
            gross: $("#gross").val(),
            health: $("#health").val(),
            contrib: $("#contrib").val(),
            housing: $("#housing").val()
        };

        $.ajax({
            type: "POST",
            url: "/payslip",
            data: formData,
            dataType: "json",
            success: function (response) {
                var formattedPayee = response.payslip.toLocaleString();
                $("#result").html("Monthly Netpay: " + "₦" + formattedPayee + "<br>" +
                    "Payee: " + "₦" + response.payee + "<br>" +
                    "Health Insurance Payment: " + "₦" + response.health + "<br>" +
                    "National Housing Fund: " + "₦" + response.housing + "<br>" +
                    "Employee Pension Contribution: " + "₦" + response.emp_pension + "<br>" +
                    "Total Pension to be Remitted: " + "₦" + response.pension);
                createPieChart(response);
            },
        })
    })
});