const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

async function request(path, options = {}) {
    let response;
    try {
        response = await fetch(`${API_URL}${path}`, {
            headers: { 'Content-Type': 'application/json' },
            ...options,
        });
    } catch {
        throw new Error('Não foi possível conectar ao servidor.');
    }
    const data = await response.json().catch(() => null);
    if (!response.ok) {
        throw new Error(data?.error || 'Algo deu errado. Tente novamente.'); 
    }
    return data;
}

export function registerUser({ name, email, password }) {
    return request('/register', {
        method: 'POST',
        body: JSON.stringify({ email, password }),
    });
}

export function loginUser({ email, password }) {
    return request('/login', {
        method: 'POST',
        body: JSON.stringify({ email, password }),
    });
}