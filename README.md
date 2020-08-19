# A MXP Go Library

Go MXP is a GoLang driven library for your MXP node. This library has received a grant from the MXP Foundation to ensure it's continuous development through 2020. 

## Installation

Get goMXP 

### Getting A Block

```
package main


func main() {
	gt, err := goMXP.New("http://127.0.0.1:8732")
	if err != nil {
		fmt.Printf("could not connect to network: %v", err)
	}

	block, err := gt.Block(1000)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(block)
}
```

### Getting a Cycle
```
	cycle, err := gt.Cycle(50)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cycle)
```

## Contributing

### The Makefile
The makefile is there as a helper to run quality code checks. To run vet and staticchecks please run: 
```
make checks
```


## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
