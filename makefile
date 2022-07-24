build:
	rm -f moogle-mod-manager.zip
	go-winres make
	go build -ldflags="-s -H=windowsgui"  -o moogle-mod-manager.exe
	upx -9 -k moogle-mod-manager.exe
	rm moogle-mod-manager.ex~
	7z a -tzip moogle-mod-manager.zip  moogle-mod-manager.exe
