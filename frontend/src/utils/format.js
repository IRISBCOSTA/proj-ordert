export function formatMoney(value) {
    return Number(value).toLocaleString('pt-BR', {
        style: 'currency',
        currency: 'BRL',
    });
}

export function formatDate(isoString) {
    return new Date(isoString).toLocaleString('pt-BR');
}