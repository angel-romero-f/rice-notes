# Frontend Environment Configuration

To run the frontend with proper API connectivity, create a `.env.local` file in this directory with:

```bash
NEXT_PUBLIC_API_URL=http://localhost:8080
```

This ensures the frontend correctly connects to the backend server running on port 8080.

## Why this is needed

Without this configuration, the frontend will default to calling `http://localhost:8081`, which will result in connection errors since the backend runs on port 8080.