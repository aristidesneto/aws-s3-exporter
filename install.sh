#!/bin/zsh

set -e

# Cores para output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Iniciando instalação do AWS S3 Metrics Exporter${NC}\n"

# Verifica se o Go está instalado
if ! command -v go &> /dev/null; then
    echo -e "${YELLOW}Go não encontrado. Por favor, instale o Go antes de continuar.${NC}"
    exit 1
fi

# Verifica se o git está instalado
if ! command -v git &> /dev/null; then
    echo -e "${YELLOW}Git não encontrado. Por favor, instale o Git antes de continuar.${NC}"
    exit 1
fi

# Define diretórios
INSTALL_DIR="/opt/aws-s3-exporter"
CONFIG_DIR="/etc/aws-s3-exporter"
SYSTEMD_DIR="/etc/systemd/system"

# Cria diretórios necessários
echo -e "${GREEN}Criando diretórios...${NC}"
sudo mkdir -p ${INSTALL_DIR}
sudo mkdir -p ${CONFIG_DIR}

# Clona o repositório
echo -e "${GREEN}Clonando repositório...${NC}"
git clone https://github.com/seu-usuario/aws-s3-exporter.git /tmp/aws-s3-exporter

# Entra no diretório
cd /tmp/aws-s3-exporter

# Build do projeto
echo -e "${GREEN}Building projeto...${NC}"
make build

# Copia os arquivos necessários
echo -e "${GREEN}Instalando arquivos...${NC}"
sudo cp bin/aws-s3-exporter ${INSTALL_DIR}/
sudo cp configs/config.example.yaml ${CONFIG_DIR}/config.yaml

# Cria o service do systemd
echo -e "${GREEN}Configurando serviço systemd...${NC}"
cat << EOF | sudo tee ${SYSTEMD_DIR}/aws-s3-exporter.service
[Unit]
Description=AWS S3 Metrics Exporter
After=network.target

[Service]
Type=simple
User=root
ExecStart=${INSTALL_DIR}/aws-s3-exporter --config ${CONFIG_DIR}/config.yaml
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

# Recarrega o systemd
sudo systemctl daemon-reload

echo -e "${GREEN}Instalação concluída!${NC}\n"
echo -e "Para configurar o exporter:"
echo -e "1. Edite o arquivo de configuração: ${YELLOW}sudo nano ${CONFIG_DIR}/config.yaml${NC}"
echo -e "2. Inicie o serviço: ${YELLOW}sudo systemctl start aws-s3-exporter${NC}"
echo -e "3. Habilite o serviço para iniciar com o sistema: ${YELLOW}sudo systemctl enable aws-s3-exporter${NC}"
echo -e "\nPara verificar o status: ${YELLOW}sudo systemctl status aws-s3-exporter${NC}"
echo -e "Para ver os logs: ${YELLOW}sudo journalctl -u aws-s3-exporter -f${NC}"
