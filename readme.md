# Gopher Chaos 🌪️

### Por que o Gopher Chaos existe?
Sistemas falham. É uma certeza, não uma probabilidade. Em arquiteturas de microserviços, o sucesso de uma requisição depende de dezenas de fatores externos: latência de rede, saturação de CPU e falhas parciais.

Muitas vezes, nossos testes de unidade e integração cobrem o "caminho feliz", mas negligenciam como o sistema se comporta sob estresse. 

O Gopher Chaos foi criado para:

1. **Validar Timeouts**: Garantir que o cliente não fique esperando infinitamente por uma resposta.
2. **Testar Circuit Breakers**: Verificar se sua aplicação para de chamar um serviço que está falhando sistematicamente.
3. **Expor Race Conditions**: A latência variável frequentemente revela bugs de concorrência que não aparecem em ambientes de baixa latência.

### O Fluxo do Caos

Diferente de ferramentas que injetam falhas na camada de rede, o Gopher Chaos atua na camada de aplicação via Interceptors gRPC. Isso permite um controle granular e visibilidade total dentro do código Go.

## Como Utilizar

1. **Configuração do Engine**

O core do projeto é o ChaosConfig, onde você define a agressividade do seu teste.

```golang
cfg := chaos.ChaosConfig{
    Probability: 0.1, // 10% das requisições sofrerão intervenção
    Latency: chaos.ChaosConfigLatency{
        Min: 100 * time.Millisecond,
        Max: 2 * time.Second,
    },
    Error: codes.Internal, // Erro gRPC a ser injetado
}
```

2. **Registro do Interceptor**

Basta adicionar o interceptor na criação do seu servidor gRPC:
```golang
chaosEngine := chaos.NewChaos(cfg)
interceptor := interceptors.NewInterceptor(chaosEngine)

s := grpc.NewServer(
    grpc.UnaryInterceptor(interceptor.UnaryInterceptor),
    grpc.StreamInterceptor(interceptor.StreamInterceptor),
)
```

## Estrutura do Projeto
`pkg/chaos`: O motor que decide quando e como falhar.

`pkg/interceptors`: A ponte entre o motor de caos e o gRPC (Unary e Stream).

`example/`: Implementação de um UserService para testes práticos de carga e streaming.

## Decisões de Arquitetura
Este trecho do documento explica o raciocínio por trás das escolhas implementadas no projeto.

### A Escolha da Camada: gRPC Interceptors
Implementar o caos via interceptores (middleware) foi uma escolha estratégica para garantir a **Separation of concerns**.

**Transparência**: A lógica de negócio no `handlers/user_grpc.go` não sabe que o caos existe.

**Abrangência**: Conseguimos interceptar tanto chamadas simples (Unary) quanto fluxos de dados contínuos (Stream), garantindo que falhas possam ocorrer no meio de uma transmissão de dados longa.

### Injeção de Latência (Sleep Progressivo)
A implementação em `pkg/interceptors/interceptor.go` utiliza um select com `time.NewTimer`.

* **Por que não apenas falhar?**

Falhas instantâneas são fáceis de tratar. O cenário mais perigoso em produção é o "Gray Failure", onde o serviço está lento, consumindo threads do chamador e causando um efeito cascata.

* **Respeito ao Contexto**: 

O uso de case ` <-ctx.Done()` no interceptor é crucial. Se o cliente desistir da requisição antes do timer do caos terminar, nós liberamos os recursos imediatamente, simulando o comportamento real de um servidor sob pressão.

## Erros gRPC Específicos
O motor permite configurar qual `codes.Code` será retornado.

* **Motivação**: Testar se o cliente diferencia um `INTERNAL` (erro de código) de um `UNAVAILABLE` (erro de infra). Isso é vital para configurar políticas de retry inteligentes — você não deve dar retry em um erro 4xx, mas deve em um 5xx.

##  Concorrência e Performance (sync.Pool)
Em `pkg/chaos/engine.go`, utilizamos um sync.Pool para gerenciar as instâncias de `rand.Rand`.

* `O Problema`: O gerador de números aleatórios global do Go (rand.Float64()) sofre de contenção de lock em sistemas de alta performance.
* `A Solução`: Ao usar um pool de geradores locais, garantimos que a injeção de caos não se torne ela mesma o gargalo do sistema durante testes de carga massivos.

## Suporte a Streaming
A implementação do `wrappedStream` em `pkg/interceptors/stream.go` é o ponto do projeto em que lidamos com streams gRPC. 

Nesse cenário, a conexão pode ficar aberta por minutos e o caos não pode ocorrer apenas no início. 

Ao interceptar `SendMsg` e `RecvMsg`, podemos simular uma conexão que começa bem e degrada no meio do processo, forçando o desenvolvedor a tratar erros dentro do loop de `Recv()`.
