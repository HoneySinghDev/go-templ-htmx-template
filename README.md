# Go Templ HTMX Template

A streamlined Go-based template integrating Echo, Templ, HTMX, SQLC, Tailwind CSS, and Alpine.js, designed for rapid development of web applications. This template serves as a foundation for full-stack development, incorporating best practices for clean architecture, security, and performance.

## Quick Start

**Clone and Install:**

```bash
git clone github.com/HoneySinghDev/go-templ-htmx-template.git
cd go-templ-htmx-template
task install
```

**Launch Development Server:**

```bash
task dev
```

## Configuration

**Environment Variables:**

Create a `.env` file based on `.env.example`. Adjust database and Supabase settings to match your local or production environments.

**PKL Configuration:**

Settings are in `pkl/app.config`, with environment-specific configurations in `pkl/local` and `pkl/prod`.

## Building for Production

```bash
task build
```

Compiles Go binary and frontend assets. Deploy `bin/app` to your server.

## Features

- **Echo Framework**: Robust foundation for web server operations.
- **Templ Templating**: Type-safe, JSX-inspired templating for Go, enhancing server-side rendering.
- **HTMX**: Dynamic HTML content updates without full page reloads.
- **SQLC**: Type-safe SQL queries with automated code generation.
- **Tailwind CSS**: Utility-first styling for rapid UI development.
- **Alpine.js**: Lightweight JavaScript for enhanced interactivity.

- [ ] Trying Experimental NilAway: A tool to help you avoid nil pointer dereferencing in Go. [Check](https://github.com/uber-go/nilaway/issues/175)

## Contributing

Contributions are welcome. Fork this repository, make your changes, and submit a pull request.

## License

MIT License. See `LICENSE` for more details.

---

Simplify your Go web development with this template, crafted for developers by developers.
