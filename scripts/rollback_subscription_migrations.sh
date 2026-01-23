#!/bin/bash
# Rollback subscription system migrations (in reverse order)
# Phase 2, Day 1

set -e  # Exit on error

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "⏪ Rolling Back Subscription System Migrations"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

DB="sneakers_marketplace"
USER="postgres"

# Check database connection
echo "🔍 Checking database connection..."
if ! psql -U $USER -d $DB -c "SELECT 1;" > /dev/null 2>&1; then
    echo "❌ Cannot connect to database $DB"
    exit 1
fi
echo "✅ Database connected"
echo ""

# Confirm rollback
read -p "⚠️  Are you sure you want to rollback subscription migrations? (yes/no): " confirm
if [ "$confirm" != "yes" ]; then
    echo "❌ Rollback cancelled"
    exit 0
fi
echo ""

# Run migrations in REVERSE order
migrations=(
    "06_add_subscription_helpers"
    "05_assign_users_to_free"
    "04_seed_subscription_plans"
    "03_create_subscription_transactions"
    "02_create_user_subscriptions"
    "01_create_subscription_plans"
)

for migration in "${migrations[@]}"; do
    file="internal/database/migrations/20260123_${migration}.down.sql"
    echo "📝 Rolling back: $migration"
    
    if [ -f "$file" ]; then
        psql -U $USER -d $DB -f "$file"
        echo "✅ Rolled back: $migration"
        echo ""
    else
        echo "❌ File not found: $file"
        exit 1
    fi
done

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ All migrations rolled back successfully!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
