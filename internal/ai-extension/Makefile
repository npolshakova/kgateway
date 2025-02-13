WEBHOOK_API_SOURCES = docs/export_openapi.py samples/app.py ai_extension/guardrails/api.py
TEST_DIR = ../../test/kubernetes/e2e/features/aiextension/

webhooks-openapi.yaml: $(WEBHOOK_API_SOURCES)
	python docs/export_openapi.py app:app --app-dir samples/ --out docs/$@

.PHONY: unit-tests
unit-tests:
	pyright .
	(cd ai_extension && python -m pytest test/ -v -o log_cli_level=DEBUG)

.PHONY: lint
lint:
	ruff check . --output-format=github --exclude '**/*_pb2.py' --exclude '**/*_pb2.pyi' --exclude '**/*_pb2_*.py' --fix
	ruff check ${TEST_DIR} --output-format=github --exclude '**/*_pb2.py' --exclude '**/*_pb2.pyi' --exclude '**/*_pb2_*.py' --fix
	ruff format . --exclude '**/*_pb2.py' --exclude '**/*_pb2.pyi' --exclude '**/*_pb2_*.py'
	ruff format ${TEST_DIR} --exclude '**/*_pb2.py' --exclude '**/*_pb2.pyi' --exclude '**/*_pb2_*.py'

.PHONY: run-webhook-server
run-webhook-server:
	python3 -m fastapi run  --host 0.0.0.0 --port 7891 samples/app.py 

