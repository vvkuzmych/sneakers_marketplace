#!/bin/bash
# Run subscription system migrations in order
# Phase 2, Day 1

set -e  # Exit on error

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🚀 Running Subscription System Migrations"
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

# Run migrations in order
migrations=(
    "01_create_subscription_plans"
    "02_create_user_subscriptions"
    "03_create_subscription_transactions"
    "04_seed_subscription_plans"
    "05_assign_users_to_free"
    "06_add_subscription_helpers"
)

for migration in "${migrations[@]}"; do
    file="internal/database/migrations/20260123_${migration}.up.sql"
    echo "📝 Running: $migration"
    
    if [ -f "$file" ]; then
        psql -U $USER -d $DB -f "$file"
        echo "✅ Completed: $migration"
        echo ""
    else
        echo "❌ File not found: $file"
        exit 1
    fi
done

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🎉 All migrations completed successfully!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Show results
echo "📊 Subscription System Summary:"
psql -U $USER -d $DB -c "SELECT name, display_name, price_monthly, buyer_fee_percent, seller_fee_percent FROM subscription_plans ORDER BY sort_order;"
echo ""
psql -U $USER -d $DB -c "SELECT COUNT(*) as total_active_subscriptions FROM user_subscriptions WHERE status = 'active';"
