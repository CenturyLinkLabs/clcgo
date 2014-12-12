package clcgo

//import (
//"encoding/json"
//"fmt"
//"testing"
//)

//func TestServerJSONUnmarshalling(t *testing.T) {
//template := `{"id": "%s", "name": "%s"}`
//id := "foo"
//name := "bar"
//j := fmt.Sprintf(template, id, name)

//s := Server{}
//err := json.Unmarshal([]byte(j), &s)

//if err != nil {
//t.Errorf("Expected no error, got '%s'", err)
//}

//if s.ID != id {
//t.Errorf("Expected ID to be '%s', was '%s'", id, s.ID)
//}

//if s.Name != name {
//t.Errorf("Expected Name to be '%s', was '%s'", name, s.Name)
//}
//}

//func TestWorkingServerURL(t *testing.T) {
//s := Server{ID: "abc123"}
//u, err := s.URL("AA")
//if err != nil {
//t.Errorf("Expected no error, got '%s'", err)
//}
//if e := APIRoot + "/servers/AA/abc123"; u != e {
//t.Errorf("Expected URL to be '%s', was '%s'", e, u)
//}
//}

//func TestErroredServerURL(t *testing.T) {
//s := Server{}
//u, err := s.URL("AA")
//if err == nil {
//t.Errorf("Expected an error, got nothing")
//} else {
//if e := "The server needs an ID attribute to generate a URL"; err.Error() != e {
//t.Errorf("Expected the error '%s', got nothing", e)
//}
//}
//if u != "" {
//t.Errorf("Expected empty URL, got '%s'", u)
//}
//}
