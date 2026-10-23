# P2 - LFG Dungeon Queue

A Go program that simulates MMORPG LFG dungeon matchmaking using goroutines, semaphores, mutexes, and wait groups.

## Configuration

Edit `config.txt`:

```text
n=3
t=5
h=5
d=15
t1=2
t2=5
```

- `n` - maximum concurrent dungeon instances
- `t` - tanks
- `h` - healers
- `d` - DPS
- `t1` - minimum clear time
- `t2` - maximum clear time (`<= 15`)

## Run

Standard:

```bash
go run .
```

Using the run script (**MacOS only**):

```bash
bash run.sh
```

`run.sh` saves outputs inside `results/` as:

```text
result-1.txt
result-2.txt
...
```

## Build

```bash
go build -o p2 .
```

Run the compiled program:

```bash
./p2
```