// main.js
const productsContainer = document.querySelector('#products');

document.addEventListener('DOMContentLoaded', (event) => {
    const searchForm = document.getElementById('search-form');
    const searchQueryInput = document.getElementById('search-query');
    let query = "";

    searchForm.addEventListener('submit', (event) => {
        event.preventDefault();
        query = searchQueryInput.value.trim(); // Update query with trimmed input value
        fetchProducts(query); // Fetch products with the updated query
    });

    // Fetch all products when the page loads
    fetchProducts();

    fetchCategories();
    checkLoginState();
});

function fetchProducts(query="") {
    fetch('http://localhost:8000/api/v1/products/search/'+query)
        .then(response => response.json())
        .then(data => {
            productsContainer.innerHTML = ''; // Clear products container
            data.forEach(product => {
                const productElement = document.createElement('div');
                productElement.className = 'product';

                const productName = document.createElement('span');
                productName.textContent = product.name;

                const productDescription = document.createElement('p');
                productDescription.textContent = product.description;

                const productPriceLabel = document.createElement('span');
                productPriceLabel.className = 'product-label'; // Add class for styling
                productPriceLabel.textContent = 'Kr ';

                const productPrice = document.createElement('p');
                productPrice.textContent = product.price;

                const productQuantityLabel = document.createElement('span');
                productQuantityLabel.className = 'product-label'; // Add class for styling

                const productQuantity = document.createElement('p');
                if (product.hasOwnProperty('qtyInStock') && product.qtyInStock > 0) {
                    productQuantityLabel.textContent = 'In Stock: ';
                    productQuantity.textContent = product.qtyInStock;
                } else {
                    productQuantity.textContent = 'Product not in stock';
                }

                productElement.appendChild(productName);
                productElement.appendChild(productDescription);
                productElement.appendChild(productPriceLabel);
                productElement.appendChild(productPrice);
                productElement.appendChild(productQuantityLabel);
                productElement.appendChild(productQuantity);

                productsContainer.appendChild(productElement);
            });
            if (productsContainer.children.length === 0) {
                const noResults = document.createElement('p');
                noResults.textContent = 'No results found.';
                productsContainer.appendChild(noResults);
            }
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

