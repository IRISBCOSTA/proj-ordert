import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { LogIn } from 'lucide-react';
import { loginUser } from '../services/api';
import './Auth.css';

export default function Login() {
    const [email, setEmail] = useState('');
    const [password, setPassword ] = useState('');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();

    async function handleSubmit(event) {
        event.preventDefault();
        setError('');
        setLoading(true);
        try {
            const data = await loginUser({ email, password });
            localStorage.setItem('token', data.token);
            navigate('/Dashboard');
        } catch (err) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    }

    return (
        <div className="auth-page">
            <div className="auth-card">
                <h1>Entrar</h1>
                <form onSubmit={handleSubmit}>
                    <label htmlFor="email">
                        Email
                        <input
                            id="email"
                            type="email"
                            required
                            value={email}
                            onChange={(event) => setEmail(event.target.value)}
                        />
                    </label>
                    <label htmlFor="password">
                        Senha
                        <input
                        id="password"
                        type="password"
                        required
                        minLength={6}
                        value={password}
                        onChange={(event) => setPassword(event.target.value)}
                        />
                    </label>

                    {error && <p className="auth-error">{error}</p>}

                    <button type="submit" className="auth-submit" disabled={loading}>
                        <LogIn size={18} />
                        {loading ? 'Entrando...' : 'Entrar'}
                    </button>
                </form>

                <p className="auth-switch">
                    Não tem conta? <Link to="/Cadastro">Criar conta</Link>
                </p>
                <Link to="/" className="auth-back">Voltar</Link>
            </div>
        </div>
    );
}