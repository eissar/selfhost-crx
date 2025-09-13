1. copy dist.crx into this directory

2. manually update version in main.go to target version
```ps1
rg --ignore-case "version" main.go ../gaafl/gaafl.json
```

```lua
print(vim.fn.system({ 'rg', '--ignore-case', 'version', 'main.go', '../gaafl/gaafl.json' }))
```

```ps1
git add dist.crx main.go
git commit -m 'version bump to VERSION'
```

... upload using az ...



