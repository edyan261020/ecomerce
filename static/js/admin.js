document.addEventListener('DOMContentLoaded', () => {
  // Verificar sesión
  if (localStorage.getItem('isLoggedIn') !== 'true') {
    alert('Acceso no autorizado. Inicie sesión.');
    window.location.href = '/';
    return;
  }

  // Cerrar sesión
  document.getElementById('logoutBtn').addEventListener('click', () => {
    localStorage.removeItem('isLoggedIn');
    window.location.href = '/';
  });

  // Simular contenido CRUDs
  document.getElementById('usuariosCrud').innerHTML = `
    <p>Lista de usuarios aquí</p>
   <a href="/users">
  <button>Administrar Usuarios</button>
</a>

  `;

  document.getElementById('productosCrud').innerHTML = `
    <p>Lista de productos aquí</p>
    <a href="/productsx">
  <button>Administrar Productos</button>
</a>
  `;

  document.getElementById('comprasCrud').innerHTML = `
    <p>Lista de compras aquí</p>
    <a href="/sales">
  <button>Dashboart Ventas</button>
</a>
  `;
});
