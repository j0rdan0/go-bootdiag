package main

var viu_binary_x86_64_linux = "https://github.com/atanunq/viu/releases/download/v1.5.1/viu-x86_64-unknown-linux-musl"
var viu_installer_windows = "https://github.com/atanunq/viu/releases/download/v1.5.1/viu-x86_64-pc-windows-msvc.exe"

type ClientSecretInfo struct {
	ClientID     string `json:"ClientID"`
	TenantID     string `json:"TenantID"`
	ClientSecret string `json:"ClientSecret`
}
