class EStore {
    constructor() {
        this.elements = {
            productsContainer: document.getElementById('products'),
            searchForm: document.getElementById('search-form'),
            searchQueryInput: document.getElementById('search-query'),
            logoutButton: document.getElementById('logoutButton'),
            loginButton: document.getElementById('loginButton'),
            closeModalButton: document.querySelector('.close'),
            sortDropdown: document.getElementById('sort-dropdown')
        }
        this.addEventListeners();
        this.checkLoginState();
        this.fetchProducts();
        this.fetchCategories();
    }

    addEventListeners() {
        this.elements.searchForm.addEventListener('submit', this.handleSearch.bind(this));
        this.elements.logoutButton.addEventListener('click', this.logoutUser.bind(this));
        this.elements.loginButton.addEventListener('click', this.openLoginModal.bind(this));
        this.elements.closeModalButton.addEventListener('click', this.closeModal.bind(this));
        this.elements.sortDropdown.addEventListener('change', this.handleSort.bind(this));
    }

    handleSearch(event) {
        event.preventDefault();
        const query = this.elements.searchQueryInput.value.trim();
        if (query !== "") {
            this.fetchProducts(query);
            this.updateBreadcrumbs([{ text: 'Home', href: '/' }, { text: 'Search Results: ' + query, href: '#' }]);
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
        this.elements.productsContainer.innerHTML = '';

        // Append sorted products
        productsArray.forEach(product => {
            this.elements.productsContainer.appendChild(product);
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
        this.updateBreadcrumbs([{ text: 'Home', href: '/' }, { text: categoryName, href: `/products/${categoryName}` }]);
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

    // Inside the EStore class
    async fetchProductDetails(productId) {
        try {
            const dataResponse = await fetch(`/api/v1/products/${productId}`);
            const productData = await dataResponse.json();
            console.log("Product Data:", productData);

            this.updateBreadcrumbs([
                { text: 'Home', href: '/' },
                { text: productData.categoryName, href: `/products/${productData.categoryName}` },
                { text: productData.brandName, href: `#` }
            ]);

            const populatedHtml = await this.populateProductTemplate(productData);
            this.insertProductDetails(populatedHtml);
            return productData;
        } catch (error) {
            console.error('Error fetching product details:', error);
            // Handle the error
            return Promise.reject(error);
        }
    }

    // Helper function to insert content
    insertProductDetails(content) {
        const container = document.getElementById('product-details-container');
        container.innerHTML = content;
        container.style.display = 'block'; // Show the product details container
    }

    async populateProductTemplate(productData) {
        try {
            const templateResponse = await fetch('http://localhost:8000/product', {credentials: 'include'});
            console.log("HTML fetch response status:", templateResponse.status);
            let newHtml = await templateResponse.text();
            console.log("HTML response text:", newHtml);

            const name = productData.name ?? 'No Name';
            const description = productData.description ?? 'No Description';
            const price = productData.price ?? 0;
            const inStock = productData.qtyInStock > 0 ? 'In Stock' : 'Out of Stock';

            // Populate the template with product data
            newHtml = newHtml.replace(/\{\{product\.name}}/g, name);
            newHtml = newHtml.replace(/\{\{product\.description}}/g, description);
            newHtml = newHtml.replace(/\{\{product\.price}}/g, price);
            newHtml = newHtml.replace(/\{\{product\.inStock}}/g, inStock);

            console.log('Final HTML:', newHtml);
            return newHtml;
        } catch (error) {
            console.error('Error populating product template:', error);
            return 'Error loading product details.';
        }
    }

    updateProductDisplay(products) {
        this.elements.productsContainer.innerHTML = '';
        if (products.length === 0) {
            const noResults = document.createElement('p');
            noResults.textContent = 'No results found.';
            this.elements.productsContainer.appendChild(noResults);
        } else {
            products.forEach(product => {
                const productElement = this.createProductElement(product);
                this.elements.productsContainer.appendChild(productElement);
            });
        }
    }

    updateBreadcrumbs(pathComponents) {
        const breadcrumbsContainer = document.querySelector('.breadcrumbs');
        breadcrumbsContainer.innerHTML = ''; // Clear previous breadcrumbs

        pathComponents.forEach((component, index) => {
            const link = document.createElement('a');
            link.textContent = component.text;
            link.href = component.href;
            breadcrumbsContainer.appendChild(link);

            if (index < pathComponents.length - 1) { // Don't add arrow on last component
                const separator = document.createTextNode(' > ');
                breadcrumbsContainer.appendChild(separator);
            }
        });
    }

    createProductElement(product) {
        // Creates the product container
        const productElement = document.createElement('div');
        productElement.className = 'product';
        productElement.dataset.productId = product.id;

        // Creates the left column for name, description, and stock
        const leftColumn = document.createElement('div');
        leftColumn.className = 'product-left-column';

        // Creates the name element
        const nameElement = document.createElement('div');
        nameElement.className = 'product-name';
        nameElement.textContent = 'Name: ' + product.name;

        // Creates the description element
        const descriptionElement = document.createElement('div');
        descriptionElement.className = 'product-description';
        descriptionElement.textContent = 'Description: ' + product.description;

        // Creates the stock element
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

        productElement.appendChild(leftColumn);
        productElement.appendChild(rightColumn);

        // The click event listener implementation starts here
        productElement.addEventListener('click', async (event) => {
            event.preventDefault();
            const productId = productElement.dataset.productId;
            console.log("ProductId clicked:", productId);

            try {
                const productData = await this.fetchProductDetails(productId);
                console.log("ProductData:", productData); // Check Product Data

                // Hide the main product list and show the product detail container
                this.elements.productsContainer.style.display = 'none';
                const productDetailsContainer = document.getElementById('product-details-container');
                productDetailsContainer.innerHTML = "";
                productDetailsContainer.style.display = 'block';

                // Add a back button to the product details container
                const backButton = document.createElement('button');
                backButton.textContent = 'Back to product list';
                backButton.addEventListener('click', () => {
                    productDetailsContainer.style.display = 'none'; // Hide the product details container
                    this.elements.productsContainer.style.display = ''; // Show the main product list
                });
                productDetailsContainer.appendChild(backButton);

                const populatedHtml = await this.populateProductTemplate(productData);
                console.log("PopulatedHTML", populatedHtml); // Check Populated HTML
                this.insertProductDetails(populatedHtml);
            } catch (error) {
                console.error('Error fetching and displaying product details:', error);
                // Display an error message to the user.
            }
        });
//...
        return productElement;
    }

    fetchProducts(query = "") {
        fetch(`/api/v1/products/search/${query}`)
            .then(response => this.handleFetchResponse(response))
            .then(data => this.updateProductDisplay(data))
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