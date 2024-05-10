// fetch product details function
async function fetchProductDetails(ProductID) {
    const response = await fetch(`/api/v1/products/${ProductID}`);
    if (!response.ok) {
        throw new Error(`Failed to fetch product details for ProductID: ${ProductID}`);
    }
    const productDetails = await response.json();
    return productDetails;
}

async function fetchAndPopulateCartSummary() {
    try {
        const cartItemsResponse = await fetch('/api/v1/cart/');
        if (!cartItemsResponse.ok) {
            throw new Error('Failed to fetch cart items');
        }

        const cartItems = await cartItemsResponse.json();

        const cartItemsElement = document.getElementById('cart-items');
        let subtotal = 0;

        for (const item of cartItems) {
            const productDetails = await fetchProductDetails(item.ProductID);

            // Create cart item elements
            const itemElement = document.createElement('div');
            itemElement.classList.add("cart-item");

            const productName = document.createElement('div');
            productName.textContent = `${productDetails.name}`;
            itemElement.appendChild(productName);

            const quantityInput = document.createElement('input');
            quantityInput.type = 'number';
            quantityInput.min = '0';
            quantityInput.value = `${item.Quantity}`;
            quantityInput.onchange = () => adjustQuantity(item.ProductID, quantityInput.value);

            const quantityContainer = document.createElement('div');
            quantityContainer.appendChild(quantityInput);
            itemElement.appendChild(quantityContainer);

            const priceElement = document.createElement('div');
            priceElement.textContent = `KR. ${productDetails.price.toFixed(2)}`;
            itemElement.appendChild(priceElement);

            const totalElement = document.createElement('div');
            totalElement.textContent = `KR. ${(productDetails.price * item.Quantity).toFixed(2)}`;
            itemElement.appendChild(totalElement);

            const removeBtn = document.createElement('span');
            removeBtn.classList.add("remove-btn");
            removeBtn.innerText = 'X';
            removeBtn.onclick = () => removeFromCart(item.ProductID);

            const removeContainer = document.createElement('div');
            removeContainer.appendChild(removeBtn);

            itemElement.appendChild(removeContainer);

            cartItemsElement.appendChild(itemElement);

            subtotal += productDetails.price * item.Quantity;
        }

        document.getElementById('subtotal').innerText = `Subtotal: KR. ${subtotal.toFixed(2)}`;

    } catch (error) {
        console.error('Error fetching cart summary:', error);
    }
}

// These functions should call your API or perform the necessary actions to adjust the cart
function removeFromCart(productId) {
    console.log("Remove from cart:", productId);
    // Implement your logic to remove the item from cart
}

function adjustQuantity(productId, newQuantity) {
    console.log("Adjust quantity:", productId, newQuantity);
    // Implement your logic to adjust the quantity of an item in the cart
}

document.addEventListener('DOMContentLoaded', fetchAndPopulateCartSummary);