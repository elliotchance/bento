# 🍱 bento

   * [Getting Started](#getting-started)
      * [Installation](#installation)
      * [Running A Program](#running-a-program)
   * [Example Use Case](#example-use-case)
   * [Language](#language)
      * [File Structure](#file-structure)
      * [Sentences](#sentences)
      * [Comments](#comments)
      * [Variables](#variables)
         * [Text](#text)
         * [Number](#number)
      * [Functions (Custom Sentences)](#functions-custom-sentences)
      * [Controlling Flow](#controlling-flow)
         * [Conditions](#conditions)
         * [Decisions (if/unless)](#decisions-ifunless)
         * [Loops (while/until)](#loops-whileuntil)
   * [Backends](#backends)
      * [System](#system)
   * [Examples](#examples)
      * [Hello, World!](#hello-world)
      * [Variables](#variables-1)
      * [Functions (Custom Sentences)](#functions-custom-sentences-1)

bento is a
[forth-generation programming language](https://en.wikipedia.org/wiki/Fourth-generation_programming_language)
that is English-based. It is designed to separate orchestration from
implementation as to provide a generic, self-documenting DSL that can be managed
by non-technical individuals.

That's a bunch of fancy talk that means that the developers are able to
setup the complex tasks and make them easily rerunnable by non-technical people.

The project is still very young, but it has a primary goal for each half of its
intended audience:

**For the developer:** Programs can be written in any language and easily
exposed through a set of specific DSLs called "sentences".

**For the user:** An English-based language that avoids any nuances that
non-technical people would find difficult to understand. The language only
contains a handful of special words required for control flow. Whitespace and
capitalization do not matter.

# Getting Started

## Installation

Bento is available for Mac, Windows and Linux. You can download the latest
release from the [Releases](https://github.com/elliotchance/bento/releases)
page.

## Running A Program

Create a file called `hello-world.bento` with the contents:

```bento
start:
	display "Hello, World!"
```

Then run this file with (you may need to replace the path to the file if you are
not in the same directory as the file):

```bash
bento hello-world.bento
```

# Example Use Case

The sales team need to be able to run customer reports against a database. They
do not understand SQL, or the reports may need to be generated from multiple
tools/steps. As the developer, your primary options are:

1. Run the reports for them. This is a waste of engineering time and a blocker
for the sales team if they need to do it frequently, or while you're on holiday.

2. Write up some documentation and give them readonly access to what they need.
The process could be highly technical. Also, allowing them to access production
systems is a bad idea. Also, it's still a waste of time for everyone when things
aren't working as expected.

3. Build a tool that runs the reports for them. This is another rabbit hole of
wasted engineering time, especially when requirements change frequently or they
need a lot of flexibility.

4. Build the reporting queries/tasks once in bento, and provide them with a
script that they can edit and rerun. An example might look like:

```
run sales report for last 7 days for customer 123
if there were sales, email report to "bob@mycompany.com"
```

The developer would implement the following sentences:

```
run sales report for last ? days for customer ?
there were sales
email report to ?
```

# Language

## File Structure

A `.bento` file consists of one or more functions. Each of those functions can
contain zero or more sentences and call other sentences before of after they are
defined.

The file must contain a `start` function. This can be located anywhere in the\
file, but it is recommended to put it as the first function.

## Sentences

A sentence contains a collection of words and values and it is terminated by a
new line. For example:

```
display "Hello"
```

## Comments

Comments start with a `#` and continue until a new line or the end of the
program is reached. The comment may be on its own line or at the end of a
sentence:

```
# A comment looks like this.
display "Hello" # It can also be here
```

## Variables

Variables can be declared with the `declare` keyword:

```
declare first-name is <type>
```

Where `<type>` is `text` or `number`. Each of the these types has a default
value (described below). However, the variable can be set with:

```
set first-name to "Bob"
```

### Text

Text variables can contain any text, including being empty (zero characters).
It's perfectly safe to use text variables before they have been given a value,
the default value will be empty.

### Number

Numbers represent exact numerical figures that have no theoretical size limit,
but have a fixed precision of 6 decimal places.

Basic mathematical operations are:

```bento
add a and b into c          # c = a + b
subtract a from b into c    # c = b - c
multiply a and b into c     # c = a * b
divide a and b into c       # c = a / b
```

## Functions (Custom Sentences)

Custom sentences can be defined by using the `:` character:

```
print everything:
	display "Hello"
	display "World"
```

The whitespace is not required. However, it is easier to read when content of
functions are indented with spaces or tabs.

Variables can be declared in the function name by specifying their names and
types in `()`, for example:

```
say greeting to persons-name (greeting is text, persons-name is text):
	display greeting
	display persons-name
```

Can be called with:

```
say "Hi" to "Bob"
```

The order in which the arguments are defined is not important.

## Controlling Flow

### Conditions

A condition is a simple comparison between two variables or values. Some
examples are:

```
name = "Bob"
counter > 10
first-name != last-name
```

All supported operators are:

- `=` - Equal.
- `!=` - Not equal.
- `>` - Greater than.
- `>=` - Greater than or equal.
- `<` - Less than.
- `<=` - Less than or equal.

Values can only be compared when they are the same type. For example the
following is not allowed, and will return an error:

```
"123" = 123
```

### Decisions (if/unless)

Sentences starting with `if` or `unless` can be used to control the flow. The
sentence takes one of the following forms (each either starting with `if` or
`unless`):

```
if/unless <condition>, <true>

if/unless <condition>, <true>, otherwise <false>
```

When `unless` is used instead of `if` the comparison is inverted, so that:

```
if "Bob" = "Bob"      # true
unless "Bob" = "Bob"  # false
```

### Loops (while/until)

Sentences starting with `while` repeat the sentence until while the condition is
true. That is, the loop will only stop once the condition becomes false.

Conversely, using `until` will repeat the sentence until the condition becomes
true.

Loops are written in one of the following forms:

```
while/until <condition>, <true>

while/until <condition>, <true>, otherwise <false>
```

# Backends

A backend performs the tasks described in the bento program.

## System

The system backend provides direct access to running programs on the host
machine.

- `run system command <command>`: Run the `command` and send all stdout and
stderr to the console.

- `run system command <command> output into <output>`: Run the `command` and
capture all of the stdout and stderr into the `output`.

- `run system command <command> status code into <status>`: Run the `command`
and discard and stdout and stderr. Instead capture the status code returned in
`status`.

- `run system command <command> output into <output> status code into <status>`:
Run the `command` and capture the stdout and stderr into `output` as well as the
status code returned into `status`.

Example:

```bento
start:
	declare echo-result is number
	run system command "echo hello" status code into echo-result
	unless echo-result = 0, display "command failed!"
```

# Examples

## Hello, World!

```
start:
	display "Hello, World!"
```

## Variables

```
start:
	declare first-name is text
	set first-name to "Bob"
	display first-name
```

## Functions (Custom Sentences)

```
start:
	print everything
	say "Hi" to "Bob"

print everything:
	display "Hello"
	display "World"

say greeting to persons-name (persons-name is text, greeting is text):
	display greeting
	display persons-name
```
