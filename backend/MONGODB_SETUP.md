# MongoDB Setup Guide

## 🗄️ MongoDB Configuration

This backend uses MongoDB as the database, connecting to `localhost:27017` with database name `investmentdb`.

## 📋 Prerequisites

- MongoDB installed and running on your system
- MongoDB running on `localhost:27017`

## 🚀 Quick Start

### 1. Install MongoDB (if not installed)

#### macOS (using Homebrew):
```bash
brew tap mongodb/brew
brew install mongodb-community
brew services start mongodb-community
```

#### Ubuntu/Debian:
```bash
sudo apt-get install -y mongodb
sudo systemctl start mongodb
sudo systemctl enable mongodb
```

#### Windows:
Download and install from: https://www.mongodb.com/try/download/community

### 2. Verify MongoDB is Running

```bash
# Check if MongoDB is running
mongosh --eval "db.version()"

# Or connect to MongoDB shell
mongosh
```

### 3. Create Database and Collections

The application will automatically create the database and collections, but you can manually create them:

```bash
mongosh
```

Then in MongoDB shell:
```javascript
use investmentdb

// Create collections
db.createCollection("investments")
db.createCollection("goals")
db.createCollection("budgets")
db.createCollection("expenses")
db.createCollection("users")

// Verify collections
show collections
```

### 4. Install Go Dependencies

```bash
cd backend
go mod tidy
```

This will download all required packages including the MongoDB driver.

### 5. Configure Environment

```bash
cp .env.example .env
```

The `.env` file should contain:
```
MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=investmentdb
PORT=8080
ALLOWED_ORIGINS=http://localhost:3000
APP_ENV=development
```

### 6. Run the Backend

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## 📊 Database Structure

### Collections:

#### 1. **investments**
```javascript
{
  _id: ObjectId,
  created_at: ISODate,
  updated_at: ISODate,
  user_id: ObjectId,
  name: String,
  type: String,  // "Stocks", "Mutual Fund", "ETF", "FD", "PPF"
  invested: Number,
  current_value: Number,
  returns: Number,  // Percentage
  status: String,  // "Growing", "Stable", "Declining"
  purchase_date: ISODate
}
```

#### 2. **goals**
```javascript
{
  _id: ObjectId,
  created_at: ISODate,
  updated_at: ISODate,
  user_id: ObjectId,
  name: String,
  target_amount: Number,
  current_amount: Number,
  deadline: ISODate,
  status: String,  // "Planned", "In Progress", "Completed"
  priority: String,  // "High", "Medium", "Low"
  description: String
}
```

#### 3. **budgets**
```javascript
{
  _id: ObjectId,
  created_at: ISODate,
  updated_at: ISODate,
  user_id: ObjectId,
  month: String,  // "2024-01"
  income: Number,
  total_expenses: Number,
  savings: Number,
  savings_goal: Number
}
```

#### 4. **expenses**
```javascript
{
  _id: ObjectId,
  created_at: ISODate,
  updated_at: ISODate,
  user_id: ObjectId,
  budget_id: ObjectId,
  category: String,  // "Food", "Transport", "Entertainment"
  amount: Number,
  description: String,
  date: ISODate
}
```

#### 5. **users**
```javascript
{
  _id: ObjectId,
  created_at: ISODate,
  updated_at: ISODate,
  name: String,
  email: String,  // unique
  password: String  // hashed
}
```

## 🔍 Useful MongoDB Commands

### View all databases:
```javascript
show dbs
```

### Switch to investmentdb:
```javascript
use investmentdb
```

### View all collections:
```javascript
show collections
```

### Query investments:
```javascript
db.investments.find().pretty()
```

### Count documents:
```javascript
db.investments.countDocuments()
```

### Insert sample data:
```javascript
db.investments.insertOne({
  name: "Nifty 50 Index Fund",
  type: "Mutual Fund",
  invested: 50000,
  current_value: 56250,
  returns: 12.5,
  status: "Growing",
  purchase_date: new Date("2024-01-15"),
  created_at: new Date(),
  updated_at: new Date()
})
```

### Delete all documents in a collection:
```javascript
db.investments.deleteMany({})
```

### Drop a collection:
```javascript
db.investments.drop()
```

### Create indexes for better performance:
```javascript
db.investments.createIndex({ user_id: 1 })
db.investments.createIndex({ status: 1 })
db.goals.createIndex({ user_id: 1 })
db.goals.createIndex({ status: 1 })
db.budgets.createIndex({ user_id: 1, month: 1 })
db.expenses.createIndex({ user_id: 1, budget_id: 1 })
```

## 🛠️ Troubleshooting

### MongoDB not starting:
```bash
# Check MongoDB status
sudo systemctl status mongodb

# Restart MongoDB
sudo systemctl restart mongodb
```

### Connection refused:
- Ensure MongoDB is running on port 27017
- Check firewall settings
- Verify MongoDB configuration in `/etc/mongod.conf`

### Permission denied:
```bash
# Fix MongoDB data directory permissions
sudo chown -R mongodb:mongodb /var/lib/mongodb
sudo chown mongodb:mongodb /tmp/mongodb-27017.sock
```

## 📚 Additional Resources

- [MongoDB Official Documentation](https://docs.mongodb.com/)
- [MongoDB Go Driver Documentation](https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo)
- [MongoDB University (Free Courses)](https://university.mongodb.com/)

## 🔐 Security Notes

For production:
1. Enable MongoDB authentication
2. Create a dedicated database user
3. Use strong passwords
4. Enable SSL/TLS
5. Configure firewall rules
6. Regular backups

Example with authentication:
```bash
MONGO_URI=mongodb://username:password@localhost:27017