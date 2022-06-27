package Modules

import(
    "strings"
)

func get_fields_whitespace(str string) ResponseStruct{
    words := strings.Fields(str)
    var res = ResponseStruct{
        Name: words[0],
        Ttl: words[1],
        Class: words[2],
        Typ: words[3],
        Reply: words[4],
    }
    return res;
}
