# Makefile for testing /users/register API endpoint

# Variables
URL = http://localhost:8080
ENDPOINT = /users/register
REGISTER_URL = $(URL)$(ENDPOINT)

# Colors for output
GREEN = \033[0;32m
RED = \033[0;31m
NC = \033[0m # No Color

# Default user data (modify as needed)
NAME ?= testuser
PHONE ?= 09110072557

.PHONY: help register register-json register-verbose register-file clean

help:
	@echo "Available targets:"
	@echo "  make register           - Register a user with default data (JSON)"
	@echo "  make register-json      - Register with custom JSON from file"
	@echo "  make register-verbose   - Register with verbose curl output"
	@echo "  make clean              - Remove temporary files"
	@echo ""
	@echo "  make register NAME=johndoe PHONE=09030072667

# Basic registration with default data
register:
	@echo "Registering user: $(NAME) ($(PHONE))"
	@curl -X POST $(REGISTER_URL) \
		-H "Content-Type: application/json" \
		-d '{"name":"$(NAME)","phone_number":"$(PHONE)"}'
	@echo ""

# Registration with verbose output (shows headers, etc.)
register-verbose:
	@echo "Registering user (verbose mode): $(PHONE)"
	@curl -X POST $(REGISTER_URL) \
		-H "Content-Type: application/json" \
		-d '{"name":"$(NAME)","phone_number":"$(PHONE)"}' \
		-v

# Registration using a JSON file
register-json:
	@echo "Registering user using JSON file"
	@curl -X POST $(REGISTER_URL) \
		-H "Content-Type: application/json" \
		-d @user_data.json

# Create example JSON file 
user_data.json:
	@echo '{"name":"testuser","phone_number":"09040072667"}' > user_data.json
	@echo "Created example user_data.json"

clean:
	@rm -f user_data.json response.json
	@echo "Cleaned up temporary files"