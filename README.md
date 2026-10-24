# goboxer
[![Go](https://github.com/jparound30/goboxer/actions/workflows/go.yml/badge.svg)](https://github.com/jparound30/goboxer/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/jparound30/goboxer.svg)](https://pkg.go.dev/github.com/jparound30/goboxer)

_goboxer_ is **UNOFFICIAL** Box API(https://developer.box.com/reference) library for golang.  

_goboxer_ is UNDER DEVELOPMENT and its API may be destructively changed.

## Changelog
See [CHANGELOG.md](CHANGELOG.md). Breaking changes are listed under a **Breaking Changes** heading for each release.

## Requirements
* Go 1.26 or later

## Installation
```sh
go get github.com/jparound30/goboxer
```

## Features
* OAuth 2.0 (authorization code / refresh token) and JWT (service account / user) authentication
* Auto refreshing access_token / refresh_token
* Builtin retry process (HTTP Status Code 429 or 500+). The number of attempts is configurable via `APIConn.MaxRequestAttempts` (default: 5)
* Batch request supported.

## Usage

### OAuth 2.0
```go
apiConn := goboxer.NewAPIConnWithRefreshToken(clientID, clientSecret, accessToken, refreshToken)

folder := goboxer.NewFolder(apiConn)
info, err := folder.GetInfo("0", nil)
```
`cmd/goboxertokengen` is a helper that runs the OAuth 2.0 authorization flow locally and prints the resulting tokens.

### Downloading a file
`DownloadFile` streams the content. Close the body when done, or use `ReadBody` to load it into memory.
```go
resp, err := goboxer.NewFile(apiConn).DownloadFile(fileID, "", "")
if err != nil {
	log.Fatal(err)
}
defer resp.Body.Close()

out, _ := os.Create("downloaded.bin")
defer out.Close()
if _, err := io.Copy(out, resp.Body); err != nil {
	log.Fatal(err)
}
```

### JWT
Load the JSON configuration file downloaded from the Box Developer Console and create a connection
as the service account (enterprise) or as a managed user.
```go
configFile, err := os.Open("config.json")
if err != nil {
	log.Fatal(err)
}
defer configFile.Close()

loader := goboxer.JwtConfigDefaultLoader{}

// as service account
apiConn, err := goboxer.NewAPIConnWithJwtConfig(configFile, loader)

// as user
// apiConn, err := goboxer.NewAPIConnWithJwtConfigForUser(configFile, loader, userID)
```
See `cmd/sandbox/jwt_example.go` for a complete example.

### Custom HTTP client
Each `APIConn` sends its requests with `APIConn.HTTPClient`. When it is nil, `goboxer.DefaultHTTPClient` is used and all connections share one connection pool.
Set it to configure a proxy, TLS, timeouts, or to inject your own `http.RoundTripper` (for tests or tracing).
```go
apiConn := goboxer.NewAPIConnWithAccessToken(accessToken)
apiConn.HTTPClient = &http.Client{
	Transport: &http.Transport{
		Proxy:                 http.ProxyURL(proxyURL),
		ResponseHeaderTimeout: 30 * time.Second,
	},
}
```
Prefer timeouts on the `Transport` over `http.Client.Timeout`: the latter also limits reading the body and would interrupt large `DownloadFile` streams. Connections that should share a pool must be given the same `*http.Client`.

### NOTICE
**About features for enterprise (Retention policy, Whitelist, etc...), it will not be implemented. ( I have no environment that testing those api! )**

### Implementation status.

| Category | SubCategory | STATUS | Priority | 
----|---- |----|----|
| Files | Get File Info | supported | - |
|  | Download File | supported | - |
|  | Upload File | supported | - |
|  | Upload File Version | supported | - |
|  | Chunked Upload | not yet | High |
|  | Update File Info | supported | - |
|  | Preflight Check | supported | - |
|  | Delete File | supported | - |
|  | Copy File | supported | - |
|  | Lock and Unlock | supported | - |
|  | Get Thumbnail | not yet | Low |
|  | Get Embed Link | not yet | Low |
|  | Get File Collaborations | supported | - |
|  | Get File Comments | not yet | Low |
|  | Get File Tasks | not yet | no plan |
|  | Representations | not yet | Low |
| File Versions | Get Versions | not yet | Low |
|  | Get File Version Info | not yet | Low |
|  | Promote Version | not yet | Low |
|  | Delete Old Version | not yet | Low |
| Folders | Get Folder Info | supported | - |
|  | Get Folder Items | supported | - |
|  | Create Folder | supported | - |
|  | Update Folder | supported | - |
|  | Delete Folder | supported | - |
|  | Copy Folder | supported | - |
|  | Get Folder Collaborations | supported | - |
| File and Folder Metadata | Get Metadata Template by Name | not yet | no plan |
|  | Get Metadata Template by ID | not yet | no plan |
|  | Create Metadata Template | not yet | no plan |
|  | Update Metadata Template | not yet | no plan |
|  | Delete Metadata Template | not yet | no plan |
|  | Get Enterprise Template | not yet | no plan |
|  | Get all Metadata on File | not yet | no plan |
|  | Get Metadata on File | not yet | no plan |
|  | Create Metadata on File | not yet | no plan |
|  | Update Metadata on File | not yet | no plan |
|  | Delete Metadata on File | not yet | no plan |
|  | Get All Metadata on Folder | not yet | no plan |
|  | Get Metadata on Folder | not yet | no plan |
|  | Create Metadata on Folder | not yet | no plan |
|  | Update Metadata on Folder | not yet | no plan |
|  | Delete Metadata on Folder | not yet | no plan |
| Metadata Cascade Policy | Get Metadata Cascade Policies | not yet | no plan |
|  | Get Metadata Cascade Policy | not yet | no plan |
|  | Create Metadata Cascade Policy | not yet | no plan |
|  | Delete Metadata Cascade Policy | not yet | no plan |
|  | Force Apply Metadata Cascade Policy | not yet | no plan |
| Search | Searching for Content | not yet | Normal |
| Trash | Get Trashed Items | not yet | Low |
|  | Get Trashed Item | not yet | Low |
|  | Restore Item | not yet | Low |
|  | Permanently Delete Item | not yet | Low |
| Shared Links | Get Shared Link | supported | - |
|  | Create or Update Shared Link | supported | - |
|  | Get Shared Item | not yet | Normal |
| Web Links | Get Web Link | not yet |  no plan |
|  | Create Web Link | not yet | no plan |
|  | Update Web Link | not yet | no plan |
|  | Delete Web Link | not yet | no plan |
| Users | Get Current User | supported | - |
|  | Get User | supported | - |
|  | Get User Avatar | not yet | no plan |
|  | Create User | supported | - |
|  | Update User | supported | - |
|  | Create App User | supported | - |
|  | Delete User | supported | - |
|  | Get Enterprise Users | supported | - |
|  | Invite User | not yet | no plan |
|  | Move Owned Items | not yet | no plan |
|  | Change User's Login | not yet | Low |
|  | Get Email Aliases | not yet | Low |
|  | Create Email Alias | not yet | Low |
|  | Delete Email Alias | not yet | Low |
| Groups | Get Group | supported | - |
|  | Create Group | supported | - |
|  | Update Group | supported | - |
|  | Delete Group | supported | - |
|  | Get Enterprise Groups | supported | - |
| Groups - Membership | Get Membership | supported | - |
|  | Create Membership | supported | - |
|  | Update Membership | supported | - |
|  | Delete Membership | supported | - |
|  | Get Memberships for Group | supported | - |
|  | Get Memberships for User | supported | - |
|  | Get Collaborations for Group | supported | - |
| Collaborations | Get Collaboration | supported | - |
|  | Create Collaboration | supported | - |
|  | Update Collaboration | supported | - |
|  | Delete Collaboration | supported | - |
|  | Pending Collaborations | supported | - |
| Comments | Get Comment | not yet | Low |
|  | Create Comment | not yet | Low |
|  | Update Comment | not yet | Low |
|  | Delete Comment | not yet | Low |
| Tasks | Get Task | not yet| no plan |
|  | Create Task | not yet | no plan |
|  | Update Task | not yet | no plan |
|  | Delete Task | not yet | no plan |
|  | Get Task Assignment | not yet | no plan |
|  | Create Task Assignment | not yet | no plan |
|  | Update Task Assignment | not yet | no plan |
|  | Delete Task Assignment | not yet | no plan |
|  | Get Assignments | not yet | no plan |
| Relay Workflow | Get List of Published Templates | not yet | no plan |
|  | Get List of Relay Workflows | not yet | no plan |
|  | Launch Relay Workflow | not yet | no plan |
| Watermarking | Get Watermark on File | not yet | no plan |
|  | Apply Watermark on File | not yet | no plan |
|  | Remove Watermark on File | not yet | no plan |
|  | Get Watermark on Folder | not yet | no plan |
|  | Apply Watermark on Folder | not yet | no plan |
|  | Remove Watermark on Folder | not yet | no plan |
| Webhooks | Get Webhooks | not yet | no plan |
|  | Get Webhook | not yet | no plan |
|  | Create Weboook | not yet | no plan |
|  | Update Webhook | not yet | no plan |
|  | Delete Webhook | not yet | no plan |
| Skills | Skill Invocation | not yet | no plan |
| Events | User Events | supported | - |
|  | Enterprise Events | supported | - |
|  | Long polling | not yet | no plan |
| Collections | Get Collections | not yet | no plan |
|  | Get Collection Items | not yet | no plan |
|  | Add or Delete Items From a Collection | not yet | no plan |
| Recent Items | Get Recent Items | not yet | no plan |
| Retention Policies | Get Retention Policy | not yet | no plan |
|  | Create Retention Policy | not yet | no plan |
|  | Update Retention Policy | not yet | no plan |
|  | Get Retention Policies | not yet | no plan |
|  | Get Retention Policy Assignment | not yet | no plan |
|  | Create Retention Policy Assignment | not yet | no plan |
|  | Get Retention Policy Assignments | not yet | no plan |
|  | Get File Version Retention | not yet | no plan |
|  | Get File Version Retentions | not yet | no plan |
| Legal Hold Object | Get Legal Hold Policy | not yet | no plan |
|  | Create Legal Hold Policy | not yet | no plan |
|  | Update Legal Hold Policy | not yet | no plan |
|  | Delete Legal Hold Policy | not yet | no plan |
|  | Get Legal Hold Policies | not yet | no plan |
|  | Get Policy Assignment | not yet | no plan |
|  | Create New Policy Assignment | not yet | no plan |
|  | Delete Policy Assignment | not yet | no plan |
|  | Get Policy Assignments | not yet | no plan |
|  | Get File Version Legal Hold | not yet | no plan |
|  | Get File Version Legal Holds | not yet | no plan |
| Device Pins | Get Device Pin | not yet | no plan |
|  | Delete Device Pin | not yet | no plan |
|  | Get Enterprise Device Pins | not yet | no plan |
| Terms of Service | Get Terms of Service | not yet | no plan |
|  | Get Terms of Service by ID | not yet | no plan |
|  | Get Terms of Service ID associated with Collaboration object | not yet | no plan |
|  | Create a Terms of Service | not yet | no plan |
|  | Update a Terms of Service | not yet | no plan |
|  | Get Terms of Service User Status | not yet | no plan |
|  | Create Terms of Service User Status | not yet | no plan |
|  | Update Terms of Service User Status | not yet | no plan |
| Collaboration Whitelist | Get Collaboration Whitelist Entries | not yet | no plan |
|  | Get Collaboration Whitelist Entry by ID | not yet | no plan |
|  | Create Collaboration Whitelist Entry | not yet | no plan |
|  | Delete Collaboration Whitelist Entry | not yet | no plan |
|  | Get Collaboration Whitelist Exempt Users | not yet | no plan |
|  | Get Collaboration Whitelist Exempt Users by ID | not yet | no plan |
|  | Create Collaboration Whitelist Exempt User | not yet | no plan |
|  | Delete Collaboration Whitelist Exempt User | not yet | no plan |
| Multi-Zones | Get Storage Policy by ID | not yet | no plan |
|  | Get Storage Policies | not yet | no plan |
|  | Get Storage Policy Assignment by ID | not yet | no plan |
|  | Get Storage Policy Assignments | not yet | no plan |
|  | Create Storage Policy Assignment | not yet | no plan |
|  | Update Storage Policy Assignment | not yet | no plan |
|  | Delete Storage Policy Assignment | not yet | no plan |

please refer box's API Reference. https://developer.box.com/reference

