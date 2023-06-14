// Global array to store the expenses.
let expenses = [];

// // Async method that runs an interval every 10 seconds to update the state of the expenses.
// async function run() {
//     // Send a request to localhost:8098/query every 10 seconds to update the state of the expenses.
//     // Also rerender the reports.
//     setInterval(() => {
//         if (expenses.length === 0) return;
//         expenses.forEach((expense, index) => {
//             fetch('http://localhost:8098/query', {
//                 method: 'POST',
//                 headers: { 'Content-Type': 'application/json' },
//                 body: JSON.stringify({ expenseID: expense.expenseID })
//             })
//             .then(response => response.json())
//             .then(data => {
//                 if (data.state.toUpperCase() == "ACCEPTED" || data.state.toUpperCase() == "REJECTED") {
//                     // Remove the expense from the list at index
//                     expenses.splice(index, 1);
//                 }
//                 expense.state = data.state;
//             }).catch((error) => {
//                 console.error('Error:', error);
//             });
//         });
//         renderReports();
//     }, 10000);
// }

// run();

// Event listener for the expense form when it is submitted
document.getElementById("expenseForm").addEventListener("submit", function(e) {
    e.preventDefault();

    const message = document.getElementById("message");
    const overlay = document.getElementById("overlay");
    const submitButton = document.getElementById("submitButton");
    const amount = document.getElementById("amount").value;
    const date = document.getElementById("date").value;

    const expense = {
        amount: amount,
        date: date,
    };

    // Display the overlay and disable the submit button
    overlay.style.display = 'flex';
    submitButton.disabled = true;

    // Add a timeout for the request
    setTimeout(() => {
        if (overlay.style.display === 'flex') {
            overlay.style.display = 'none';
            submitButton.disabled = false;
            message.style.display = 'block';
            message.textContent = 'Error: Request timed out!';
            message.className = 'message error';
        }
    }, 10000);

    fetch('http://localhost:8097/expense', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: `{"amount": ${expense.amount}, "date": "${expense.date}"}`,
    })
    .then(response => response.json())
    .then(data => {
        const { expenseID, state } = data;
        expense.expenseID = expenseID;
        expense.state = state;

        renderReports();
        // Hide the overlay, enable the submit button, push the expense, and display success message
        expenses.push(expense);
        overlay.style.display = 'none';
        submitButton.disabled = false;
        message.style.display = 'block';
        message.textContent = 'Expense submitted successfully!';
        message.className = 'message success';
    })
    .catch((error) => {
        console.error('Error:', error);

        // Hide the overlay, enable the submit button, and display error message
        overlay.style.display = 'none';
        submitButton.disabled = false;
        message.style.display = 'block';
        message.textContent = 'Error submitting expense!';
        message.className = 'message error';
    });
});

// Ran when the user wants to view the details of an expense
function showExpenseDetails(index) {
    const expense = expenses[index];
    const detailsDiv = document.createElement('div');
    detailsDiv.id = `details-${index}`;
    
    const stateButton = document.createElement('button');
    stateButton.innerHTML = "Query state";
    stateButton.addEventListener('click', function() {
        // Send a request to localhost:8098/query
        fetch('http://localhost:8098/query', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
        })
        .then(response => response.json())
        .then(data => {
            expense.id = data.expenseID;
            expense.date = data.date;
            expense.state = data.state;
            detailsDiv.querySelector('#state').innerHTML = `State: ${data.state}`;
        });
    });
    
    const stateLabel = document.createElement('p');
    stateLabel.id = 'state';
    stateLabel.innerHTML = `State: ${expense.state}`;

    detailsDiv.appendChild(stateButton);
    detailsDiv.appendChild(stateLabel);
    document.getElementById('reportList').appendChild(detailsDiv);
}

// Event listener for the expand button. Toggles the display of the report list.
document.getElementById('expandButton').addEventListener('click', function(e) {
    const reportList = document.getElementById('reportList');
    reportList.style.display = reportList.style.display === 'none' ? 'block' : 'none';
    renderReports();
});

// Render the report to display the most current expense details
function renderReports() {
    const reportsContainer = document.getElementById('reportList');
    reportsContainer.innerHTML = ''; // Clear the container
    expenses.forEach((expense, index) => {
        const expenseDiv = document.createElement('div');
        expenseDiv.classList.add('expense');
        
        const expenseIDLink = document.createElement('a');
        expenseIDLink.href = '#';
        expenseIDLink.innerHTML = `Expense ID: ${expense.expenseID}`;
        expenseIDLink.addEventListener('click', function(e) {
            e.preventDefault();
            showExpenseDetails(index);
        });

        const amountLabel = document.createElement('p');
        amountLabel.innerHTML = `Amount: $${expense.amount}`;

        const dateLabel = document.createElement('p');
        dateLabel.innerHTML = `Date: ${expense.date}`;

        expenseDiv.appendChild(expenseIDLink);
        expenseDiv.appendChild(amountLabel);
        expenseDiv.appendChild(dateLabel);
        reportsContainer.appendChild(expenseDiv);
    });
}

