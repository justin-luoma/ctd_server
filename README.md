**GOPATH configuration**
Two entries need to be added to the GOPATH
{path to ctd_server}/ctd_server/_lib
{path to ctd_server}/ctd_server/include

_lib is for external libraries
include is for ctd_server packages

`go get` will pull into the first path in GOPATH
`go env` will output you configured GOPATH

if adding to _lib ensure that .git and .gitignore files and directories are removed 
**DO NOT COMMIT ANY WITH A .GIT FOLDER ANYWHERE BUT ctd_server/**

**Naming conventions for CTD_Server**

Item              | Value
----------------- | ----------
Packages          | snake_case
Public Functions  | Snake_Case
Private Functions | snake_case
Structs           | PascalCase
Variables         | camelCase