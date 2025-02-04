package profile

import (
	"reflect"
	"rwa/internal/models"
	"testing"
)

var testdb *Store

func TestMain(m *testing.M) {
	testdb = NewStore()
	db := &testdb.db
	m.Run()
	_ = db
}

func TestStore(t *testing.T) {
	testCase := []*models.Profile{
		&models.Profile{
			ID:        1,
			Follow:    make(map[uint64]struct{}),
			Followers: make(map[uint64]struct{}),
		},
		&models.Profile{
			ID:        2,
			Follow:    make(map[uint64]struct{}),
			Followers: make(map[uint64]struct{}),
		},
		&models.Profile{
			ID:        3,
			Follow:    make(map[uint64]struct{}),
			Followers: make(map[uint64]struct{}),
		},
	}

	t.Run("AddProfile", func(t *testing.T) {
		// Добавление
		for _, tt := range testCase {
			err := testdb.AddProfile(tt)
			if err != nil {
				t.Fatal(err)
			}
		}

		// Проверка
		for _, tt := range testCase {
			if profile, ok := testdb.db[tt.ID]; !ok {
				t.Fatal("no add db: ", tt.ID)
			} else {
				if !reflect.DeepEqual(profile, tt) {
					t.Fatal("not equal profile: ", tt.ID, tt.Follow)
				}
			}
		}
	})

	t.Run("GetProfile", func(t *testing.T) {
		for _, tt := range testCase {
			if profile, err := testdb.GetProfile(tt.ID); err != nil {
				t.Fatal(err)
			} else {
				if !reflect.DeepEqual(profile, tt) {
					t.Fatal("not equal profile: ", tt.ID, tt.Follow)
				}
			}
		}
	})

	t.Run("Follow", func(t *testing.T) {
		// 1 -> 3
		// 1, 3 -> 2
		err := testdb.Follow(1, 3)
		if err != nil {
			t.Fatal(err)
		}

		err = testdb.Follow(1, 2)
		if err != nil {
			t.Fatal(err)
		}

		err = testdb.Follow(3, 2)
		if err != nil {
			t.Fatal(err)
		}

		expect := []struct {
			id        uint64
			follow    map[uint64]struct{}
			followers map[uint64]struct{}
		}{
			{
				id: 1,
				follow: map[uint64]struct{}{
					3: {},
					2: {},
				},
				followers: map[uint64]struct{}{},
			},
			{
				id:     2,
				follow: map[uint64]struct{}{},
				followers: map[uint64]struct{}{
					1: {},
					3: {},
				},
			},
			{
				id: 3,
				follow: map[uint64]struct{}{
					2: {},
				},
				followers: map[uint64]struct{}{
					1: {},
				},
			},
		}

		for i, tt := range testCase {
			follow := testdb.db[tt.ID].Follow
			followers := testdb.db[tt.ID].Followers
			if !reflect.DeepEqual(follow, expect[i].follow) {
				t.Fatal("not equal follow: ", tt.ID, tt.Follow)
			}

			if !reflect.DeepEqual(followers, expect[i].followers) {
				t.Fatal("not equal followers: ", tt.ID, tt.Followers)
			}
		}

	})

	t.Run("Unfollow", func(t *testing.T) {
		err := testdb.Unfollow(1, 3)
		if err != nil {
			t.Fatal(err)
		}

		expectFollow := map[uint64]struct{}{2: {}}
		expectFollowers := map[uint64]struct{}{}

		if !reflect.DeepEqual(testdb.db[1].Follow, expectFollow) {
			t.Fatal("not equal follow: ", testdb.db[1].Follow, expectFollow)
		}

		if !reflect.DeepEqual(testdb.db[3].Followers, expectFollowers) {
			t.Fatal("not equal followers: ", testdb.db[3].Followers, expectFollowers)
		}

	})

}
