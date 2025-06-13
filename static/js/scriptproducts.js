document.addEventListener('DOMContentLoaded', function() {
    // Mobile Menu Toggle
    const mobileMenuBtn = document.querySelector('.mobile-menu-btn');
    const mainNav = document.querySelector('.main-nav');
    
    mobileMenuBtn.addEventListener('click', function() {
        mainNav.classList.toggle('active');
    });
    
    // Hero Slider
    const slides = document.querySelectorAll('.slide');
    const prevBtn = document.querySelector('.slider-controls .prev');
    const nextBtn = document.querySelector('.slider-controls .next');
    let currentSlide = 0;
    
    function showSlide(index) {
        slides.forEach(slide => slide.classList.remove('active'));
        slides[index].classList.add('active');
    }
    
    function nextSlide() {
        currentSlide = (currentSlide + 1) % slides.length;
        showSlide(currentSlide);
    }
    
    function prevSlide() {
        currentSlide = (currentSlide - 1 + slides.length) % slides.length;
        showSlide(currentSlide);
    }
    
    nextBtn.addEventListener('click', nextSlide);
    prevBtn.addEventListener('click', prevSlide);
    
    // Auto slide change every 5 seconds
    setInterval(nextSlide, 5000);
    
    // Testimonial Slider
    const testimonials = document.querySelectorAll('.testimonial');
    const dots = document.querySelectorAll('.dot');
    let currentTestimonial = 0;
    
    function showTestimonial(index) {
        testimonials.forEach(testimonial => testimonial.classList.remove('active'));
        dots.forEach(dot => dot.classList.remove('active'));
        
        testimonials[index].classList.add('active');
        dots[index].classList.add('active');
    }
    
    dots.forEach((dot, index) => {
        dot.addEventListener('click', () => {
            currentTestimonial = index;
            showTestimonial(currentTestimonial);
        });
    });
    
    // Auto testimonial change every 5 seconds
    setInterval(() => {
        currentTestimonial = (currentTestimonial + 1) % testimonials.length;
        showTestimonial(currentTestimonial);
    }, 5000);
    
    // Shopping Cart Functionality
    const cartIcon = document.getElementById('cart-icon');
    const cartSidebar = document.getElementById('cart-sidebar');
    const closeCart = document.getElementById('close-cart');
    const cartOverlay = document.getElementById('cart-overlay');
    const cartItemsContainer = document.getElementById('cart-items');
    const cartCount = document.querySelector('.cart-count');
    const totalPriceElement = document.querySelector('.total-price');
    
    let cart = JSON.parse(localStorage.getItem('cart')) || [];
    
    // Sample product data
    const products = [
        {
            id: 1,
            name: 'Wireless Headphones',
            category: 'Electronics',
            price: 99.99,
            oldPrice: 129.99,
            image: 'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=1470&q=80',
            badge: 'Sale'
        },
        {
            id: 2,
            name: 'Smart Watch',
            category: 'Electronics',
            price: 199.99,
            oldPrice: 249.99,
            image: 'https://images.unsplash.com/photo-1523275335684-37898b6baf30?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=1399&q=80',
            badge: 'Popular'
        },
        {
            id: 3,
            name: 'Running Shoes',
            category: 'Men\'s Fashion',
            price: 79.99,
            oldPrice: 89.99,
            image: 'https://images.unsplash.com/photo-1542291026-7eec264c27ff?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=1470&q=80',
            badge: 'New'
        },
        {
            id: 4,
            name: 'Denim Jacket',
            category: 'Women\'s Fashion',
            price: 59.99,
            oldPrice: 69.99,
            image: 'https://images.unsplash.com/photo-1591047139829-d91aecb6caea?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=1636&q=80',
            badge: 'Trending'
        },
        {
            id: 5,
            name: 'Leather Backpack',
            category: 'Accessories',
            price: 89.99,
            oldPrice: 99.99,
            image: 'https://images.unsplash.com/photo-1553062407-98eeb64c6a62?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=1374&q=80',
            badge: 'Sale'
        },
        {
            id: 6,
            name: 'Bluetooth Speaker',
            category: 'Electronics',
            price: 69.99,
            oldPrice: 79.99,
            image: 'https://images.unsplash.com/photo-1572569511254-d8f925fe2cbb?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=1377&q=80',
            badge: 'Popular'
        },
        {
            id: 7,
            name: 'Sunglasses',
            category: 'Accessories',
            price: 49.99,
            oldPrice: 59.99,
            image: 'https://images.unsplash.com/photo-1511499767150-a48a237f0083?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=1480&q=80',
            badge: 'New'
        },
        {
            id: 8,
            name: 'Fitness Tracker',
            category: 'Electronics',
            price: 79.99,
            oldPrice: 89.99,
            image: 'https://images.unsplash.com/photo-1576243345690-4e4b79b63288?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=1470&q=80',
            badge: 'Trending'
        }
    ];
    
    // Display products
    function displayProducts() {
        const productGrid = document.getElementById('product-grid');
        productGrid.innerHTML = '';
        
        products.forEach(product => {
            const productCard = document.createElement('div');
            productCard.className = 'product-card';
            
            productCard.innerHTML = `
                ${product.badge ? `<span class="product-badge">${product.badge}</span>` : ''}
                <div class="product-thumb">
                    <img src="${product.image}" alt="${product.name}">
                    <div class="product-actions">
                        <button class="quick-view"><i class="far fa-eye"></i></button>
                        <button class="add-to-wishlist"><i class="far fa-heart"></i></button>
                    </div>
                </div>
                <div class="product-info">
                    <span class="product-category">${product.category}</span>
                    <h3 class="product-title">${product.name}</h3>
                    <div class="product-price">
                        <span class="current-price">$${product.price.toFixed(2)}</span>
                        ${product.oldPrice ? `<span class="old-price">$${product.oldPrice.toFixed(2)}</span>` : ''}
                    </div>
                    <button class="add-to-cart" data-id="${product.id}">Add to Cart</button>
                </div>
            `;
            
            productGrid.appendChild(productCard);
        });
    }
    
    // Add to cart
    function addToCart(productId) {
        const product = products.find(p => p.id === productId);
        
        if (!product) return;
        
        const existingItem = cart.find(item => item.id === productId);
        
        if (existingItem) {
            existingItem.quantity += 1;
        } else {
            cart.push({
                id: product.id,
                name: product.name,
                price: product.price,
                image: product.image,
                quantity: 1
            });
        }
        
        updateCart();
        saveCartToLocalStorage();
    }
    
    // Remove from cart
    function removeFromCart(productId) {
        cart = cart.filter(item => item.id !== productId);
        updateCart();
        saveCartToLocalStorage();
    }
    
    // Update quantity
    function updateQuantity(productId, newQuantity) {
        const item = cart.find(item => item.id === productId);
        
        if (item) {
            item.quantity = parseInt(newQuantity) || 1;
        }
        
        updateCart();
        saveCartToLocalStorage();
    }
    
    // Update cart UI
    function updateCart() {
        // Update cart count
        const totalItems = cart.reduce((total, item) => total + item.quantity, 0);
        cartCount.textContent = totalItems;
        
        // Update cart items
        if (cart.length === 0) {
            cartItemsContainer.innerHTML = `
                <div class="empty-cart">
                    <i class="fas fa-shopping-cart"></i>
                    <p>Your cart is empty</p>
                </div>
            `;
            totalPriceElement.textContent = '$0.00';
        } else {
            cartItemsContainer.innerHTML = '';
            
            let totalPrice = 0;
            
            cart.forEach(item => {
                const itemTotal = item.price * item.quantity;
                totalPrice += itemTotal;
                
                const cartItem = document.createElement('div');
                cartItem.className = 'cart-item';
                
                cartItem.innerHTML = `
                    <div class="cart-item-img">
                        <img src="${item.image}" alt="${item.name}">
                    </div>
                    <div class="cart-item-details">
                        <h4 class="cart-item-title">${item.name}</h4>
                        <span class="cart-item-price">$${item.price.toFixed(2)}</span>
                        <button class="cart-item-remove" data-id="${item.id}">Remove</button>
                        <div class="cart-item-quantity">
                            <button class="quantity-btn minus" data-id="${item.id}">-</button>
                            <input type="number" class="quantity-input" value="${item.quantity}" min="1" data-id="${item.id}">
                            <button class="quantity-btn plus" data-id="${item.id}">+</button>
                        </div>
                    </div>
                `;
                
                cartItemsContainer.appendChild(cartItem);
            });
            
            totalPriceElement.textContent = `$${totalPrice.toFixed(2)}`;
        }
        
        // Add event listeners to cart items
        document.querySelectorAll('.cart-item-remove').forEach(button => {
            button.addEventListener('click', (e) => {
                const productId = parseInt(button.getAttribute('data-id'));
                removeFromCart(productId);
            });
        });
        
        document.querySelectorAll('.quantity-btn.minus').forEach(button => {
            button.addEventListener('click', (e) => {
                const productId = parseInt(button.getAttribute('data-id'));
                const item = cart.find(item => item.id === productId);
                
                if (item && item.quantity > 1) {
                    updateQuantity(productId, item.quantity - 1);
                }
            });
        });
        
        document.querySelectorAll('.quantity-btn.plus').forEach(button => {
            button.addEventListener('click', (e) => {
                const productId = parseInt(button.getAttribute('data-id'));
                const item = cart.find(item => item.id === productId);
                
                if (item) {
                    updateQuantity(productId, item.quantity + 1);
                }
            });
        });
        
        document.querySelectorAll('.quantity-input').forEach(input => {
            input.addEventListener('change', (e) => {
                const productId = parseInt(input.getAttribute('data-id'));
                updateQuantity(productId, input.value);
            });
        });
    }
    
    // Save cart to localStorage
    function saveCartToLocalStorage() {
        localStorage.setItem('cart', JSON.stringify(cart));
    }
    
    // Load cart from localStorage
    function loadCartFromLocalStorage() {
        const savedCart = localStorage.getItem('cart');
        if (savedCart) {
            cart = JSON.parse(savedCart);
            updateCart();
        }
    }
    
    // Toggle cart sidebar
    function toggleCartSidebar() {
        cartSidebar.classList.toggle('active');
        cartOverlay.classList.toggle('active');
        document.body.style.overflow = cartSidebar.classList.contains('active') ? 'hidden' : 'auto';
    }
    
    // Event listeners
    cartIcon.addEventListener('click', toggleCartSidebar);
    closeCart.addEventListener('click', toggleCartSidebar);
    cartOverlay.addEventListener('click', toggleCartSidebar);
    
    // Initialize
    displayProducts();
    loadCartFromLocalStorage();
    
    // Add event listeners to all "Add to Cart" buttons
    document.addEventListener('click', function(e) {
        if (e.target.classList.contains('add-to-cart')) {
            const productId = parseInt(e.target.getAttribute('data-id'));
            addToCart(productId);
            
            // Show cart sidebar when adding an item
            if (!cartSidebar.classList.contains('active')) {
                toggleCartSidebar();
            }
        }
    });
    
    // Smooth scrolling for navigation links
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function(e) {
            e.preventDefault();
            
            const targetId = this.getAttribute('href');
            if (targetId === '#') return;
            
            const targetElement = document.querySelector(targetId);
            if (targetElement) {
                targetElement.scrollIntoView({
                    behavior: 'smooth'
                });
                
                // Close mobile menu if open
                if (mainNav.classList.contains('active')) {
                    mainNav.classList.remove('active');
                }
            }
        });
    });
});