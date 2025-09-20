#  String-Decompression

The project decompresses strings containing character repetitions and escape sequences.  
It reads a string and expands numbers after characters into repeated sequences.  
Escape sequences (`\`) allow digits and slashes to be treated as literal characters.

---
## ⚙️ Prerequisites

- Go **1.21+**  

---

##  Installation

Clone the repository and set up the Go module:

```bash
git clone https://github.com/veranemoloko/wb-tech-l2.git
cd wb-tech-l2/string_decompression
go mod tidy
```
---
## Examples
```go
unpackString("a3")      // "aaa"
unpackString("abc")     // "abc"
unpackString("a\\3")    // "a3"
unpackString("\\\\")    // "\"
unpackString("a10")     // "aaaaaaaaaa"
```
---
## Limitations

- Maximum repeat: MaxRepeat = 1_000_000
- Strings cannot start with a digit
- Strings cannot end with a single backslash
---
## Advantages of solution ✨
- The `unpackString` function was implemented using a **finite state machine** approach.  
- Using explicit states makes it easy to follow the logic.  
- It prevents excessive memory allocation by limiting repetitions with `MaxRepeat`.  
- Easily extendable to support more escape rules or other encoding schemes.  
- Uses `strings.Builder` to build the result, avoiding unnecessary string concatenations.
- The function uses strings.Builder with an initial preallocation equal to the input string length.
- Unit tests cover normal cases, escapes, repetitions, and error conditions.
