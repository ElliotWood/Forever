# Forever talent checks

Ported from [ElliotWood/Forever](https://github.com/ElliotWood/Forever), where the talent trees were built by hand
before the beta client existed. Here the trees in `ui/sim/talents/trees` are generated from the client, so these
scripts are cross-checks, not the source of the trees.

- `export_beta.py <build> <dir> [class...]` reads the Trait tables of a client build from wago.tools and writes one
  `<class>.json` per class (cached under `<dir>/.cache`). Descriptions carry per-rank numbers taken from the
  TraitDefinitionEffectPoints curves.
- `diff_trees.py <dir-or-file>` reports what a dataset changes against `ui/sim/talents/trees`, per talent: added,
  removed, renamed, moved, rank counts, prerequisites and tooltip numbers. `--tree-schema` reads a dataset in the tree
  layout (for example the trees of another commit or another repo). `--json` for machine output. Exits 1 when
  anything changed.

Typical use after a new client build:

    tools/forever_talents/export_beta.py 1.60.1.69913 beta/
    tools/forever_talents/diff_trees.py beta/

`import_talents.py` only supplies helpers to the two scripts above. Its command line writes the old repo's tree and
proto layout; do not run it here.
