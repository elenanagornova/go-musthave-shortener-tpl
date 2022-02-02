package crypto

import (
	"crypto/aes"
	"fmt"
	"testing"
)

func TestCrypto(t *testing.T) {
	key, err := generateRandom(aes.BlockSize) // ключ шифрования
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	type args struct {
		key  []byte
		text []byte
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "simple",
			args: args{
				key:  key,
				text: []byte("hello, world"),
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, nonce, err := Encrypt(tt.args.key, tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			decrypted, err := Decrypt(tt.args.key, nonce, encrypted)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			isEqual := string(tt.args.text) == string(decrypted)
			if tt.want != isEqual {
				t.Errorf("Expected equal = %v, got %v", tt.want, isEqual)
				return
			}
		})
	}
}
