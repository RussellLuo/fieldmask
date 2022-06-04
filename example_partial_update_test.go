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

func (a *Address2) Update(other Address2, fm fieldmask.FieldMask) {
	if len(fm) == 0 {
		// Clear the entire address.
		*a = other
		return
	}

	if fm.Has("country") {
		a.Country = other.Country
	}
	if fm.Has("province") {
		a.Province = other.Province
	}
	if fm.Has("city") {
		a.City = other.City
	}
}

type Person2 struct {
	Name    string   `json:"name,omitempty"`
	Age     int      `json:"age,omitempty"`
	Address Address2 `json:"address,omitempty"`
}

func (p *Person2) Update(other Person2, fm fieldmask.FieldMask) {
	if len(fm) == 0 {
		// Clear the entire person.
		*p = other
		return
	}

	if fm.Has("name") {
		p.Name = other.Name
	}
	if fm.Has("age") {
		p.Age = other.Age
	}
	if addressFM, ok := fm.FieldMask("address"); ok {
		p.Address.Update(other.Address, addressFM)
	}
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

	person.Update(req.Person2, req.FieldMask)
	fmt.Printf("updated: %#v\n", person)

	// Output:
	// initial: fieldmask_test.Person2{Name:"foo", Age:20, Address:fieldmask_test.Address2{Country:"X", Province:"Y", City:"Z"}}
	// updated: fieldmask_test.Person2{Name:"foo", Age:10, Address:fieldmask_test.Address2{Country:"X", Province:"Y", City:"ZZ"}}
}
