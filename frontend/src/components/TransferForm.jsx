import { useState } from 'react';
import { transfer } from '../services/api';

export default function TransferForm({ onSuccess }) {
    const [toAccountId, setToAccountId] = useState('');
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
            await transfer({
                to_account_id: Number(toAccountId),
                amount: Number(amount),
            });
            setToAccountId('');
            setAmount('');
            setMessage('Transferência realizada com sucesso!');
            onSuccess();
        } catch (err) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    }

    return (
        <form className="card" onSubmit={handleSubmit}>
            <h2>Transferência</h2>
            <label htmlFor="to-account">
                ID da conta de destino
                <input
                    id="to-account"
                    type="number"
                    min="1"
                    required
                    value={toAccountId}
                    onChange={(event) => setToAccountId(event.target.value)}
                />
            </label>
            <label htmlFor="transfer-amount">
                Valor (R$)
                <input
                    id="transfer-amount"
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
                {loading ? 'Enviando...' : 'Transferir'}
            </button>
        </form>
    );
}