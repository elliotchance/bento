# 🍱 bento

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
         * [Decisions](#decisions)
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

### Decisions

Sentences starting with `if` can be used to control the flow. An `if` sentence
takes on of the following forms:

```
if <condition>, <true>

if <condition>, <true>, otherwise <false>
```

Where `<condition>` is a simple comparison between two variables or values. Some
examples are:

```
name = "Bob"
counter > 10
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
