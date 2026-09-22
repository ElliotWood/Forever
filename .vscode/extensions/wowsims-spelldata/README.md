# WoWSims Spelldata hover

Hover a spell id in Go, JSON or TypeScript and read the client's row: what the spell is, what the
header columns say and what each effect does, humanised where the shape is known and always followed by
the client's own literal.

The ids it recognises are the ones hand-written code states: `spelldata.MustFind(11574)`,
`spelldata.Find(116)`, `spellData.Rend.ByID(11574)`, `core.ActionID{SpellID: 11574}`, an APL file's
`"spellId": 11574`, `ActionId.fromSpellId(23563)` and a TS `spellId: 23563`.

In a Go file it also answers where no id is written. The family in `spellData.Execute` hovers as the
whole ladder - every rank with the call that reaches it, then the highest rank's row - and a name bound
to one rank, `executeRank` in `var executeRank = spellData.Execute.Highest()`, hovers as the rank it
picks, wherever in the package that name is used. Any other Go name gives no hover.

## How it answers

The hover runs `go run ./tools/spelldata -json <id>` in the folder holding `go.mod`, so the row always
matches the store in the checkout being read. A file outside a Go module gets no hover, and so does an
id the store does not carry. A family asks `-family <package>/<Family>` and a bound name asks
`-expr <call> -package <package>`, both with the file's own folder as the package.

A name is looked up in the ladder picks the file's folder declares, which the extension reads out of the
`.go` files there once and rereads when one of them is saved. A name that is not one of those picks
costs nothing: no process is started.

VS Code has no way to show a hover and then fill it in, so the first hover of a session waits for the
tool to compile - a second or two. Every later one is a cached read, and the ids are cached per
repository for the session.

## Build

```
cd .vscode/extensions/wowsims-spelldata
npm install
npm run compile      # tsc, writes out/extension.js
npm test             # compiles, then runs the matcher tests under node --test
```

## Install

The folder lives under `.vscode/extensions/`, which VS Code (1.89 or later) treats as a workspace
extension: opening the repository shows an *Install Workspace Extension* prompt once, and after that
the hover loads for this workspace only. The compiled `out/extension.js` is committed so the prompt
works from a fresh clone; rebuild it with `npm run compile` after changing `src/`.

For development, open this folder in VS Code and press F5, or from a shell:

```
code --extensionDevelopmentPath=$PWD
```

To install it for good, build a `.vsix` and hand it to VS Code:

```
npx @vscode/vsce package
code --install-extension wowsims-spelldata-0.1.0.vsix
```

It compiles first and packs `package.json`, this README and `out/`. The one warning it prints - no
LICENSE beside the manifest - is expected: the repository's sits at its root.

## Settings

|                              |                                                                                                                       |
| ---------------------------- | --------------------------------------------------------------------------------------------------------------------- |
| `wowsims-spelldata.goBinary` | the `go` binary the hover runs, `go` by default. Set it to a full path where `go` is not on the PATH VS Code inherits |

Failures are written to the `WoWSims Spelldata` output channel rather than shown: a tool that exits
nonzero means no hover, and a missing `go` binary warns once per session.
