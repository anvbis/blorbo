# Specification
...

### Operator precedence

| Precedence | Operator  | Description                                                        | Associativity |
|------------|-----------|--------------------------------------------------------------------|---------------|
| 1          | () [] .   | Function call, array index, property access                        | Left-to-right |
| 2          | + - ! ~   | Unary plus, unary minus, logical not, bitwise not                  | Right-to-left |
| 3          | * / %     | Multiplication, division, modulus                                  | Left-to-right |
| 4          | + -       | Addition Subtraction                                               | Left-to-right |
| 5          | >> <<     | Bitwise right shift, bitwise left shift                            | Left-to-right |
| 6          | > >= < <= | Greater than, greater than or equal, less than, less than or equal | Left-to-right |
| 7          | == !=     | Equal, not equal                                                   | Left-to-right |
| 8          | &         | Bitwise and                                                        | Left-to-right |
| 9          | ^         | Bitwise xor                                                        | Left-to-right |
| 10         | \|        | Bitwise or                                                         | Left-to-right |
| 11         | and       | Logical and                                                        | Left-to-right |
| 12         | or        | Logical or                                                         | Left-to-right |
| 13         | =         | Assignment                                                         | Right-to-left |

