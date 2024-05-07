document.addEventListener('DOMContentLoaded', () => {
    const productsContainer = document.getElementById('products');
    const searchForm = document.getElementById('search-form');
    const searchQueryInput = document.getElementById('search-query');

    searchForm.addEventListener('submit', (event) => {
        event.preventDefault();
        const query = searchQueryInput.value.trim();
        if (query !== "") { // Check if the query is not empty
            fetchProducts(query);
        } else {
            // Optionally, you can provide feedback to the user about the empty search query
            alert("Please enter a search query.");
        }
    });

    // Event listener for the logout button
    const logoutButton = document.getElementById('logoutButton');
    logoutButton.addEventListener('click', logoutUser);

    // Event listener for the login button
    const loginButton = document.getElementById('loginButton');
    loginButton.addEventListener('click', login);

    // Event listener for the close button in the modal
    const closeModalButton = document.querySelector('.close');
    closeModalButton.addEventListener('click', closeModal);
    fetchProducts();
    checkLoginState();

    fetchCategories();

    const sortDropdown = document.getElementById('sort-dropdown');
    sortDropdown.addEventListener('change', (e) => {
        const sortingOption = e.target.value;
        // Get the products currently displayed
        const products = document.querySelectorAll('.product');

        // Convert NodeList to Array
        const productsArray = Array.from(products);

        // Apply sorting logic based on the selected option
        switch (sortingOption) {
            case 'price-asc':
                productsArray.sort((a, b) => {
                    const priceA = parseFloat(a.querySelector('.product-price').textContent);
                    const priceB = parseFloat(b.querySelector('.product-price').textContent);
                    return priceA - priceB;
                });
                break;
            case 'price-desc':
                productsArray.sort((a, b) => {
                    const priceA = parseFloat(a.querySelector('.product-price').textContent);
                    const priceB = parseFloat(b.querySelector('.product-price').textContent);
                    return priceB - priceA;
                });
                break;
            case 'quantity-asc':
                productsArray.sort((a, b) => {
                    const quantityA = a.querySelector('.product-quantity').textContent === 'Product not in stock' ? Infinity : parseInt(a.querySelector('.product-quantity').textContent);
                    const quantityB = b.querySelector('.product-quantity').textContent === 'Product not in stock' ? Infinity : parseInt(b.querySelector('.product-quantity').textContent);
                    // If both products are out of stock, or both are in stock, sort them based on their quantities
                    if (isNaN(quantityA) === isNaN(quantityB)) {
                        return quantityA - quantityB;
                    } else {
                        // If one product is out of stock, it should be considered greater (placed after) than the other
                        // Product with quantity (in stock) should be considered less (placed before) than the out of stock product
                        return isNaN(quantityA) ? 1 : -1;
                    }
                });
                break;
            case 'quantity-desc':
                productsArray.sort((a, b) => {
                    const quantityA = a.querySelector('.product-quantity').textContent === 'Product not in stock' ? -Infinity : parseInt(a.querySelector('.product-quantity').textContent);
                    const quantityB = b.querySelector('.product-quantity').textContent === 'Product not in stock' ? -Infinity : parseInt(b.querySelector('.product-quantity').textContent);
                    // If both products are out of stock, or both are in stock, sort them based on their quantities
                    if (isNaN(quantityA) === isNaN(quantityB)) {
                        return quantityB - quantityA;
                    } else {
                        // If one product is out of stock, it should be considered greater (placed after) than the other
                        // Product with quantity (in stock) should be considered less (placed before) than the out of stock product
                        return isNaN(quantityA) ? 1 : -1;
                    }
                });
                break;
            default:
                break;
        }


        // Clear the products container
        productsContainer.innerHTML = '';

        // Append sorted products to the products container
        productsArray.forEach(product => {
            productsContainer.appendChild(product);
        });
    });
});

function createProductElement(product) {
    const productElement = document.createElement('div');
    productElement.className = 'product';
    productElement.dataset.productId = product.id;

    const productName = document.createElement('span');
    productName.textContent = product.name;

    const productDescription = document.createElement('span');
    productDescription.textContent = product.description;

    const priceContainer = document.createElement('span');
    const productPriceLabel = document.createElement('span');
    productPriceLabel.className = 'product-label';
    productPriceLabel.textContent = 'Kr ';

    const productPrice = document.createElement('span');
    productPrice.textContent = product.price;
    productPrice.className = 'product-price';

    priceContainer.appendChild(productPriceLabel);
    priceContainer.appendChild(productPrice);

    const quantityContainer = document.createElement('span');
    const productQuantityLabel = document.createElement('span');
    productQuantityLabel.className = 'product-label';

    const productQuantity = document.createElement('span');
    if (product.hasOwnProperty('qtyInStock') && product.qtyInStock > 0) {
        productQuantityLabel.textContent = 'In Stock: ';
        productQuantity.textContent = product.qtyInStock;
    } else {
        productQuantity.textContent = 'Product not in stock';
    }
    productQuantity.className = 'product-quantity'; // Add class for quantity

    quantityContainer.appendChild(productQuantityLabel);
    quantityContainer.appendChild(productQuantity);

    productElement.appendChild(productName);
    productElement.appendChild(document.createElement('br'));
    productElement.appendChild(productDescription);
    productElement.appendChild(document.createElement('br'));
    productElement.appendChild(priceContainer);
    productElement.appendChild(document.createElement('br'));
    productElement.appendChild(quantityContainer);

    return productElement;
}

