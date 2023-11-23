# Teamwork test task

## Task description

 Package customerimporter reads from the given customers.csv file and returns a
 sorted (data structure of your choice) of email domains along with the number
 of customers with e-mail addresses for each domain.  Any errors should be
 logged (or handled). Performance matters (this is only ~3k lines, but *could*
 be 1m lines or run on a small machine).

## Benchmarks

I decided to verify my solution using various configurations. The single-threaded implementation proves to be the most efficient, as the time spent on communication between goroutines (such as writing to a channel, reading from a channel, and utilizing wait groups and mutex) exceeds the time it takes to execute the data processing itself. The single-threaded solution is in the main branch, while the concurrent implementation is in the concurrency-realisation branch.


Below is a table with the average execution times:
| |  single-threaded | 1 goroutine for reading from file, 1 goroutine for extracting domain, no buffered channel | 1 goroutine for reading from file, 1 goroutine for extracting domain, channel with buffer 5 | 1 goroutine for reading from file, 4 goroutines for extracting domain, channel with buffer 5 |
|----------|----------|----------|----------|----------|
| 3К rows  | 33,8 ms  | 35,6 ms   | **32,3 ms**   | 35,1 ms   |
| 1M rows  | **3,46 s**   | 4,01 s   | 3,90 s   | 3,75 s   |

## Tests
Below is a picture with the test coverage of the solution:
<img width="913" alt="Screenshot 2023-11-23 at 18 47 12" src="https://github.com/ulalawell/teamwork-test-task/assets/77236479/e4bf793d-4bba-45e5-b0bf-c42dca75295e">
