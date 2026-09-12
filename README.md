# Generic Functional Programming Language
A functional programming language with as few types as (realistically) possible.

![demo video](demo.gif)   

# Installation

### Linux and OS-X

Simply run this command in a shell:
```
curl -sL https://raw.githubusercontent.com/Yoyodotpy/programming-language/main/install.sh | sh
```

### Windows

Run this command in a PowerShell window:
```
irm https://raw.githubusercontent.com/Yoyodotpy/programming-language/main/install.ps1 | iex
```

# Quick Start

Once you've installed the interpretter binary in the current folder, you can easily make custom scripts like a "hello world" or "fibonacci" program. When you install it, you also get two files containing example code, which can be useful references for your first script. If you feel like reading the thoughts behind the syntax and design, please read the next section (overview).

# Overview
Goal of the language: To make a usable functional programming language with as few types as possible.

### 1. way to write a lambda
    (x:x)
    Supports sugar for currying:
    (x,y:(x y))
### 2. way to apply value to lambda
    ((x:x) y)
### 3. way to assign a value to a variable
    (p = (x:x))
    (p x)
### 4. Two datatypes
    - Concell -> [x y]
        - Concells are "lists" of two values, being one value and a second one that points to it. Can be stacked forming larger "lists". Language will feature a way to write concell stacks in a smaller form. They also accept newlines to seperate entries.
        - Ex. : [1 [2 [3 4]]] == [1 2 3 4]

    - HEX -> '6E'
        - Hex numbers will be used for pretty much everything, with the compiler automatically translating numbers and other values to hex before adding them to thne AST.
        - Ex. : 12 -> 'C'

While these two(three?) datatypes aren't strictly required for the programming language to be turing complete and theoretically usable with a perfect cpu and infinite ram, modern cpus aren't designed to run complex pure untyped functional programs, so I decided to add these datatypes as a compromise to allow the programming language to run moderately quickly on modern cpus but still follow the base logic of lambda calculus.

Example of use to simulate other datatypes:
"hello world" -> ['68' '65' '6C' '6C' '6F' '20' '77' '6F' '72' '6C' '64']

Language features a "standard library" that translates common programming words into their proper lambda calculus functions, like True and False or nil.
```
(True = (x,y:x))
(False = (x,y:y))
(nil = (x:x))
```
Example programming excerpt:
```
(running = true)

(print_if_running = (text,times:(running (do
        ((c times) print text)
        (print 4)
    )))
)
(print_if_running "hello world",20)
```
