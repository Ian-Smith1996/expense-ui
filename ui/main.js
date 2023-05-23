let expenses = [];

document.getElementById("expenseForm").addEventListener("submit", function(e) {
    e.preventDefault();

    const message = document.getElementById("message");
    const overlay = document.getElementById("overlay");
    const submitButton = document.getElementById("submitButton");
    const amount = document.getElementById("amount").value;
    const date = document.getElementById("date").value;

    const expense = {
        amount: parseFloat(amount),
        date: date
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

    fetch('http://localhost:8098/create', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(expense),
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
            body: JSON.stringify({ expenseID: expense.expenseID })
        })
        .then(response => response.json())
        .then(data => {
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


document.getElementById('expandButton').addEventListener('click', function(e) {
    const reportList = document.getElementById('reportList');
    reportList.style.display = reportList.style.display === 'none' ? 'block' : 'none';
    renderReports();
});

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

