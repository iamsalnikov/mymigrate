# mymigrate

Yet another migration library for golang.

- [Why](#why)
- [How to use it](#how-to-use-it)
  - [Installation](#installation)
  - [Setup a database connection](#setup-a-database-connection)
  - [Add migrations](#add-migrations)
  - [Apply, Down, View history with direct commands](#apply-down-view-history-with-direct-commands)
  - [Cobra commands](#cobra-commands)

## Why
<img align="left" width="150" height="130" src="https://www.meme-arsenal.com/memes/7d70b26fc3a93cd663768d1c52a445b5.jpg">
Sometimes we need to perform some complex logic during migration, such as receiving data from an external resource, complex data mapping, and so on. 

I've tried to find a tool that will allow me to do these things, but I've failed (maybe because I didn't look for it well enough).
So I've decided to write a simple migration tool that will allow you to write migrations in golang and work with an app's DB connection.

It works with golang's SQL package, so, theoretically, it may work with any SQL DB (I've tested it only with MySQL).

## How to use it

You can use direct functions to work with migrations or add cobra's commands to your app and use them. I'll try to describe all these ways.

### Installation

To install this package you need to run:

```bash
go get github.com/iamsalnikov/mymigrate
```

### Setup a database connection

To work with migrations we need to know a database connection. After opening a connection with DB we need to pass the connection to `mymigrate` package via:

```golang
mymigrate.SetDatabase(db)
```

Example:

```golang
import (
    "database/sql"
    "log"

    _ "github.com/go-sql-driver/mysql"
    "github.com/iamsalnikov/mymigrate"
	"github.com/iamsalnikov/mymigrate/provider/mysql"
)

func main() {
    db, err := sql.Open("mysql", getConnString(name))
    if err != nil {
        log.Fatalln(err)
    }

	provider := mysql.NewMysqlProvider(db)
    mymigrate.SetDatabaseProvider(provider)
```

### Add migrations

To add a new migration to a migration pool we need to call the method `Add` and pass the name of the migration, a function to UP the migration, a function to DOWN the migration. Example:

```golang
mymigrate.Add(
    "mig_001",
    func (db *sql.DB) error {
        // TODO: implemet up logic
        panic("Implement me!")
    },
    func (db *sql.DB) error {
        // TODO: implemet down logic
        panic("Implement me!")
    },
)
```

We can create a package `migrations` inside a project and put all migrations here. And then just import in to the entrypoint of project. Example:

Project structure:

```
- app/
    - migrations/
        - mig_001.go
        - mig_002.go
    main.go
```

Content of `app/migrations/mig_001.go`:

```golang
package migrations

import (
    "database/sql"

    "github.com/iamsalnikov/mymigrate"
)

func init() {
    mymigrate.Add(
        "mig_001",
        func (db *sql.DB) error {
            // TODO: implemet up logic
            panic("Implement me!")
        },
        func (db *sql.DB) error {
            // TODO: implemet down logic
            panic("Implement me!")
        },
    )
}
```

Content of `app/main.go`: 

```golang
import (
    "database/sql"
    "log"

    _ "github.com/go-sql-driver/mysql"
    "github.com/iamsalnikov/mymigrate"

    // Import project migrations
    _ "app/migrations"
)

func main() {
    db, err := sql.Open("mysql", getConnString(name))
    if err != nil {
        log.Fatalln(err)
    }

    mymigrate.SetDatabase(db)
    appliedMigrations, err := mymigrate.Apply()
    // TODO: work with it
)
```

### Apply, Down, View history with direct commands

To Apply migrations with direct command we need to run `mymigrate.Apply()` function. It will return a list of applied migrations and an error.

To Down migrations with direct command we need to run `mymigrate.Down(int)` function and pass number of migrations to be downed. It will return a list of downed migrations and an error.

To view a history of applied migrations with direct command we need to run `mymigrate.History()`. It will return a list of applied migrations and an error.

### Cobra commands

If you use [spf13/cobra](https://github.com/spf13/cobra) package to build nice CLI app then there is one piece of good new for you: this package has commands to work with migrations. You can find them at [cobracmd](cobracmd) directory.

Package `github.com/iamsalnikov/mymigrate/cobracmd` export next commands:
- [ApplyCmd](cobracmd/apply_cmd.go) - command to apply new migrations
- [CreateCmd](cobracmd/create_cmd.go) - command to create new migration
- [DownCmd](cobracmd/down_cmd.go) - command to down applied migrations
- [HistoryCmd](cobracmd/history_cmd.go) - command to view a list of applied migrations
- [NewListCmd](cobracmd/new_cmd.go) - command to view a list of new migrations
- [MigrateCmd](cobracmd/cmd.go) - root command to work with migrations

Also you can find at `github.com/iamsalnikov/mymigrate/cobracmd` functions for cobra commands if you want to configure commands by yourself.
