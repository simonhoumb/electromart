// fetch product details function (unchanged)
async function fetchProductDetails(ProductID) {
    const response = await fetch(`/api/v1/products/${ProductID}`);
    if (!response.ok) {
        throw new Error(`Failed to fetch product details for ProductID: ${ProductID}`);
    }
    const productDetails = await response.json();
    return productDetails;
}

async function fetchAndPopulateCartSummary() {
    const cartItemsElement = document.getElementById('cart-items');
    const subtotalElement = document.getElementById('subtotal');

    try {
        // 1. Fetch Cart Items from API
        const cartItemsResponse = await fetch('/api/v1/cart/');
        if (!cartItemsResponse.ok) {
            throw new Error('Failed to fetch cart items');
        }

        // 2. Parse response, handling potential null value
        let cartItems = await cartItemsResponse.json();
        cartItems = cartItems || [];  // If cartItems is null, set it to an empty array

        // 3. Handle Empty Cart
        cartItemsElement.innerHTML = ''; // Always clear existing content
        if (cartItems.length === 0) {
            cartItemsElement.innerHTML = '<div class="no-items-message">No items in cart.</div>';
            subtotalElement.textContent = 'Subtotal: KR. 0.00';
            return; // No need to fetch product details if cart is empty
        }

        // 4. Fetch Product Details Concurrently
        const productDetailsMap = {};
        const allProductDetailsPromises = cartItems.map(item => fetchProductDetails(item.ProductID));
        const allProductDetails = await Promise.allSettled(allProductDetailsPromises);

        allProductDetails.forEach((result, index) => {
            if (result.status === 'fulfilled') {
                productDetailsMap[cartItems[index].ProductID] = result.value;
            } else {
                console.error(`Failed to fetch details for ProductID: ${cartItems[index].ProductID}`);
            }
        });

        // 5. Populate Cart Items and Calculate Subtotal
        let subtotal = 0;
        cartItems.forEach(item => {
            const productDetails = productDetailsMap[item.ProductID];
            if (productDetails) {
                const itemElement = createCartItemElement(productDetails, item.Quantity);
                cartItemsElement.appendChild(itemElement);
                subtotal += productDetails.price * item.Quantity;
            }
        });

        // 6. Update Subtotal Display
        subtotalElement.innerText = `Subtotal: KR. ${subtotal.toFixed(2)}`;

    } catch (error) {
        console.error('Error fetching cart summary:', error);
        cartItemsElement.innerHTML = `<div class="error-message">Error loading cart. Please try again.</div>`;
    } finally {
        // 7. Remove Loading Indicator (if it exists)
        const loadingIndicator = cartItemsElement.querySelector('.loading-indicator');
        if (loadingIndicator) {
            loadingIndicator.remove();
        }
    }
}

function createCartItemElement(productDetails, quantity) {
    // Create the container for the cart item
    const itemElement = document.createElement('div');
    itemElement.classList.add('cart-item');
    itemElement.dataset.productId = productDetails.id;

    // Create elements for each piece of information
    const nameElement = document.createElement('div');
    nameElement.textContent = productDetails.name;

    const quantityInput = document.createElement('input');
    quantityInput.type = 'number';
    quantityInput.min = '1';
    quantityInput.value = quantity;
    quantityInput.addEventListener('input', () => {
        let newQuantity = parseInt(quantityInput.value, 10);
        if (isNaN(newQuantity) || newQuantity < 1) {
            newQuantity = 1;
            quantityInput.value = 1; // Reset to 1 if invalid input
        }
        adjustQuantity(productDetails.id, newQuantity);
    });

    const priceElement = document.createElement('div');
    priceElement.textContent = 'KR. ' + productDetails.price.toFixed(2);

    const totalElement = document.createElement('div');
    totalElement.textContent = 'KR. ' + (productDetails.price * quantity).toFixed(2);

    const removeButton = document.createElement('button');
    removeButton.textContent = 'Remove';
    removeButton.classList.add('remove-btn');
    removeButton.onclick = () => removeFromCart(productDetails.id);

    // Assemble the cart item
    itemElement.appendChild(nameElement);
    itemElement.appendChild(quantityInput);
    itemElement.appendChild(priceElement);
    itemElement.appendChild(totalElement);
    itemElement.appendChild(removeButton);

    return itemElement;
}



async function removeFromCart(productId) {
    try {
        const response = await fetch(`/api/v1/cart?productID=${productId}`, {
            method: 'DELETE',
        });

        if (response.ok) {
            // Item successfully deleted

            // Clear the cart display and repopulate it
            const cartItemsElement = document.getElementById('cart-items');
            cartItemsElement.innerHTML = ''; // Clear existing items

            // Re-fetch and populate the cart summary
            await fetchAndPopulateCartSummary();
            await updateCartBadge();
            ;
        } else {
            console.error('Failed to remove item from cart:', response.statusText);
            // Handle error, e.g., show a message to the user
        }
    } catch (error) {
        console.error('Error:', error);
        // Handle network or other errors
    }
}

async function adjustQuantity(productId, newQuantity) {
    try {
        // 1. Input Validation and Sanitization:
        newQuantity = parseInt(newQuantity, 10);
        if (isNaN(newQuantity) || newQuantity < 1) {
            newQuantity = 1;
            const quantityInput = document.querySelector(`[data-product-id="${productId}"] .quantity-input`);
            if (quantityInput) {
                quantityInput.value = newQuantity;
            }
        }

        if (newQuantity === 0) {
            // 2. Handle Removal When Quantity is Zero:
            await removeFromCart(productId);
        } else {
            // 3. Update Quantity on Server:
            const response = await fetch(`/api/v1/cart?productID=${productId}`, {
                method: 'PATCH',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ newQuantity }),
            });

            if (response.ok) {
                // 4. Update Cart Item Element:
                const itemElement = document.querySelector(`[data-product-id="${productId}"]`);
                if (itemElement) {
                    const totalElement = itemElement.querySelector('div:nth-child(4)');

                    // Fetch updated product details (in case price has changed):
                    const productDetails = await fetchProductDetails(productId);
                    totalElement.textContent = `KR. ${(productDetails.price * newQuantity).toFixed(2)}`;

                    // 5. Update Subtotal:
                    await fetchAndPopulateCartSummary(); // Refresh the entire cart summary
                }
            } else {
                console.error('Failed to adjust quantity:', response.statusText);
                // ... (error handling)
            }
        }
    } catch (error) {
        console.error('Error:', error);
        // ... (error handling)
    }
}

