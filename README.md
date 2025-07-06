# AWS S3 Metrics Exporter

Este é um exportador Prometheus que coleta métricas de buckets S3 da AWS, fornecendo informações sobre contagem de arquivos e tamanho total por prefixo de data.

> **Nota sobre o Caso de Uso**: Este exporter foi desenvolvido para atender uma necessidade específica de monitoramento de buckets S3 onde os arquivos são organizados por data (exemplo: 2025/05/26). A lógica de agrupamento por prefixo de data é um requisito específico, portanto, pode necessitar adaptações para outros casos de uso.

## Funcionalidades

- Coleta métricas de buckets S3 configurados
- Agrupa métricas por prefixo de data
- Exposição de métricas no formato Prometheus
- Configuração flexível via arquivo YAML
- Suporte a múltiplos buckets
- Atualização periódica das métricas

## Métricas Exportadas

- `s3_file_count`: Número de arquivos por bucket e prefixo de data
- `s3_total_size`: Tamanho total dos arquivos por bucket e prefixo de data
- `s3_backup_last_upload_timestamp`: Timestamp do último upload por prefixo (ano/mês/dia)

## Pré-requisitos

- Go 1.x ou superior
- Credenciais AWS configuradas
- make (opcional, para usar os comandos do Makefile)

## Instalação

Crie um arquivo de configuração chamado `config.yaml` com as seguintes configurações.

> Atualize conforme necessidade

```yaml
# AWS S3 Exporter Configuration
# =============================
# Este arquivo define as configurações do exporter de métricas S3.
# Siga o padrão YAML e consulte a documentação para detalhes de cada campo.

aws:
  # Profile utilizado para acessar a AWS (opcional, pode ser sobrescrito por variáveis de ambiente ou flags)
  # Em ambientes Docker ou cloud, recomenda-se utilizar as variáveis de ambiente padrão da AWS
  # (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_REGION), pois elas têm prioridade sobre o profile configurado aqui.
  # Deixe em branco para ser usado a default
  profile: ""
  # Região AWS (opcional, pode ser sobrescrito por variáveis de ambiente ou flags)
  region: ""

s3:
  # Lista de buckets que serão monitorados
  buckets: []

# Intervalo em minutos para verificar os buckets
interval: 10
```

Escolha como deseja rodar sua aplicação:

1. **Docker**

```bash
docker run -p 2112:2112 \
  -v ./config.yaml:/etc/aws-s3-exporter/config.yaml \
  -e AWS_DEFAULT_REGION=us-east-1 \
  -e AWS_ACCESS_KEY_ID=<access-key> \
  -e AWS_SECRET_ACCESS_KEY=<secret-key> \
  aristidesbneto/aws-s3-exporter:v1 \
    --config /etc/aws-s3-exporter/config.yaml
```

2. **Docker Compose**

```yaml
services:
  aws-s3-exporter:
    image: aristidesbneto/aws-s3-exporter:v1
    environment:
      AWS_DEFAULT_REGION: ${AWS_DEFAULT_REGION}
      AWS_ACCESS_KEY_ID: ${AWS_ACCESS_KEY_ID}
      AWS_SECRET_ACCESS_KEY: ${AWS_SECRET_ACCESS_KEY}
    # command: "--config /etc/aws-s3-exporter/config.yaml"
    volumes:
      - ./configs/config.yaml:/etc/aws-s3-exporter/config.yaml:ro
    ports:
      - "2112:2112"
```

3. **Linux**

Para instalar o exporter em um servidor Linux com systemd, você pode usar o script de instalação:

```bash
curl -sSL https://raw.githubusercontent.com/aristidesneto/aws-s3-exporter/main/install.sh | sudo bash
```

O script irá:

1. Verificar os pré-requisitos (Go e Git)
2. Criar os diretórios necessários
3. Clonar o repositório
4. Buildar o projeto
5. Instalar o binário e arquivos de configuração
6. Configurar o serviço systemd

Após a instalação:

