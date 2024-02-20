package client

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"math/big"

	"github.com/huandu/xstrings"
)

// IntToHex32 Normalize to a 32-byte hex string without 0x prefix.
func IntToHex32(x *big.Int) string {
	str := x.Text(16)
	return xstrings.RightJustify(str, 64, "0")
}

// SerializeSignature Convert a Sign from an r, s pair to a 32-byte hex string.
func SerializeSignature(r, s *big.Int) string {
	return IntToHex32(r) + IntToHex32(s)
}


func Sign(starkPrivateKey string, hashStr string) (string, error) {
	r, s1 := doSign(starkPrivateKey, hashStr)
	log.Println("\nr: ", r.Text(16))
	log.Println("\ns1: ", s1.Text(16))

	// no need to serialize directly use this r and s
	return SerializeSignature(r, s1), nil
}

// helper function usedfor signing
func GenerateKRfc6979(msgHash, priKey *big.Int, seed int) *big.Int {
	msgHash = big.NewInt(0).Set(msgHash) // copy
	bitMod := msgHash.BitLen() % 8
	if bitMod <= 4 && bitMod >= 1 && msgHash.BitLen() > 248 {
		msgHash.Mul(msgHash, big.NewInt(16))
	}

	var extra []byte
	if seed > 0 {
		buf := new(bytes.Buffer)
		var data interface{}
		if seed < 256 {
			data = uint8(seed)
		} else if seed < 65536 {
			data = uint16(seed)
		} else if seed < 4294967296 {
			data = uint32(seed)
		} else {
			data = uint64(seed)
		}
		_ = binary.Write(buf, binary.BigEndian, data)
		extra = buf.Bytes()
	}

	return generateSecret(EC_ORDER, priKey, sha256.New, msgHash.Bytes(), extra)
}

// N_ELEMENT_BITS_ECDSA math.floor(math.log(FIELD_PRIME, 2))
var N_ELEMENT_BITS_ECDSA = big.NewInt(251)

func doSign(starkPrivateKey string, hash string) (*big.Int, *big.Int) {
	priKey, _ := new(big.Int).SetString(starkPrivateKey, 16)


	// convert from hex instead of decimal
	msgHash, _ := new(big.Int).SetString(hash, 16)


	seed := 0
	EcGen := pedersenCfg.ConstantPoints[1]
	alpha := pedersenCfg.ALPHA
	nBit := big.NewInt(0).Exp(big.NewInt(2), N_ELEMENT_BITS_ECDSA, nil)
	for {
		k := GenerateKRfc6979(msgHash, priKey, seed)
		//	Update seed for next iteration in case the value of k is bad.
		if seed == 0 {
			seed = 1
		} else {
			seed += 1
		}
		// Cannot fail because 0 < k < EC_ORDER and EC_ORDER is prime.
		x := ecMult(k, EcGen, alpha, FIELD_PRIME)[0]
		// !(1 <= x < 2 ** N_ELEMENT_BITS_ECDSA)
		if !(x.Cmp(one) > 0 && x.Cmp(nBit) < 0) {
			continue
		}
		// msg_hash + r * priv_key
		x1 := big.NewInt(0).Add(msgHash, big.NewInt(0).Mul(x, priKey))
		// (msg_hash + r * priv_key) % EC_ORDER == 0
		if big.NewInt(0).Mod(x1, EC_ORDER).Cmp(zero) == 0 {
			continue
		}
		// w = div_mod(k, msg_hash + r * priv_key, EC_ORDER)
		w := divMod(k, x1, EC_ORDER)
		// not (1 <= w < 2 ** N_ELEMENT_BITS_ECDSA)
		if !(w.Cmp(one) > 0 && w.Cmp(nBit) < 0) {
			continue
		}
		s1 := divMod(one, w, EC_ORDER)
		return x, s1
	}
}