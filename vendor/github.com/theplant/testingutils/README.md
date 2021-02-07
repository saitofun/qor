




* [Pretty Json Diff](#pretty-json-diff)
* [Println Json](#println-json)




## Pretty Json Diff
``` go
func PrettyJsonDiff(expected interface{}, actual interface{}) (r string)
```
It convert the two objects into pretty json, and diff them, output the result.


```go
	type Company struct {
	    Name string
	}
	type People struct {
	    Name    string
	    Age     int
	    Company Company
	}
	
	p1 := People{
	    Name: "Felix",
	    Age:  20,
	    Company: Company{
	        Name: "The Plant",
	    },
	}
	p2 := People{
	    Name: "Tom",
	    Age:  21,
	    Company: Company{
	        Name: "Microsoft",
	    },
	}
	
	fmt.Println(PrettyJsonDiff(p1, p2))
	//Output:
	// --- Expected
	// +++ Actual
	// @@ -1,7 +1,7 @@
	//  {
	// -	"Name": "Felix",
	// -	"Age": 20,
	// +	"Name": "Tom",
	// +	"Age": 21,
	//  	"Company": {
	// -		"Name": "The Plant"
	// +		"Name": "Microsoft"
	//  	}
	//  }
```

## Println Json
``` go
func PrintlnJson(vals ...interface{})
```





