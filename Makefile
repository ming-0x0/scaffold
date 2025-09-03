.PHONY: proto.gen proto.lint proto.clean proto.dep

# Generate proto files
proto.gen:
	@cd api && buf generate

# Lint proto files
proto.lint:
	@cd api && buf lint

# Clean generated files
proto.clean:
	@find . -name "*.pb.go" -delete
	@find . -name "*.pb.gw.go" -delete
	@find . -name "openapi.yaml" -delete

# Update dependencies
proto.dep:
	@cd api && buf dep update
