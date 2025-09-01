
we can update the extension for a new version with
<https://developer.chrome.com/docs/extensions/how-to/distribute/host-on-linux#package_through_command_line>
```ps1
chrome.exe --pack-extension=./dist --pack-extension-key=<selfhost-crx-pack-ext-key.pem>
```
NOTE: make sure to version bump, otherwise this may do nothing?

dttlc.com/installExtension/GAAFL/dist.crx

dttlc.com/install/GAAFL@latest

dttlc.com/extension/install/GAAFL
dttlc.com/extension/install/GAAFL@latest

dttlc.com/extension/install/GAAFL/dist.crx
