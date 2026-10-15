create-table:
	@if [ -z "$(name)" ]; then \
		echo "❌ Error: Missing migration name."; \
		echo "Usage: make create-table name=your_migration_name"; \
		exit 1; \
	fi
	@goose -dir=db/migrations create $(name) sql

migration-up: 
	@goose -dir=db/migrations sqlite db/login_example.db up

migration-down: 
	@goose -dir=db/migrations sqlite db/login_example.db down

migration-fresh:
	@echo "🧹 Dropping all migrations..."
	@goose -dir=db/migrations sqlite db/login_example.db reset
	@echo "🚀 Applying all migrations..."
	@goose -dir=db/migrations sqlite db/login_example.db up
	@echo "✅ Database is fresh!"

%:
	@: