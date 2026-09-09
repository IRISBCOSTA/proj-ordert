import { Link } from 'react-router-dom';
import { Landmark, LogIn, UserPlus} from 'lucide-react';
import './Landing.css';

export default function Landing() {
    return (
        <div className="landing">
            <header className="landing-header">
                <div className="landing-logo">
                    <Landmark size={22} strokeWidth={2.4} />
                    Ordet Bank
                    </div>
                    <Link to="/Login" className="btn-nav">Entrar</Link>
            </header>
            
            <main className="landing-hero">
                <h1>ORDET BANK</h1>
                <span className="landing-rule" aria-hidden="true" />
                <p className="landing-subtitle">Sua conta digital, simples e segura.</p>

                <div className="landing-cta">
                    <Link to="/Login" className="btn-primary">
                        <LogIn size={20} />
                        Entrar
                    </Link>
                    <Link to="/Cadastro" className="btn-primary">
                        <UserPlus size={20} />
                        Criar Conta
                    </Link>
                </div>
            </main>
        </div>
    );
}