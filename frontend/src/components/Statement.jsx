import { formatMoney, formatDate } from '../utils/format';

const LABELS = {
    deposito: 'Depósito',
    saque: 'Saque',
    transferencia_enviada: 'Transferência enviada',
    transferencia_recebida: 'Transferência recebida',
};

const DEBITS = ['saque', 'transferencia_enviada'];

export default function Statement({ entries }) {
    return (
        <section className="card statement">
            <h2>Extrato</h2>

            {entries.length === 0 ? (
                <p className="statement-empty">Nenhuma movimentação ainda.</p>
            ) : (
                <ul className="statement-list">
                    {entries.map((entry, index) => {
                        const isDebit = DEBITS.includes(entry.type);
                        return (
                            <li key={index} className="statement-item">
                                <div>
                                    <p className="statement-type">
                                        {LABELS[entry.type] || entry.type}
                                    </p>
                                    <p className="statement-date">
                                        {formatDate(entry.created_at)}
                                    </p>
                                </div>
                                <p className={isDebit ? 'amount-debit' : 'amount-credit'}>
                                    {isDebit ? '- ' : '+ '}
                                    {formatMoney(entry.amount)}
                                </p>
                            </li>
                        );
                    })}
                </ul>
            )}
        </section>
    );
}