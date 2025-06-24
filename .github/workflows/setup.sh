apt update
apt install -y build-essential curl git
git clone https://github.com/arp242/bsdgrep.git
cd bsdgrep/
./update.sh
sed -i 's/#error.*getprogname.*/return \"grep\";/' progname.c
sed -i 's/warnc(/warn(/g' util.c
sed -i 's/warn(p->fts_errno,/warn(/g' util.c
rm freebsd.c
make
make install
which grep
grep --version
mv /usr/bin/grep /usr/bin/grep.gnu && ln -sf /usr/local/bin/grep /usr/bin/grep
grep --version
history