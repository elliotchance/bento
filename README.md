# 🍱 bento

   * [Example Use Case](#example-use-case)
   * [Language](#language)
      * [Sentences](#sentences)
      * [Comments](#comments)
      * [Variables](#variables)
      * [Functions (Custom Sentences)](#functions-custom-sentences)
   * [Examples](#examples)
      * [Hello, World!](#hello-world)
      * [Variables](#variables-1)

bento is a
[forth-generation programming language](https://en.wikipedia.org/wiki/Fourth-generation_programming_language)
that is English-based. It is designed to separate orchestration from
implementation as to provide a generic, self-documenting DSL that can be managed
by non-technical individuals.

That's a bunch of fancy talk to mean that the developers should be able to setup
the complex tasks that are easily rerunnable by non-technical people.

The project is still very young, but it has a primary goal for each half of its
intended audience:

**For the developer:** Programs can be written in any language and easily
exposed through a set of specific DSLs called "sentences".

**For the user:** An English-based language that avoids any nuances that
non-technical people would find it difficult to understand. The language only
contains a handful of special words required for control flow. Whitespace and
capitalization do not matter. Running incorrect programs gives automatic and
helpful suggestions where possible (see *Example Use Case*).

# Example Use Case

The sales team need to be able to run customer reports against a database. They
do not understand SQL, or the reports may need to be generated from multiple
tools/steps. As the developer, your main options are:

1. Run the reports for them. This is a waste of engineering time and a blocker
for the sales team if they need to do it frequently, or while you're on holiday.

2. Write up some documentation and give them readonly access to what they need.
The process could be highly technical, allow them to access production systems
(bad or just insecure) and is still a waste of time for everyone when things
aren't working.

3. Build a tool that runs the reports for them. This is another rabbit hole of
wasted engineering time, especially when requirements change.

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

Since bento is designed for non-technical people, it provides automatic feedback
when things go wrong (if possible) to help them diagnose the problem. Let's
consider they try to run:

```
run sales report for last month for customer 123
```

This will produce an error, but suggest an alternative:

```
Error: On line 1, I do not understand:

  run sales report for last month for customer 123

Did you mean this instead?

  run sales report for last ? days for customer ?
```

Hopefully, they will be able to adjust the script accordingly and resolve it
immediately. Hooray!

# Language

## Sentences

A sentence contains a collection of words and values and it is terminated by a
new line. For example:

```
display "Hello"
```

## Comments

Comments start with a `#` and continue until a new line or the end of the
program is reached. The comment may be on its own line or at the end of a
sentence.

```
# A comment looks like this.
display "Hello" # It can also be here
```

## Variables

Variables can be declared with the `declare` keyword:

```
declare my-var as text
```

At the moment, only `text` type is supported.

The default value will be empty, you can set the value with:

```
set my-var to "hello"
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

The order in which the variables are defined is not important.

# Examples

## Hello, World!

```
display "Hello, World!"
```

## Variables

```
declare first-name as text
set first-name to "Bob"
display first-name
```
