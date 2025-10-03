.PHONY: docker-up docker-down test-e2e

docker-up:
	docker compose up -d --build
	@echo "Waiting for services to be ready..."
	@until curl -s http://localhost:8080/healthz >/dev/null; do sleep 1; done
	@echo "Services are ready!"

docker-down:
	docker compose down

test-e2e:
	./e2e-test.sh