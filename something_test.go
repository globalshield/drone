package main

import (
    "crypto/md5"
    "crypto/sha1"
    "crypto/sha256"
    "crypto/sha512"
    "hash"
    "hash/fnv"
    "testing"
)

const (
    K       = 1024
    DATALEN = 1 * K
)

func runHash(b *testing.B, h hash.Hash, n int) {
    var data = make([]byte, n)
    for i := 0; i < n; i++ {
        data[i] = byte(i * i)
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        h.Reset()
        h.Sum(data)
    }
}

func BenchmarkFNV32(b *testing.B) {
    runHash(b, fnv.New32(), DATALEN)
}

func BenchmarkFNV64(b *testing.B) {
    runHash(b, fnv.New64(), DATALEN)
}

func BenchmarkFNV128(b *testing.B) {
    runHash(b, fnv.New128(), DATALEN)
}

func BenchmarkMD5(b *testing.B) {
    runHash(b, md5.New(), DATALEN)
}

func BenchmarkSHA1(b *testing.B) {
    runHash(b, sha1.New(), DATALEN)
}

func BenchmarkSHA224(b *testing.B) {
    runHash(b, sha256.New224(), DATALEN)
}

func BenchmarkSHA256(b *testing.B) {
    runHash(b, sha256.New(), DATALEN)
}

func BenchmarkSHA512(b *testing.B) {
    runHash(b, sha512.New(), DATALEN)
}
