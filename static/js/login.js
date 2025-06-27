document.addEventListener('DOMContentLoaded', () => {
  const form       = document.getElementById('loginForm');
  const btnCancelar= document.getElementById('btnCancelar');
  const btnClients = document.getElementById('btnClients');
  const errorMsg   = document.getElementById('errorMsg');

  form.addEventListener('submit', async evt => {
    evt.preventDefault();
    errorMsg.style.display = 'none';

    const user = document.getElementById('username').value.trim();
    const pass = document.getElementById('password').value;

    try {
      const res = await fetch('/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ user, pass })
      });

      if (!res.ok) {
        // 401 o cualquier otro error
        errorMsg.style.display = 'block';
        return;
      }

      // Leer la respuesta con role, cedula y nombres
      const { role, cedula, nombres } = await res.json();

      // Guardar en localStorage
      localStorage.setItem('isLoggedIn', 'true');
      localStorage.setItem('role', role);

      if (role === 'client') {
        localStorage.setItem('clienteCedula', cedula);
        localStorage.setItem('clienteNombres', nombres);
      }

      // Redirigir según el rol
      if (role === 'admin') {
        window.location.href = '/admin';
      } else {
        window.location.href = '/products';
      }

    } catch (err) {
      console.error('Error al hacer login:', err);
      errorMsg.style.display = 'block';
    }
  });

  btnCancelar.addEventListener('click', () => {
    window.location.href = '/';
  });
  btnClients.addEventListener('click', () => {
    window.location.href = '/clients';
  });
});
