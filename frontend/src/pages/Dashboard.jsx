import { useNavigate } from 'react-router-dom';
import './Dashboard.css';

export default function Dashboard() {
    const navigate = useNavigate();

    function handleLogout() {
        localStorage.removeItem('token');
        navigate('/');
    }

    return (
        <div className="dashboard-stub">
            <h1>Dashboard</h1>
            <p>Em construção...</p>
            <button className="btn-nav" onClick={handleLogout}>Sair</button>
        </div>
    );
}