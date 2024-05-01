// Fetching the user profile
var xhr = new XMLHttpRequest();
xhr.open('GET', '/api/profile', true);
xhr.setRequestHeader('Content-type', 'application/json');
xhr.onload = function () {
    if (this.status == 200) {
        var respJson = JSON.parse(this.responseText);
        console.log(respJson); // to log the whole object
        console.log(respJson.username); // to log just the username
        document.getElementById('username').value = respJson.username;
        document.getElementById('firstName').value = respJson.firstName;
        document.getElementById('lastName').value = respJson.lastName;
        document.getElementById('address').value = respJson.address;
    } else if (this.status == 401) {
        alert('Please login to access your profile.');
        // Here you might want to redirect your user to the login page depending on your UX design
    } else {
        alert('Failed to get profile information.');
    }
};
xhr.send();