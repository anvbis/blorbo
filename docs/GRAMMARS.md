# Grammars

## Expressions

```
expr -> assign 
```

### Assignment

```
assign -> IDENT "=" assign 
        | log_or
```

### Binary expressions

```
log_or -> log_and ( "or" log_and )*
```

```
log_and -> bit_or ( "and" bit_or )*
```

```
bit_or -> bit_xor ( "|" bit_xor )*
```

```
bit_xor -> bit_and ( "^" bit_and )*
```

```
bit_and -> equality ( "&" equality )* 
```

```
equality -> comparison ( ( "==" | "!=" ) comparison )* 
```

```
comparison -> shift ( ( ">" | ">=" | "<" | "<=" ) shift )* 
```

```
shift -> term ( ( ">>" | "<<" ) term )*
```

```
term -> factor ( ( "+" | "-" ) factor )*
```

```
factor -> unary ( ( "*" | "/" | "%" ) unary )*
```

### Unary expressions

```
unary -> ( "+" | "-" | "!" | "~" ) unary
       | call
```

```
call -> primary ( "(" args? ")" )*
args -> expr ( "," expr )*
```

### Primary expressions

```
primary -> IDENT
         | INTEGER
         | FLOAT
         | "true"
         | "false"
         | "null"
         | "(" expr ")"
```

