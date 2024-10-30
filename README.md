## Godb
Basicamente uma tentativa de criação de um clone do sqlite para aprender mais sobre Golang

Feito seguindo:
* [Medium](https://medium.com/felixklauke/database-i-developing-your-own-data-storage-engine-aka-create-your-own-database-ed4560c8d80a)
* [Build your own](https://build-your-own.org/database/90_end)
* [Sqlite docs](https://www.sqlite.org/arch.html)
* [Cstack](https://cstack.github.io/db_tutorial/)

As partes do projeto são:

<details>
<summary>Interface</summary>
        Funções relacionadas ao REPL, entrada e saída de dados
</details>

<details>
<summary>Processor(Tokenizer e Parser)</summary>
        Tem basicamente as funções de Parser para as instruções e Tokenizer,
        que são passadas para o core, atráves do Interface
</details>

<details>
<summary>Core</summary>
        Implementação das execuções dos comandos, além das outras estrutura para
        auxiliar, no momento está armazenando tudo em um slice de []byte
        Em breve vou trocar para alguma estrutura de dados mais complexa
</details>
