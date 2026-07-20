.PHONY: build test bench demo clean

BINARY := cmd/bloom-dedup-demo/bloom-dedup-demo.exe
CONTROL := testdata/control
OUTPUT := output

build:
	go build -o $(BINARY) ./cmd/bloom-dedup-demo

test:
	go test ./...

bench: build
	$(BINARY) bench --in $(CONTROL)/demo_events.jsonl --config $(CONTROL)/demo_config.json

demo: build
	$(BINARY) run --in $(CONTROL)/demo_events.jsonl --config $(CONTROL)/demo_config.json --out $(OUTPUT)/demo_result.jsonl --report $(OUTPUT)/demo_report.md
clean:
	@if exist "$(BINARY)" del "$(BINARY)"

