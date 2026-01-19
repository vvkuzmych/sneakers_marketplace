╔══════════════════════════════════════════════════════════════════╗
║                                                                  ║
║        📝 SHELL SCRIPTS - QUICK REFERENCE                        ║
║                                                                  ║
╚══════════════════════════════════════════════════════════════════╝

📦 FILES:
  run_comparison.sh  ⭐ Main script (interactive menu)
  demo.sh           Quick demo launcher
  bench.sh          Quick benchmark launcher

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🚀 USAGE:

1) INTERACTIVE MENU:
   ./run_comparison.sh
   
   Choose:
   1) Demo - Full demonstration
   2) Benchmarks - Performance tests
   3) ALL - Demo + Benchmarks
   4) Build - Just compile

2) DIRECT COMMANDS:
   ./run_comparison.sh demo
   ./run_comparison.sh bench
   ./run_comparison.sh all
   ./run_comparison.sh build

3) SHORTCUTS:
   ./demo.sh
   ./bench.sh

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✨ FEATURES:

✅ Database connection check
✅ Automatic project build
✅ Beautiful terminal formatting
✅ Colored output (GREEN=success, RED=error, YELLOW=warning)
✅ Benchmark analysis (shows overhead %)
✅ Performance comparison summary
✅ Error handling

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📊 EXAMPLE OUTPUT:

Demo:
  🔹 Raw SQL (pgx):
  ✅ Created user ID: 42 (took 2.5ms)

  🔹 GORM:
  ✅ Created user ID: 43 (took 3.8ms)

  📊 Performance: Raw SQL 2.5ms vs GORM 3.8ms (1.5x)

Benchmarks:
  CREATE:
    Raw SQL:  1250000 ns/op
    GORM:     1875000 ns/op
    Overhead: +50%

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🛠️ TROUBLESHOOTING:

Problem: "Permission denied"
Solution: chmod +x run_comparison.sh demo.sh bench.sh

Problem: "Database connection failed"
Solution: docker-compose up -d (start PostgreSQL)

Problem: "Command not found: psql"
Solution: Script will skip DB check (not critical)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📚 MORE INFO:
  README.md - Full documentation
  ../../docs/GORM_QUICKSTART.md - Quick start guide
  ../../docs/GORM_INVESTIGATION.md - Detailed analysis

Created: 2026-01-19