function updateProductDisplay(products) {
    const productsContainer = document.getElementById('products');
    productsContainer.innerHTML = '';

    if (products.length === 0) {
        const noResults = document.createElement('p');
        noResults.textContent = 'No results found.';
        productsContainer.appendChild(noResults);
    } else {
        products.forEach(product => {
            const productElement = createProductElement(product);
            productsContainer.appendChild(productElement);
        });
    }
}

function fetchProducts(query = "") {
    fetch(`/api/v1/products/search/${query}`)
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            const productsContainer = document.getElementById('products');
            productsContainer.innerHTML = ''; // Clear previous content

            if (data && data.length > 0) { // Check if data is not null and has at least one element
                updateProductDisplay(data);
                // Update breadcrumbs with search query
                updateBreadcrumbsBySearch(query);
            } else {
                const noResults = document.createElement('p');
                noResults.textContent = `Could not find anything with "${query}"`;
                productsContainer.appendChild(noResults);
            }
        })
        .catch(error => console.error('Error fetching products:', error));
}

function fetchCategories() {
    fetch('/api/v1/categories/')
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to fetch categories');
            }
            return response.json();
        })
        .then(categories => {
            const dropdown = document.querySelector('#categories');
            dropdown.innerHTML = '';
            categories.forEach(category => {
                const a = document.createElement('a');
                a.text = category.name;
                const categoryName = encodeURIComponent(category.name); // Encode category name
                a.href = '/products/' + categoryName; // Set the URL
                a.dataset.categoryName = category.name; // Store the category name
                a.addEventListener('click', (event) => {
                    event.preventDefault();
                    const categoryName = event.target.dataset.categoryName; // Get category name
                    fetchProductsByCategory(categoryName);
                    history.pushState(null, null, `/products/${categoryName}`); // Update URL
                });
                dropdown.appendChild(a); // Append the created <a> element to the dropdown
            });
        })
        .catch(error => console.error('Error fetching categories:', error));
}

function fetchProductsByCategory(categoryName) {
    // Fetch products based on the category name
    fetch(`http://localhost:8000/api/v1/products/search/${categoryName}`)
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(products => {
            updateProductDisplay(products); // Update display with products received from backend
            updateBreadcrumbs(categoryName); // Update breadcrumbs after fetching products
        })
        .catch(error => {
            console.error('Error fetching products:', error);
            displayErrorMessage('Sorry, there was an error loading products.');
        });
}

function updateBreadcrumbsBySearch(query) {
    const breadcrumbs = document.querySelector('.breadcrumbs');
    breadcrumbs.innerHTML = ''; // Clear existing breadcrumbs
    const homeLink = document.createElement('a');
    homeLink.textContent = 'Home';
    homeLink.href = '/';
    breadcrumbs.appendChild(homeLink);

    // Add separator if there are more breadcrumbs
    if (query !== "") {
        breadcrumbs.appendChild(document.createTextNode(' * '));

        // Add search query
        const searchLink = document.createElement('span');
        searchLink.textContent = `Search: ${query}`;
        breadcrumbs.appendChild(searchLink);
    }
}




function updateBreadcrumbs(categoryId) {
    const breadcrumbs = document.querySelector('.breadcrumbs');
    breadcrumbs.innerHTML = ''; // Clear existing breadcrumbs
    const homeLink = document.createElement('a');
    homeLink.textContent = 'Home';
    homeLink.href = '/';
    breadcrumbs.appendChild(homeLink);

    // Fetch category details
    fetch(`http://localhost:8000/api/v1/categories/${categoryId}`)
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to fetch category details');
            }
            return response.json();
        })
        .then(category => {
            const categoryLink = document.createElement('a');
            categoryLink.textContent = category.name;
            categoryLink.href = `#${categoryId}`;
            categoryLink.addEventListener('click', (event) => {
                event.preventDefault();
                fetchProductsByCategory(categoryId);
            });
            breadcrumbs.appendChild(document.createTextNode(' * '));
            breadcrumbs.appendChild(categoryLink);
        })
        .catch(error => console.error('Error fetching category details:', error));
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
                logoutButton.style.display = 'block';
            } else {
                userNotLogged.style.display = 'block';
                userLogged.style.display = 'none';
                logoutButton.style.display = 'none';
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
    const userMenu = document.getElementById("user-menu");
    userMenu.classList.toggle("show");
}

document.getElementById("logged-username").addEventListener("click", () => {
    openUserMenu();
});

window.onclick = function(event) {
    if (!event.target.matches('#logged-username')) {
        const userMenu = document.getElementById("user-menu");
        if (userMenu.classList.contains('show')) {
            userMenu.classList.remove('show');
        }
    }
};

function login() {
    const modal = document.getElementById("loginModal");
    modal.style.display = "block";
}

function closeModal() {
    const modal = document.getElementById("loginModal");
    modal.style.display = "none";
}