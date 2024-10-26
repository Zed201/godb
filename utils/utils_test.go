package utils

import (
	"testing"
	"time"

	"golang.org/x/exp/rand"
)

var num_test = 100

func TestVarchar(t *testing.T) {
	tests := []string{"Test iwdiwd", "Luiz Gustavo", "wibfiwbef", "a         o"}
	dec := Decoders[0]
	enc := Encoders[0]
	Emsg := "Test de string %s deu erro"
	for _, s := range tests {
		res, e := enc(s)
		if e != nil {
			t.Fatalf(Emsg, s)
		}
		S, e := dec(res)
		if S != s {
			t.Fatalf(Emsg, s)
		}
	}
}

func TestInt(t *testing.T) {
	tests := randomInt32Slice(num_test)
	dec := Decoders[1]
	enc := Encoders[1]
	Emsg := "Test de Int, esperava %v, consegui %v, %T != %T"
	for _, s := range tests {
		res, e := enc(s)
		if e != nil {
			t.Fatalf("Erro de nil")
		}
		S, e := dec(res)
		// Conversão de interface{} para int32 e depois int
		if v, o := S.(int32); o {
			i := int32(v)
			if i != s {
				t.Fatalf(Emsg, s, i, s, i)
			}
		}
	}
}

func TestFloat(t *testing.T) {
	tests := randomFloatSlice(num_test, 50)
	dec := Decoders[2]
	enc := Encoders[2]
	Emsg := "Test de Float, esperava %v, consegui %v"

	for _, s := range tests {
		res, e := enc(s)
		if e != nil {
			t.Fatalf("Erro de nil")
		}
		S, e := dec(res)
		if v, o := S.(float32); o {
			i := float32(v)
			if i != s {
				t.Fatalf(Emsg, s, i, s, i)
			}
		}
	}
}

func TestBool(t *testing.T) {
	tests := randomBoolSlice(num_test)
	dec := Decoders[3]
	enc := Encoders[3]

	Emsg := "Test de Bool, esperava %s, consegui %s"
	for _, s := range tests {
		res, e := enc(s)
		if e != nil {
			t.Fatalf("Erro de nil")
		}
		S, e := dec(res)
		if S != s {
			t.Fatalf(Emsg, s, S)
		}
	}
}

func randomInt32Slice(size int) []int32 {
	rand.Seed(uint64(time.Now().UnixNano())) // Convertendo UnixNano() para uint64
	slice := make([]int32, size)
	for i := range slice {
		slice[i] = rand.Int31() // Gera números aleatórios de 0 até max - 1
	}
	return slice
}

// Função para gerar um slice de floats aleatórios
func randomFloatSlice(size int, max float32) []float32 {
	rand.Seed(uint64(time.Now().UnixNano()))
	slice := make([]float32, size)
	for i := range slice {
		slice[i] = rand.Float32() * max // Gera floats de 0 até max
	}
	return slice
}

// Função para gerar um slice de bools aleatórios
func randomBoolSlice(size int) []bool {
	rand.Seed(uint64(time.Now().UnixNano()))
	slice := make([]bool, size)
	for i := range slice {
		slice[i] = rand.Intn(2) == 1 // 50% de chance de ser true ou false
	}
	return slice
}
