const productId = window.location.search.slice(1).split('&')[0].split('=')[1]; // Assuming product ID is passed in the URL like ?id=123

function displayProduct(product) {
    console.log('Starting to display product', product);
    document.getElementById("product-title").textContent = product.name;
    document.getElementById("product-brand").textContent = product.brand;
    document.getElementById("product-description").textContent = product.description;
    document.getElementById("product-price").textContent = "$" + product.price;

    // Breadcrumbs
    document.getElementById("product-category").textContent = product.categoryName;
    document.getElementById("product-brand").textContent = product.brandName;

    // Availability
    const availabilityElement = document.getElementById('product-availability');
    if (product.qtyInStock > 0) {
        availabilityElement.textContent = "In Stock";
        availabilityElement.classList.add('in-stock');
    } else {
        availabilityElement.textContent = "Out of Stock";
        availabilityElement.classList.add('out-of-stock');
    }
    console.log('Finished displaying product');
}

fetch(`/api/v1/products/${productId}`)  // Assuming your API endpoint is like this
    .then(response => {
        if (!response.ok) {
            // Handle API errors gracefully (e.g., display an error message)
            console.error('Error fetching product:', response.statusText);
            return;
        }
        return response.json();
    })
    .then(product => {
        displayProduct(product);
    })
    .catch(error => console.error('Error:', error));
// }
//
