const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

async function request(path, options = {}) {
    // Se houver token salvo, envia junto em toda requisição
    const token = localStorage.getItem('token');
    const headers = { 'Content-Type': 'application/json' };
    if (token) {
        headers.Authorization = `Bearer ${token}`;
    }

    let response;
    try {
        response = await fetch(`${API_URL}${path}`, { ...options, headers });
    } catch {
        throw new Error('Não foi possível conectar ao servidor.');
    }

    const data = await response.json().catch(() => null);
    if (!response.ok) {
        const error = new Error(data?.error || 'Algo deu errado. Tente novamente.');
        error.status = response.status; // usado para detectar 401 (token inválido)
        throw error;
    }
    return data;
}

export function registerUser({ name, email, password }) {
    return request('/register', {
        method: 'POST',
        body: JSON.stringify({ name, email, password }),
    });
}

export function loginUser({ email, password }) {
    return request('/login', {
        method: 'POST',
        body: JSON.stringify({ email, password }),
    });
}

export function getMe() {
    return request('/me');
}

export function getAccount() {
    return request('/account');
}

export function deposit(amount) {
    return request('/deposit', {
        method: 'POST',
        body: JSON.stringify({ amount }),
    });
}

export function withdraw(amount) {
    return request('/withdraw', {
        method: 'POST',
        body: JSON.stringify({ amount }),
    });
}

export function transfer({ to_account_id, amount }) {
    return request('/transfer', {
        method: 'POST',
        body: JSON.stringify({ to_account_id, amount }),
    });
}

export function getTransactions() {
    return request('/transactions');
}