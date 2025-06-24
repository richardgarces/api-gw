// --- JWT generation in JS ---
function base64url(input) {
  return btoa(input)
    .replace(/=/g, '')
    .replace(/\+/g, '-')
    .replace(/\//g, '_');
}

async function hmacSHA256(key, msg) {
  const enc = new TextEncoder();
  const cryptoKey = await window.crypto.subtle.importKey(
    "raw", enc.encode(key), { name: "HMAC", hash: "SHA-256" }, false, ["sign"]
  );
  const sig = await window.crypto.subtle.sign("HMAC", cryptoKey, enc.encode(msg));
  let str = '';
  new Uint8Array(sig).forEach(b => { str += String.fromCharCode(b); });
  return base64url(str);
}

async function generateJWT(secret) {
  const header = base64url(JSON.stringify({ alg: "HS256", typ: "JWT" }));
  const payload = base64url(JSON.stringify({
    sub: "admin",
    exp: Math.floor(Date.now() / 1000) + 365 * 24 * 60 * 60
  }));
  const data = header + "." + payload;
  const signature = await hmacSHA256(secret, data);
  return data + "." + signature;
}

function saveToken() {
  localStorage.setItem('jwt_token', document.getElementById('token').value);
  let box = document.getElementById('error-box');
  if (box) box.textContent = '';
  alert('Token guardado');
}

function getTokenSync() {
  return localStorage.getItem('jwt_token') || '';
}

// Obtiene el token del backend si no existe en localStorage
async function ensureToken() {
  let token = getTokenSync();
  if (!token) {
    const res = await fetch('/admin/token');
    const data = await res.json();
    token = data.token;
    localStorage.setItem('jwt_token', token);
  }
  document.getElementById('token').value = token;
}

(async function() {
  await ensureToken();
})();

function showError(msg) {
  let box = document.getElementById('error-box');
  if (!box) {
    box = document.createElement('div');
    box.id = 'error-box';
    box.style.color = 'red';
    box.style.margin = '1em 0';
    document.body.insertBefore(box, document.body.firstChild.nextSibling);
  }
  box.textContent = msg;
}

async function loadServices() {
  const res = await fetch('/admin/services', {
    headers: {
      'Authorization': 'Bearer ' + getTokenSync()
    }
  });
  if (res.status === 401 || res.status === 403) {
    showError('Token inválido o expirado. Por favor, actualiza el token.');
    return;
  }
  const data = await res.json();
  const tbody = document.querySelector('#services tbody');
  tbody.innerHTML = '';
  data.forEach(s => {
    tbody.innerHTML += `<tr>
      <td>${s.id}</td>
      <td>${s.name}</td>
      <td>${(s.targets||[]).join('<br>')}</td>
      <td><pre>${JSON.stringify(s.plugins, null, 2)}</pre></td>
    </tr>`;
  });
}

async function loadRoutes() {
  const res = await fetch('/admin/routes', {
    headers: {
      'Authorization': 'Bearer ' + getTokenSync()
    }
  });
  if (res.status === 401 || res.status === 403) {
    showError('Token inválido o expirado. Por favor, actualiza el token.');
    return;
  }
  const data = await res.json();
  const tbody = document.querySelector('#routes tbody');
  tbody.innerHTML = '';
  data.forEach(r => {
    tbody.innerHTML += `<tr>
      <td>${r.id}</td>
      <td>${r.path}</td>
      <td>${r.service_id}</td>
      <td><pre>${JSON.stringify(r.plugins, null, 2)}</pre></td>
    </tr>`;
  });
}
(async function() {
  await ensureToken();
  loadServices();
  loadRoutes();
})();