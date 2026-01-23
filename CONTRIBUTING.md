# 🤝 Contributing Guidelines

Thank you for contributing to Sneakers Marketplace! Please follow these guidelines.

---

## 📋 Quick Rules

### 1. 🚨 File Operations

**ALWAYS use `git mv` for moving/renaming files:**

```bash
# ✅ Correct
git mv file.md docs/file.md

# ❌ Wrong - loses Git history!
mv file.md docs/file.md
```

**Why?** Preserves Git history, makes `git blame` and `git log --follow` work correctly.

📖 **Full guide:** [docs/GIT_BEST_PRACTICES.md](docs/GIT_BEST_PRACTICES.md)

---

### 2. 📝 Commit Messages

Follow the [Conventional Commits](https://www.conventionalcommits.org/) format:

```bash
feat: add user authentication
fix: resolve database connection issue
docs: update README with new features
refactor: rename UserService to UserHandler
test: add tests for bidding service
chore: update dependencies
```

**Format:** `<type>: <description>`

**Types:**
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation changes
- `refactor` - Code refactoring
- `test` - Adding tests
- `chore` - Maintenance tasks
- `style` - Code style changes

---

### 3. 🔧 Development Workflow

```bash
# 1. Pull latest changes
git pull origin main

# 2. Create feature branch
git checkout -b feature/your-feature-name

# 3. Make changes
# ... code ...

# 4. Use Makefile commands
make lint           # Check linting
make test           # Run tests
make lint-fix       # Auto-fix issues

# 5. Commit changes
git add .
git commit -m "feat: add your feature"

# 6. Push and create PR
git push origin feature/your-feature-name
```

---

### 4. 🧪 Before Committing

```bash
# Always run these before committing:
make lint           # Backend linting
make test           # Run all tests

# Frontend:
cd frontend
npm run lint        # Frontend linting
npm run test        # Frontend tests
```

---

### 5. 📂 Project Structure

```
sneakers_marketplace/
├── cmd/              # Service entry points
├── internal/         # Private application code
├── pkg/              # Shared packages
├── frontend/         # React application
├── docs/             # 📚 All documentation
├── scripts/          # Helper scripts
└── migrations/       # Database migrations
```

**Rules:**
- New docs → `docs/` directory
- Backend code → `internal/` or `pkg/`
- Frontend code → `frontend/src/`
- Scripts → `scripts/`

---

### 6. 🗄️ Database Migrations

```bash
# Always use migration files
# Don't manually modify database!

# Create migration
make db-migrate

# Rollback migration
make db-rollback

# See: internal/database/migrations/
```

---

### 7. 🛠️ Use Makefile

```bash
make help           # Show all commands
make dev            # Start everything
make build          # Build all services
make test           # Run tests
make lint           # Run linters
```

📖 **Full guide:** [docs/MAKEFILE_GUIDE.md](docs/MAKEFILE_GUIDE.md)

---

## 🐛 Reporting Issues

1. Check if issue already exists
2. Use issue templates
3. Provide:
   - Clear description
   - Steps to reproduce
   - Expected vs actual behavior
   - Environment (OS, Go version, etc.)

---

## 💡 Pull Request Process

1. **Fork** the repository
2. **Create** a feature branch
3. **Make** your changes
4. **Test** thoroughly (`make test`)
5. **Lint** your code (`make lint`)
6. **Commit** with clear messages
7. **Push** to your fork
8. **Open** a Pull Request

### PR Checklist:

- [ ] Tests pass (`make test`)
- [ ] Linting passes (`make lint`)
- [ ] Documentation updated (if needed)
- [ ] Commit messages follow convention
- [ ] No merge conflicts
- [ ] Code reviewed by yourself first

---

## 📖 Documentation

- Always update docs when adding features
- Keep README.md up to date
- Add inline comments for complex logic
- Update [docs/INDEX.md](docs/INDEX.md) for new docs

---

## 🎯 Code Style

### Backend (Go)

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` (automatically via `make lint-fix`)
- Use meaningful variable names
- Add comments for exported functions

### Frontend (TypeScript/React)

- Follow [TypeScript guidelines](https://www.typescriptlang.org/docs/handbook/declaration-files/do-s-and-don-ts.html)
- Use ESLint (`npm run lint`)
- Use functional components with hooks
- Use TypeScript types (avoid `any`)

---

## ⚠️ Common Mistakes to Avoid

1. ❌ Using `mv` instead of `git mv`
2. ❌ Committing without running tests
3. ❌ Vague commit messages ("fix stuff", "update")
4. ❌ Large commits (break into smaller ones)
5. ❌ Pushing directly to `main` branch
6. ❌ Not updating documentation
7. ❌ Ignoring linter warnings

---

## 🚀 Quick Start for Contributors

```bash
# 1. Clone repository
git clone https://github.com/vvkuzmych/sneakers_marketplace
cd sneakers_marketplace

# 2. Install dependencies
make install

# 3. Setup database
make db-setup

# 4. Start development
make dev

# 5. Run tests
make test

# 6. Check linting
make lint
```

---

## 📞 Need Help?

- 📖 Check [docs/](docs/) directory
- 💬 Open a Discussion on GitHub
- 🐛 Create an Issue
- 📧 Contact maintainers

---

## 🙏 Thank You!

Your contributions make this project better! 🎉

---

**Remember:** Quality over quantity. Small, well-tested PRs are better than large, untested ones.

---

*Happy coding! 🚀*
