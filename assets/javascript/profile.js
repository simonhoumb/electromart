document.addEventListener('DOMContentLoaded', function () {
    // Fetch user profile when the page is loaded
    fetchUserProfile();

    // Add event listener to the form for updating profile
    document.getElementById('profileForm').addEventListener('submit', function (event) {
        event.preventDefault(); // Prevent default form submission behavior
        updateUserProfile();
    });
});

// Function to fetch user profile
function fetchUserProfile() {
    console.log("Fetching user profile..."); // Check if function is being called
    var xhr = new XMLHttpRequest();
    xhr.open('GET', '/api/profile', true);
    xhr.setRequestHeader('Content-type', 'application/json');
    xhr.onload = function () {
        console.log("Received response:", xhr.responseText); // Check response from server
        if (xhr.status === 200) {
            var respJson = JSON.parse(xhr.responseText);
            fillUserProfileForm(respJson);
        } else if (xhr.status === 401) {
            alert('Please login to access your profile.');
            // Redirect the user to the login page
            window.location.href = '/login';
        } else {
            showError('Failed to get profile information. Please try again later.');
        }
    };
    xhr.onerror = function () {
        showError('Failed to fetch profile. Please check your internet connection and try again.');
    };
    xhr.send();
}

// Function to fill user profile form fields
function fillUserProfileForm(user) {
    console.log("Filling user profile form..."); // Check if function is being called
    document.getElementById('username').value = user.Username;
    document.getElementById('email').value = user.Email;
    document.getElementById('firstName').value = user.FirstName;
    document.getElementById('lastName').value = user.LastName;
    document.getElementById('phone').value = user.Phone;
    if (user.Address && user.Address.Valid) {
        document.getElementById('address').value = user.Address.String;
    }
    if (user.PostCode && user.PostCode.Valid) {
        document.getElementById('postCode').value = user.PostCode.String;
    }
}



// Function to update user profile
function updateUserProfile() {
    // Retrieve input values
    var username = document.getElementById('username').value.trim();
    var email = document.getElementById('email').value.trim();
    var firstName = document.getElementById('firstName').value.trim();
    var lastName = document.getElementById('lastName').value.trim();
    var phone = document.getElementById('phone').value.trim();
    var address = document.getElementById('address').value.trim();
    var postCode = document.getElementById('postCode').value.trim();

    // Check if required fields are empty
    if (username === '' || email === '' || firstName === '' || lastName === '' || address === '' || postCode === '') {
        alert('Please fill in all required fields.');
        return;
    }

    // Prepare data for PATCH request
    var data = {
        "Username": username,
        "Email": email,
        "FirstName": firstName,
        "LastName": lastName,
        "Phone": phone,
        "Address": {"String": address, "Valid": true},
        "PostCode": {"String": postCode, "Valid": true}
    };

    // Send PATCH request with console logging
    var xhr = new XMLHttpRequest();
    xhr.open('PATCH', '/api/profile', true);
    xhr.setRequestHeader('Content-type', 'application/json');
    console.log("Sending PATCH request with data:", data);  // Log data before sending
    xhr.onload = function () {
        console.log("PATCH request response:", xhr.responseText);  // Log response
        console.log("Status code:", xhr.status);                     // Log status code
        if (this.status == 200) {
            alert('Profile updated successfully');
            // Fetch updated profile after successful update
            fetchUserProfile();
        } else if (this.status == 400) {
            alert('Bad request: The server could not understand the request due to invalid syntax.');
        } else if (this.status == 401) {
            alert('Unauthorized: Please login to update your profile.');
            // Redirect the user to the login page
            window.location.href = '/login';
        } else if (this.status == 403) {
            alert('Forbidden: You are not allowed to update this profile.');
        } else if (this.status == 404) {
            alert('Not found: The requested resource could not be found.');
        } else if (this.status == 500) {
            alert('Internal Server Error: Failed to update profile. Please try again later.');
        } else {
            alert('Unknown error occurred. Status code: ' + this.status);
        }
    };
    xhr.onerror = function () {
        alert('Failed to update profile. Please check your internet connection and try again.');
    }
    xhr.send(JSON.stringify(data));
}

document.addEventListener('DOMContentLoaded', function () {
    // Add event listener to the delete user button
    document.getElementById('deleteUserBtn').addEventListener('click', function () {
        deleteUser();
    });
});

// Function to send DELETE request to delete user
function deleteUser() {
    var xhr = new XMLHttpRequest();
    xhr.open('DELETE', '/api/profile', true);
    xhr.setRequestHeader('Content-type', 'application/json');
    xhr.onload = function () {
        if (xhr.status === 204) {
            alert('User deleted successfully');
            // Redirect the user to the login page
            window.location.href = '/';
        } else {
            alert('Failed to delete user. Please try again later.');
        }
    };
    xhr.onerror = function () {
        alert('Failed to delete user. Please check your internet connection and try again.');
    };
    xhr.send();
}


