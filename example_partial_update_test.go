package fieldmask_test

import (
	"encoding/json"
	"fmt"

	"github.com/RussellLuo/fieldmask"
	"github.com/mitchellh/mapstructure"
)

type Address2 struct {
	Country  string `json:"country,omitempty"`
	Province string `json:"province,omitempty"`
	City     string `json:"city,omitempty"`
}

type Person2 struct {
	Name    string   `json:"name,omitempty"`
	Age     int      `json:"age,omitempty"`
	Address Address2 `json:"address,omitempty"`
}

type UpdatePersonRequest struct {
	Person2
	FieldMask fieldmask.FieldMask `json:"-"`
}

func (req *UpdatePersonRequest) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &req.FieldMask); err != nil {
		return err
	}
	return mapstructure.Decode(req.FieldMask, &req.Person2)
}

func Example_partialUpdate() {
	person := Person2{
		Name: "foo",
		Age:  20,
		Address: Address2{
			Country:  "X",
			Province: "Y",
			City:     "Z",
		},
	}
	fmt.Printf("initial: %#v\n", person)

	blob := []byte(`{"age": 10, "address": {"city": "ZZ"}}`)
	req := new(UpdatePersonRequest)
	if err := json.Unmarshal(blob, req); err != nil {
		fmt.Printf("err: %#v\n", err)
	}

	// Update name if needed.
	if req.FieldMask.Has("name") {
		person.Name = req.Name
	}

	// Update age if needed.
	if req.FieldMask.Has("age") {
		person.Age = req.Age
	}

	// Update address if needed.
	if req.FieldMask.Has("address") {
		fm, _ := req.FieldMask.FieldMask("address")
		if len(fm) == 0 {
			// Clear the entire address.
			person.Address = req.Address
			fmt.Printf("updated: %#v\n", person)
			return
		}

		if fm.Has("country") {
			person.Address.Country = req.Address.Country
		}
		if fm.Has("province") {
			person.Address.Province = req.Address.Province
		}
		if fm.Has("city") {
			person.Address.City = req.Address.City
		}
		fmt.Printf("updated: %#v\n", person)
	}

	// Output:
	// initial: fieldmask_test.Person2{Name:"foo", Age:20, Address:fieldmask_test.Address2{Country:"X", Province:"Y", City:"Z"}}
	// updated: fieldmask_test.Person2{Name:"foo", Age:10, Address:fieldmask_test.Address2{Country:"X", Province:"Y", City:"ZZ"}}
}
