#!/bin/bash

# Script to switch between development and production environments

if [ "$1" == "dev" ] || [ "$1" == "development" ]; then
    echo "Switching to DEVELOPMENT environment..."
    if [ -f .env.local ]; then
        cp .env .env.production.backup
        cp .env.local .env
        echo "✅ Switched to development (.env.local → .env)"
        echo "   - APP_ENV=development"
        echo "   - FRONTEND_URL=http://localhost:3000"
        echo "   - GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback"
    else
        echo "❌ .env.local not found!"
    fi
elif [ "$1" == "prod" ] || [ "$1" == "production" ]; then
    echo "Switching to PRODUCTION environment..."
    if [ -f .env.production.backup ]; then
        cp .env.production.backup .env
        echo "✅ Switched to production (.env.production.backup → .env)"
    else
        echo "⚠️  No backup found. Please manually configure .env for production:"
        echo "   - APP_ENV=production"
        echo "   - FRONTEND_URL=https://trackerapp-livid.vercel.app"
        echo "   - GOOGLE_REDIRECT_URL=https://trackerbackend-ao16.onrender.com/auth/google/callback"
    fi
else
    echo "Usage: ./switch-env.sh [dev|prod]"
    echo ""
    echo "Examples:"
    echo "  ./switch-env.sh dev   - Switch to development environment"
    echo "  ./switch-env.sh prod  - Switch to production environment"
fi

# 
