# Conductor

Conductor is an application super-orchestrator and scheduler for the Nalej platform.

## Musician

Musician are cluster allocated schedulers that collect information at a cluster level.

To run:

```bash
./bin/musician  run --prometheus="http://192.168.99.100:31080" --sleep=10000
```

or use a config file:

```bash
./bin/musician --config=config/musician.yaml run
```

