# Makefile for raspi-live

ffmpeg.tar.bz2:
	sudo apt-get install libomxil-bellagio-dev -y; sudo dd if=/dev/zero of=/var/swap.img bs=1024k count=1000;
	sudo mkswap /var/swap.img;
	sudo swapon /var/swap.img;
	wget -O ffmpeg.tar.bz2 https://ffmpeg.org/releases/ffmpeg-snapshot-git.tar.bz2
		
ffmpeg-src: ffmpeg.tar.bz2
	tar xvjf ffmpeg.tar.bz2

ffmpeg: ffmpeg-src
	cd ffmpeg;
	sudo ./configure --arch=arm --target-os=linux --enable-gpl --enable-omx --enable-omx-rpi --enable-nonfree;
	sudo make -j$(grep -c ^processor /proc/cpuinfo);
	sudo make install
clean:
	sudo rm ffmpeg.tar.bz2
	sudo rm -r ffmpeg
    
