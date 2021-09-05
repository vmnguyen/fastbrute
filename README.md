# fastbrute
HTTP brute forcer written in Go

# Usage
fastbrute:
  -c int
        Concurrent (default 50)
  -m int
        Mode to scan, 1 for stdin, 2 for request file (default 1)
  -r string
        Path to request file (default "/path/to/request/file")
  -t string
        Target to brute force (default "https://example.com")
  -w string
        Path to wordlist file (default "/path/to/wordlist/")
