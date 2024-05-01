document.addEventListener('DOMContentLoaded', (event) => {
    const products = [
        {name: 'Product 1', description: 'This is product 1', price: '$100'},
        {name: 'Product 2', description: 'This is product 2', price: '$200'},
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

function login() {
    var modal = document.getElementById("loginModal");
    modal.style.display = "block";
}

function closeModal() {
    var modal = document.getElementById("loginModal");
    modal.style.display = "none";
}

function checkLoginState() {
    fetch('/api/check_login', { credentials: 'include' })
        .then(response => response.json())
        .then(respJson => {
            if (respJson.logged_in) {
                document.getElementById('logged-username').textContent = respJson.username;
                document.getElementById('user-not-logged').style.display = 'none';
                document.getElementById('user-logged').style.display = 'block';
            } else {
                document.getElementById('user-not-logged').style.display = 'block';
                document.getElementById('user-logged').style.display = 'none';
            }
        }).catch(error => console.error('Error checking login:', error));
}

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
            console.log('fetch response:', response);
            return response.json();
        })
        .then(categories => {
            console.log('fetch data:', categories);
            const dropdown = document.querySelector('#categories');
            categories.forEach(category => {
                const a = document.createElement('a');
                a.text = category.name;
                a.href = '/category/' + category.id;
                dropdown.appendChild(a);
            });
        })
        .catch(error => console.error('Error fetching categories:', error));

    checkLoginState();
}

window.onclick = function(event) {
    var modal = document.getElementById("loginModal");
    if (event.target == modal) {
        modal.style.display = "none";
    }
}
