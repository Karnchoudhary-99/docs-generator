# User Guide

Write your own Markdown pages next to the generated API pages: pass
their directory with `--docs` and docsgen copies them into the site,
links them from the index, and adds them to `search.json`.

Search works offline — `search.json` is a plain list of title, URL,
and text entries, so any static host or tiny script can query it.
