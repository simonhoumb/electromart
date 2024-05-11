// Function to handle form submission
function handleFormSubmit(event) {
    event.preventDefault(); // Prevent the default form submission

    // Retrieve the values of the username and password fields
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    // Log the retrieved values for debugging
    console.log("Username:", username);
    console.log("Password:", password);

    // Perform the fetch request with the retrieved values
    fetch('/api/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            username: username,
            password: password,
        }),
    })
        .then(response => response.json())
        .then(data => {
            // Close the modal and update the UI
            closeModal();
            return new EStore().checkLoginState();
        })
        .catch((error) => {
            console.error('Error:', error);
        });
}

// Function to close the modal
function closeModal() {
    var modal = document.getElementById("loginModal");
    modal.style.display = "none";
}

// Event listener for form submission
document.addEventListener('DOMContentLoaded', function () {
    document.getElementById('loginForm').addEventListener('submit', handleFormSubmit);

    // Event listener for closing the modal
    document.querySelector('.close').addEventListener('click', closeModal);

    // Event listener for redirecting to the registration page
    document.getElementById('registerLink').addEventListener('click', function(event) {
        event.preventDefault();  // Prevent default link behavior
        window.location.href = '/register'; // Redirect to the registration page
    });
});

