---
name: 🐞 Bug Report (Agile-style)
about: Report a bug following Agile best practices
title: "[Bug] "
labels: [bug]
---

## 🐛 Bug Summary

_Briefly describe the bug._

> Example:
> When submitting the user form without a Clerk ID, the server crashes with a panic.

---

## 🔁 Steps to Reproduce

_List minimal steps to reproduce the bug._

1. Go to `/admin/users`
2. Click "Create User"
3. Leave Clerk ID empty
4. Click "Submit"

---

## 📌 Expected Behavior

_What did you expect to happen?_

> Example:
> The user should be created without Clerk ID and no error should occur.

---

## 💥 Actual Behavior

_What actually happened?_

> Example:
> The server panics with a nil pointer exception.

---

## 🖥 Environment

- OS: (e.g., Ubuntu 22.04 / macOS 14)
- Browser/Postman: (e.g., Chrome 126 / curl)
- Backend Version: (e.g., `v0.3.1`)
- Branch: (e.g., `main` / `dev`)

---

## ✅ Acceptance Criteria

- [ ] Bug is reproducible with tests
- [ ] Missing Clerk ID does not crash server
- [ ] Proper validation or fallback applied

---

## 🧠 Possible Fix / Notes

_If you have any ideas for the fix, write here._

> Example:
> Missing null check on `clerk_id` field in `CreateUserHandler`.

---

## 🔗 Related Issues / PR

_Link related issues or PRs, e.g., #123, #45_

---

## 🏷 Labels and Priority

_(Optional: e.g., `priority: high`, `area: backend`)_
