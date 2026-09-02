# Skill: Dev Workflow

## First-time setup

```bash
npm install
cp .env.example .env
# Edit .env — set VITE_API_URL to your backend base URL
```

## Local development

```bash
npm run dev
# Opens http://localhost:5173
```

## Build (verify before any commit / PR)

```bash
npm run build
# Must exit 0. Fails on TypeScript errors and Vite build errors.
```

## Preview production build

```bash
npm run preview
# Serves the dist/ folder locally
```

## Common issues

| Symptom | Fix |
|---|---|
| White screen after login | Check `VITE_API_URL` is set and backend is running |
| TS error in `apiClient` | Confirm you `await` the call before casting |
| Tailwind classes not applying | Check `tailwind.config.js` `content` includes your file |
