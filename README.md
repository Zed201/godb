## Godb
Basicamente uma tentativa de criação de um clone do sqlite para aprender mais sobre Golang

Feito seguindo mais ou meno:
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
        Depois irei trocar para uma implementação de Btree, mas ai teria que mudar como salva
        o binário então deixarei para estudos futuros
</details>

#### Detalhes
Basicamente é como um sqlite3, suportando insert, select e delet, 
com apenas os tipos de Int, Float, varchar e bool; Enquanto não trocar a forma de 
armazenamento, basicamente transforma tudo em um slice de bytes e salva uma struct que 
representar o banco de dados, salva tudo num binário com o nome do banco, aí usa o .use 
para "abrir" esse arquivo e usa os comandos sql suportados
O insert é o normal, o select ele funciona com o where e passando os nomes as colunas e o delete funciona normalmente
