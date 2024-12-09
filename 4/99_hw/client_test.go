package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

var (
	testServer  *httptest.Server
	errorServer *httptest.Server
	mockUsers   *[]User
)

type mockTransport struct{}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("custom error") // При любом запросе вернет "custom error"??
}

func TestMain(m *testing.M) {
	testServer = httptest.NewServer(http.HandlerFunc(SearchServer))
	defer testServer.Close()

	mockUsers, _ = Search(&SearchRequest{
		Limit:      99,
		Offset:     0,
		Query:      "",
		OrderField: "",
		OrderBy:    0,
	})

	// Сервер для вызова ошибок
	errorServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request, err := parseRequest(r)
		if err != nil {
			w.WriteHeader(err.Code())
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		switch request.Query {
		case "timeout_test":
			time.Sleep(2 * time.Second)
			w.WriteHeader(http.StatusOK)
		case "token_test":
			if r.Header.Get("AccessToken") == "your token here" {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"error":"Unauthorized"}`))
				return
			}
			_, _ = w.Write([]byte(`"{"error": "BAD TEST"}"`))
		case "internal_error_test":
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"Internal server error"}`))
		case "json_invalid_test":
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`"Error":"Internal server error"}`)) // missed "{"
		case "invalid_response_json":
			w.WriteHeader(http.StatusOK)
			// ID is not int
			invalidJsonResponse := `{
				"Id":     "abc",
				"Name":   "Firstname and Lastname",
				"Age":    0,
				"About":  "about text",
				"Gender": "male"
			}`
			_, _ = w.Write([]byte(invalidJsonResponse))
		}
	}))
	defer errorServer.Close()

	ret := m.Run()
	os.Exit(ret)
}

