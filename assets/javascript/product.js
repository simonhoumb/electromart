document.addEventListener("DOMContentLoaded", function() {
    const params = new URLSearchParams(window.location.search);
    const productId = params.get('id');

    if (productId) {
        fetchProduct(productId);
    }
});

function fetchProduct(productId) {
    fetch("/api/v1/products?id=" + productId)
        .then(response => {
            if (!response.ok) {
                throw new Error("Network response was not ok");
            }
            return response.json();
        })
        .then(product => {
            displayProduct(product);
        })
        .catch(error => {
            console.error("Error fetching product:", error);
        });
}

function displayProduct(product) {
    // Update HTML elements with product details
    document.getElementById("product-title").textContent = product.Name;
    document.getElementById("product-category").textContent = product.Category;
    document.getElementById("product-name").textContent = product.Name;
    document.getElementById("product-description").textContent = product.Description;
    document.getElementById("product-brand").textContent = "Brand: " + product.BrandName;
    document.getElementById("product-price").textContent = "$" + product.Price;
    // Add more elements as needed
}
