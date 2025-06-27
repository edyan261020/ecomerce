/* scriptproducts.js */

// Array global vacío para productos

let products = [];
let cart = JSON.parse(localStorage.getItem('cart')) || [];
updateCartUI();
/**
 * Renderiza un array de productos en el grid
 * @param {Array} list - lista de productos a mostrar
 */
function displayProducts(list) {
  const productGrid = document.getElementById('product-grid');
  if (!productGrid) {
    console.error('#product-grid no existe en el DOM');
    return;
  }

  productGrid.innerHTML = '';
  list.forEach(product => {
    const card = document.createElement('div');
    card.className = 'product-card';
    card.innerHTML = `
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
    productGrid.appendChild(card);
  });
}

// Carrito de compras: actualizar UI, guardar, acciones
function updateCartUI() {
  const cartCount = document.querySelector('.cart-count');
  const totalPriceEl = document.querySelector('.total-price');
  const cartItemsContainer = document.getElementById('cart-items');

  if (!cartCount || !totalPriceEl || !cartItemsContainer) {
    console.error('Elementos del carrito no encontrados en el DOM');
    return;
  }

  const totalItems = cart.reduce((sum, item) => sum + item.quantity, 0);
  cartCount.textContent = totalItems;

  if (cart.length === 0) {
    cartItemsContainer.innerHTML = `
      <div class="empty-cart">
        <i class="fas fa-shopping-cart"></i>
        <p>Tu carrito está vacío</p>
      </div>
    `;
    totalPriceEl.textContent = '$0.00';
    return;
  }

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
        <h4>${item.name}</h4>
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

  totalPriceEl.textContent = `$${totalPrice.toFixed(2)}`;
}

function saveCart() {
  localStorage.setItem('cart', JSON.stringify(cart));
}

function addToCart(id) {
  const prod = products.find(p => p.id === id);
  if (!prod) return;
  const existing = cart.find(i => i.id === id);
  if (existing) existing.quantity++;
  else cart.push({ id: prod.id, name: prod.name, price: prod.price, image: prod.image, quantity: 1 });
  saveCart();
  updateCartUI();
}

function removeFromCart(id) {
  cart = cart.filter(i => i.id !== id);
  saveCart();
  updateCartUI();
}

function changeQuantity(id, qty) {
  const item = cart.find(i => i.id === id);
  if (!item) return;
  item.quantity = Math.max(1, qty);
  saveCart();
  updateCartUI();
}

function toggleCart() {
  const sidebar = document.getElementById('cart-sidebar');
  const overlay = document.getElementById('cart-overlay');
  if (!sidebar || !overlay) return;
  sidebar.classList.toggle('active');
  overlay.classList.toggle('active');
  document.body.style.overflow = sidebar.classList.contains('active') ? 'hidden' : 'auto';
}

function setupMobileMenu() {
  const btn = document.querySelector('.mobile-menu-btn');
  const nav = document.querySelector('.main-nav');
  if (btn && nav) btn.addEventListener('click', () => nav.classList.toggle('active'));
}

// Inicialización al cargar ventana
window.addEventListener('load', async () => {
  // 1) Cargar productos desde API
  try {
    const res = await fetch('/api/productos');
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    products = await res.json();
  } catch (e) {
    console.error('Error cargando productos:', e);
    const productGrid = document.getElementById('product-grid');
    if (productGrid) productGrid.innerHTML = '<p>No se pudieron cargar los productos.</p>';
    return;
  }

  // 2) Crear combo de categorías
  const sectionHeader = document.querySelector('.section-header');
  if (sectionHeader) {
    const select = document.createElement('select');
    select.id = 'category-filter';
    select.innerHTML = '<option value="">Todas las categorías</option>';
    const categories = [...new Set(products.map(p => p.category))];
    categories.forEach(cat => {
      const opt = document.createElement('option');
      opt.value = cat;
      opt.textContent = cat;
      select.appendChild(opt);
    });
    select.addEventListener('change', () => {
      const filtered = select.value ? products.filter(p => p.category === select.value) : products;
      displayProducts(filtered);
    });
    sectionHeader.appendChild(select);
  }

  // 3) Mostrar productos y configurar UI
  displayProducts(products);
  setupMobileMenu();

  // Carrito\  updateCartUI();
  document.getElementById('cart-icon')?.addEventListener('click', toggleCart);
  document.getElementById('close-cart')?.addEventListener('click', toggleCart);
  document.getElementById('cart-overlay')?.addEventListener('click', toggleCart);

  // Delegación de eventos
  document.addEventListener('click', e => {
    const addBtn = e.target.closest('.add-to-cart');
    if (addBtn) {
      addToCart(+addBtn.dataset.id);
      toggleCart();
      return;
    }
    if (e.target.matches('.cart-item-remove')) removeFromCart(+e.target.dataset.id);
    if (e.target.matches('.quantity-btn.plus')) changeQuantity(+e.target.dataset.id, cart.find(i => i.id === +e.target.dataset.id).quantity + 1);
    if (e.target.matches('.quantity-btn.minus')) changeQuantity(+e.target.dataset.id, cart.find(i => i.id === +e.target.dataset.id).quantity - 1);
  });

  document.addEventListener('change', e => {
    if (e.target.matches('.quantity-input')) changeQuantity(+e.target.dataset.id, +e.target.value);
  });
});



document.querySelector('.checkout-btn')?.addEventListener('click', async () => {


  try {
 const cedula = localStorage.getItem('clienteCedula');
  const nombre = localStorage.getItem('clienteNombres');

  // 2) Si no hay sesión de cliente, redirigir al login
  if (!cedula || !nombre) {
     alert('Debes iniciar sesión antes de realizar un pedido.');
    return window.location.href = '/login';
  }else{await sendOrderToAPI();}
    
    alert('Orden generada correctaemente. Opcional puedes continuar por WhatsApp para finalizar el pedido.');
  } catch (err) {
    return alert('No se pudo guardar el pedido en el servidor.');
  }



  const modal = document.getElementById('invoice-modal');
  const content = document.getElementById('invoice-content');

  if (!cart.length) {
    content.innerHTML = `<p>Tu carrito está vacío.</p>`;
    modal.style.display = 'flex';
    return;
  }

  const fecha = new Date().toLocaleString();
  let rows = '';
  let total = 0;
  cart.forEach(item => {
    const subtotal = item.price * item.quantity;
    total += subtotal;
    rows += `
      <tr>
        <td>${item.name}</td>
        <td style="text-align:right;">${item.quantity}</td>
        <td style="text-align:right;">$${subtotal.toFixed(2)}</td>
      </tr>`;
  });

  const mensajeWhatsApp = encodeURIComponent(
    `Hola, quiero hacer un pedido:\n\n` +
    cart.map(item => `- ${item.name} x${item.quantity} ($${item.price.toFixed(2)})`).join('\n') +
    `\n\nTotal: $${total.toFixed(2)}`
  );

  content.innerHTML = `
    <h2 style="text-align:center;">Factura de Pedido</h2>
    <p style="text-align:center;">Fecha: ${fecha}</p>
    <table>
      <thead>
        <tr>
          <th>Producto</th>
          <th style="text-align:right;">Cantidad</th>
          <th style="text-align:right;">Subtotal</th>
        </tr>
      </thead>
      <tbody>${rows}</tbody>
      <tfoot>
        <tr>
          <td colspan="2" style="text-align:right;">Total:</td>
          <td style="text-align:right;">$${total.toFixed(2)}</td>
        </tr>
      </tfoot>
    </table>
    <div style="text-align:center; margin-top: 20px;">
      <a href="https://wa.me/593987498260?text=${mensajeWhatsApp}" target="_blank"
        style="background:#25D366; color:white; padding:10px 20px; border-radius:5px; text-decoration:none; font-weight:bold;">
        Seguir por WhatsApp
      </a>
    </div>
  `;

  modal.style.display = 'flex';
});




















// 1) Agrega esto arriba del listener
async function sendOrderToAPI() {
   const cedula = localStorage.getItem('clienteCedula');
  const nombre = localStorage.getItem('clienteNombres');
 
  // 1.1) Crear la venta
  const ventaRes = await fetch('/api/ventas', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ 
      cedula_cliente: cedula,
      nombre_cliente: nombre,
      fecha: new Date().toISOString().slice(0, 19).replace('T', ' '),
      total: cart.reduce((sum, i) => sum + i.price * i.quantity, 0)
    })
  });
  if (!ventaRes.ok) throw new Error('Error creando venta');

  // 1.2) Obtener el ID recién creado (asumiendo GET /api/ventas/:id con el max)
  const { id } = await fetch('/api/ventas/max') // o la ruta que devuelva el último ID
    .then(r => r.json());
    console.log('Último ID de venta recuperado:', id);

  // 1.3) Enviar cada detalle
  await Promise.all(cart.map(item =>
    fetch('/api/detalles_ventas', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        id_venta: id,
        id_producto: item.id,
        detalle: item.name,        // o el texto que quieras
        cantidad: item.quantity,
        precio_uni: item.price,
        sub_total: item.price * item.quantity
      })
    })
  ));
}








 
// Botón cerrar modal y limpiar carrito
document.querySelector('#invoice-modal .close-btn')?.addEventListener('click', () => {
  // Cerrar modal
  document.getElementById('invoice-modal').style.display = 'none';
  // Vaciar carrito y localStorage
  cart = [];
  saveCart();
  updateCartUI();
});


 