func TestSearchClient_FindUsers(t *testing.T) {
	testSrv := &SearchClient{
		AccessToken: "your token here",
		URL:         testServer.URL,
	}
	tests := []struct {
		name    string
		req     SearchRequest
		want    *SearchResponse
		wantErr bool
	}{
		// limit
		{
			name: "limit: limit > 25 (100)",
			req: SearchRequest{
				Limit:      100,
				Offset:     0,
				Query:      "",
				OrderField: "",
				OrderBy:    OrderByAsIs,
			},
			want: &SearchResponse{
				Users:    append([]User{}, (*mockUsers)[:25]...),
				NextPage: true,
			},
			wantErr: false,
		},
		{
			name: "limit: limit + nextPage",
			req: SearchRequest{
				Limit:      10,
				Offset:     0,
				Query:      "",
				OrderField: "",
				OrderBy:    OrderByAsIs,
			},
			want: &SearchResponse{
				Users:    append([]User{}, (*mockUsers)[:10]...),
				NextPage: true,
			},
			wantErr: false,
		},
		{
			name: "limit: limit < 0",
			req: SearchRequest{
				Limit:      -1,
				Offset:     0,
				Query:      "",
				OrderField: "",
				OrderBy:    OrderByAsIs,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "limit+offset: limit+offset == end == len(users)",
			req: SearchRequest{
				Limit:      35,
				Offset:     10,
				Query:      "",
				OrderField: "",
				OrderBy:    OrderByAsIs,
			},
			want: &SearchResponse{
				Users:    append([]User{}, (*mockUsers)[10:35]...),
				NextPage: false,
			},
			wantErr: false,
		},
		// offset
		{
			name: "offset: offset < 0",
			req: SearchRequest{
				Limit:      25,
				Offset:     -1,
				Query:      "",
				OrderField: "",
				OrderBy:    OrderByAsIs,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "offset: offset == 10",
			req: SearchRequest{
				Limit:      25,
				Offset:     10,
				Query:      "",
				OrderField: "",
				OrderBy:    OrderByAsIs,
			},
			want: &SearchResponse{
				Users:    append([]User{}, (*mockUsers)[10:35]...),
				NextPage: false,
			},
			wantErr: false,
		},
		{
			name: "offset: offset == 100",
			req: SearchRequest{
				Limit:      25,
				Offset:     100,
				Query:      "",
				OrderField: "",
				OrderBy:    OrderByAsIs,
			},
			want: &SearchResponse{
				Users:    append([]User{}),
				NextPage: false,
			},
			wantErr: false,
		},
		{
			name: "offset: offset == end == len(users)",
			req: SearchRequest{
				Limit:      25,
				Offset:     35,
				Query:      "",
				OrderField: "",
				OrderBy:    OrderByAsIs,
			},
			want: &SearchResponse{
				Users:    append([]User{}),
				NextPage: false,
			},
			wantErr: false,
		},
		// Query
		{
			name: "query: query == none",
			req: SearchRequest{
				Limit:      25,
				Offset:     0,
				Query:      "",
				OrderField: "",
				OrderBy:    OrderByAsIs,
			},
			want: &SearchResponse{
				Users:    append([]User{}, (*mockUsers)[:25]...),
				NextPage: true,
			},
			wantErr: false,
		},
		{
			name: "query: query == Ex",
			req: SearchRequest{
				Limit:      25,
				Offset:     0,
				Query:      "Ex",
				OrderField: "",
				OrderBy:    OrderByAsIs,
			},
			want: &SearchResponse{
				Users: append([]User{}, (*mockUsers)[4],
					(*mockUsers)[6],
					(*mockUsers)[7],
					(*mockUsers)[9],
					(*mockUsers)[10],
					(*mockUsers)[15],
					(*mockUsers)[26],
					(*mockUsers)[28],
				),
				NextPage: false,
			},
			wantErr: false,
		},
		{
			name: "query: query == ex + nextpage",
			req: SearchRequest{
				Limit:      5,
				Offset:     0,
				Query:      "Ex",
				OrderField: "",
				OrderBy:    OrderByAsIs,
			},
			want: &SearchResponse{
				Users: append([]User{}, (*mockUsers)[4],
					(*mockUsers)[6],
					(*mockUsers)[7],
					(*mockUsers)[9],
					(*mockUsers)[10],
				),
				NextPage: true,
			},
			wantErr: false,
		},
		{
			name: "query: query == ex + nextpage + offset (pagination)",
			req: SearchRequest{
				Limit:      5,
				Offset:     5,
				Query:      "Ex",
				OrderField: "",
				OrderBy:    OrderByAsIs,
			},
			want: &SearchResponse{
				Users: append([]User{},
					(*mockUsers)[15],
					(*mockUsers)[26],
					(*mockUsers)[28],
				),
				NextPage: false,
			},
			wantErr: false,
		},
		{
			name: "query: find in name == Henders (found Henderson)",
			req: SearchRequest{
				Limit:      25,
				Offset:     0,
				Query:      "Henders",
				OrderField: "",
				OrderBy:    OrderByAsIs,
			},
			want: &SearchResponse{
				Users:    append([]User{}, (*mockUsers)[10]),
				NextPage: false,
			},
			wantErr: false,
		},
		// OrderField + OrderBy:
		{
			name: "OrderField: invalid orderField (want nil)",
			req: SearchRequest{
				Limit:      25,
				Offset:     0,
				Query:      "J",
				OrderField: "abc",
				OrderBy:    OrderByAsIs,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "OrderField: Name AsIs",
			req: SearchRequest{
				Limit:      25,
				Offset:     0,
				Query:      "J",
				OrderField: "Name",
				OrderBy:    OrderByAsIs,
			},
			want: &SearchResponse{
				Users: append([]User{},
					(*mockUsers)[6],
					(*mockUsers)[8],
					(*mockUsers)[21],
					(*mockUsers)[25],
				),
				NextPage: false,
			},
			wantErr: false,
		},
		{
			name: "OrderField: Age AsIs",
			req: SearchRequest{
				Limit:      25,
				Offset:     0,
				Query:      "J",
				OrderField: "Age",
				OrderBy:    OrderByAsIs,
			},
			want: &SearchResponse{
				Users: append([]User{},
					(*mockUsers)[6],
					(*mockUsers)[8],
					(*mockUsers)[21],
					(*mockUsers)[25],
				),
				NextPage: false,
			},
			wantErr: false,
		},
		{
			name: "OrderField: Id AsIs",
			req: SearchRequest{
				Limit:      25,
				Offset:     0,
				Query:      "J",
				OrderField: "Age",
				OrderBy:    OrderByAsIs,
			},
			want: &SearchResponse{
				Users: append([]User{},
					(*mockUsers)[6],
					(*mockUsers)[8],
					(*mockUsers)[21],
					(*mockUsers)[25],
				),
				NextPage: false,
			},
			wantErr: false,
		},
		{
			name: "OrderField: Empty orderField (default empty: Name) OrderByDesc",
			req: SearchRequest{
				Limit:      25,
				Offset:     0,
				Query:      "J",
				OrderField: "",
				OrderBy:    OrderByDesc,
			},
			want: &SearchResponse{
				Users: append([]User{},
					(*mockUsers)[8],
					(*mockUsers)[6],
					(*mockUsers)[21],
					(*mockUsers)[25],
				),
				NextPage: false,
			},
			wantErr: false,
		},
		{
			name: "OrderField: Name OrderByDesc",
			req: SearchRequest{
				Limit:      25,
				Offset:     0,
				Query:      "J",
				OrderField: "Name",
				OrderBy:    OrderByDesc,
			},
			want: &SearchResponse{
				Users: append([]User{},
					(*mockUsers)[8],
					(*mockUsers)[6],
					(*mockUsers)[21],
					(*mockUsers)[25],
				),
				NextPage: false,
			},
			wantErr: false,
		},
		{
			name: "OrderField: Id OrderByDesc",
			req: SearchRequest{
				Limit:      25,
				Offset:     0,
				Query:      "J",
				OrderField: "Id",
				OrderBy:    OrderByDesc,
			},
			want: &SearchResponse{
				Users: append([]User{},
					(*mockUsers)[6],
					(*mockUsers)[8],
					(*mockUsers)[21],
					(*mockUsers)[25],
				),
				NextPage: false,
			},
			wantErr: false,
		},
		{
			name: "OrderField: Age OrderByDesc",
			req: SearchRequest{
				Limit:      25,
				Offset:     0,
				Query:      "J",
				OrderField: "Age",
				OrderBy:    OrderByDesc,
			},
			want: &SearchResponse{
				Users: append([]User{},
					(*mockUsers)[21],
					(*mockUsers)[8],
					(*mockUsers)[25],
					(*mockUsers)[6],
				),
				NextPage: false,
			},
			wantErr: false,
		},
		{
			name: "OrderBy: Name OrderByAsc",
			req: SearchRequest{
				Limit:      25,
				Offset:     0,
				Query:      "J",
				OrderField: "Name",
				OrderBy:    OrderByAsc,
			},
			want: &SearchResponse{
				Users: append([]User{},
					(*mockUsers)[25],
					(*mockUsers)[21],
					(*mockUsers)[6],
					(*mockUsers)[8],
				),
				NextPage: false,
			},
			wantErr: false,
		},
		{
			name: "OrderBy: Invalid OrderBy",
			req: SearchRequest{
				Limit:      25,
				Offset:     0,
				Query:      "J",
				OrderField: "Name",
				OrderBy:    -2,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testSrv.FindUsers(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindUsers() \nerror = %v, \nwantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindUsers() \ngot  = %v, \nwant = %v", got, tt.want)
			}
		})
	}
}

func TestSearchClient_FindUsersError(t *testing.T) {
	srvError := SearchClient{
		AccessToken: "your token here",
		URL:         errorServer.URL,
	}
	tests := []struct {
		name    string
		req     SearchRequest
		want    *SearchResponse
		wantErr bool
		err     error
	}{
		//{
		//	name:    "Error: timeout",
		//	req:     SearchRequest{},
		//	want:    &SearchResponse{},
		//	wantErr: true,
		//},
		{
			name: "Error: timeout",
			req: SearchRequest{
				Limit:      25,
				Offset:     0,
				Query:      "timeout_test",
				OrderField: "",
				OrderBy:    0,
			},
			want:    nil,
			wantErr: true,
			// limit 26 в ответе, так как запрашивает на 1 больше для пагинации
			err: fmt.Errorf("timeout for limit=26&offset=0&order_by=0&order_field=&query=timeout_test"),
		},
		{
			name: "Error: token_test",
			req: SearchRequest{
				Limit:      25,
				Offset:     0,
				Query:      "token_test",
				OrderField: "",
				OrderBy:    0,
			},
			want:    nil,
			wantErr: true,
			err:     fmt.Errorf("Bad AccessToken"),
		},
		{
			// TODO: Попробовать вызвать internal без имитации
			name: "Error: internal_error_test",
			req: SearchRequest{
				Limit:      25,
				Offset:     0,
				Query:      "internal_error_test",
				OrderField: "",
				OrderBy:    0,
			},
			want:    nil,
			wantErr: true,
			err:     fmt.Errorf("SearchServer fatal error"),
		},
		{
			name: "Error: json_invalid_test",
			req: SearchRequest{
				Limit:      25,
				Offset:     0,
				Query:      "json_invalid_test",
				OrderField: "",
				OrderBy:    0,
			},
			want:    nil,
			wantErr: true,
			err:     fmt.Errorf("cant unpack error json: invalid character ':' after top-level value"),
		},
		{
			name: "Error: invalid user json",
			req: SearchRequest{
				Limit:      25,
				Offset:     0,
				Query:      "invalid_response_json",
				OrderField: "",
				OrderBy:    0,
			},
			want:    nil,
			wantErr: true,
			err:     fmt.Errorf("cant unpack result json: json: cannot unmarshal object into Go value of type []main.User"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := srvError.FindUsers(tt.req)
			if (err != nil) != tt.wantErr {
				t.Log("err: ", err)
				t.Log("tt.err: ", tt.err)
				t.Errorf("FindUsers() \nerror = %v, \nwantErr %v", err, tt.wantErr)
				return
			}
			if err.Error() != tt.err.Error() {
				t.Errorf("err == tt.err: FindUsers() \ngot = %v\nwant = %v", err, tt.err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindUsers() \ngot  = %v, \nwant = %v", got, tt.want)
			}
		})
	}
}

func TestSearchClient_FindUsers_invalidUrl(t *testing.T) {
	srvError := SearchClient{
		AccessToken: "your token here",
		URL:         "http://invalidurl.comcom",
	}
	client.Transport = &mockTransport{}
	tests := []struct {
		name    string
		req     SearchRequest
		want    *SearchResponse
		wantErr bool
		err     error
	}{
		{
			name: "Error: invalid url",
			req: SearchRequest{
				Limit:      25,
				Offset:     0,
				Query:      "panic",
				OrderField: "",
				OrderBy:    0,
			},
			want:    nil,
			wantErr: true,
			err:     fmt.Errorf("unknown error Get \"http://invalidurl.comcom?limit=26&offset=0&order_by=0&order_field=&query=panic\": custom error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := srvError.FindUsers(tt.req)
			if (err != nil) != tt.wantErr {
				t.Log("err: ", err)
				t.Log("tt.err: ", tt.err)
				t.Errorf("FindUsers() \nerror = %v, \nwantErr %v", err, tt.wantErr)
				return
			}
			if err.Error() != tt.err.Error() {
				t.Errorf("err == tt.err: FindUsers() \ngot = %v\nwant = %v", err, tt.err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindUsers() \ngot  = %v, \nwant = %v", got, tt.want)
			}
		})
	}
}
