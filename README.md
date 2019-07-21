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
         * [Blackhole](#blackhole)
         * [Text](#text)
         * [Number](#number)
            * [Mathematical Operations](#mathematical-operations)
      * [Functions](#functions)
         * [Arguments](#arguments)
         * [Questions](#questions)
      * [Controlling Flow](#controlling-flow)
         * [Conditions](#conditions)
         * [Decisions (if/unless)](#decisions-ifunless)
         * [Loops (while/until)](#loops-whileuntil)
   * [Backends](#backends)
      * [System](#system)
   * [Examples](#examples)
      * [Hello, World!](#hello-world)
      * [Variables](#variables-1)
      * [Functions (Custom Sentences)](#functions-custom-sentences)

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

Long sentences can be broken up into multiple lines using `...` at the end of
each line, excluding the last line:

```
this is a really long...
	sentence that should go...
	over multiple lines
```

Indentation between lines does not make an difference. However, it is easier to
read when following lines are indented.

## Comments

Comments start with a `#` and continue until a new line or the end of the
program is reached. The comment may be on its own line or at the end of a
sentence:

```
# A comment looks like this.
display "Hello" # It can also be here
```

## Variables

```
declare first-name is text
declare counter is a number
```

1. Only `text` and `number` is supported. See specific documentation below.
2. The word `a` or `an` may appear before the type. This can make it easier to
read: "is a number" rather than "is number". However, the "a" or "an" does not
have any affect on the program.
3. All variables for a function must be declare before any other sentences.
4. The same variable name cannot be defined twice within the same function.
However, the same variable name can appear in different functions. These are
unrelated to each other.
5. You cannot declare a variable with the same name as one of the function
parameters.
6. All types have a default value which is safe to use before it is set to
another value.
7. There is a special variable called `_` which is called the blackhole.
Explained in more detail below.

Variables can be set with:

```
set first-name to "Bob"
```

### Blackhole

The blackhole variable is a single underscore (`_`). It can be used as a
placeholder when the value can be ignored. For example:

```bento
divide 1.23 by 7.89 into _
```

If `_` is used in a place where the value would be read it will return a zero
value (the same as the default value for that type).

The following lines would function the same. However, you should not rely on the
blackhole variable as a readable value and it may be made illegal in the future:

```bento
add _ and 7.89 into result
add 0 and 7.89 into result
```

### Text

```bento
my-variable is text
```

1. A `text` variable can contain any text, including being empty (zero
characters).
2. It's perfectly safe to use text variables before they have been given a
value, the default value will be empty.

### Number

```bento
my-variable is number
my-variable is number with 1 decimal place
my-variable is number with 3 decimal places
```

1. A number variable is exact and has a maximum number of decimal places (this
is also called the precision).
2. If the number of decimal places is not specified it will use 6.
3. For integers you should use `number with 0 decimal places`.
4. The number of decimal places cannot be negative.
5. A number has no practical minimum (negative) or maximum (positive) value. You
can process incredibly large numbers with absolute precision.
6. Any calculated value will be rounded at the end of the operation so that it
never contains more precision than what is allowed. For example if the number
has one decimal place, `5.5 * 6.5 * 11` evaluates to `393.8` because
`5.5 * 6.5 = 35.75 => 35.8`, `35.8 * 11 = 393.8`.
7. Numbers are always displayed without trailing zeroes after the decimal point.
For example, `12.3100` is displayed as `12.31` as long as the number of decimal
places is at least 2.
8. The words `places` and `place` mean the same thing. However, it is easier to
read when `place` is reserved for when there is only one decimal place.
9. The default value of a `number` is `0`. This is safe to use use before it has
been set.

#### Mathematical Operations

```bento
add a and b into c          # c = a + b
subtract a from b into c    # c = b - c
multiply a and b into c     # c = a * b
divide a and b into c       # c = a / b
```

Note: Be careful with `subtract` as the operands are in the reverse order of the
others.

## Functions

Functions (custom sentences) can be defined by using the `:` character:

```
print everything:
	display "Hello"
	display "World"
```

The whitespace is not required. However, it is easier to read when content of
functions are indented with spaces or tabs.

### Arguments

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

### Questions

A question is a special type of function that is defined with a `?` instead of a
`:`:

```bento
it is ok?
	yes
```

A question is answered with the `yes` or `no` sentences. Once a question is
answered it will return immediately.

If a question is not explicitly answered by the end, it's assumed to be `no`.

Questions can be asked in conditionals:

```bento
start:
	if it is ok, display "All good!"
```

Questions can also take arguments in the same way that functions do:

```bento
start:
	declare x is number

	set x to 123
    if x is over 100, display "It's over 100", otherwise display "Not yet"

the-number is over threshold (the-number is number, threshold is number)?
	if the-number > threshold, yes
```

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
