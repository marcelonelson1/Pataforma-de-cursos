  GNU nano 8.4                                                                                                                                                                                                                                                                        deploy_vps.sh                                                                                                                                                                                                                                                                                 
  GNU nano 8.2                                                                                                                                                                                                               deploy_vps
#!/bin/bash

# Configuración
echo "=== Automatización Apache en VPS ==="
read -p "Tu dominio (ej: midominio.com): " dominio
read -p "Ruta completa del archivo .zip en la VPS (ej: /tmp/web.zip): " zip_path
read -p "Email para SSL (ej: tu@email.com): " email_ssl

# Verificar que el archivo .zip existe
if [ ! -f "$zip_path" ]; then
    echo "❌ Error: El archivo .zip no existe en la ruta especificada."
    exit 1
fi

# Extraer el nombre del archivo .zip
zip_name=$(basename "$zip_path")

# Crear directorio para el sitio web
sudo mkdir -p "/var/www/$dominio"

# Descomprimir el archivo .zip
echo "⬆️ Descomprimiendo archivo..."
sudo unzip "$zip_path" -d "/var/www/$dominio"

# Configurar permisos
sudo chown -R www-data:www-data "/var/www/$dominio"

# Instalar Apache y dependencias
echo "🔧 Instalando Apache y dependencias..."
sudo apt update
sudo apt install -y apache2 unzip

# Configurar Virtual Host
echo "📄 Configurando Virtual Host..."
vhost_content="<VirtualHost *:80>
    ServerName $dominio
    ServerAlias www.$dominio
    DocumentRoot /var/www/$dominio
    <Directory /var/www/$dominio>
        Options -Indexes +FollowSymLinks
        AllowOverride All
        Require all granted
    </Directory>
</VirtualHost>"

sudo echo "$vhost_content" | sudo tee "/etc/apache2/sites-available/$dominio.conf" > /dev/null

# Habilitar el sitio y reiniciar Apache
echo "🚀 Habilitando sitio y reiniciando Apache..."
sudo a2dissite 000-default.conf
sudo a2ensite "$dominio.conf"
sudo systemctl restart apache2

# Instalar Certbot y obtener SSL
echo "🔒 Configurando SSL..."
sudo apt install -y certbot python3-certbot-apache
sudo certbot --apache --non-interactive --agree-tos --email "$email_ssl" -d "$dominio" -d "www.$dominio"

# Abrir puertos en el firewall (si UFW está instalado)
if command -v ufw &> /dev/null; then
    echo "🔥 Configurando firewall..."
    sudo ufw allow 80/tcp
    sudo ufw allow 443/tcp
    sudo ufw reload
fi

echo "🎉 ¡Apache configurado correctamente!"
echo "Accede a: https://$dominio"
