# nute

To run nute, use the following commands:
go mod download

Build common components:
make mashupsdk

Install g3n support libraries:
sudo apt-get install xorg-dev libgl1-mesa-dev libopenal1 libopenal-dev libvorbis0a libvorbis-dev libvorbisfile3

Generate self signed certs:
./mashupsdk/tls/certs_gen.sh
mkdir examples/helloworld/hellocustos/tls
mv mashupsdk.* examples/helloworld/hellocustos/tls

Make example:
make hellocustosworld

Add nute hello custos world example to your $PATH:
examples/helloworld/bin

Run worldg3n example:
worldg3n -custos -insecure

Mac users have indicated some problems here.  You'll need these if you don't have them yet.

brew install libvorbis 
brew install openal-soft

This works:
worldg3n -toruslayout -headless

This doesn't:
worldg3n -custos -tls-skip-validation -toruslayout

Run Hello world gio:
make helloworldgio

hellogio -insecure

Run Hello World fyne:
make helloworldfyne

hellofyne -insecure

Run Hello World mobile (do not build...)
