# Distributed System Challenges

## Challenges:

### Echo (Completed):
```bash
./maelstrom test -w echo --bin ./main --node-count 1 --time-limit 10
```

### Unique Ids (Completed):
```bash
./maelstrom test -w unique-ids --bin ./main --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition
```

### Broadcast:
#### a: Single-Node Broadcast (Completed)
```bash
./maelstrom test -w broadcast --bin ./main --node-count 1 --time-limit 20 --rate 10
```
#### b: Multi-Node Broadcast (Completed)
```bash
./maelstrom test -w broadcast --bin ./main --node-count 5 --time-limit 20 --rate 10
```
#### c: Fault Tolerant Broadcast (Completed)
```bash
./maelstrom test -w broadcast --bin ./main --node-count 5 --time-limit 20 --rate 10
```

#### d: Efficient Broadcast (In Progress)
```bash
./maelstrom test -w broadcast --bin ./main --node-count 25 --time-limit 20 --rate 100 --topology total
```