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
        Está no Inter também, pois no fim é uma parte da interface, basicamente 
        vai interpretar os comandos e passar para o core
</details>

<details>
<summary>Core</summary>
        Implementação das execuções dos comandos, além das outras estruturas 
        para armazenar dados e outros
</details>
