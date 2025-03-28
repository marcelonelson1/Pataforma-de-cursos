import React, { useState } from 'react';
import { motion } from 'framer-motion';
import { 
  Mail, 
  Phone, 
  MapPin, 
  Send, 
  User, 
  MessageCircle,
  Loader2,
  CheckCircle,
  AlertCircle
} from 'lucide-react';
import axios from 'axios';
import './ContactoPage.css';

const ContactPage = () => {
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    phone: '',
    message: ''
  });
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitStatus, setSubmitStatus] = useState({
    success: false,
    error: null
  });

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsSubmitting(true);
    setSubmitStatus({ success: false, error: null });

    try {
      const response = await axios.post('http://localhost:5000/api/contact', formData);
      
      if (response.data.success) {
        setSubmitStatus({ success: true, error: null });
        setFormData({
          name: '',
          email: '',
          phone: '',
          message: ''
        });
      }
    } catch (error) {
      const errorMsg = error.response?.data?.error || 'Error al enviar el mensaje. Por favor intente nuevamente.';
      setSubmitStatus({ success: false, error: errorMsg });
      console.error('Error:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="contact-page">
      <motion.div 
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
        className="contact-container"
      >
        <div className="contact-info">
          <h2>Contáctanos</h2>
          
          <div className="contact-details">
            <div className="contact-detail">
              <Mail className="icon" size={20} />
              <div>
                <p>Correo Electrónico</p>
                <a href="mailto:contacto@ejemplo.com" className="email-link">
                  contacto@ejemplo.com
                </a>
              </div>
            </div>

            <div className="contact-detail">
              <Phone className="icon" size={20} />
              <div>
                <p>Teléfono</p>
                <p>+54 351 388 2695</p>
              </div>
            </div>

            <div className="contact-detail">
              <MapPin className="icon" size={20} />
              <div>
                <p>Ubicación</p>
                <p>Córdoba, Argentina</p>
              </div>
            </div>
          </div>
        </div>

        <form onSubmit={handleSubmit} className="contact-form">
          {submitStatus.success && (
            <div className="alert success">
              <CheckCircle className="icon" size={18} />
              <span>¡Mensaje enviado con éxito! Nos pondremos en contacto contigo pronto.</span>
            </div>
          )}
          
          {submitStatus.error && (
            <div className="alert error">
              <AlertCircle className="icon" size={18} />
              <span>{submitStatus.error}</span>
            </div>
          )}

          <div className="form-group">
            <label htmlFor="name">
              <User className="icon" size={18} /> Nombre Completo
            </label>
            <input
              type="text"
              id="name"
              name="name"
              required
              value={formData.name}
              onChange={handleChange}
              className="form-input"
              placeholder="Tu nombre completo"
              disabled={isSubmitting}
            />
          </div>

          <div className="form-group">
            <label htmlFor="email">
              <Mail className="icon" size={18} /> Correo Electrónico
            </label>
            <input
              type="email"
              id="email"
              name="email"
              required
              value={formData.email}
              onChange={handleChange}
              className="form-input"
              placeholder="tu.correo@ejemplo.com"
              disabled={isSubmitting}
            />
          </div>

          <div className="form-group">
            <label htmlFor="phone">
              <Phone className="icon" size={18} /> Teléfono (Opcional)
            </label>
            <input
              type="tel"
              id="phone"
              name="phone"
              value={formData.phone}
              onChange={handleChange}
              className="form-input"
              placeholder="123 456 7890"
              disabled={isSubmitting}
            />
          </div>

          <div className="form-group">
            <label htmlFor="message">
              <MessageCircle className="icon" size={18} /> Mensaje
            </label>
            <textarea
              id="message"
              name="message"
              required
              value={formData.message}
              onChange={handleChange}
              className="form-textarea"
              placeholder="Escribe tu mensaje aquí..."
              rows="5"
              disabled={isSubmitting}
            />
          </div>

          <button
            type="submit"
            className="submit-btn"
            disabled={isSubmitting}
          >
            {isSubmitting ? (
              <>
                <Loader2 className="icon spin" size={18} />
                Enviando...
              </>
            ) : (
              <>
                <Send className="icon" size={18} />
                Enviar Mensaje
              </>
            )}
          </button>
        </form>
      </motion.div>
    </div>
  );
};

export default ContactPage;