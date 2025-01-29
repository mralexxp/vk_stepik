package password

import "testing"

func TestHash(t *testing.T) {
	testCase := []struct {
		name       string
		FirstPass  string
		SecondPass string
		want       bool
	}{
		{
			name:       "valid pass",
			FirstPass:  "123456",
			SecondPass: "123456",
			want:       true,
		},
		{
			name:       "invalid pass",
			FirstPass:  "123456",
			SecondPass: "12345",
			want:       false,
		},
		{
			name:       "short pass",
			FirstPass:  "1",
			SecondPass: "1",
			want:       true,
		},
		{ // true, так как проверяем пароль непосредственно при валидации структуры
			name:       "empty pass",
			FirstPass:  "",
			SecondPass: "",
			want:       true,
		},
		{
			name:       "Very long valid pass",
			FirstPass:  "123456jfdskfdkkkdfskaljalfsdjlksjaflkdjf*FAFd898af981274*$**#$*8iieriiueriuiewfjlkjalkdsfjklksjdjfsahjvvnxvjjhsdf",
			SecondPass: "123456jfdskfdkkkdfskaljalfsdjlksjaflkdjf*FAFd898af981274*$**#$*8iieriiueriuiewfjlkjalkdsfjklksjdjfsahjvvnxvjjhsdf",
			want:       true,
		},
		{
			name:       "Very long invalid pass",
			FirstPass:  "23456jfdskfdkkkdfskaljalfsdjlksjaflkdjf*FAFd898af981274*$**#$*8iieriiueriuiewfjlkjalkdsfjklksjdjfsahjvvnxvjjhsdf",
			SecondPass: "123456jfdskfdkkkdfskaljalfsdjlksjaflkdjf*FAFd898af981274*$**#$*8iieriiueriuiewfjlkjalkdsfjklksjdjfsahjvvnxvjjhsdf",
			want:       false,
		},
	}

	for _, tt := range testCase {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := Hash(tt.FirstPass)
			if err != nil {
				t.Error(err)
			}

			ok := Check(tt.SecondPass, hash)
			if ok != tt.want {
				t.Errorf("want %v, got %v", tt.want, ok)
			}
		})
	}
}

func BenchmarkHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Hash("12345678HFjdfjhhdfshjfdjkh&@&#762371uifijafj787958437987543298")
	}
}

func BenchmarkCheck(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Check("123456", "hfjsdhfjsdhdjskf")
	}
}
