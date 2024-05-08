class EStore {
    constructor() {
        this.productsContainer = document.getElementById('products');
        this.searchForm = document.getElementById('search-form');
        this.searchQueryInput = document.getElementById('search-query');
        this.logoutButton = document.getElementById('logoutButton');
        this.loginButton = document.getElementById('loginButton');
        this.closeModalButton = document.querySelector('.close');
        this.sortDropdown = document.getElementById('sort-dropdown');

        this.searchForm.addEventListener('submit', this.handleSearch.bind(this));
        this.logoutButton.addEventListener('click', this.logoutUser.bind(this));
        this.loginButton.addEventListener('click', this.openLoginModal.bind(this));
        this.closeModalButton.addEventListener('click', this.closeModal.bind(this));
        this.sortDropdown.addEventListener('change', this.handleSort.bind(this));

        this.checkLoginState();
        this.fetchProducts();
        this.fetchCategories();
    }

    handleSearch(event) {
        event.preventDefault();
        const query = this.searchQueryInput.value.trim();
        if (query !== "") {
            this.fetchProducts(query);
            this.updateBreadcrumbs('Search Results: '); // Update breadcrumbs for search
        } else {
            alert("Please enter a search query.");
        }
    }


    handleSort(event) {
        const sortingOption = event.target.value;
        const products = document.querySelectorAll('.product');
        const productsArray = Array.from(products);

        switch (sortingOption) {
            case 'price-asc':
                productsArray.sort((a, b) => this.getProductPrice(a) - this.getProductPrice(b));
                break;
            case 'price-desc':
                productsArray.sort((a, b) => this.getProductPrice(b) - this.getProductPrice(a));
                break;
            case 'quantity-asc':
                productsArray.sort((a, b) => this.getProductQuantity(a) - this.getProductQuantity(b));
                break;
            case 'quantity-desc':
                productsArray.sort((a, b) => this.getProductQuantity(b) - this.getProductQuantity(a));
                break;
            default:
                break;
        }

        // Clear products container
        this.productsContainer.innerHTML = '';

        // Append sorted products
        productsArray.forEach(product => {
            this.productsContainer.appendChild(product);
        });
    }

    getProductPrice(productElement) {
        const priceString = productElement.querySelector('.product-price').textContent;
        // Extract the numeric part of the price string
        const priceNumeric = parseFloat(priceString.replace(/[^\d.-]/g, ''));
        return isNaN(priceNumeric) ? -1 : priceNumeric;
    }


    getProductQuantity(productElement) {
        const stockString = productElement.querySelector('.product-stock').textContent;
        const quantity = parseInt(stockString.split(': ')[1]);
        return isNaN(quantity) ? -1 : quantity;
    }

    fetchCategories() {
        fetch('/api/v1/categories/')
            .then(this.handleFetchResponse)
            .then(this.renderCategories.bind(this))
            .catch(error => console.error('Error fetching categories:', error));
    }

    renderCategories(categories) {
        const dropdown = document.querySelector('#categories');
        dropdown.innerHTML = '';
        categories.forEach(category => {
            const a = document.createElement('a');
            a.textContent = category.name;
            const categoryName = encodeURIComponent(category.name);
            a.href = '/products/' + categoryName;
            a.dataset.categoryName = category.name;
            a.addEventListener('click', this.handleCategoryClick.bind(this));
            dropdown.appendChild(a);
        });
    }

    handleCategoryClick(event) {
        event.preventDefault();
        const categoryName = event.target.dataset.categoryName;
        this.fetchProductsByCategory(categoryName);
        this.updateBreadcrumbs(categoryName); // Update breadcrumbs
        history.pushState(null, null, `/products/${categoryName}`);
    }

    fetchProductsByCategory(categoryName) {
        fetch(`/api/v1/products/search/${categoryName}`)
            .then(this.handleFetchResponse)
            .then(this.updateProductDisplay.bind(this))
            .catch(error => console.error('Error fetching products:', error));
    }

    handleFetchResponse(response) {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    }

    updateProductDisplay(products) {
        this.productsContainer.innerHTML = '';
        if (products.length === 0) {
            const noResults = document.createElement('p');
            noResults.textContent = 'No results found.';
            this.productsContainer.appendChild(noResults);
        } else {
            products.forEach(product => {
                const productElement = this.createProductElement(product);
                this.productsContainer.appendChild(productElement);
            });
        }
    }

    updateBreadcrumbs(categoryName) {
        const breadcrumbsContainer = document.querySelector('.breadcrumbs');
        breadcrumbsContainer.innerHTML = ''; // Clear previous breadcrumbs

        // Create a link to the home page
        const homeLink = document.createElement('a');
        homeLink.textContent = 'Home';
        homeLink.href = '/';
        breadcrumbsContainer.appendChild(homeLink);

        // Create a separator
        const separator = document.createElement('span');
        separator.textContent = ' > ';
        breadcrumbsContainer.appendChild(separator);

        // Create a link to the category page
        const categoryLink = document.createElement('a');
        categoryLink.textContent = categoryName; // Use the name of the category
        categoryLink.href = `/products/${categoryName}`;
        breadcrumbsContainer.appendChild(categoryLink);
    }


    createProductElement(product) {
        // Create the product container
        const productElement = document.createElement('div');
        productElement.className = 'product';
        productElement.dataset.productId = product.id;

        // Create the left column for name, description, and stock
        const leftColumn = document.createElement('div');
        leftColumn.className = 'product-left-column';

        // Create the name element
        const nameElement = document.createElement('div');
        nameElement.className = 'product-name';
        nameElement.textContent = 'Name: ' + product.name;

        // Create the description element
        const descriptionElement = document.createElement('div');
        descriptionElement.className = 'product-description';
        descriptionElement.textContent = 'Description: ' + product.description;

        // Create the stock element
        const stockElement = document.createElement('div');
        stockElement.className = 'product-stock';
        if (product.hasOwnProperty('qtyInStock') && product.qtyInStock > 0) {
            stockElement.textContent = 'In Stock: ' + product.qtyInStock;
        } else {
            stockElement.textContent = 'Product not in stock';
        }

        // Append name, description, and stock to the left column
        leftColumn.appendChild(nameElement);
        leftColumn.appendChild(descriptionElement);
        leftColumn.appendChild(stockElement);

        // Create the right column for price and button
        const rightColumn = document.createElement('div');
        rightColumn.className = 'product-right-column';

        // Create the price element
        const priceElement = document.createElement('div');
        priceElement.className = 'product-price';
        priceElement.textContent = 'Price: Kr ' + product.price;

        // Create the button element
        const buttonElement = document.createElement('button');
        buttonElement.className = 'add-to-cart-button';
        buttonElement.dataset.productId = product.id;
        buttonElement.textContent = 'Add to Cart';

        // Append price and button to the right column
        rightColumn.appendChild(priceElement);
        rightColumn.appendChild(buttonElement);

        // Append left column to the product container
        productElement.appendChild(leftColumn);

        // Append right column to the product container
        productElement.appendChild(rightColumn);

        return productElement;
    }



    fetchProducts(query = "") {
        fetch(`/api/v1/products/search/${query}`)
            .then(this.handleFetchResponse)
            .then(this.updateProductDisplay.bind(this))
            .catch(error => console.error('Error fetching products:', error));
    }

    checkLoginState() {
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

    logoutUser() {
        fetch('/api/logout', { credentials: 'include' })
            .then(response => response.text())
            .then(text => {
                this.checkLoginState();
                window.location.href = "/";
            })
            .catch(error => console.error('Error:', error));
    }

    openLoginModal() {
        const modal = document.getElementById("loginModal");
        modal.style.display = "block";
    }

    closeModal() {
        const modal = document.getElementById("loginModal");
        modal.style.display = "none";
    }
}

// Instantiate the EStore class
new EStore();
