package sessions

import "testing"

func TestGenerateSession(t *testing.T) {
	public, private, _ := GenerateSession("username")

	if private != GetPrivateKey(public) {
		t.Errorf("GenerateSession got private %v, want %v", private, GetPrivateKey(public))
	}
}

func BenchmarkGenerateSession(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateSession("username")
	}
}
