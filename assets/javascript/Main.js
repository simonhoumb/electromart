// main.js
document.addEventListener('DOMContentLoaded', (event) => {
    // Fetch products and categories, and check login state
    fetchProducts();
    fetchCategories();
    checkLoginState();
});

function fetchProducts() {
    fetch('/products')
        .then(response => response.json())
        .then(data => {
            const productsContainer = document.querySelector('#products');
            data.forEach(product => {
                const productElement = document.createElement('div');
                productElement.className = 'product';

                const productName = document.createElement('h3');
                productName.textContent = product.name;

                const productDescription = document.createElement('p');
                productDescription.textContent = product.description;

                const productPrice = document.createElement('p');
                productPrice.textContent = product.price;

                productElement.appendChild(productName);
                productElement.appendChild(productDescription);
                productElement.appendChild(productPrice);

                productsContainer.appendChild(productElement);
            });
        })
        .catch(error => console.error('Error fetching products:', error));
}

function fetchCategories() {
    fetch('/api/categories')
        .then(response => response.json())
        .then(categories => {
            const dropdown = document.querySelector('#categories');
            categories.forEach(category => {
                const a = document.createElement('a');
                a.text = category.name;
                a.href = '/category/' + category.id;
                dropdown.appendChild(a);
            });
        })
        .catch(error => console.error('Error fetching categories:', error));
}

function checkLoginState() {
    fetch('/api/check_login', { credentials: 'include' })
        .then(response => response.json())
        .then(respJson => {
            const userNotLogged = document.getElementById('user-not-logged');
            const userLogged = document.getElementById('user-logged');
            const logoutButton = document.getElementById('logoutButton');

            if (respJson.logged_in) {
                document.getElementById('logged-username').textContent = respJson.username;
                userNotLogged.style.display = 'none';
                userLogged.style.display = 'block';
                logoutButton.style.display = 'block'; // Show logout button
            } else {
                userNotLogged.style.display = 'block';
                userLogged.style.display = 'none';
                logoutButton.style.display = 'none'; // Hide logout button
            }
        })
        .catch(error => console.error('Error checking login:', error));
}


function logoutUser() {
    fetch('/api/logout', { credentials: 'include' })
        .then(response => response.text())
        .then(text => {
            checkLoginState();
            window.location.href = "/";
        })
        .catch(error => console.error('Error:', error));
}

function openUserMenu() {
    var userMenu = document.getElementById("user-menu");
    userMenu.classList.toggle("show");
}

// Add event listener to the username to open user menu
document.getElementById("logged-username").addEventListener("click", function() {
    openUserMenu();
});

// Close the user menu if the user clicks outside of it
window.onclick = function(event) {
    if (!event.target.matches('#logged-username')) {
        var userMenu = document.getElementById("user-menu");
        if (userMenu.classList.contains('show')) {
            userMenu.classList.remove('show');
        }
    }
}

function login() {
    var modal = document.getElementById("loginModal");
    modal.style.display = "block"; // Display the modal when login button is clicked
}

function closeModal() {
    var modal = document.getElementById("loginModal");
    modal.style.display = "none"; // Hide the modal when close button is clicked
}

