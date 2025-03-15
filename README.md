# nute

To run nute, use the following commands:
go mod download

Build common components:
make mashupsdk

Install g3n support libraries:
sudo apt-get install xorg-dev libgl1-mesa-dev libopenal1 libopenal-dev libvorbis0a libvorbis-dev libvorbisfile3

Generate self signed certs:
cd ./mashupsdk/tls/
./certs_gen.sh
cd ../..
mkdir examples/helloworld/hellocustos/tls
mv mashupsdk.* examples/helloworld/hellocustos/tls

Make example:
make hellocustosworld

Add nute hello custos world example to your $PATH:
examples/helloworld/bin

Run worldg3n example:
worldg3n -custos -tls-skip-validation

Mac users have indicated some problems here.  You'll need these if you don't have them yet.

brew install libvorbis 
brew install openal-soft
brew install glfw

This works:
worldg3n -toruslayout -headless

This doesn't:
worldg3n -custos -tls-skip-validation -toruslayout

Run Hello world gio:
make helloworldgio

hellogio -tls-skip-validation

Run Hello World fyne:
make helloworldfyne

hellofyne -tls-skip-validation

Run Hello World mobile (do not build...)
