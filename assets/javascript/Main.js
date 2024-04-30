document.addEventListener('DOMContentLoaded', (event) => {
    const products = [
        {name: 'Product 1', description: 'This is product 1', price: '$100'},
        {name: 'Product 2', description: 'This is product 2', price: '$200'},
        // Add more products as required
    ];

    const productsContainer = document.querySelector('#products');

    products.forEach(product => {
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
});

window.onload = function() {
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
        });

    fetch('/api/categories')
        .then(response => {
            console.log('fetch response:', response); // to check the response
            return response.json();
        })
        .then(categories => {
            console.log('fetch data:', categories); // to check the data
            // rest of your code
        })
        .catch(error => console.error('Error fetching categories:', error));
    fetch('/api/categories')
        .then(response => response.json())
        .then(categories => {
            const dropdown = document.querySelector('#categories');
            categories.forEach(category => {
                const a = document.createElement('a');
                a.text = category.name;
                a.href = '/category/' + category.id; // adjust this line as necessary
                dropdown.appendChild(a);
            });
        });
};


function login() {
    // Get the modal
    var modal = document.getElementById("loginModal");
    // Display the modal
    modal.style.display = "block";
}

// When the user clicks on the close button, close the modal
function closeModal() {
    var modal = document.getElementById("loginModal");
    modal.style.display = "none";
}

// When the user clicks anywhere outside of the modal, close it
window.onclick = function(event) {
    var modal = document.getElementById("loginModal");
    if (event.target == modal) {
        modal.style.display = "none";
    }
}