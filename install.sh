#!/bin/bash

set -e

# Cores para output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Função para backup
backup_config() {
    if [ -f "$1" ]; then
        backup_file="${1}.backup-$(date +%Y%m%d_%H%M%S)"
        echo -e "${YELLOW}Fazendo backup de $1 para $backup_file${NC}"
        sudo cp "$1" "$backup_file"
    fi
}

# Função para parar o serviço se estiver rodando
stop_service() {
    if systemctl is-active --quiet aws-s3-exporter; then
        echo -e "${YELLOW}Parando serviço aws-s3-exporter...${NC}"
        sudo systemctl stop aws-s3-exporter
    fi
}

# Verifica se é uma atualização
UPDATE_MODE=0
if [ -d "$INSTALL_DIR" ]; then
    echo -e "${YELLOW}Instalação existente detectada. Iniciando modo de atualização...${NC}"
    UPDATE_MODE=1
fi

echo -e "${GREEN}Iniciando $([ $UPDATE_MODE -eq 1 ] && echo "atualização" || echo "instalação") do AWS S3 Metrics Exporter${NC}\n"

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

# Limpa instalação anterior temporária se existir
rm -rf /tmp/aws-s3-exporter

# Cria diretórios necessários
echo -e "${GREEN}Criando/verificando diretórios...${NC}"
sudo mkdir -p ${INSTALL_DIR}
sudo mkdir -p ${CONFIG_DIR}

# Clona o repositório
echo -e "${GREEN}Baixando última versão...${NC}"
git clone --depth 1 https://github.com/aristidesneto/aws-s3-exporter.git /tmp/aws-s3-exporter

# Entra no diretório
cd /tmp/aws-s3-exporter

# Build do projeto
echo -e "${GREEN}Building projeto...${NC}"
make build

# Para o serviço se estiver rodando
stop_service

# Backup e copia dos arquivos
echo -e "${GREEN}Instalando arquivos...${NC}"
if [ $UPDATE_MODE -eq 1 ]; then
    # Backup da configuração existente
    backup_config "${CONFIG_DIR}/config.yaml"
    # Atualiza apenas o binário
    sudo cp bin/aws-s3-exporter ${INSTALL_DIR}/
else
    # Instalação inicial
    sudo cp bin/aws-s3-exporter ${INSTALL_DIR}/
    # Copia arquivo de configuração exemplo apenas se não existir
    if [ ! -f "${CONFIG_DIR}/config.yaml" ]; then
        sudo cp configs/config.example.yaml ${CONFIG_DIR}/config.yaml
    fi
fi

# Cria ou atualiza o service do systemd
echo -e "${GREEN}Configurando serviço systemd...${NC}"
# Backup do arquivo de serviço existente
backup_config "${SYSTEMD_DIR}/aws-s3-exporter.service"
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

# Limpa arquivos temporários
rm -rf /tmp/aws-s3-exporter

if [ $UPDATE_MODE -eq 1 ]; then
    echo -e "${GREEN}Atualização concluída!${NC}\n"
    echo -e "Para reiniciar o exporter:"
    echo -e "1. Verifique o arquivo de configuração: ${YELLOW}sudo nano ${CONFIG_DIR}/config.yaml${NC}"
    echo -e "2. Reinicie o serviço: ${YELLOW}sudo systemctl restart aws-s3-exporter${NC}"
    echo -e "\nBackups foram criados com timestamp dos arquivos modificados"
else
    echo -e "${GREEN}Instalação concluída!${NC}\n"
    echo -e "Para configurar o exporter:"
    echo -e "1. Edite o arquivo de configuração: ${YELLOW}sudo nano ${CONFIG_DIR}/config.yaml${NC}"
    echo -e "2. Inicie o serviço: ${YELLOW}sudo systemctl start aws-s3-exporter${NC}"
    echo -e "3. Habilite o serviço para iniciar com o sistema: ${YELLOW}sudo systemctl enable aws-s3-exporter${NC}"
fi

echo -e "\nComandos úteis:"
echo -e "Status do serviço: ${YELLOW}sudo systemctl status aws-s3-exporter${NC}"
echo -e "Logs do serviço: ${YELLOW}sudo journalctl -u aws-s3-exporter -f${NC}"
echo -e "Métricas exportadas: ${YELLOW}curl -i http://localhost:2112/metrics${NC}"

if [ $UPDATE_MODE -eq 1 ]; then
    echo -e "\nArquivos de backup:"
    find "${CONFIG_DIR}" "${SYSTEMD_DIR}" -name "*.backup-*" -type f 2>/dev/null
fi
