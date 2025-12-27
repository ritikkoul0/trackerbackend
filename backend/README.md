# Investment Tracker Backend API

A RESTful API built with Go and Gin framework for managing investments, budgets, goals, and expenses.

## 🚀 Tech Stack

- **Go** 1.21+
- **Gin** - Web framework
- **GORM** - ORM library
- **SQLite** - Database
- **CORS** - Cross-origin resource sharing

## 📁 Project Structure

```
backend/
├── config/
│   └── database.go          # Database configuration
├── controllers/
│   ├── investment_controller.go
│   ├── goal_controller.go
│   ├── budget_controller.go
│   ├── expense_controller.go
│   └── dashboard_controller.go
├── models/
│   ├── user.go
│   ├── investment.go
│   ├── goal.go
│   ├── budget.go
│   └── expense.go
├── routes/
│   └── routes.go            # API routes
├── main.go                  # Application entry point
├── go.mod                   # Go modules
└── .env.example             # Environment variables template
```

## 🛠️ Installation

### Prerequisites

- Go 1.21 or higher
- Git

### Steps

1. **Navigate to backend directory:**
   ```bash
   cd backend
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Create .env file:**
   ```bash
   cp .env.example .env
   ```

4. **Run the server:**
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080`

## 📡 API Endpoints

### Health Check
- `GET /api/v1/health` - Check server status

### Investments
- `GET /api/v1/investments` - Get all investments
- `GET /api/v1/investments/:id` - Get single investment
- `POST /api/v1/investments` - Create investment
- `PUT /api/v1/investments/:id` - Update investment
- `DELETE /api/v1/investments/:id` - Delete investment

### Goals
- `GET /api/v1/goals` - Get all goals
- `GET /api/v1/goals/:id` - Get single goal
- `POST /api/v1/goals` - Create goal
- `PUT /api/v1/goals/:id` - Update goal
- `DELETE /api/v1/goals/:id` - Delete goal

### Budgets
- `GET /api/v1/budgets` - Get all budgets
- `GET /api/v1/budgets/:id` - Get single budget
- `POST /api/v1/budgets` - Create budget
- `PUT /api/v1/budgets/:id` - Update budget
- `DELETE /api/v1/budgets/:id` - Delete budget

### Expenses
- `GET /api/v1/expenses` - Get all expenses
- `GET /api/v1/expenses/:id` - Get single expense
- `POST /api/v1/expenses` - Create expense
- `PUT /api/v1/expenses/:id` - Update expense
- `DELETE /api/v1/expenses/:id` - Delete expense

### Dashboard
- `GET /api/v1/dashboard` - Get dashboard summary

## 📝 Request Examples

### Create Investment
```bash
curl -X POST http://localhost:8080/api/v1/investments \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Nifty 50 Index Fund",
    "type": "Mutual Fund",
    "invested": 50000,
    "current_value": 56250,
    "purchase_date": "2024-01-15T00:00:00Z"
  }'
```

### Create Goal
```bash
curl -X POST http://localhost:8080/api/v1/goals \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Emergency Fund",
    "target_amount": 500000,
    "current_amount": 150000,
    "deadline": "2024-12-31T00:00:00Z",
    "priority": "High"
  }'
```

### Create Budget
```bash
curl -X POST http://localhost:8080/api/v1/budgets \
  -H "Content-Type: application/json" \
  -d '{
    "month": "2024-01",
    "income": 75000,
    "total_expenses": 45000,
    "savings_goal": 35000
  }'
```

## 🗄️ Database Models

### Investment
- ID, Name, Type, Invested, CurrentValue, Returns, Status, PurchaseDate

### Goal
- ID, Name, TargetAmount, CurrentAmount, Deadline, Status, Priority, Description

### Budget
- ID, Month, Income, TotalExpenses, Savings, SavingsGoal

### Expense
- ID, Category, Amount, Description, Date, BudgetID

## 🔧 Development

### Run in development mode:
```bash
go run main.go
```

### Build for production:
```bash
go build -o investment-tracker-api
./investment-tracker-api
```

### Run tests:
```bash
go test ./...
```

## 🌐 CORS Configuration

The API is configured to accept requests from `http://localhost:3000` by default. Update the CORS settings in `main.go` if needed.

## 📊 Database

The application uses SQLite for simplicity. The database file `investment_tracker.db` will be created automatically on first run.

### Auto-Migration

Database tables are automatically created/updated based on the models when the server starts.

## 🚀 Deployment

### Build binary:
```bash
go build -o investment-tracker-api
```

### Run binary:
```bash
./investment-tracker-api
```

## 📄 License

This project is open source and available for educational purposes.

## 🤝 Contributing

Feel free to submit issues and pull requests!