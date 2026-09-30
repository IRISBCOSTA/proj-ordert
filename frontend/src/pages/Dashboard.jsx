import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Landmark, LogOut } from 'lucide-react';
import {
    getMe,
    getAccount,
    getTransactions,
    deposit,
    withdraw,
} from '../services/api';
import { formatMoney } from '../utils/format';
import AmountForm from '../components/AmountForm';
import TransferForm from '../components/TransferForm';
import Statement from '../components/Statement';
import './Dashboard.css';

export default function Dashboard() {
    const navigate = useNavigate();
    const [user, setUser] = useState(null);
    const [account, setAccount] = useState(null);
    const [transactions, setTransactions] = useState([]);
    const [error, setError] = useState('');

    function handleLogout() {
        localStorage.removeItem('token');
        navigate('/');
    }

    async function loadData() {
        try {
            const [me, acc, txs] = await Promise.all([
                getMe(),
                getAccount(),
                getTransactions(),
            ]);
            setUser(me);
            setAccount(acc);
            setTransactions(txs);
            setError('');
        } catch (err) {
            if (err.status === 401) {
                //token expirado ou invalido: volta para o login
                handleLogout();
                return;
            }
            setError(err.message);
        }
    }

    useEffect(() => {
        if (!localStorage.getItem('token')) {
            navigate('/Login');
            return;
        }
        loadData();
    }, []);

    return (
        <div className="dashboard">
            <header className="dash-header">
                <div className="dash-logo">
                    <Landmark size={22} strokeWidth={2.4} />
                    Ordet Bank
                </div>
                <div className="dash-user">
                    <span>Olá, {user?.name}</span>
                    <button className="dash-logout" onClick={handleLogout}>
                        <LogOut size={16} />
                        Sair
                    </button>
                </div>
            </header>

            <main className="dash-main">
                {error && <p className="form-error">{error}</p>}

                <section className="balance-card">
                    <p className="balance-label">Saldo disponível</p>
                    <p className="balance-value">
                        {account ? formatMoney(account.balance) : '...'}
                    </p>
                    <p className="balance-id">
                        ID da conta: <strong>{account?.id}</strong>
                    </p>
                </section>

                <div className="dash-grid">
                    <AmountForm
                        title="Depósito"
                        buttonLabel="Depositar"
                        action={deposit}
                        onSuccess={loadData}
                    />
                    <AmountForm
                        title="Saque"
                        buttonLabel="Sacar"
                        action={withdraw}
                        onSuccess={loadData}
                    />
                    <TransferForm onSuccess={loadData} />
                </div>

                <Statement entries={transactions} />
            </main>
        </div>
    );
}