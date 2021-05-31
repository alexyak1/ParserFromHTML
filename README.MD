# Parser for table from web site
url for parsed web site: https://candystore.zimpler.net/
## Result will be data from table, sorted by total eated amount of candy and show favirite candy.

# Expected result
```json
[
    {
        name: "Jonas",
        favouriteSnack: "Geisha",
        totalSnacks: 1982,
    },
    {
        name: "Annika",
        favouriteSnack: "Geisha",
        totalSnacks: 208,
    },
    {
        name: "Jane",
        favouriteSnack: "Nötchoklad",
        totalSnacks: 22,
    },
    {
        name: "Aadya",
        favouriteSnack: "Center",
        totalSnacks: 11,
    },

]
```

# How to run

* Clone repository: `git clone https://github.com/alexyak1/ZimplerParser.git`
* Run command `go run main.go` <br>
    or build project `go build` and `./ZimplerParser`  
* open http://127.0.0.1:10000/ in a browser 
