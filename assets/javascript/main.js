class EStore {
    constructor() {
        this.elements = {
            productsContainer: document.getElementById('products'),
            searchForm: document.getElementById('search-form'),
            searchQueryInput: document.getElementById('search-query'),
            logoutButton: document.getElementById('logoutButton'),
            loginButton: document.getElementById('loginButton'),
            closeModalButton: document.querySelector('.close'),
            sortDropdown: document.getElementById('sort-dropdown'),
        }
        this.checkLoginState().then(this.fetchProducts.bind(this));
        this.addEventListeners();
        this.brandClicked = '';
        this.categoryClicked = '';
        this.fetchCategories();
        this.addEventListenersProduct();
        this.fetchBrands();
        this.updateCartBadge();
        this.allBrands = []; // Inside the EStore class constructor
    }

    addEventListeners() {
        this.elements.searchForm.addEventListener('submit', this.handleSearch.bind(this));
        this.elements.logoutButton.addEventListener('click', this.logoutUser.bind(this));
        this.elements.loginButton.addEventListener('click', this.openLoginModal.bind(this));
        this.elements.closeModalButton.addEventListener('click', this.closeModal.bind(this));
        this.elements.closeModalButton.addEventListener('click', this.closeModal.bind(this));
        this.elements.sortDropdown.addEventListener('change', this.handleSort.bind(this));

        // Event delegation for "Add to Cart" buttons on both product listings and details
        document.addEventListener('click', (event) => {
            const target = event.target.closest('.add-to-cart-button');
            if (target) {
                event.preventDefault();
                this.handleAddToCart(target.dataset.productId); // Pass product ID directly
            }
        });

        // Cart Icon Click Handler
        document.querySelector('#cart a').addEventListener('click', (event) => {
            this.handleCartIconClick(event); // Pass the event object
        });

    }

    handleSearch(event) {
        event.preventDefault();
        const query = this.elements.searchQueryInput.value.trim();
        if (query !== "") {
            this.fetchProducts(query);
            this.updateBreadcrumbs([
                { text: 'Home', href: '/', isCategory: false },
                { text: 'Search Results: ' + query, href: '#', isCategory: false }
            ]);
        } else {
            alert("Please enter a search query.");
        }
    }

    async handleCartIconClick() {
        const productArea = document.getElementById('products').parentElement;  // Get parent of #products
        const productDetailsContainer = document.getElementById('product-details-container');
        const productFilters = document.getElementById('product-filters');
        const productsHeader = document.querySelector('.products-header'); // Get the products header

        try {
            // 1. Show loading indicator immediately
            productArea.innerHTML = '<div class="loading-indicator">Loading cart...</div>';

            const response = await fetch('/cart');
            if (!response.ok) {
                throw new Error('Failed to fetch cart content');
            }
            const cartHtmlContent = await response.text();
            const existingCart = document.getElementById('cart-summary');
            if (existingCart) {
                // Toggle existing cart's visibility
                existingCart.style.display = existingCart.style.display === 'none' ? 'block' : 'none';
                if (existingCart.style.display === 'block') {
                    await fetchAndPopulateCartSummary(existingCart.querySelector('#cart-items'));
                    productsHeader.style.display = 'none'; // Hide products-header when cart is shown
                } else {
                    // Remove cart and re-fetch products
                    existingCart.remove();
                    await this.fetchProducts();

                    // Re-create the product-details-container if it's missing
                    const productDetailsContainer = document.getElementById('product-details-container');
                    if (!productDetailsContainer) {
                        const newProductDetailsContainer = document.createElement('section');
                        newProductDetailsContainer.id = 'product-details-container';
                        newProductDetailsContainer.style.display = 'none'; // Initially hidden
                        productArea.appendChild(newProductDetailsContainer);
                    }
                    if (productDetailsContainer) {
                        productDetailsContainer.style.display = 'block';
                    }
                    if (productFilters) {
                        productFilters.style.display = 'block';
                    }
                }
            } else {
                // Create a new container for the cart
                const cartContainer = document.createElement('div');
                cartContainer.id = 'cart-summary';

                // Set innerHTML before appending to the DOM
                cartContainer.innerHTML = cartHtmlContent;

                // Add cartContainer to the productArea
                productArea.innerHTML = '';
                productArea.appendChild(cartContainer);

                // Hide other sections (if needed)
                if (productDetailsContainer) productDetailsContainer.style.display = 'none';
                if (productFilters) productFilters.style.display = 'none';
                productsHeader.style.display = 'none'; // Hide the product header when the cart is shown


                // Wait for the cart items to be added to the DOM
                await new Promise(resolve => {
                    const checkForCartItems = () => {
                        const cartItemsElement = document.getElementById('cart-items');
                        if (cartItemsElement) {
                            resolve(cartItemsElement);
                        } else {
                            setTimeout(checkForCartItems, 10);
                        }
                    };
                    checkForCartItems();
                }).then(async (cartItemsElement) => {
                    await fetchAndPopulateCartSummary(cartItemsElement);
                });
            }

        } catch (error) {
            console.error('Error fetching cart content:', error);
            productArea.innerHTML = '<div class="error-message">Error loading cart. Please try again.</div>';
        }
        this.updateCartBadge();
    }

    async updateCartBadge() {
        try {
            const response = await fetch('/api/v1/cart/');
            if (!response.ok) {
                throw new Error('Failed to fetch cart items');
            }
            const cartItems = await response.json();

            const count = cartItems.length;  // Get the length of the array
            const cartBadge = document.getElementById('cart-badge');
            if (cartBadge) {
                cartBadge.textContent = count;
                cartBadge.style.display = count > 0 ? 'block' : 'none';
            }
        } catch (error) {
            console.error('Error updating cart badge:', error);
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

    fetchBrands() {
        fetch('/api/v1/brands/')
            .then(this.handleFetchResponse)
            .then(brands => {
                this.allBrands = brands;
                this.renderBrandLinks(); // Populate the dropdown after fetching brands
            })
            .catch(error => console.error('Error fetching brands:', error));
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

    renderBrandLinks(brands) {
        const dropdown = document.querySelector('#brands');
        dropdown.innerHTML = '';
        this.allBrands.forEach(brand => {
            const a = document.createElement('a');
            a.textContent = brand.name;
            const brandName = encodeURIComponent(brand.name);
            a.href = '/brands/' + brandName;
            a.dataset.brandName = brand.name;
            a.addEventListener('click', this.handleBrandClick.bind(this));
            dropdown.appendChild(a);
        });
    }



    async handleCategoryClick(event) {
        event.preventDefault();

        const categoryName = event.target.dataset.categoryName;

        // 1. Show Product Section, Update Content, and Hide Cart and Details (if visible)
        await this.toggleCartAndProducts('products'); // Toggle to products

        // 2. Fetch products by category
        this.fetchProductsByCategory(categoryName);

        // 3. Update breadcrumbs and history
        this.updateBreadcrumbs([
            { text: 'Home', href: '/', type: 'home' },
            { text: categoryName, href: `/${categoryName}`, type: 'category' }
        ]);
        history.pushState(null, null, `/${categoryName}`);
    }

    async toggleCartAndProducts(showElementId) {
        const productsSection = document.getElementById('products');
        const existingCart = document.getElementById('cart-summary');
        const productDetailsContainer = document.getElementById('product-details-container');
        const productFilters = document.getElementById('product-filters');
        const productsHeader = document.querySelector('.products-header');

        if (showElementId === 'cart-summary') {
            // ... (Logic for showing the cart remains the same)
        } else { // showElementId is 'products'
            if (existingCart) {
                existingCart.remove();
            }

            // If productArea does not exist, re-create it
            if (!productsSection) {
                const productAreaParent = document.querySelector('.container');
                const productAreaElement = document.createElement('div');
                productAreaElement.id = 'products';
                productAreaElement.className = 'product-grid';
                productAreaParent.appendChild(productAreaElement);
            }

            // 1. Fetch products FIRST before updating references
            await this.fetchProducts();

            // Update the reference to the new productArea
            this.elements.productsContainer = document.getElementById('products'); // Update the reference

            // 2. Update breadcrumbs and history (if needed)
            this.updateBreadcrumbs([
                { text: 'Home', href: '/', type: 'home' },
                // ... (rest of your breadcrumb logic)
            ]);
            history.pushState(null, null, `/`);

            document.querySelector('.products-header').style.display = 'block';
            productDetailsContainer.style.display = 'block';
        }
    }









    handleBrandClick(event) {
        event.preventDefault();
        const brandName = event.target.dataset.brandName; // Assuming 'brandName' dataset
        const productDetailsContainer = document.getElementById('product-details-container');
        const product = document.getElementById('products');
        productDetailsContainer.style.display = 'none'; // Hide product details if visible
        productDetailsContainer.innerHTML = '';
        product.style.display = 'none';
        product.innerHTML = '';
        this.elements.productsContainer.style.display = 'block';


        this.fetchProductsByBrand(brandName);


        // Update breadcrumbs
        this.updateBreadcrumbs([
            { text: 'Home', href: '/', type: 'home' },
            { text: brandName, href: `/products/brand/${brandName}`, type: 'brand' }
        ]);
        history.pushState(null, null, `/products/brand/${brandName}`);
    }


    fetchProductsByBrandAndCategory(brandName, categoryName) {
        fetch(`/api/v1/products/brand/${brandName}/category/${categoryName}`)
            .then(this.handleFetchResponse)
            .then(products => {
                // Filter products to include only those from the selected brand and category
                const filteredProducts = products.filter(product => product.brandName === brandName && product.categoryName === categoryName);
                this.updateProductDisplay(filteredProducts);
            })
            .catch(error => console.error('Error fetching products:', error));
    }


    async fetchProductDetails(productId) { // Removed 'category' from parameters
        try {
            const dataResponse = await fetch(`/api/v1/products/category/${this.categoryClicked}/brand/${productId}`);
            const productsData = await dataResponse.json();
            const productData = productsData.find(product => product.id === productId);
            if (!productData) throw new Error(`Product with id ${productId} was not found`);
            console.log("Product Data:", productData);

            // Fetch category from product data instead of function parameter
            const category = productData.categoryName;
            this.brandClicked = productData.brandName; // Store the brand name

            this.updateBreadcrumbs([
                { text: 'Home', href: '/', type: 'home' },
                { text: productData.brandName, href: `/products/brand/${productData.brandName}`, type: 'brand' },
                { text: productData.categoryName, href: `/products/category/${productData.categoryName}`, type: 'category' },
                { text: productData.name, href: `/products/category/${productData.categoryName}/brand/${productData.brandName}/product/${productId}`, type: 'product' }
            ]);

            const populatedHtml = await this.populateProductTemplate(productData);
            this.insertProductDetails(populatedHtml);
            return productData;
        } catch (error) {
            console.error('Error fetching product details:', error);
            return Promise.reject(error);
        }
    }

    fetchProductsByCategory(categoryName) {
        fetch(`/api/v1/products/category/${categoryName}`)
            .then(this.handleFetchResponse)
            .then(this.updateProductDisplay.bind(this))
            .catch(error => console.error('Error fetching products:', error));
    }

    fetchProductsByBrand(brandName) {
        fetch(`/api/v1/products/brand/${brandName}`)
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

    // Helper function to insert content
    insertProductDetails(content) {
        const container = document.getElementById('product-details-container');
        container.innerHTML = ''; // Clear existing content
        container.appendChild(content);
        container.style.display = 'block';

        // Update button text after the product details have been rendered
        const addToCartButton = container.querySelector('.add-to-cart-button');
        if (addToCartButton) {
            addToCartButton.innerHTML = this.isUserLoggedIn ? 'Add To Cart' : 'Login to Add To Cart';
        }
    }


    async populateProductTemplate(productData) {
        try {
            const templateResponse = await fetch('http://localhost:8000/product', {credentials: 'include'});
            console.log("HTML fetch response status:", templateResponse.status);
            let newHtml = await templateResponse.text();
            console.log("HTML response text:", newHtml);

            // Create elements dynamically
            const section = document.createElement('section');
            section.className = 'product-details';
            const div = document.createElement('div');
            div.className = 'product-info';
            const title = document.createElement('h1');
            title.id = 'product-title';
            title.textContent = productData.name;
            const availability = document.createElement('p');
            availability.className = "availability";
            availability.id = "product-availability";
            availability.textContent = productData.qtyInStock > 0 ? 'In Stock' : 'Out of Stock';
            const price = document.createElement('p');
            price.className = "price";
            price.id = "product-price";
            price.textContent = productData.price;
            const description = document.createElement('p');
            description.className = "description";
            description.id = "product-description";
            description.textContent = productData.description;

            const button = document.createElement('button');
            button.className = 'add-to-cart-button';
            button.dataset.productId = productData.id;
            // Use the `isUserLoggedIn` property here to set the button text dynamically
            button.textContent = this.isUserLoggedIn ? 'Add To Cart' : 'Login to Add To Cart';

            // Append elements to create structure
            div.appendChild(title);
            div.appendChild(availability);
            div.appendChild(price);
            div.appendChild(description);
            div.appendChild(button);
            section.appendChild(div);

            return section; // Return the dynamically created section

        } catch (error) {
            console.error('Error populating product template:', error);
            return 'Error loading product details.';
        }
    }

    updateProductDisplay(products) {
        this.elements.productsContainer.innerHTML = "";
        products.forEach(product => {
            const productElement = this.createProductElement(product);
            this.elements.productsContainer.appendChild(productElement);
        });
    }

    async updateBreadcrumbs(pathComponents) {
        const breadcrumbsContainer = document.querySelector('.breadcrumbs');

        if (pathComponents.length > 1) {
            breadcrumbsContainer.innerHTML = '';

            for (const [index, component] of pathComponents.entries()) {
                const link = document.createElement('a');
                link.textContent = component.text;
                link.href = component.href;


                link.addEventListener("click", async (event) => {
                    event.preventDefault();

                    const productDetailsContainer = document.getElementById('product-details-container');
                    productDetailsContainer.style.display = 'none';
                    this.elements.productsContainer.style.display = 'block';

                    if (component.type === 'home') {
                        await this.fetchProducts();
                    } else if (component.type === 'brand') {
                        const brand = component.text;
                        await this.fetchProductsByBrand(brand);
                    }         else if (component.type === 'category') {
                        const category = component.text;
                        if (this.brandClicked) {
                            await this.fetchProductsByBrandAndCategory(this.brandClicked, category);
                        } else {
                            await this.fetchProductsByCategory(category);
                        }

                        // After fetching products by that category,
                        // fetch all products and filter them by selected category to get unique brands.
                        const allProducts = await this.fetchProducts();
                        const brands = this.getUniqueBrands(allProducts, category);
                        this.renderBrandLinks(brands);
                    }

                    // Update breadcrumbs to only include up to current level
                    this.updateBreadcrumbs(pathComponents.slice(0, index + 1));
                });

                breadcrumbsContainer.appendChild(link);
                if (index !== pathComponents.length - 1) {
                    breadcrumbsContainer.appendChild(document.createTextNode(' > '));
                }
            }
        } else {
            breadcrumbsContainer.innerHTML = '';
        }
    }

    getUniqueBrands(products, category) {
        // Add error handling to check if 'products' is undefined
        if (!products) {
            console.error("'products' is undefined in getUniqueBrands function");
            return [];
        }
        const categoryProducts = products.filter(product => product.categoryName === category);
        const uniqueBrands = [];
        categoryProducts.forEach(product => {
            if (!uniqueBrands.find(brand => brand.name === product.brandName)) {
                uniqueBrands.push({ name: product.brandName });
            }
        });
        return uniqueBrands;
    }


    async handleAddToCart(productId) {  // Change argument to directly receive productId
        try {
            if (!this.isUserLoggedIn) {
                this.openLoginModal();
                return;
            }
            console.log('Product ID to add:', productId);
            const quantity = 1;

            const data = {
                productID: productId,
                quantity: quantity
            };
            console.log('Data to be sent:', data);

            const response = await fetch('/api/v1/cart/', {
                method: 'POST',
                mode: 'cors',
                credentials: 'same-origin',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(data)
            });

            console.log('Response status:', response.status);

            if (!response.ok) {
                const errorData = await response.json();
                const errorMessage = errorData.error || "Failed to add to cart.";
                throw new Error(errorMessage);
            } else {
                this.updateCartBadge();
            }
        } catch (error) {
            console.error('Error adding to cart:', error);
            alert("An error occurred while adding the product to your cart.");
        }
    }


    createProductElement(product) {
        // Creates the product container
        const productElement = document.createElement('div');
        productElement.className = 'product';
        productElement.dataset.productId = product.id;
        productElement.dataset.productCategory = product.categoryName;

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
        buttonElement.innerHTML = this.isUserLoggedIn ? 'Add To Cart' : 'Login to Add To Cart';

        // Append price and button to the right column
        rightColumn.appendChild(priceElement);
        rightColumn.appendChild(buttonElement);

        productElement.appendChild(leftColumn);
        productElement.appendChild(rightColumn);

        return productElement;
    }

    addEventListenersProduct() {
        this.elements.productsContainer.addEventListener('click', async (event) => {
            // Getting the closest product from the clicked target
            const productElement = event.target.closest('.product');
            const clickedButton = event.target.closest('button');

            // If a product was clicked and it was not a button inside a product
            if (productElement && !clickedButton) {
                event.preventDefault();

                // If the clicked product is the same as the previously clicked product, we return
                if (this.clickedProduct && this.clickedProduct.id === productElement.dataset.productId) {
                    return;
                }

                const productId = productElement.dataset.productId;
                const productCategory = productElement.dataset.productCategory;
                console.log(productCategory);

                // Fetch and display product details
                try {
                    const productData = await this.fetchProductDetails(productId, productCategory);

                    // Hide the main product list and show the product detail container
                    this.elements.productsContainer.style.display = 'none';
                    const productDetailsContainer = document.getElementById('product-details-container');
                    productDetailsContainer.innerHTML = "";
                    productDetailsContainer.style.display = 'block';

                    const populatedHtml = await this.populateProductTemplate(productData);
                    console.log("PopulatedHTML", populatedHtml); // Check Populated HTML
                    this.insertProductDetails(populatedHtml);
                } catch (error) {
                    console.error('Error fetching and displaying product details:', error);
                    // Display an error message to the user.
                }
            }
        });
    }

    fetchProducts(query = "") {
        return fetch(`/api/v1/products/search/${query}`)
            .then(response => this.handleFetchResponse(response))
            .then(data => {
                this.updateProductDisplay(data);
                return data;
            })
            .catch(error => {
                console.error('Error fetching products:', error);
                return [];
            });
    }

    async checkLoginState() {
        try {
            const response = await fetch('/api/check_login', { credentials: 'include' });
            const respJson = await response.json();
            this.isUserLoggedIn = respJson.logged_in;

            const userNotLogged = document.getElementById('user-not-logged');
            const userLogged = document.getElementById('user-logged');
            const logoutButton = document.getElementById('logoutButton');
            const cartIcon = document.getElementById('cart');
            console.log("isUserLoggedIn:", this.isUserLoggedIn);

            if (this.isUserLoggedIn) {
                document.getElementById('logged-username').textContent = respJson.username;
                userNotLogged.style.display = 'none';
                userLogged.style.display = 'block';
                logoutButton.style.display = 'block';
                cartIcon.style.display = 'block';
            } else {
                userNotLogged.style.display = 'block';
                userLogged.style.display = 'none';
                logoutButton.style.display = 'none';
                cartIcon.style.display = 'none';
            }
        } catch (error) {
            console.error('Error checking login:', error);
        }
    }

    logoutUser() {
        fetch('/api/logout', { credentials: 'include' })
            .then(response => response.text())
            .then(text => {
                this.checkLoginState();
                window.location.href = "/";
                // Remove event listener for login button
                this.elements.loginButton.removeEventListener('click', this.openLoginModal.bind(this));
                // Add event listener for login button
                this.elements.loginButton.addEventListener('click', this.openLoginModal.bind(this));
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

async function updateCartBadge() {
    try {
        const response = await fetch('/api/v1/cart/');
        if (!response.ok) {
            throw new Error('Failed to fetch cart items');
        }
        const cartItems = await response.json();

        const count = cartItems.length;  // Get the length of the array
        const cartBadge = document.getElementById('cart-badge');
        if (cartBadge) {
            cartBadge.textContent = count;
            cartBadge.style.display = count > 0 ? 'block' : 'none';
        }
    } catch (error) {
        console.error('Error updating cart badge:', error);
    }
}
