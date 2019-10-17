//
// Copyright 2018 Malin Yamato Lääkkö --  All rights reserved.
// https://github.com/MalinYamato
//
// MIT License
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Rakuen. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
)

func encode64(b []byte) string {
	return base64.URLEncoding.EncodeToString(b)
}

func decode64(s string) []byte {
	data, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}
func decodeHex(data []byte) ([]byte, error) {
	decoded := make([]byte, hex.DecodedLen(len(data)))
	_, err := hex.Decode(decoded, data)
	if err != nil {
		return nil, err
	}
	fmt.Println("decoded ", decoded)
	return decoded, nil
}
func encodeHex(data []byte) []byte {
	encoded := make([]byte, hex.EncodedLen(len(data)))
	hex.Encode(encoded, data)

	fmt.Printf("enoded %s\n", encoded)
	return encoded
}

func generateRandomBytes(n int) ([]byte, error) {
	// generatee 64 bit key
	key := make([]byte, 64)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		panic(err.Error())
	}
	return key, nil
}

func createHash(key string) []byte {
	hash := sha256.Sum256([]byte(key))
	return hash[:]
}

func encrypt(data []byte, pvtKey []byte) []byte {
	key := make([]byte, hex.DecodedLen(len(pvtKey)))
	_, err := hex.Decode(key, pvtKey)
	if err != nil {
		log.Fatal(err)
	}

	block, _ := aes.NewCipher(key)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func decrypt(data []byte, pvtKey []byte) ([]byte, error) {
	key := make([]byte, hex.DecodedLen(len(pvtKey)))
	_, err := hex.Decode(key, pvtKey)
	if err != nil {
		log.Fatal(err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func readKeyFile(path string) []byte {
	file, err := os.Open(path) // For read access.
	if err != nil {
		panic(err)
	}
	pvtKey := make([]byte, 64)
	count, err := file.Read(pvtKey)
	if err != nil {
		panic(err)
	}
	fmt.Printf("read %d bytes: %q\n", count, pvtKey[:count])
	return pvtKey
}

func test_main() {

	clearText := "[{'Hello' : 'World'}][{'Hello' : 'World'}]"
	fmt.Println(clearText)

	pvtKey := readKeyFile("private.key")

	ciphertext := encrypt([]byte(clearText), pvtKey)
	fmt.Println(ciphertext)

	encoded := encodeHex(ciphertext)
	// <--------------network ------------->
	decoded, _ := decodeHex(encoded)
	plaintext, _ := decrypt(decoded, pvtKey)
	fmt.Println(string(plaintext))
}
