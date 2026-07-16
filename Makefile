.PHONY: build test bench demo clean

APP=C:\Users\Alexandra\GolandProjects\bloom-dedup-demo\cmd\bloom-dedup-demo\bloom-dedup-demo.exe
ROOT=C:\Users\Alexandra\GolandProjects\bloom-dedup-demo
CONTROL=C:\Users\Alexandra\GolandProjects\bloom-dedup-demo\testdata\control
OUTPUT=C:\Users\Alexandra\GolandProjects\bloom-dedup-demo\output

#build:
#	...

#test:
#	...

#bench:
#	...

demo:
	$(APP) run --in $(CONTROL)\demo_events.jsonl --config $(CONTROL)\demo_config.json --out $(OUTPUT)\demo_result.jsonl --report $(OUTPUT)\demo_report.json
