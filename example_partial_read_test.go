package fieldmask_test

import (
	"encoding/json"
	"fmt"

	"github.com/RussellLuo/fieldmask"
	"github.com/fatih/structs"
)

type Address1 struct {
	Country  string `json:"country,omitempty"`
	Province string `json:"province,omitempty"`
	City     string `json:"city,omitempty"`
}

type Person1 struct {
	Name    string   `json:"name,omitempty"`
	Age     int      `json:"age,omitempty"`
	Address Address1 `json:"address,omitempty"`
}

type GetPersonResponse struct {
	Person1
	FieldMask []string `json:"-"`
}

func (resp *GetPersonResponse) MarshalJSON() ([]byte, error) {
	s := structs.New(resp.Person1)
	s.TagName = "json"
	m := s.Map()

	fm := fieldmask.From(m, resp.FieldMask...)
	return json.Marshal(fm)
}

func Example_partialRead() {
	person := Person1{
		Name: "foo",
		Age:  20,
		Address: Address1{
			Country:  "X",
			Province: "Y",
			City:     "Z",
		},
	}
	fmt.Printf("original: %#v\n", person)

	resp := &GetPersonResponse{
		Person1:   person,
		FieldMask: []string{"name", "address.city"},
	}
	blob, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("err: %#v\n", err)
	}

	fmt.Printf("marshaled: %s\n", blob)

	// Output:
	// original: fieldmask_test.Person1{Name:"foo", Age:20, Address:fieldmask_test.Address1{Country:"X", Province:"Y", City:"Z"}}
	// marshaled: {"address":{"city":"Z"},"name":"foo"}
}
