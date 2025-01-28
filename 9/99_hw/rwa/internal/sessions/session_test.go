package sessions

import (
	"reflect"
	"testing"
	"time"
)

func TestGenerateSession(t *testing.T) {
	public, _ := GenerateSession("username")
	private := GetPrivateKey(public)

	if private != GetPrivateKey(public) {
		t.Errorf("GenerateSession got private %v, want %v", private, GetPrivateKey(public))
	}
}

func BenchmarkGenerateSession(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateSession("username")
	}
}

func TestSessionManager(t *testing.T) {
	sm := NewSessionManager()

	// Создание сессий
	TestsCaseCreate := []struct {
		name      string
		username  string
		wantExist bool
	}{
		{
			name:      "ValidUser1",
			username:  "ValidUser1",
			wantExist: true,
		},
		{
			name:      "ValidUser2",
			username:  "ValidUser2",
			wantExist: true,
		},
		{
			name:      "ExistUser",
			username:  "ValidUser1",
			wantExist: true,
		},
	}

	testStore := make(map[string]*Session)

	publicKeys := make(map[string]string)

	for _, tt := range TestsCaseCreate {
		t.Run(tt.name, func(t *testing.T) {
			publicKey, err := sm.Create(tt.username)
			if err != nil {
				t.Errorf("sm.Create got err %v, want nil", err)
			}

			privateKey := GetPrivateKey(publicKey)

			sess, ok := sm.store[privateKey]
			if ok != tt.wantExist {
				t.Errorf("sm.Create got ok %v, want ok %v", ok, tt.wantExist)
			}

			if sess.Username != tt.username {
				t.Errorf("sm.Create got sess.Username %v, want %v", sess.Username, tt.username)
			}

			testStore[privateKey] = sess
			publicKeys[publicKey] = sess.Username
		})
	}

	if !reflect.DeepEqual(sm.store, testStore) {
		t.Errorf("store error: want %v, got %v", testStore, sm.store)
	}

	// SessionManager.Check OK
	for publicKey, WantUsername := range publicKeys {
		t.Run("check test: "+WantUsername, func(t *testing.T) {
			gotUsername, ok := sm.Check(publicKey)
			if !ok {
				t.Errorf("sm.Check got ok %v, want true", ok)
			}

			if gotUsername != WantUsername {
				t.Errorf("sm.Check got %v, want %v", gotUsername, WantUsername)
			}
		})
	}

	// SessionManager.Check invalid PrivateKey
	t.Run("invalid privatekey", func(t *testing.T) {
		gotUsername, ok := sm.Check("invalidkey!!@#3432")
		if ok {
			t.Errorf("sm.Check got ok %v, want false", ok)
		}

		if gotUsername != "" {
			t.Errorf(`sm.Check got %v, want ""`, gotUsername)
		}
	})

	// SessionManager.Check expired
	publicKey, err := sm.Create("ExpireUser")
	if err != nil {
		t.Errorf("sm.Create got err %v, want nil", err)
	}

	sess := sm.store[GetPrivateKey(publicKey)]
	if sess.Expire > time.Now().Unix()+ExpirationSession || sess.Expire < time.Now().Unix() {
		t.Errorf("Expiration error: session expired %v", sess.Expire)
	}

	GotUsername, expirationOK := sm.Check(publicKey)
	if expirationOK != true || GotUsername != sess.Username {
		t.Errorf("Check error: %v, %v", GotUsername, expirationOK)
	}

	sess.Expire = time.Now().Unix() - 2
	sm.store[GetPrivateKey(publicKey)] = sess

	GotUsername, expirationOK = sm.Check(publicKey)
	if expirationOK != false || GotUsername != "" {
		t.Errorf("Check error: %v, %v", GotUsername, expirationOK)
	}

	// SessionManager.Destroy
	t.Run("DestroyTest", func(t *testing.T) {
		// Создаем тестовую сессию
		destroyToken, err := sm.Create("DestroyUser")
		if err != nil {
			t.Errorf("sm.Create got err %v, want nil", err)
		}

		// проверяем успешное создание
		destroyUsername, ok := sm.Check(destroyToken)
		if ok == false || destroyUsername != "DestroyUser" {
			t.Errorf("Check error: %v, %v", ok, destroyUsername)
		}

		// уничтожаем сессию по токену
		username, err := sm.DestroyByToken(destroyToken)
		if err != nil || username != "DestroyUser" {
			t.Errorf("sm.DestroyTokenSession got err %v, want nil", err)
		}

		// Проверяем наличие удаленной сессии
		destroyUsername, ok = sm.Check(destroyToken)
		if ok == true || destroyUsername != "" {
			t.Errorf("Check error: %v, %v", ok, destroyUsername)
		}

		// Создаем 10 сессия от одного пользователя
		testKeys := make([]string, 10)
		for i := 0; i < 10; i++ {
			pkey, err := sm.Create("DestroyByUsername")
			if err != nil {
				t.Errorf("sm.Create got err %v, want nil", err)
			}

			testKeys[i] = pkey
		}

		// проверяем, что успешно создали
		for _, key := range testKeys {
			un, ok := sm.Check(key)
			if ok != true || un == "" {
				t.Errorf("Check error: %v, %v", ok, un)
			}
		}

		deleted, err := sm.DestroyByUsername("DestroyByUsername")
		if deleted != 10 || err != nil {
			t.Errorf("DestroyByUsername got deleted %v, want %v", deleted, 10)
		}

		// Проверяем наличие сессий
		for _, key := range testKeys {
			un, ok := sm.Check(key)
			if ok != false || un != "" {
				t.Errorf("Check error: %v, %v", ok, un)
			}
		}
	})

	t.Run("ClearExpiredSession", func(t *testing.T) {
		sm.store = map[string]*Session{}

		expiredCases := []struct {
			username  string
			deltaTime int64
		}{
			{
				username:  "delta1000",
				deltaTime: 1000,
			},
			{
				username:  "delta100",
				deltaTime: 100,
			},
			{
				username:  "delta2",
				deltaTime: 2,
			},
			{
				username:  "deltaReverse100",
				deltaTime: -100,
			},
			{
				username:  "deltaReverse10",
				deltaTime: -10,
			},
		}

		for _, expiredCase := range expiredCases {
			public, err := sm.Create(expiredCase.username)
			if err != nil {
				t.Errorf("sm.Create got err %v, want nil", err)
			}

			private := GetPrivateKey(public)

			sm.store[private].Expire = time.Now().Unix() + expiredCase.deltaTime
		}

		if len(sm.store) != len(expiredCases) {
			t.Errorf("Lens store not matched do clearExpired func")
		}

		deleted := sm.ClearExpired()
		if deleted != 2 {
			t.Errorf("deleted got %v, want %v", deleted, 2)
		}

		if len(sm.store) != 3 {
			t.Errorf("unexpected result after len store: %d", len(sm.store))
		}
	})

}
