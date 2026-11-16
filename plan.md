Go Practice: learning plan
A small event-ticketing API ("EventFlow") that I build to learn Go. Standard library only (net/http), seats kept in memory, everything in main.go for now.
For the AI assistant: read this first
I'm a JavaScript/TypeScript developer (Node.js, Express, React). I'm new to Go.
Compare Go with plain Node.js/Express: plain objects and functions. I don't use classes or NestJS, so don't explain with them.
Go one small step at a time. Give me hints and a skeleton with // TODO comments. I type the code myself, so don't write the whole solution.
When I paste code, review it: point out only what really matters, not tiny style details.
When there's a compiler error, explain in plain words what it means and where the problem is.
Keep answers short. Explain in Bulgarian if you can.
Where the project stands
Done:
Mission 1: HTTP server with http.NewServeMux, GET /events/{id} and GET /health, JSON with struct tags, log.Fatal(http.ListenAndServe(...))
Mission 2: seats in a slice, GET /events/{id}/seats?status=free filter
Mission 3, first half:
memStore struct holding seats map[string]Seat, keyed by seat ID
newMemStore() creates 20 free seats for event "20"
method SeatsByEventID(eventID string) []Seat returns that event's seats sorted with slices.SortFunc + cmp.Compare
the GET seats handler uses it
Start here: compare main.go with the list above. If something is missing, finish it first. Then continue with Mission 3 below.
How to run and check
go run .                              # start the server (Ctrl+C, then run again after every change)
curl -i localhost:8080/events/20/seats
go vet ./...                          # catches mistakes like broken JSON tags
gofmt -w .                            # formats the code (the editor also does it on save)
Rules
One step at a time. A mission is done when its "Done when" passes.
Commit after each step, push at the end of the day.
The repo is public: nothing from work goes in.
Mission 3 · finish: errors and a POST request
Learn: errors are normal return values, sentinel errors, errors.Is, the "comma ok" map lookup, reading a JSON request body.
Create the errors once, at the top of the file, like const NOT_FOUND = new Error(...) in JS:
var ErrNotFound = errors.New("not found")
var ErrSeatTaken = errors.New("seat taken")
Method func (s *memStore) Hold(seatID string) error:
seat, ok := s.seats[seatID]: ok is false when the key doesn't exist → return ErrNotFound
status isn't "free" → return ErrSeatTaken
set seat.Status = "held", then write it back: s.seats[seatID] = seat (the value you got is a copy)
return nil means "no error"
Quick check in main: fmt.Println(store.Hold("7"), store.Hold("7"), store.Hold("999")) → <nil> seat taken not found
POST /reservations with body {"seatId": "7"}:
a small struct for the body, with a JSON tag
json.NewDecoder(r.Body).Decode(&body): broken JSON → 400
call Hold, then pick the status with errors.Is(err, ErrNotFound) → 404, errors.Is(err, ErrSeatTaken) → 409
success → 201
Error wrapping: inside Hold, return fmt.Errorf("hold seat %s: %w", seatID, ErrSeatTaken) instead of the bare error. Check that errors.Is still finds it, and understand why (%w keeps the original inside).
Remember: set the status (w.WriteHeader or http.Error) before writing the body. Once the body starts, the status is locked at 200.
Done when: with curl -i -X POST localhost:8080/reservations -d '{"seatId":"7"}' you get 201, then 409 for the same seat, 404 for "999", 400 for broken JSON, and GET shows seat 7 as "held".
Mission 4 · functions
Learn: multiple return values, helper functions, functions that take and return functions, defer.
Method func (s *memStore) Seat(id string) (Seat, error): returns the seat or ErrNotFound. Add GET /seats/{id} that uses it: 200 or 404.
Helper func writeJSON(w http.ResponseWriter, status int, v any): sets Content-Type, writes the status, encodes v. Use it in every JSON handler, so the handlers get shorter.
Method Release(seatID string) error and DELETE /reservations/{seatId}: held → free (204), not held → 409, unknown → 404.
A function that wraps a handler, like Express middleware:
func logRequests(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // TODO: remember the start time, call next(w, r),
        //       log method, path and how long it took
    }
}
Then try the same log with defer instead.
Done when: every route works with curl, and you can explain what logRequests receives and what it returns.
Mission 5 · goroutines and the mutex
Learn: go f(), sync.WaitGroup, race conditions, sync.Mutex, the race detector.
Read Go by Example: Goroutines, WaitGroups, Mutexes.
Two buyers, one seat: in main (before the server starts), start 100 goroutines that all call store.Hold("5") at the same time. Each goroutine writes its result into its own slot of results := make([]error, 100). Wait with a WaitGroup, then count how many got nil.
Put time.Sleep(time.Millisecond) in Hold between the "is it free?" check and the write. Run go run -race . and read the report. You'll see several winners, a "DATA RACE" report, or a "concurrent map writes" crash.
Fix it: add mu sync.Mutex to memStore, and at the start of every method do s.mu.Lock() and defer s.mu.Unlock().
A mutex must never be copied: make newMemStore() return *memStore (return &memStore{...}). go vet warns about copied locks.
This is the same bug as two users signing up with the same username at the same moment.
Done when: exactly one winner, -race is quiet, and you can explain why the mutex fixes it.
Mission 6 · channels
Learn: a channel is a pipe between goroutines. Send ch <- v, receive v := <-ch, buffered vs unbuffered, close + range, select with a timeout.
Read Go by Example: Channels, Channel Buffering, Closing Channels, Range over Channels, Select, Timeouts.
Redo Mission 5's counting with a channel: results := make(chan error, 100). Each goroutine sends its result, and main receives 100 times and counts. No slice needed.
Background worker, an in-process queue (like RabbitMQ, but inside the program):
in main, create jobs := make(chan string, 10) and start one goroutine that loops over jobs and, for each seat ID, sleeps 2 seconds and logs "confirmation sent for seat X"
after a successful hold, the POST handler sends the seat ID into jobs and answers 201 right away
curl should return at once, and the log line appears 2 seconds later
select with a timeout: a function that waits for a result from a channel, but gives up after 1 second with time.After.
Done when: you can explain buffered vs unbuffered, who closes a channel (the sender), and why the handler doesn't wait for the "email".
Optional, if there's time
Context: r.Context(), context.WithTimeout, and a slow handler that stops when the client presses Ctrl+C (select on ctx.Done()). In real services almost every function that does I/O takes ctx context.Context as its first parameter.
Tests: a table-driven test for Hold (free seat, seat already held, unknown ID), one handler test with httptest that expects 409, then go test ./... and go test -race ./....
Structure: move the seat code into a package (internal/seats), add a SeatStore interface, and a NewServer(store) http.Handler so main only wires things together.
TypeScript / Node → Go
Node / TypeScript
Go
package.json, npm install x
go.mod, go get x
object { id: "1" }
struct value Seat{ID: "1"} (the type name is required)
object used as a dictionary
map[string]Seat, created with make
array
slice []Seat, grows with append
export
a capital first letter: Seat is public, seat isn't
null / undefined
zero values (0, "", false) and nil
try / catch, throw
return err, then if err != nil { ... } at the caller
arr.sort((a, b) => a.n - b.n)
slices.SortFunc(arr, func(a, b Seat) int { return cmp.Compare(a.N, b.N) })
`event ${id}`
fmt.Sprintf("event %s", id)
console.log
fmt.Println (terminal); fmt.Fprintln(w, ...) writes to the response
res.json(x)
json.NewEncoder(w).Encode(x)
req.params.id
r.PathValue("id")
req.query.status
r.URL.Query().Get("status")
const x = ... (new variable)
x := ... (= only changes an existing one)
Go quirks I've met so far
for _, s := range seats gives a copy of each item; change the original with seats[i], or write it back into the map.
A map has no order; sort a slice if you need one.
A map declared with var is nil and crashes on write; create it with make.
An empty slice declared with var becomes null in JSON; []Seat{} becomes [].
A JSON tag without quotes (json:name) is silently ignored; go vet catches it.
Unused variables and imports are compile errors.
Links
Go by Example: short runnable examples for every topic above
pkg.go.dev: docs for every package
Working with Errors in Go 1.13: %w, errors.Is, errors.As