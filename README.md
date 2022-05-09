# nute

To run nute, use the following commands:
go mod download

Build common components:
make mashupsdk

Generate self signed certs:
./mashupsdk/tls/certs_gen.sh

Mac users have indicated some problems here.  You'll need these if you don't have them yet.
brew install libvorbis openal-soft

Run Hello world gio:
make helloworldgio

hellogio -insecure

Run Hello World fyne:
make helloworldfyne

hellofyne -insecure

Run Hello World mobile (do not build...)
