PHONY: build-docs
build-docs:
	pip install -r docs/requirements.txt
	mkdocs build



