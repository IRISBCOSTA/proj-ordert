import { useState } from 'react';

export default function AmountForm({ title, buttonLabel, action, onSuccess }) {
    const [amount, setAmount] = useState('');
    const [error, setError] = useState('');
    const [message, setMessage] = useState('');
    const [loading, setLoading] = useState(false);

    async function handleSubmit(event) {
        event.preventDefault();
        setError('');
        setMessage('');
        setLoading(true);
        try {
            await action(Number(amount));
            setAmount('');
            setMessage('Operação realizada com sucesso!');
            onSuccess(); // pede ao Dashboard para recarregar saldo e extrato
        } catch (err) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    }

    return (
        <form className="card" onSubmit={handleSubmit}>
            <h2>{title}</h2>
            <label htmlFor={`amount-${title}`}>
                Valor (R$)
                <input
                    id={`amount-${title}`}
                    type="number"
                    step="0.01"
                    min="0.01"
                    required
                    value={amount}
                    onChange={(event) => setAmount(event.target.value)}
                />
            </label>

            {error && <p className="form-error">{error}</p>}
            {message && <p className="form-success">{message}</p>}

            <button type="submit" className="card-submit" disabled={loading}>
                {loading ? 'Enviando...' : buttonLabel}
            </button>
        </form>
    );
}