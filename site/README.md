# GitHub Pages site

This folder builds the published site from the Markdown files in the repo. It holds no content of its own.

- **Engine:** [Hugo](https://gohugo.io) (extended, v0.146+), theme [Hextra](https://imfing.github.io/hextra/) as a Hugo module (`go.mod`).
- **Content:** `hugo.toml` mounts `../*.md`, `teams/`, `shared/` and `pilot/` as they are. Each folder's `README.md` becomes that section's landing page; the repo `README.md` is the home page. `tools/` is never published.
- **Sidebar labels and order:** the `[[cascade]]` entries in `hugo.toml`. Add one when you add a file; without it the page still appears, sorted by name.
- **Page headings:** taken from each file's own `# ...` line (`layouts/docs/` hides the theme's extra title).
- **Publishing:** `.github/workflows/pages.yml` builds and deploys on every push to `main` that touches Markdown or `site/`.

Preview locally (needs Go and Hugo extended):

```sh
cd site
hugo server
```
