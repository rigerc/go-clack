# `go-clack/prompts`

Effortlessly build beautiful command-line apps 🪄 [Try the demo](https://stackblitz.com/edit/clack-prompts?file=index.js)

![clack-prompt](https://github.com/orochaa/go-clack/blob/master/.github/assets/clack-demo.gif)

---

`go-clack/prompts` is an opinionated, pre-styled wrapper around [`go-clack/core`](https://github.com/orochaa/go-clack/blob/master/core/README.md).

- 💎 Beautiful, minimal UI
- ✅ Simple API
- 🧱 Comes with `Text`, `Confirm`, `Select`, `MultiSelect`, `Spinner`, and more.

## Get Started

To start using `go-clack/prompts`, follow these steps:

### 1. Install the Package

Add the `go-clack/prompts` package to your Go project:

```bash
go get github.com/orochaa/go-clack/prompts
```

### 2. Create a Prompt

To create and run a simple text prompt, you can use the following code:

```go
// main.go
package main

import (
  "github.com/orochaa/go-clack/prompts"
)

func main() {
  name, err :=  prompts.Text(prompts.TextParams{
    Message: "What's your name?",
  })
  prompts.ExitOnError(err)

  fmt.Printf("Hello, %s!\n", name)
}
```

### 3. Run Your Application

Compile and run your application:

```bash
go run main.go
```

## Basics

### Setup

The `Intro` and `Outro` functions will print a message to begin or end a prompt session, respectively.

```go
prompts.Intro("create-my-app")
// Do stuff
prompts.Outro("You're all set!")
```

### Cancellation

An `error` is returned when a user cancels a prompt with `CTRL + C`.

You should handle this situation for each prompt, optionally providing a nice cancellation message with the `Cancel` utility.

```go
value, err := prompts.Text(/* TODO */)
if (err != nil) {
  prompts.Cancel("Operation cancelled.")
  os.Exit(0)
}

// Do stuff with `value`
```

Or just exiting using the `ExitOnError` utility.

```go
value, err := prompts.Text(/* TODO */)
prompts.ExitOnError(err)

// Do stuff with `value`
```

## Components

### Text

The `Text` component accepts a single line of text.

```go
meaning, err := prompts.Text(prompts.TextParams{
  Message: "What is the meaning of life?",
  Placeholder: "Not sure",
  InitialValue: "42",
  Validate: func(value string) error {
    if value.length === 0 {
      return errors.New("Value is required!")
    }
    return nil
  },
})
```

### Password

The `Password` component accepts a password input, masking the characters.

```go
password, err := prompts.Password(prompts.PasswordParams{
  Message: "Enter your password:",
  Validate: func(value string) error {
    if value.length < 8 {
      return errors.New("Password must be at least 8 characters long!")
    }
    return nil
  },
})
```

### Path

The `Path` component accepts a file or directory path.

```go
path, err := prompts.Path(prompts.PathParams{
  Message: "Enter the file path:",
  Validate: func(value string) error {
    if !fileExists(value) {
      return errors.New("File does not exist!")
    }
    return nil
  },
})
```

### Confirm

The `Confirm` component accepts a yes or no answer. The result is a boolean value of `true` or `false`.

```go
confirmed, err := prompts.Confirm(prompts.ConfirmParams{
  Message: "Are you sure?",
  InitialValue: true,
})
```

### Select

The `Select` component allows the user to choose a single option from a list.

```go
project, err := prompts.Select(prompts.SelectParams{
  Message: "Pick a project type:",
  Options: []*prompts.SelectOption[string]{
    {Label: "TypeScript", Value: "ts"},
    {Label: "JavaScript", Value: "js"},
    {Label: "CoffeeScript", Value: "coffee", Hint: "oh no"},
  },
})
```

### Reorder

The `Reorder` form lets users arrange a list without losing any items. Arrow keys navigate; press **Space** to pick up or drop an item, and **Enter** to confirm the final order.

```go
order, err := prompts.Reorder(prompts.ReorderParams[string]{
  Message: "Order the deployment steps:",
  Options: []*prompts.ReorderOption[string]{
    {Label: "Install", Value: "install"},
    {Label: "Build", Value: "build"},
    {Label: "Deploy", Value: "deploy"},
  },
})
prompts.ExitOnError(err)
// order contains all three values in the order chosen by the user
```

### Banner

`Banner` prints FIGlet ASCII art using the renderer in `prompts/banner`. It supports fonts, colors, custom gradients, random presets, and plain or ANSI output.

```go
import "github.com/orochaa/go-clack/prompts/banner"

err := prompts.Banner("CLACK", prompts.BannerOptions{
  Config: banner.Config{
    Font: "standard",
    Gradient: &banner.GradientOcean,
  },
})
prompts.ExitOnError(err)
```

### Persistent components above forms

Use `SetPersistent` to keep a banner or other component above successive interactive forms. When the next form starts, it replaces the previous form while the persistent header stays in place. This mode supports one active form at a time and components are rendered in order.

```go
header, err := prompts.BannerComponent(banner.Config{
  Text: "CLACK",
  Font: "standard",
  Gradient: &banner.GradientOcean,
})
prompts.ExitOnError(err)

prompts.SetPersistent(
  header,
  prompts.StaticComponent("Project setup"),
)
defer prompts.ClearPersistent()

name, err := prompts.Text(prompts.TextParams{Message: "Project name?"})
prompts.ExitOnError(err)
_, err = prompts.Confirm(prompts.ConfirmParams{Message: "Continue?"})
prompts.ExitOnError(err)
_ = name
```

Persistent components may also be functions that return fresh text. Call `prompts.RefreshPersistent()` to redraw the active form with updated content. `ClearPersistent()` removes the registered header for future renders; it does not erase terminal history.

### MultiSelect

The `MultiSelect`component allows the user to choose multiple options from a list.

```go
additionalTools, err := prompts.MultiSelect(prompts.MultiSelectParams[string]{
  Message: "Select additional tools:",
  Options: []*prompts.MultiSelectOption[string]{
    {Value: "eslint", Label: "ESLint", Hint: "recommended"},
    {Value: "prettier", Label: "Prettier"},
    {Value: "gh-action", Label: "GitHub Action"},
  },
})
```

### GroupMultiSelect

The `GroupMultiSelect` component allows the user to choose multiple options from grouped lists.

```go
groupChoices, err := prompts.GroupMultiSelect(prompts.GroupMultiSelectParams[string]{
  Message: "Select additional tools:",
  Options: map[string][]prompts.MultiSelectOption[string]{
    "Group 1": {
      {Label: "Option 1", Value: "1"},
      {Label: "Option 2", Value: "2"},
    },
    "Group 2": {
      {Label: "Option 3", Value: "3"},
      {Label: "Option 4", Value: "4"},
    },
  },
})
```

### SelectPath

The `SelectPath` component allows the user to select a file/folder on a tree based select with free navigation by arrow keys.

```go
selectedPath, err := prompts.SelectPath(prompts.SelectPathParams{
  Message: "Select a path:",
})
```

### MultiSelectPath

The `MultiSelectPath` component allows the user to select multiple files/folders on a tree based select with free navigation by arrow keys.

```go
selectedPaths, err := prompts.MultiSelectPath(prompts.MultiSelectPathParams{
  Message: "Select paths",
})
```

### SelectKey

The `SelectKey` component allows the user to choose a option associated to a key.

```go
selectedKey, err := prompts.SelectKey(prompts.SelectKeyParams{
  Message: "Select a key",
  Options: []prompts.SelectKeyOption[string]{
    {Key: "f", Label: "Foo"},
    {Key: "b", Label: "Bar"},
    {Key: "Enter", Label: "Baz"},
  },
})
```

### Spinner

The spinner component surfaces a pending action, such as a long-running download or dependency installation.

```go
s := prompts.Spinner()
s.Start("Installing via npm")
// Do installation here
s.Stop("Installed via npm")
```

## Utilities

### Tasks

Execute multiple tasks in spinners.

```go
prompts.Tasks(
  []prompts.Task{
    {
      Title: "Installing via npm",
      Task: func(message func(msg string)) (string, error) {
        // Do installation here
        return "Installed via npm", nil
      }
    },
  },
  prompts.SpinnerOptions{}
)
```

### Logs

```go
prompts.Info("Info!")
prompts.Success("Success!")
prompts.Step("Step!")
prompts.Warn("Warn!")
prompts.Error("Error!")
prompts.Message("Hello, World", prompts.MessageOptions{
  FirstLine: prompts.MessageLineOptions{
    Start: picocolors.Cyan("~"),
  },
  Default: prompts.MessageLineOptions{
    Start: picocolors.Gray(symbols.BAR),
  },
})
```

![clack-log-prompts](https://github.com/orochaa/go-clack/blob/master/.github/assets/clack-logs.png)
