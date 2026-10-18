# GitHub Pages site

This folder builds the published site from the Markdown files in the repo. It holds no content of its own.

- **Engine:** [Hugo](https://gohugo.io) (extended, v0.146+), theme [Hextra](https://imfing.github.io/hextra/) as a Hugo module (`go.mod`).
- **Content:** `hugo.toml` mounts `../*.md`, `shared/` and `pilot/` as they are; `teams/` is left out until it is ready. Each folder's `README.md` becomes that section's landing page; the repo `README.md` is the home page. `tools/` is never published.
- **Sidebar labels, order and "Edit this page" targets:** the `[[cascade]]` entries in `hugo.toml`. Add one when you add a file; without it the page still appears, sorted by name, but its edit link points at the repo root.
- **Links between files:** write them as normal relative repo links (`../00-plan.md#41-cutting-review-wait`, `teams/`, `pilot/cass-operator/README.md`) so they work on GitHub. `layouts/_markup/render-link.html` maps them to the published pages, and sends links to unpublished paths such as `tools/` to GitHub. Section anchors use GitHub's heading IDs, which Hugo generates the same way.
- **Page headings:** taken from each file's own `# ...` line (`layouts/docs/` hides the theme's extra title).
- **Publishing:** `.github/workflows/pages.yml` builds and deploys on every push to `main` that touches Markdown or `site/`.

Preview locally (needs Go and Hugo extended):

```sh
cd site
hugo server
```
