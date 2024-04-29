var data = JSON.stringify({
    "username": document.getElementById('username').value,
    "password": document.getElementById('password').value,
    "email": document.getElementById('email').value,
    "firstName": document.getElementById('firstName').value,
    "lastName": document.getElementById('lastName').value,
    "phone": document.getElementById('phone').value
});

document.getElementById('registerForm').addEventListener('submit', function(event) {
    event.preventDefault();
    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/register', true);
    xhr.setRequestHeader('Content-type', 'application/json');
    xhr.onload = function () {
        if (this.status == 201) {
            alert('Registration successful.');
            window.location.replace('/login.html');  // Redirect to login page
        } else {
            alert('Registration failed');
        }
    };
    xhr.send(data);
});