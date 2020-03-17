# e-TF1 go-config

Package e-TF1 go-config allow you to load your application config from multiple source (dotenv, env, flags...).  
This package is build on top of [confita](https://github.com/heetch/confita) 

The main difference with confita is that the last loader that matches the config key will override the previous match.

# Quick Usage

Create your own struct with your configuration values

```go
package main

import (
	"context"
	"fmt"

	"github.com/etf1/config"
)

type MyConfig struct {
	Debug    bool   `config:"DEBUG"`
	HTTPAddr string `config:"HTTP_ADDR"`
	BDDPassword string `config:"BDD_PASSWORD" print:"-"` // with print:"-" it will be print as "*** Hidden value ***"
}

func New(context context.Context) *MyConfig {
	cfg := &MyConfig{
		HTTPAddr: ":8001",
        BDDPassword: "my_fake_password",
	}

	config.LoadOrFatal(context, cfg) // It use the DefaultConfigLoader
	return cfg
}

func main() {
	cfg := New(context.Background())
	
	fmt.Println(config.TableString(cfg))
}
```

It will print something similar to

```
-----------------------------------
       Debug|                false|bool `config:"DEBUG"`
    HTTPAddr|                :8001|string `config:"HTTP_ADDR"`
 BDDPassword| *** Hidden value ***|string `config:"BDD_PASSWORD" print:"-"`
-----------------------------------
```

# Loaders 

The library provider provide a DefaultConfigLoader that load from

- .env (if file ./.env was found)
- environment variables
- flags 

You can create your own loader chain

```go
// create your own chain loader
cl := config.NewConfigLoader(
    file.NewBackend("myfile.yml"),
    flags.NewBackend(),
)
```

```go
// you can even append multiple backends to it
cl.AppendBackends(
    file.NewBackend("myfile.json"),
)
```

```go
// and even prepend multiple backends to it
f := ".env"
if _, err := os.Stat(f); err == nil {
    cl.PrependBackends(dotenv.NewBackend(f))
}
```

# Printer

The library provide you a config table printer. A special tag `print:"-"` prevent config value to be printed.

```
-----------------------------------
       Debug|                false|bool `config:"DEBUG"`
    HTTPAddr|                :8001|string `config:"HTTP_ADDR"`
 BDDPassword| *** Hidden value ***|string `config:"BDD_PASSWORD" print:"-"`
-----------------------------------
```