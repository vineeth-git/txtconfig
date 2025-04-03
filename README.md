A Go (golang) project inspired by godotenv (which loads env vars from a .env file) that loads a config object from a .txt file

There are multiple file formats that we can store config files such as json, yaml, txt and etc. We found that loading .txt files might be very faster since it is very straight forward.
Loading the values from the .txt file could be a very basic thing and we can use this library to do the same.

There's a basic test coverage available in the repo

## Installation  
```shell
go get github.com/vineeth-git/txtconfig
```

## Usage
Add your application configuration to your `.txt` file

eg: config.txt
```
user_name=txtuser
token=dummy
max_timeout=3000
```

Then in your Go app, define the struct that is going to hold these config values
```
package main

type AppConfig struct {
    UserName string
    Token string
    MaxTimeout int64
}
```

Then on your Go file, where you need to load the config

```
appConfig := AppConfig{}
err := txtconfig.Load('/path/to/file.txt', &appConfig)
if err != nil {
    log.Fatal("Error loading Config)
}
fmt.Println(appConfig)
```

## Convensions and how it works
1. Config variable `UserName string` will look for its camel case `user_name` in the `.txt` file.
1. You can also give custom name to look up in the `.txt` file using the field Tags `key` like this

    ```UserName string `key:"permanent_address"` ```
1. You can also define the default value in your config object, in case the config is not available in `.txt` file like this
    ```UserName string `default:"root"` ```
1. You can also mark certain fields required using `required:"true"`. You will get error `Required field is empty: fieldName` if required field is empty. By default, it is false.

    ```UserName string `required:"true"` ```
