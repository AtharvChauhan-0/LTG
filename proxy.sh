#!/bin/bash

# Script to fetch free HTTP proxies for testing
# Usage: ./get-proxies.sh

echo "════════════════════════════════════════"
echo "  FREE PROXY FETCHER"
echo "════════════════════════════════════════"
echo ""

# Check if curl is available
if ! command -v curl &> /dev/null; then
    echo "❌ Error: curl is not installed"
    echo "   Please install curl first"
    exit 1
fi

echo "📥 Fetching free proxies from GitHub..."
echo ""

# Try multiple sources
SOURCES=(
    "https://raw.githubusercontent.com/TheSpeedX/PROXY-List/master/http.txt"
    "https://raw.githubusercontent.com/clarketm/proxy-list/master/proxy-list-raw.txt"
    "https://raw.githubusercontent.com/ShiftyTR/Proxy-List/master/http.txt"
)

> proxies.txt  # Clear file

success=0
for source in "${SOURCES[@]}"; do
    echo "Trying: $source"
    if curl -s -f "$source" | grep -E "^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+:[0-9]+$" >> proxies.txt; then
        echo "✅ Success"
        success=1
    else
        echo "⚠️  Failed"
    fi
    echo ""
done

if [ $success -eq 0 ]; then
    echo "❌ Could not fetch proxies from any source"
    echo ""
    echo "Manual options:"
    echo "1. Visit https://free-proxy-list.net/"
    echo "2. Copy proxies and paste into proxies.txt"
    echo "3. Format: IP:PORT (one per line)"
    exit 1
fi

# Add http:// prefix and remove duplicates
sed -i 's|^|http://|' proxies.txt 2>/dev/null || sed -i '' 's|^|http://|' proxies.txt
sort -u proxies.txt -o proxies.txt

count=$(wc -l < proxies.txt | tr -d ' ')

echo "════════════════════════════════════════"
echo "✅ Success! Downloaded $count proxies"
echo "════════════════════════════════════════"
echo ""
echo "File: proxies.txt"
echo ""
echo "Preview (first 5):"
head -5 proxies.txt
echo ""
echo "Next steps:"
echo "1. Run your load test: go run main.go logger.go"
echo "2. The tool will test each proxy before using it"
echo ""
echo "Note: Free proxies are often slow/unreliable"
echo "      Consider premium proxies for production"