# AWS S3 Metrics Exporter

Este é um exportador Prometheus que coleta métricas de buckets S3 da AWS, fornecendo informações sobre contagem de arquivos e tamanho total por prefixo de data.

> **Nota sobre o Caso de Uso**: Este exporter foi desenvolvido para atender uma necessidade específica de monitoramento de buckets S3 onde os arquivos são organizados por data (exemplo: 2025/05/26). A lógica de agrupamento por prefixo de data é um requisito específico da empresa onde foi originalmente desenvolvido, portanto, pode necessitar adaptações para outros casos de uso.

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

## Pré-requisitos

- Go 1.x ou superior
- Credenciais AWS configuradas
- make (opcional, para usar os comandos do Makefile)

## Configuração

1. Clone o repositório
2. Copie o arquivo de exemplo de configuração:
   ```bash
   cp config.yaml.example.yaml config.yaml
   ```
3. Configure o arquivo `config.yaml` com suas configurações:
   ```yaml
   aws_profile: "seu-profile" # Profile AWS a ser usado
   aws_region: "us-east-1" # Região AWS
   interval: 5 # Intervalo de coleta em minutos
   buckets: # Lista de buckets para monitorar
     - "bucket-1"
     - "bucket-2"
   ```

## Comandos Disponíveis

O projeto utiliza um Makefile para facilitar as operações comuns:

- `make help`: Exibe a lista de comandos disponíveis
- `make build`: Compila o projeto
- `make run`: Executa a aplicação
- `make test`: Executa os testes
- `make clean`: Remove os binários gerados
- `make lint`: Executa o linter (golangci-lint)
- `make fmt`: Formata o código fonte
- `make dev`: Executa formatação, lint, build e roda a aplicação

## Como Executar

1. Build do projeto:

   ```bash
   make build
   ```

2. Executar o exportador:

   ```bash
   # Usando o arquivo de configuração padrão em ./configs/config.yaml
   make run

   # Ou especificando um arquivo de configuração diferente
   ./bin/aws-s3-exporter --config /caminho/para/config.yaml
   ```

O exportador estará disponível em `http://localhost:2112/metrics`

## Instalação Rápida

Para instalar o exporter em um novo servidor, você pode usar o script de instalação:

```bash
curl -sSL https://raw.githubusercontent.com/seu-usuario/aws-s3-exporter/main/install.sh | sudo zsh
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

Para desenvolvimento, você pode usar o comando:

```bash
make dev
```

Este comando irá:

1. Formatar o código
2. Executar o linter
3. Compilar o projeto
4. Executar a aplicação

## Estrutura do Projeto

```
aws-s3-exporter/
├── config/          # Configurações
├── internal/        # Funções internas
├── prometheus/      # Definições das métricas
├── config.yaml      # Arquivo de configuração
├── main.go         # Ponto de entrada
└── Makefile        # Comandos úteis
```

## Contribuindo

1. Faça o fork do projeto
2. Crie sua feature branch (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -am 'Adicionando nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Crie um Pull Request

## Licença

Este projeto está sob a licença MIT. Veja o arquivo `LICENSE` para mais detalhes.
