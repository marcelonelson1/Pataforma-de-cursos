import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Header from './components/Header';
import Footer from './components/Footer';
import Home from './pages/Home';
import CursosPage from './pages/CursosPage';
import SobreMiPage from './pages/SobreMiPage';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import CursoDetalle from './components/CursoDetalle';
import { AuthProvider } from './context/AuthContext';
import 'animate.css';
import './App.css';

function App() {
  return (
    <AuthProvider>
      <Router>
        <div className="app">
          <Header />
          <main>
            <Routes>
              <Route path="/" element={<Home />} />
              <Route path="/cursos" element={<CursosPage />} />
              <Route path="/sobre-mi" element={<SobreMiPage />} />
              <Route path="/login" element={<LoginPage />} />
              <Route path="/registro" element={<RegisterPage />} />
              <Route path="/curso/:id" element={<CursoDetalle />} />
            </Routes>
          </main>
          <Footer />
        </div>
      </Router>
    </AuthProvider>
  );
}

export default App;