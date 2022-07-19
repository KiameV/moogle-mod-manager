build:
	go build -ldflags="-s -H=windowsgui"  -o moogle-mod-manager.exe
	upx -9 -k moogle-mod-manager.exe
	rm moogle-mod-manager.ex~
