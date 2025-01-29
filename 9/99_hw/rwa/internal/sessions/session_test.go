package sessions

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func BenchmarkGenerateSession(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateSession()
	}
}

func TestSessionManager(t *testing.T) {
	sm := NewSessionManager()

	// Создание сессий
	TestsCaseCreate := []struct {
		name      string
		id        uint64
		wantExist bool
	}{
		{
			name:      "ValidUser1",
			id:        1,
			wantExist: true,
		},
		{
			name:      "ValidUser2",
			id:        2,
			wantExist: true,
		},
		{
			name:      "ExistUser",
			id:        1,
			wantExist: true,
		},
	}

	testStore := make(map[string]*Session)

	publicKeys := make(map[string]uint64)

	for _, tt := range TestsCaseCreate {
		t.Run(tt.name, func(t *testing.T) {
			publicKey, err := sm.Create(tt.id)
			if err != nil {
				t.Fatalf("sm.Create got err %v, want nil", err)
			}

			privateKey := GetPrivateKey(publicKey)

			sess, ok := sm.store[privateKey]
			if ok != tt.wantExist {
				t.Fatalf("sm.Create got ok %v, want ok %v", ok, tt.wantExist)
			}

			if sess.ID != tt.id {
				t.Fatalf("sm.Create got sess.Username %v, want %v", sess.ID, tt.id)
			}

			testStore[privateKey] = sess
			publicKeys[publicKey] = sess.ID
		})
	}

	if !reflect.DeepEqual(sm.store, testStore) {
		t.Fatalf("store error: want %v, got %v", testStore, sm.store)
	}

	// SessionManager.Check OK
	for publicKey, WantID := range publicKeys {
		t.Run(fmt.Sprintf("Check test for id: %d", WantID), func(t *testing.T) {
			gotUsername, ok := sm.Check(publicKey)
			if !ok {
				t.Fatalf("sm.Check got ok %v, want true", ok)
			}

			if gotUsername != WantID {
				t.Fatalf("sm.Check got %v, want %v", gotUsername, WantID)
			}
		})
	}

	// SessionManager.Check invalid PrivateKey
	t.Run("invalid privatekey", func(t *testing.T) {
		gotUsername, ok := sm.Check("invalidkey!!@#3432")
		if ok {
			t.Fatalf("sm.Check got ok %v, want false", ok)
		}

		if gotUsername != 0 {
			t.Fatalf(`sm.Check got %v, want ""`, gotUsername)
		}
	})

	// SessionManager.Check expired
	publicKey, err := sm.Create(1)
	if err != nil {
		t.Fatalf("sm.Create got err %v, want nil", err)
	}

	sess := sm.store[GetPrivateKey(publicKey)]
	if sess.Expire > time.Now().Unix()+ExpirationSession || sess.Expire < time.Now().Unix() {
		t.Fatalf("Expiration error: session expired %v", sess.Expire)
	}

	GotID, expirationOK := sm.Check(publicKey)
	if expirationOK != true || GotID != sess.ID {
		t.Fatalf("Check error: %v, %v", GotID, expirationOK)
	}

	sess.Expire = time.Now().Unix() - 2
	sm.store[GetPrivateKey(publicKey)] = sess

	GotID, expirationOK = sm.Check(publicKey)
	if expirationOK != false || GotID != 0 {
		t.Fatalf("Check error: %v, %v", GotID, expirationOK)
	}

	// SessionManager.Destroy
	t.Run("DestroyTest", func(t *testing.T) {
		var DestroyUserID uint64 = 11

		// Создаем тестовую сессию
		destroyToken, err := sm.Create(DestroyUserID)
		if err != nil {
			t.Fatalf("sm.Create got err %v, want nil", err)
		}

		// проверяем успешное создание
		destroyID, ok := sm.Check(destroyToken)
		if ok == false || destroyID != DestroyUserID {
			t.Fatalf("Check error: %v, %v", ok, destroyID)
		}

		// уничтожаем сессию по токену
		id, err := sm.DestroyByToken(destroyToken)
		if err != nil || id != DestroyUserID {
			t.Fatalf("sm.DestroyTokenSession got err %v, want nil", err)
		}

		// Проверяем отсутствие удаленной сессии
		destroyID, ok = sm.Check(destroyToken)
		if ok || destroyID == DestroyUserID {
			t.Fatalf("Check error: %v, %v", ok, destroyID)
		}

		// Создаем 10 сессий от одного пользователя
		testKeys := make([]string, 10)
		var destroyByIDUser uint64 = 22

		for i := 0; i < 10; i++ {
			pkey, err := sm.Create(destroyByIDUser)
			if err != nil {
				t.Fatalf("sm.Create got err %v, want nil", err)
			}

			testKeys[i] = pkey
		}

		// проверяем, что успешно создали
		for _, key := range testKeys {
			uid, ok := sm.Check(key)
			if ok != true || uid == 0 {
				t.Fatalf("Check error: %v, %v", ok, uid)
			}
		}

		deleted, err := sm.DestroyByID(destroyByIDUser)
		if deleted != 10 || err != nil {
			t.Fatalf("DestroyByUsername got deleted %v, want %v", deleted, 10)
		}

		// Проверяем наличие сессий
		for _, key := range testKeys {
			uid, ok := sm.Check(key)
			if ok != false || uid != 0 {
				t.Fatalf("Check error: %v, %v", ok, uid)
			}
		}
	})

	t.Run("ClearExpiredSession", func(t *testing.T) {
		sm.store = map[string]*Session{}

		expiredCases := []struct {
			id        uint64
			deltaTime int64
		}{
			{
				id:        1000,
				deltaTime: 1000,
			},
			{
				id:        100,
				deltaTime: 100,
			},
			{
				id:        2,
				deltaTime: 2,
			},
			{
				id:        100100,
				deltaTime: -100,
			},
			{
				id:        1010,
				deltaTime: -10,
			},
		}

		for _, expiredCase := range expiredCases {
			public, err := sm.Create(expiredCase.id)
			if err != nil {
				t.Fatalf("sm.Create got err %v, want nil", err)
			}

			private := GetPrivateKey(public)

			sm.store[private].Expire = time.Now().Unix() + expiredCase.deltaTime
		}

		if len(sm.store) != len(expiredCases) {
			t.Fatalf("Lens store not matched do clearExpired func")
		}

		deleted := sm.ClearExpired()
		if deleted != 2 {
			t.Fatalf("deleted got %v, want %v", deleted, 2)
		}

		if len(sm.store) != 3 {
			t.Fatalf("unexpected result after len store: %d", len(sm.store))
		}
	})

	t.Run("DestroyByID", func(t *testing.T) {
		sm.store = map[string]*Session{}

		// Создаем 3 сессии:
		// id1: 1 шт
		// id2: 2 шт
		publicKey1, err := sm.Create(1)
		if err != nil {
			t.Fatalf("sm.Create got err %v, want nil", err)
		}

		privateKey1 := GetPrivateKey(publicKey1)

		if _, ok := sm.store[privateKey1]; !ok {
			t.Fatalf("Store not created")
		}

		publicKey2, err := sm.Create(2)
		if err != nil {
			t.Fatalf("sm.Create got err %v, want nil", err)
		}

		privateKey2 := GetPrivateKey(publicKey2)

		if _, ok := sm.store[privateKey2]; !ok {
			t.Fatalf("Store not created")
		}

		publicKey3, err := sm.Create(1)
		if err != nil {
			t.Fatalf("sm.Create got err %v, want nil", err)
		}

		privateKey3 := GetPrivateKey(publicKey3)

		if _, ok := sm.store[privateKey3]; !ok {
			t.Fatalf("Store not created")
		}

		// Удаляем сначала id2, проверяем, delete id1
		deleted, err := sm.DestroyByID(1)
		if deleted != 2 || err != nil {
			t.Fatalf("DestroyByID got deleted %v, want %v", deleted, 2)
		}

		if len(sm.store) != 1 {
			t.Fatalf("unexpected result after len store: %d", len(sm.store))
		}

		deleted, err = sm.DestroyByID(2)
		if deleted != 1 || err != nil {
			t.Fatalf("DestroyByID got deleted %v, want %v", deleted, 1)
		}

		if len(sm.store) != 0 {
			t.Fatalf("unexpected result after len store: %d", len(sm.store))
		}
	})
}
