<img src="https://avatars2.githubusercontent.com/u/24763891?s=400&u=c1150e7da5667f47159d433d8e49dad99a364f5f&v=4"  width="256px" height="256px" align="right" alt="Multiverse OS Logo">

## Multiverse OS: `radix` trie library
**URL** [multiverse-os.org](https://multiverse-os.org)

A simple prefix sorting radix tree designed specifically for working with words, providing fuzzy search functionality, prefix filtering for autocomplete.

### Examples
A simple example using the `.String()` function provided which prints out the
current state of the trie. 

``` go
package main

import (
	"fmt"

	radix "github.com/multiverse-os/cli/radix"
)

func main() {
	fmt.Println("radix trie example")
	fmt.Println("==================")

	tree := radix.New()
	tree.Add("romane", 112)
	tree.Add("romanus", 11)
	tree.Add("tuber", 44)
	tree.Add("romulus", 4)
	tree.Add("ruber", 8)
	tree.Add("tubular", 99)
	tree.Add("rubens", 9)
	tree.Add("tub", 3)
	tree.Add("rubicon", 44)
	tree.Add("rubicundus", 71)

	fmt.Printf("%s", tree.String())
}
```


Setting and querying flags is simple.

``` go
package main

import (
  "fmt"
  "os"

  cli "github.com/hackwave/cli-framework"
)

func main() {
  cmd := cli.New(&cli.CLI{
  Flags: []cli.Flag {
    cli.StringFlag{
      Name: "lang",
      Value: "english",
      Usage: "language for the greeting",
    },
  },
  Action: func(c *cli.Context) error {
    name := "Nefertiti"
    if c.NArg() > 0 {
      name = c.Args().Get(0)
    }
    if c.String("lang") == "spanish" {
      fmt.Println("Hola", name)
    } else {
      fmt.Println("Hello", name)
    }
    return nil
  })

  cmd.Run(os.Args)
}
```

You can also set a destination variable for a flag, to which the content will be
scanned.

``` go
package main

import (
  "os"
  "fmt"

  cli "github.com/multiverse-os/cli-framework"
)

func main() {
  var language string

  cmd := cli.New(&cli.CLI{
  Flags: []cli.Flag {
    cli.StringFlag{
      Name:        "lang",
      Value:       "english",
      Usage:       "language for the greeting",
      Destination: &language,
    },
  },
  Action: func(c *cli.Context) error {
    name := "someone"
    if c.NArg() > 0 {
      name = c.Args()[0]
    }
    if language == "spanish" {
      fmt.Println("Hola", name)
    } else {
      fmt.Println("Hello", name)
    }
    return nil
  })

  cmd.Run(os.Args)
}
```

A full list of flags can be found within the source code under the file `flag.go`.

#### Placeholder Values

Sometimes it's useful to specify a flag's value within the usage string itself.
Such placeholders are indicated with back quotes.

For example this:

```go
package main

import (
  "os"

  cli "github.com/hackwave/cli"
)

func main() {
  cmd := cli.New(&cli.CLI{
    Flags: []cli.Flag{
    cli.StringFlag{
      Name:  "config, c",
      Usage: "Load configuration from `FILE`",
    },
  })

  cmd.Run(os.Args)
}
```

Will result in help output like:

```
--config FILE, -c FILE   Load configuration from FILE
```

Note that only the first placeholder is used. Subsequent back-quoted words will
be left as-is.

#### Flag Name Aliasing

You can set alternate (or short) names for flags by providing a comma-delimited
list for the `Name`. e.g.

``` go
package main

import (
  "os"

  cli "github.com/hackwave/cli-framework"
)

func main() {
  cmd := cli.New(&cli.CLI{
  Flags: []cli.Flag {
    cli.StringFlag{
      Name: "lang, l",
      Value: "english",
      Usage: "language for the greeting",
    },
  })

  cmd.Run(os.Args)
}
```

That flag can then be set with `--lang spanish` or `-l spanish`. Note that
giving two different forms of the same flag in the same command invocation is an
error.

#### Ordering

Flags for the application and commands are shown in the order they are defined.
However, it's possible to sort them from outside this library by using `FlagsByName`
or `CommandsByName` with `sort`.

For example this:

``` go
package main

import (
  "os"
  "sort"

  cli "github.com/hackwave/cli-framework"
)

func main() {
  cmd := cli.New(&cli.CLI{
  Flags: []cli.Flag {
    cli.StringFlag{
      Name: "lang, l",
      Value: "english",
      Usage: "Language for the greeting",
    },
    cli.StringFlag{
      Name: "config, c",
      Usage: "Load configuration from `FILE`",
    },
  },
  Commands: []cli.Command{
    {
      Name:    "complete",
      Aliases: []string{"c"},
      Usage:   "complete a task on the list",
      Action:  func(c *cli.Context) error {
        return nil
      },
    },
    {
      Name:    "add",
      Aliases: []string{"a"},
      Usage:   "add a task to the list",
      Action:  func(c *cli.Context) error {
        return nil
      },
    },
  })

  sort.Sort(cli.FlagsByName(cmd.Flags))
  sort.Sort(cli.CommandsByName(cmd.Commands))

  cmd.Run(os.Args)
}
```

Will result in help output like:

```
--config FILE, -c FILE  Load configuration from FILE
--lang value, -l value  Language for the greeting (default: "english")
```

#### Values from the Environment

You can also have the default value set from the environment via `EnvVar`.  e.g.

``` go
package main

import (
  "os"

  cli "github.com/hackwave/cli-framework"
)

func main() {
  cmd := cli.New(&cli.CLI{
  Flags: []cli.Flag {
    cli.StringFlag{
      Name: "lang, l",
      Value: "english",
      Usage: "language for the greeting",
      EnvVar: "APP_LANG",
    },
  })

  cmd.Run(os.Args)
}
```

The `EnvVar` may also be given as a comma-delimited "cascade", where the first
environment variable that resolves is used as the default.

``` go
package main

import (
  "os"

  cli "github.com/multiverse-os/cli-framework"
)

func main() {
  cmd := cli.New(&cli.CLI{
  Flags: []cli.Flag {
    cli.StringFlag{
      Name: "lang, l",
      Value: "english",
      Usage: "language for the greeting",
      EnvVar: "LEGACY_COMPAT_LANG,APP_LANG,LANG",
    },
  })

  cmd.Run(os.Args)
}
```

#### Order of operations

The order of operations to assign flag value is (highest to lowest):

1. Command line flag value from user
2. Environment variable (if specified)
3. Configuration file (if specified)
4. Default defined on the flag

### Subcommands

Subcommands can be defined for a more git-like command line app.

```go
package main

import (
  "fmt"
  "os"

  cli "github.com/multiverse-os/cli-framework"
)

func main() {
  cmd := cli.New(&cli.CLI{
  Commands: []cli.Command{
    {
      Name:    "add",
      Aliases: []string{"a"},
      Usage:   "add a task to the list",
      Action:  func(c *cli.Context) error {
        fmt.Println("added task: ", c.Args().First())
        return nil
      },
    },
    {
      Name:    "complete",
      Aliases: []string{"c"},
      Usage:   "complete a task on the list",
      Action:  func(c *cli.Context) error {
        fmt.Println("completed task: ", c.Args().First())
        return nil
      },
    },
    {
      Name:        "template",
      Aliases:     []string{"t"},
      Usage:       "options for task templates",
      Subcommands: []cli.Command{
        {
          Name:  "add",
          Usage: "add a new template",
          Action: func(c *cli.Context) error {
            fmt.Println("new task template: ", c.Args().First())
            return nil
          },
        },
        {
          Name:  "remove",
          Usage: "remove an existing template",
          Action: func(c *cli.Context) error {
            fmt.Println("removed task template: ", c.Args().First())
            return nil
          },
        },
      },
    },
  })

  cmd.Run(os.Args)
}
```

### Subcommands categories

For additional organization in apps that have many subcommands, you can
associate a category for each command to group them together in the help
output.

E.g.

```go
package main

import (
  "os"

  cli "github.com/multiverse-os/cli-framework"
)

func main() {
  cmd := cli.New(&cli.CLI{
  Commands: []cli.Command{
    {
      Name: "noop",
    },
    {
      Name:     "add",
      Category: "Template actions",
    },
    {
      Name:     "remove",
      Category: "Template actions",
    },
  })

  cmd.Run(os.Args)
}
```

Will include:

```
COMMANDS:
    noop

  Template actions:
    add
    remove
```

### Exit code

Calling `App.Run` will not automatically call `os.Exit`, which means that by
default the exit code will "fall through" to being `0`.  An explicit exit code
may be set by returning a non-nil error that fulfills `cli.ExitCoder`, *or* a
`cli.MultiError` that includes an error that fulfills `cli.ExitCoder`, e.g.:

``` go
package main

import (
  "os"

  cli "github.com/multiverse-os/cli-framework"
)

func main() {
  cmd := cli.New(&cli.CLI{
  Flags: []cli.Flag{
    cli.BoolTFlag{
      Name:  "ginger-crouton",
      Usage: "is it in the soup?",
    },
  }
  Action: func(ctx *cli.Context) error {
    if !ctx.Bool("ginger-crouton") {
      return cli.NewExitError("it is not in the soup", 86)
    }
    return nil
  })

  cmd.Run(os.Args)
}
```

### Bash Completion
Bash completion is enabled by setting the `boolean` option `BashCompletion` to `true`.

``` go
package main

import (
  "fmt"
  "os"

  cli "github.com/multiverse-os/cli-framework"
)

func main() {
  tasks := []string{"cook", "clean", "laundry", "eat", "sleep", "code"}

  cmd := cli.New(&cli.CLI{
    BashCompletion: true,
    Commands: []cli.Command{
    {
      Name:  "complete",
      Aliases: []string{"c"},
      Usage: "complete a task on the list",
      Action: func(c *cli.Context) error {
         fmt.Println("completed task: ", c.Args().First())
         return nil
      },
      BashComplete: func(c *cli.Context) {
        // This will complete if no args are passed
        if c.NArg() > 0 {
          return
        }
        for _, t := range tasks {
          fmt.Println(t)
        }
      },
    },
  })

  cmd.Run(os.Args)
}
```

Source the `autocomplete/bash_autocomplete` file in your `.bashrc` file while
setting the `PROG` variable to the name of your program:

`PROG=myprogram source /.../cli/autocomplete/bash_autocomplete`

#### Distribution

Copy `autocomplete/bash_autocomplete` into `/etc/bash_completion.d/` and rename
it to the name of the program you wish to add autocomplete support for (or
automatically install it there if you are distributing a package). Don't forget
to source the file to make it active in the current shell.

```
sudo cp src/bash_autocomplete /etc/bash_completion.d/<myprogram>
source /etc/bash_completion.d/<myprogram>
```

Alternatively, you can just document that users should source the generic
`autocomplete/bash_autocomplete` in their bash configuration with `$PROG` set
to the name of their program (as above).

#### Customization

The default bash completion flag (`--generate-bash-completion`) is defined as
`cli.BashCompletionFlag`, and may be redefined if desired, e.g.:

``` go
package main

import (
  "os"

  cli "github.com/multiverse-os/cli-framework"
)

func main() {
  cli.BashCompletionFlag = cli.BoolFlag{
    Name:   "compgen",
    Hidden: true,
  }

  cmd := cli.New(&cli.CLI{
    EnableBashCompletion: true,
    Commands: []cli.Command{
      {
       Name: "wat",
      },
  })

  cmd.Run(os.Args)
}
```

### Version Flag

The default version flag (`-v/--version`) is defined as `cli.VersionFlag`, which
is checked by the cli internals in order to print the `App.Version` via
`cli.VersionPrinter` and break execution.

#### Customization

The default flag may be customized to something other than `-v/--version` by
setting `cli.VersionFlag`, e.g.:

<!-- {
  "args": ["&#45;&#45print-version"],
  "output": "partay version 19\\.99\\.0"
} -->
``` go
package main

import (
  "os"

  cli "github.com/multiverse-os/cli-framework"
)

func main() {
  cli.VersionFlag = cli.BoolFlag{
    Name: "print-version, V",
    Usage: "print only the version",
  }

  cmd := cli.New(nil)
  cmd.Name = "partay"
  cmd.Version = Version{
                  Major: 19,
                  Minor: 99,
                  Path:  0,
                }
  cmd.Run(os.Args)
}
```

Alternatively, the version printer at `cli.VersionPrinter` may be overridden, e.g.:

``` go
package main

import (
  "fmt"
  "os"

  cli "github.com/multiverse-os/cli-framework"
)

var (
  Revision = "fafafaf"
)

func main() {
  cli.VersionPrinter = func(c *cli.Context) {
    fmt.Printf("version=%s revision=%s\n", c.App.Version, Revision)
  }

  cmd := cli.New(nil)
  cmd.Name = "partay"
  cmd.Version = Version{
                  Major: 19,
                  Minor: 99,
                  Path:  0,
                }
  cmd.Run(os.Args)
}
```

#### Full API Example

**Notice**: This is a contrived (functioning) example meant strictly for API
demonstration purposes.  Use of one's imagination is encouraged.

``` go
package main

import (
  "errors"
  "flag"
  "fmt"
  "io"
  "io/ioutil"
  "os"
  "time"

  cli "github.com/multiverse-os/cli-framework"
)

func init() {
  cli.AppHelpTemplate += "\nCUSTOMIZED: you bet ur muffins\n"
  cli.CommandHelpTemplate += "\nYMMV\n"
  cli.SubcommandHelpTemplate += "\nor something\n"

  cli.HelpFlag = cli.BoolFlag{Name: "halp"}
  cli.BashCompletionFlag = cli.BoolFlag{Name: "compgen", Hidden: true}
  cli.VersionFlag = cli.BoolFlag{Name: "print-version, V"}

  cli.HelpPrinter = func(w io.Writer, templ string, data interface{}) {
    fmt.Fprintf(w, "best of luck to you\n")
  }
  cli.VersionPrinter = func(c *cli.Context) {
    fmt.Fprintf(c.App.Writer, "version=%s\n", c.App.Version)
  }
  cli.OsExiter = func(c int) {
    fmt.Fprintf(cli.ErrWriter, "refusing to exit %d\n", c)
  }
  cli.ErrWriter = ioutil.Discard
  cli.FlagStringer = func(fl cli.Flag) string {
    return fmt.Sprintf("\t\t%s", fl.GetName())
  }
}

type hexWriter struct{}

func (w *hexWriter) Write(p []byte) (int, error) {
  for _, b := range p {
    fmt.Printf("%x", b)
  }
  fmt.Printf("\n")

  return len(p), nil
}

type genericType struct{
  s string
}

func (g *genericType) Set(value string) error {
  g.s = value
  return nil
}

func (g *genericType) String() string {
  return g.s
}

func main() {
  cmd := cli.New(&cli.CLI{
    Name: "program-cli"
    Version: Version{Major: 0, Minor: 1, Patch: 0},
    HelpName: "contrive",
    Usage: "demonstrate available API",
    UsageText: "contrive - demonstrating the available API",
    ArgsUsage: "[args and such]",
    Commands: []cli.Command{
      cli.Command{
        Name:        "doo",
        Aliases:     []string{"do"},
        Category:    "motion",
        Usage:       "do the doo",
        UsageText:   "doo - does the dooing",
        Description: "no really, there is a lot of dooing to be done",
        ArgsUsage:   "[arrgh]",
        Flags: []cli.Flag{
          cli.BoolFlag{Name: "forever, forevvarr"},
        },
        Subcommands: cli.Commands{
          cli.Command{
            Name:   "wop",
            Action: wopAction,
          },
        },
        SkipFlagParsing: false,
        HideHelp:        false,
        Hidden:          false,
        HelpName:        "doo!",
        BashComplete: func(c *cli.Context) {
          fmt.Fprintf(c.App.Writer, "--better\n")
        },
        Before: func(c *cli.Context) error {
          fmt.Fprintf(c.App.Writer, "brace for impact\n")
          return nil
        },
        After: func(c *cli.Context) error {
          fmt.Fprintf(c.App.Writer, "did we lose anyone?\n")
          return nil
        },
        Action: func(c *cli.Context) error {
          c.Command.FullName()
          c.Command.HasName("wop")
          c.Command.Names()
          c.Command.VisibleFlags()
          fmt.Fprintf(c.App.Writer, "dodododododoodododddooooododododooo\n")
          if c.Bool("forever") {
            c.Command.Run(c)
          }
          return nil
        },
        OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
          fmt.Fprintf(c.App.Writer, "for shame\n")
          return err
        },
      },
    }
    Flags: []cli.Flag{
      cli.BoolFlag{Name: "fancy"},
      cli.BoolTFlag{Name: "fancier"},
      cli.DurationFlag{Name: "howlong, H", Value: time.Second * 3},
      cli.Float64Flag{Name: "howmuch"},
      cli.GenericFlag{Name: "wat", Value: &genericType{}},
      cli.Int64Flag{Name: "longdistance"},
      cli.Int64SliceFlag{Name: "intervals"},
      cli.IntFlag{Name: "distance"},
      cli.IntSliceFlag{Name: "times"},
      cli.StringFlag{Name: "dance-move, d"},
      cli.StringSliceFlag{Name: "names, N"},
      cli.UintFlag{Name: "age"},
      cli.Uint64Flag{Name: "bigage"},
    }
    BashCompletion: true,
    HideHelp: false,
    HideVersion: false,
    BashComplete: func(c *cli.Context) {
      fmt.Fprintf(c.App.Writer, "lipstick\nkiss\nme\nlipstick\nringo\n")
    }
    Before: func(c *cli.Context) error {
      fmt.Fprintf(c.App.Writer, "HEEEERE GOES\n")
      return nil
    }
    After: func(c *cli.Context) error {
      fmt.Fprintf(c.App.Writer, "Phew!\n")
      return nil
    }
    CommandNotFound: func(c *cli.Context, command string) {
      fmt.Fprintf(c.App.Writer, "Thar be no %q here.\n", command)
    }
    OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
      if isSubcommand {
        return err
      }

      fmt.Fprintf(c.App.Writer, "WRONG: %#v\n", err)
      return nil
    }
    Action: func(c *cli.Context) error {
      cli.DefaultAppComplete(c)
      cli.HandleExitCoder(errors.New("not an exit coder, though"))
      cli.ShowAppHelp(c)
      cli.ShowCommandCompletions(c, "nope")
      cli.ShowCommandHelp(c, "also-nope")
      cli.ShowCompletions(c)
      cli.ShowSubcommandHelp(c)
      cli.ShowVersion(c)

      categories := c.App.Categories()
      categories.AddCommand("sounds", cli.Command{
        Name: "bloop",
      })

      for _, category := range c.App.Categories() {
        fmt.Fprintf(c.App.Writer, "%s\n", category.Name)
        fmt.Fprintf(c.App.Writer, "%#v\n", category.Commands)
        fmt.Fprintf(c.App.Writer, "%#v\n", category.VisibleCommands())
      }

      fmt.Printf("%#v\n", c.App.Command("doo"))
      if c.Bool("infinite") {
        c.App.Run([]string{"cmd", "doo", "wop"})
      }

      if c.Bool("forevar") {
        c.App.RunAsSubcommand(c)
      }
      c.App.Setup()
      fmt.Printf("%#v\n", c.App.VisibleCategories())
      fmt.Printf("%#v\n", c.App.VisibleCommands())
      fmt.Printf("%#v\n", c.App.VisibleFlags())

      fmt.Printf("%#v\n", c.Args().First())
      if len(c.Args()) > 0 {
        fmt.Printf("%#v\n", c.Args()[1])
      }
      fmt.Printf("%#v\n", c.Args().Present())
      fmt.Printf("%#v\n", c.Args().Tail())

      set := flag.NewFlagSet("contrive", 0)
      nc := cli.NewContext(c.App, set, c)

      fmt.Printf("%#v\n", nc.Args())
      fmt.Printf("%#v\n", nc.Bool("nope"))
      fmt.Printf("%#v\n", nc.BoolT("nerp"))
      fmt.Printf("%#v\n", nc.Duration("howlong"))
      fmt.Printf("%#v\n", nc.Float64("hay"))
      fmt.Printf("%#v\n", nc.Generic("bloop"))
      fmt.Printf("%#v\n", nc.Int64("bonk"))
      fmt.Printf("%#v\n", nc.Int64Slice("burnks"))
      fmt.Printf("%#v\n", nc.Int("bips"))
      fmt.Printf("%#v\n", nc.IntSlice("blups"))
      fmt.Printf("%#v\n", nc.String("snurt"))
      fmt.Printf("%#v\n", nc.StringSlice("snurkles"))
      fmt.Printf("%#v\n", nc.Uint("flub"))
      fmt.Printf("%#v\n", nc.Uint64("florb"))
      fmt.Printf("%#v\n", nc.GlobalBool("global-nope"))
      fmt.Printf("%#v\n", nc.GlobalBoolT("global-nerp"))
      fmt.Printf("%#v\n", nc.GlobalDuration("global-howlong"))
      fmt.Printf("%#v\n", nc.GlobalFloat64("global-hay"))
      fmt.Printf("%#v\n", nc.GlobalGeneric("global-bloop"))
      fmt.Printf("%#v\n", nc.GlobalInt("global-bips"))
      fmt.Printf("%#v\n", nc.GlobalIntSlice("global-blups"))
      fmt.Printf("%#v\n", nc.GlobalString("global-snurt"))
      fmt.Printf("%#v\n", nc.GlobalStringSlice("global-snurkles"))

      fmt.Printf("%#v\n", nc.FlagNames())
      fmt.Printf("%#v\n", nc.GlobalFlagNames())
      fmt.Printf("%#v\n", nc.GlobalIsSet("wat"))
      fmt.Printf("%#v\n", nc.GlobalSet("wat", "nope"))
      fmt.Printf("%#v\n", nc.NArg())
      fmt.Printf("%#v\n", nc.NumFlags())
      fmt.Printf("%#v\n", nc.Parent())

      nc.Set("wat", "also-nope")

      ec := cli.NewExitError("ohwell", 86)
      fmt.Fprintf(c.App.Writer, "%d", ec.ExitCode())
      fmt.Printf("made it!\n")
      return ec
    }

    if os.Getenv("HEXY") != "" {
      cmd.Writer = &hexWriter{}
      cmd.ErrWriter = &hexWriter{}
    }

    Metadata: map[string]interface{}{
      "layers":     "many",
      "explicable": false,
      "whatever-values": 19.99,
  })

  cmd.Run(os.Args)
}

func wopAction(c *cli.Context) error {
  fmt.Fprintf(c.App.Writer, ":wave: over here, eh\n")
  return nil
}
```

### Combining short Bool options

Traditional use of boolean options using their shortnames look like this:
```
# cmd foobar -s -o
```

Suppose you want users to be able to combine your bool options with their shortname.  This
can be done using the **UseShortOptionHandling** bool in your commands.  Suppose your program
has a two bool flags such as *serve* and *option* with the short options of *-o* and
*-s* respectively. With **UseShortOptionHandling** set to *true*, a user can use a syntax
like:
```
# cmd foobar -so
```

If you enable the **UseShortOptionHandling*, then you must not use any flags that have a single
leading *-* or this will result in failures.  For example, **-option** can no longer be used.  Flags
with two leading dashes (such as **--options**) are still valid.
