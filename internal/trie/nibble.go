package trie

import "github.com/crypto2lab/ground/lib/comparable"

func NibblesToKey(nibbles []byte) []byte {
	if len(nibbles)%2 == 0 {
		key := make([]byte, len(nibbles)/2)
		for idx := 0; idx < len(nibbles); idx += 2 {
			key[idx/2] = (nibbles[idx] << 4 & 0xf0) | (nibbles[idx] & 0xf)
		}
		return key
	}

	key := make([]byte, len(nibbles)/2+1)
	key[0] = nibbles[0]
	for idx := 2; idx < len(nibbles); idx += 2 {
		key[idx/2] = (nibbles[idx-1] << 4 & 0xf0) | (nibbles[idx] & 0xf)
	}

	return key
}

func KeyToNibbles(key []byte) []byte {
	if len(key) == 0 {
		return []byte{}
	} else if len(key) == 1 && key[0] == 0 {
		return []byte{0, 0}
	}

	length := len(key) * 2
	nibbles := make([]byte, length)

	for idx, b := range key {
		nibbles[2*idx] = b / 16
		nibbles[2*idx+1] = b % 16
	}

	return nibbles
}

func nibblesEqualsUntil(nibblesA, nibblesB []byte) int {
	maxIterations := comparable.Min(len(nibblesA), len(nibblesB))
	for idx := 0; idx < maxIterations; idx++ {
		if nibblesA[idx] == nibblesB[idx] {
			continue
		}

		return idx
	}

	return maxIterations
}
