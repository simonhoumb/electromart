document.getElementById('loginForm').addEventListener('submit', function (event) {
    event.preventDefault();

    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    fetch('/login', {
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
            document.getElementById('user-not-logged').style.display = 'none';
            document.getElementById('logged-username').textContent = data.username;
            document.getElementById('user-logged').style.display = 'block';
        })
        .catch((error) => {
            console.error('Error:', error);
        });
});

function logout() {
    fetch('/logout', {
        method: 'POST',
    })
        .then(() => {
            // Clear the username and update the UI
            document.getElementById('user-not-logged').style.display = 'block';
            document.getElementById('logged-username').textContent = '';
            document.getElementById('user-logged').style.display = 'none';
            // Now, let's redirect user to login page
            window.location.href = '/';
        })
        .catch((error) => {
            console.error('Error:', error);
        });
}