#!/bin/bash

# Kanban Reports - Setup Script
# This script helps new users get started quickly

set -e

echo "🔄 Kanban Reports - Setup Script"
echo "================================="
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+ first."
    echo "   Visit: https://golang.org/dl/"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | cut -d' ' -f3 | sed 's/go//')
REQUIRED_VERSION="1.21"

echo "✅ Go found: version $GO_VERSION"

# Build the application
echo ""
echo "🔨 Building application..."
make build

# Create data directory if it doesn't exist
if [ ! -d "data" ]; then
    echo "📁 Creating data directory..."
    mkdir -p data
fi

# Check if sample data exists
if [ ! -f "data/sample.csv" ]; then
    echo "📄 Creating sample CSV file..."
    cat > data/sample.csv << 'EOF'
id,name,type,estimate,is_completed,completed_at,owners,epic,team,product_area,created_at,started_at,labels
1,User Authentication,Feature,5,TRUE,2024/05/07 10:30:00,john@example.com,User Management,Team Alpha,Backend,2024/05/01 09:00:00,2024/05/03 11:00:00,feature
2,Fix Login Bug,Bug,2,TRUE,2024/05/08 15:45:00,jane@example.com,User Management,Team Alpha,Frontend,2024/05/02 14:00:00,2024/05/05 10:00:00,bug
3,Dashboard UI,Feature,8,TRUE,2024/05/10 16:30:00,john@example.com;jane@example.com,Analytics,Team Beta,Frontend,2024/05/03 08:00:00,2024/05/06 09:00:00,feature
4,Database Migration,Task,3,FALSE,,bob@example.com,Infrastructure,Team Beta,Backend,2024/05/04 11:00:00,2024/05/08 14:00:00,task
5,Urgent Hotfix,Feature,1,TRUE,2024/05/09 12:00:00,alice@example.com,Support,Team Alpha,Backend,2024/05/08 10:00:00,2024/05/09 11:00:00,ad-hoc-request
EOF
    echo "✅ Sample data created: data/sample.csv"
fi

echo ""
echo "🎉 Setup complete!"
echo ""
echo "📖 Next steps:"
echo "   1. Try the interactive mode:"
echo "      ./bin/kanban-reports --interactive"
echo ""
echo "   2. Or run with sample data:"
echo "      ./bin/kanban-reports --csv data/sample.csv --type contributor --last 30"
echo ""
echo "   3. Get help:"
echo "      ./bin/kanban-reports --help"
echo ""
echo "   4. See examples:"
echo "      ./bin/kanban-reports --examples"
echo ""
echo "📁 Place your own CSV files in the 'data/' directory"