1. Configure o arquivo em `/etc/aws-s3-exporter/config.yaml`
2. Inicie o serviço: `sudo systemctl start aws-s3-exporter`
3. Habilite o início automático: `sudo systemctl enable aws-s3-exporter`

## Desenvolvimento

Se deseja contribuir, você pode executar o projeto localmente:

1. Clone o repositório

2. Copie o arquivo de exemplo de configuração:

   ```bash
   cp configs/config.example.yaml configs/config.yaml
   ```

3. Atualize o arquivo `./configs/config.yaml` com as suas configurações.
4. Ao executar, sua aplicação estará disponíveis em `http://localhost:2112/metrics`.

### Usando Makefile

Executar o exportador:

```bash
# O arquivo de configuração padrão será usado: ./configs/config.yaml
make run
```

### Usando Go CLI

Para executar via linha de comando

```bash
go run cmd/aws-s3-exporter/main.go --config=./configs/config.yaml
```

Você pode sobrescrever qualquer valor do arquivo de configuração usando variáveis de ambiente. Exemplo:

```
AWS_DEFAULT_REGION=us-east-1 \
AWS_ACCESS_KEY_ID=<access-key> \
AWS_SECRET_ACCESS_KEY=<secret-key> \
S3_BUCKETS=bucket1,bucket2 \
INTERVAL=10 \
go run cmd/aws-s3-exporter/main.go
```

### Usando Docker Compose (recomendado)

1. Crie o arquivo `.env` e adicione suas credenciais da AWS.

```bash
cp .env.example .env
```

2. Edite o arquivo `./configs/config.yaml` conforme necessário.

3. O arquivo `compose.yaml` já está configurado para ler seu `.env` e substiuir as variáveis de ambiente corretamente.

4. Execute `docker compose build && docker compose up -d` para buildar a aplicação e executar o container.

## Estrutura do Projeto

```
aws-s3-exporter/
├── config/          # Configurações
├── internal/        # Funções internas
├── prometheus/      # Definições das métricas
├── config.yaml      # Arquivo de configuração
├── main.go          # Ponto de entrada
└── Makefile         # Comandos úteis
```

## Comandos Disponíveis Makefile

O projeto utiliza um Makefile para facilitar as operações comuns:

- `make help`: Exibe a lista de comandos disponíveis
- `make build`: Compila o projeto
- `make run`: Executa a aplicação
- `make clean`: Remove os binários gerados
- `make lint`: Executa o linter (golangci-lint)
- `make fmt`: Formata o código fonte
- `make dev`: Executa formatação, lint, build e roda a aplicação

## Permissões AWS IAM Necessárias

O exporter precisa das seguintes permissões IAM mínimas para funcionar:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:PutObject",
        "s3:ListBucket",
        "s3:GetBucketLocation",
        "s3:GetObject"
      ],
      "Resource": [
        "arn:aws:s3:::nome-do-seu-bucket",
        "arn:aws:s3:::nome-do-seu-bucket/*"
      ]
    }
  ]
}
```

Substitua `nome-do-seu-bucket` pelos nomes dos buckets que você deseja monitorar. Para múltiplos buckets, adicione suas ARNs na lista de `Resource`.

Exemplo para múltiplos buckets:

```json
"Resource": [
    "arn:aws:s3:::bucket1",
    "arn:aws:s3:::bucket1/*",
    "arn:aws:s3:::bucket2",
    "arn:aws:s3:::bucket2/*"
]
```

### Detalhamento das Permissões

As permissões necessárias são:

- `s3:PutObject`: Permite salvar os objetos no bucket
- `s3:ListBucket`: Permite listar os objetos dentro do bucket, necessário para coletar informações sobre arquivos e diretórios
- `s3:GetBucketLocation`: Permite verificar se temos acesso ao bucket e sua localização
- `s3:GetObject`: Permite obter metadados dos objetos, como tamanho e última modificação

Estas são as permissões mínimas necessárias para o funcionamento do exporter.
