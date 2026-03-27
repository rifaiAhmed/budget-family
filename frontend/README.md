# Family Budget Control (Frontend)

Mobile-first React 18 + TypeScript web app using MUI, React Router, React Query, React Hook Form, Zod, Recharts.

## Requirements

- Node.js 18+
- Backend running on `http://localhost:8080`

## Setup

1) Configure environment

Copy `.env.example` to `.env` and set the API URL:

```bash
# frontend/.env
VITE_API_URL=http://localhost:8080
VITE_DEFAULT_FAMILY_ID=
```

2) Install dependencies

```bash
npm install
```

3) Run dev server

```bash
npm run dev
```

Open `http://localhost:5173`

## Notes

- JWT access/refresh tokens are stored in `localStorage`.
- The app expects a `family_id` for list/create calls.
  - Temporarily you can set it in `Settings` page (stored in `localStorage`).
  - Or set `VITE_DEFAULT_FAMILY_ID` in `.env`.

## Suggested flow (first time)

- Register
- Create family from backend/API
- Put the returned `family_id` into `Settings`
- Create wallets/categories/budgets using backend endpoints
- Use `Add` tab to create transactions
