# Cookie Authentication Fix for Cross-Domain Setup

## Problem
The `/me` endpoint returns 200 OK, but the page redirects to login because cookies are not being sent from the frontend to the backend.

## Root Cause
**Cross-domain cookie issue**: Your backend (`trackerbackend-ao16.onrender.com`) and frontend (`trackerapp-livid.vercel.app`) are on different domains. Cookies require special configuration to work across domains.

## Backend Changes Made ✅

1. **Updated `GetUserInfo` function** to return proper 401 status codes when authentication fails
2. **Cookie settings already configured** with:
   - `SameSite=None` in production
   - `Secure=true` in production
   - `HttpOnly=true` for security
   - `AllowCredentials=true` in CORS

## Frontend Changes Required ⚠️

### 1. Include Credentials in All API Requests

**Every fetch request** to the backend must include `credentials: 'include'`:

```javascript
// Example: Checking authentication status
fetch('https://trackerbackend-ao16.onrender.com/me', {
  method: 'GET',
  credentials: 'include',  // ← CRITICAL: This sends cookies
  headers: {
    'Content-Type': 'application/json',
  },
})
  .then(response => response.json())
  .then(data => {
    if (data.loggedIn) {
      // User is authenticated
      console.log('User:', data.email);
    } else {
      // Redirect to login
      window.location.href = '/login';
    }
  });
```

### 2. Update All API Calls

Apply `credentials: 'include'` to **ALL** API requests:

```javascript
// GET request
fetch('https://trackerbackend-ao16.onrender.com/api/v1/investments', {
  credentials: 'include',
  headers: { 'Content-Type': 'application/json' },
});

// POST request
fetch('https://trackerbackend-ao16.onrender.com/api/v1/expenses', {
  method: 'POST',
  credentials: 'include',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify(expenseData),
});

// PUT request
fetch('https://trackerbackend-ao16.onrender.com/api/v1/budgets/123', {
  method: 'PUT',
  credentials: 'include',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify(budgetData),
});

// DELETE request
fetch('https://trackerbackend-ao16.onrender.com/api/v1/goals/456', {
  method: 'DELETE',
  credentials: 'include',
});
```

### 3. If Using Axios

```javascript
import axios from 'axios';

// Create axios instance with credentials
const api = axios.create({
  baseURL: 'https://trackerbackend-ao16.onrender.com',
  withCredentials: true,  // ← This is equivalent to credentials: 'include'
});

// Use it for all requests
api.get('/me');
api.post('/api/v1/expenses', expenseData);
api.put('/api/v1/budgets/123', budgetData);
api.delete('/api/v1/goals/456');
```

### 4. Update Authentication Check

```javascript
// In your auth context or authentication check
const checkAuth = async () => {
  try {
    const response = await fetch('https://trackerbackend-ao16.onrender.com/me', {
      credentials: 'include',  // ← MUST include this
    });
    
    const data = await response.json();
    
    if (response.ok && data.loggedIn) {
      // User is authenticated
      setUser({ email: data.email, userId: data.user_id });
      return true;
    } else {
      // Not authenticated - redirect to login
      setUser(null);
      return false;
    }
  } catch (error) {
    console.error('Auth check failed:', error);
    return false;
  }
};
```

## Testing the Fix

1. **Clear browser cookies** for both domains
2. **Login again** via Google OAuth
3. **Check browser DevTools**:
   - Network tab → Click on `/me` request
   - Check "Cookies" tab - should see `auth_token` being sent
   - Response should show `{"loggedIn": true, "email": "...", "user_id": "..."}`

## Why This Happens

- **Same-origin policy**: Browsers don't send cookies to different domains by default
- **Security measure**: Prevents malicious sites from accessing your cookies
- **Solution**: Explicitly tell the browser to include credentials with `credentials: 'include'`

## Backend Configuration (Already Done ✅)

```go
// CORS allows credentials
AllowCredentials: true

// Cookies use SameSite=None and Secure in production
if appEnv == "production" {
    secure = true
    sameSite = http.SameSiteNoneMode
}
c.SetCookie("auth_token", jwtToken, 3600, "/", domain, secure, true)
```

## Common Mistakes to Avoid

❌ **Don't do this:**
```javascript
fetch('https://trackerbackend-ao16.onrender.com/me')  // Missing credentials
```

✅ **Do this:**
```javascript
fetch('https://trackerbackend-ao16.onrender.com/me', {
  credentials: 'include'
})
```

## Summary

The backend is working correctly and returning 200 OK. The issue is that **cookies aren't being sent from the frontend**. Add `credentials: 'include'` to all fetch requests or `withCredentials: true` for axios.

---
**Status**: Backend fixed ✅ | Frontend changes required ⚠️