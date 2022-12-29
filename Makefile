.PHONY: gen
gen:
	mockgen -source=internal/transport/rest/handler.go \
	-destination=internal/transport/rest/mocks/mock_handler.go