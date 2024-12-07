package main

import (
	"reflect"
	"testing"
)

//
//import (
//	"reflect"
//	"testing"
//)
//
//func TestSearchServer(t *testing.T) {
//	type args struct {
//		query      string
//		orderField string
//		orderBy    int
//		limit      int
//		offset     int
//	}
//	tests := []struct {
//		name string
//		args args
//		want []Usr
//	}{
//		{
//			name: "1: Henderson",
//			args: args{
//				query:      "Henderson",
//				orderField: "",
//				orderBy:    0,
//				limit:      0,
//				offset:     0,
//			},
//			want: []Usr{
//				{
//					Id:        10,
//					Name:      "Henderson Maxwell",
//					FirstName: "Henderson",
//					LastName:  "Maxwell",
//					Age:       30,
//					About:     "Ex et excepteur anim in eiusmod. Cupidatat sunt aliquip exercitation velit minim aliqua ad ipsum cillum dolor do sit dolore cillum. Exercitation eu in ex qui voluptate fugiat amet.",
//					Gender:    "male",
//				},
//			},
//		},
//		{
//			name: "2: Velit about",
//			args: args{
//				query:      "Velit",
//				orderField: "",
//				orderBy:    0,
//				limit:      0,
//				offset:     0,
//			},
//			want: []Usr{
//				{
//					Id:        2,
//					Name:      "Brooks Aguilar",
//					FirstName: "Brooks",
//					LastName:  "Aguilar",
//					Age:       25,
//					About:     "Velit ullamco est aliqua voluptate nisi do. Voluptate magna anim qui cillum aliqua sint veniam reprehenderit consectetur enim. Laborum dolore ut eiusmod ipsum ad anim est do tempor culpa ad do tempor. Nulla id aliqua dolore dolore adipisicing.",
//					Gender:    "male",
//				},
//				{
//					Id:        12,
//					Name:      "Cruz Guerrero",
//					FirstName: "Cruz",
//					LastName:  "Guerrero",
//					Age:       36,
//					About:     "Sunt enim ad fugiat minim id esse proident laborum magna magna. Velit anim aliqua nulla laborum consequat veniam reprehenderit enim fugiat ipsum mollit nisi. Nisi do reprehenderit aute sint sit culpa id Lorem proident id tempor. Irure ut ipsum sit non quis aliqua in voluptate magna. Ipsum non aliquip quis incididunt incididunt aute sint. Minim dolor in mollit aute duis consectetur.",
//					Gender:    "male",
//				},
//			},
//		},
//		{
//			name: "3: Velit offset",
//			args: args{
//				query:      "Velit",
//				orderField: "",
//				orderBy:    0,
//				limit:      0,
//				offset:     1,
//			},
//			want: []Usr{
//				{
//					Id:        12,
//					Name:      "Cruz Guerrero",
//					FirstName: "Cruz",
//					LastName:  "Guerrero",
//					Age:       36,
//					About:     "Sunt enim ad fugiat minim id esse proident laborum magna magna. Velit anim aliqua nulla laborum consequat veniam reprehenderit enim fugiat ipsum mollit nisi. Nisi do reprehenderit aute sint sit culpa id Lorem proident id tempor. Irure ut ipsum sit non quis aliqua in voluptate magna. Ipsum non aliquip quis incididunt incididunt aute sint. Minim dolor in mollit aute duis consectetur.",
//					Gender:    "male",
//				},
//			},
//		},
//		{
//			name: "2: Velit limit",
//			args: args{
//				query:      "Velit",
//				orderField: "",
//				orderBy:    0,
//				limit:      1,
//				offset:     0,
//			},
//			want: []Usr{
//				{
//					Id:        2,
//					Name:      "Brooks Aguilar",
//					FirstName: "Brooks",
//					LastName:  "Aguilar",
//					Age:       25,
//					About:     "Velit ullamco est aliqua voluptate nisi do. Voluptate magna anim qui cillum aliqua sint veniam reprehenderit consectetur enim. Laborum dolore ut eiusmod ipsum ad anim est do tempor culpa ad do tempor. Nulla id aliqua dolore dolore adipisicing.",
//					Gender:    "male",
//				},
//			},
//		},
//		{
//			name: "2: Sunt offset limit",
//			args: args{
//				query:      "Velit",
//				orderField: "",
//				orderBy:    0,
//				limit:      1,
//				offset:     1,
//			},
//			want: []Usr{
//				{
//					Id:        12,
//					Name:      "Cruz Guerrero",
//					FirstName: "Cruz",
//					LastName:  "Guerrero",
//					Age:       36,
//					About:     "Sunt enim ad fugiat minim id esse proident laborum magna magna. Velit anim aliqua nulla laborum consequat veniam reprehenderit enim fugiat ipsum mollit nisi. Nisi do reprehenderit aute sint sit culpa id Lorem proident id tempor. Irure ut ipsum sit non quis aliqua in voluptate magna. Ipsum non aliquip quis incididunt incididunt aute sint. Minim dolor in mollit aute duis consectetur.",
//					Gender:    "male",
//				},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := SearchServer(tt.args.query, tt.args.orderField, tt.args.orderBy, tt.args.limit, tt.args.offset); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("SearchServer() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func TestLoadUsers(t *testing.T) {
	type args struct {
		xmlFilename string
	}
	tests := []struct {
		name    string
		args    args
		want    *[]Usr
		wantErr bool
	}{
		{
			name: "TestLoadUsers",
			args: args{
				xmlFilename: "mock_dataset.xml",
			},
			want: &[]Usr{
				{
					Id:        0,
					Name:      "",
					FirstName: "Boyd",
					LastName:  "Wolf",
					Age:       22,
					About:     "Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n",
					Gender:    "male",
				},
				{
					Id:        1,
					Name:      "",
					FirstName: "Hilda",
					LastName:  "Mayer",
					Age:       21,
					About:     "Sit commodo consectetur minim amet ex. Elit aute mollit fugiat labore sint ipsum dolor cupidatat qui reprehenderit. Eu nisi in exercitation culpa sint aliqua nulla nulla proident eu. Nisi reprehenderit anim cupidatat dolor incididunt laboris mollit magna commodo ex. Cupidatat sit id aliqua amet nisi et voluptate voluptate commodo ex eiusmod et nulla velit.\n",
					Gender:    "female",
				},
				{
					Id:        2,
					Name:      "",
					FirstName: "Brooks",
					LastName:  "Aguilar",
					Age:       25,
					About:     "Velit ullamco est aliqua voluptate nisi do. Voluptate magna anim qui cillum aliqua sint veniam reprehenderit consectetur enim. Laborum dolore ut eiusmod ipsum ad anim est do tempor culpa ad do tempor. Nulla id aliqua dolore dolore adipisicing.\n",
					Gender:    "male",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadUsers(tt.args.xmlFilename)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadUsers() got = %v, want %v", got, tt.want)
			}
		})
	}
}
