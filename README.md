# gomemcache
Simple in-memory, non persistent key-value cache utilizing the standard map for microservices in go

# Installation
Install this package simply by requiring it with:

```bash
go get github.com/suchadean/gomemcache
```

# Usage
```go
	// Create a new cache instance
	cache := gomemcache.New()

	// Set a key-value pair
	cache.SetValue("key", []byte("value"), time.Second*30)

	// Get the value by key
	result, err := cache.GetValue(key)
	if err != nil {
	    // handle err
	}

	// Check if a key exists
	exists := cache.KeyExists(key) 

	// Delete a key value pair
	cache.DeleteKey(key)
```


# License
This package is licensed under the [MIT License](https://opensource.org/license/mit)